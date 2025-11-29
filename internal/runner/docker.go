package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// DockerRunner executes code in Docker containers
type DockerRunner struct {
	client   *client.Client
	configs  map[string]LanguageConfig
	imageDir string // Directory containing Dockerfiles
}

// NewDockerRunner creates a new Docker-based runner
func NewDockerRunner(imageDir string) (*DockerRunner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &DockerRunner{
		client:   cli,
		configs:  DefaultLanguageConfigs,
		imageDir: imageDir,
	}, nil
}

// Run executes code in a container
func (r *DockerRunner) Run(ctx context.Context, config RunConfig) (*RunResult, error) {
	langConfig, ok := r.configs[config.Language]
	if !ok {
		return &RunResult{
			Verdict: VerdictSystemError,
			Error:   fmt.Errorf("unsupported language: %s", config.Language),
		}, nil
	}

	startTime := time.Now()

	// Check if image exists
	imageExists, err := r.imageExists(ctx, langConfig.Image)
	if err != nil {
		return &RunResult{
			Verdict: VerdictSystemError,
			Error:   fmt.Errorf("failed to check image: %w", err),
		}, nil
	}

	if !imageExists {
		return &RunResult{
			Verdict: VerdictSystemError,
			Error:   fmt.Errorf("%w: %s (run 'make docker-build' first)", ErrImageNotFound, langConfig.Image),
		}, nil
	}

	// Prepare the command
	cmd := r.buildCommand(langConfig.RunCmd, "/sandbox/solution"+langConfig.FileExtension)

	// Get absolute path for source file
	absSourcePath, err := filepath.Abs(config.SourcePath)
	if err != nil {
		return &RunResult{
			Verdict: VerdictSystemError,
			Error:   fmt.Errorf("failed to get absolute path: %w", err),
		}, nil
	}

	// Container configuration
	containerConfig := &container.Config{
		Image:        langConfig.Image,
		Cmd:          cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    true,
		WorkingDir:   "/sandbox",
		// Run as non-root user (created in Dockerfile)
		User: "runner",
	}

	// Host configuration with resource limits
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   absSourcePath,
				Target:   "/sandbox/solution" + langConfig.FileExtension,
				ReadOnly: true,
			},
		},
		// Security options
		NetworkMode: "none", // No network access
		AutoRemove:  true,   // Clean up after exit
		// Resource limits
		Resources: container.Resources{
			Memory:     config.MemoryLimit,
			MemorySwap: config.MemoryLimit, // Disable swap
			CPUPeriod:  100000,
			CPUQuota:   100000, // 1 CPU
			PidsLimit:  func() *int64 { v := int64(64); return &v }(),
		},
	}

	// Create container
	resp, err := r.client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return &RunResult{
			Verdict: VerdictSystemError,
			Error:   fmt.Errorf("failed to create container: %w", err),
		}, nil
	}

	containerID := resp.ID

	// Ensure cleanup even on error
	defer func() {
		// Try to remove container (may already be removed due to AutoRemove)
		_ = r.client.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true})
	}()

	// Attach to container for stdin/stdout
	attachResp, err := r.client.ContainerAttach(ctx, containerID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return &RunResult{
			Verdict: VerdictSystemError,
			Error:   fmt.Errorf("failed to attach to container: %w", err),
		}, nil
	}
	defer attachResp.Close()

	// Start container
	if err := r.client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return &RunResult{
			Verdict: VerdictSystemError,
			Error:   fmt.Errorf("failed to start container: %w", err),
		}, nil
	}

	// Write stdin
	go func() {
		defer attachResp.CloseWrite()
		io.WriteString(attachResp.Conn, config.Stdin)
	}()

	// Create a context with timeout for the execution
	execCtx, cancel := context.WithTimeout(ctx, config.TimeLimit)
	defer cancel()

	// Read stdout and stderr
	var stdout, stderr bytes.Buffer
	outputDone := make(chan error, 1)
	go func() {
		_, err := stdcopy.StdCopy(&stdout, &stderr, attachResp.Reader)
		outputDone <- err
	}()

	// Wait for container to finish
	statusCh, errCh := r.client.ContainerWait(execCtx, containerID, container.WaitConditionNotRunning)

	var result RunResult
	result.Duration = time.Since(startTime)

	select {
	case <-execCtx.Done():
		// Timeout - kill container
		_ = r.client.ContainerKill(context.Background(), containerID, "KILL")
		result.Verdict = VerdictTimeLimitExceeded
		result.Error = ErrTimeLimitExceeded
		result.Duration = config.TimeLimit
	case err := <-errCh:
		result.Verdict = VerdictSystemError
		result.Error = fmt.Errorf("container wait error: %w", err)
	case status := <-statusCh:
		result.Duration = time.Since(startTime)
		result.ExitCode = int(status.StatusCode)

		// Wait for output to be fully read
		<-outputDone
		result.Stdout = stdout.String()
		result.Stderr = stderr.String()

		// Determine verdict based on exit code
		if status.StatusCode == 0 {
			result.Verdict = VerdictAccepted // Will be compared later
		} else if status.StatusCode == 137 { // SIGKILL (often OOM)
			result.Verdict = VerdictMemoryLimitExceeded
			result.Error = ErrMemoryLimitExceeded
		} else {
			result.Verdict = VerdictRuntimeError
			result.Error = fmt.Errorf("%w: exit code %d", ErrRuntimeError, status.StatusCode)
		}
	}

	return &result, nil
}

// Supported returns the list of supported languages
func (r *DockerRunner) Supported() []string {
	languages := make([]string, 0, len(r.configs))
	for lang := range r.configs {
		languages = append(languages, lang)
	}
	return languages
}

// Cleanup releases Docker client resources
func (r *DockerRunner) Cleanup() error {
	return r.client.Close()
}

// BuildImage builds the Docker image for a language
func (r *DockerRunner) BuildImage(ctx context.Context, language string) error {
	langConfig, ok := r.configs[language]
	if !ok {
		return fmt.Errorf("unsupported language: %s", language)
	}

	dockerfilePath := filepath.Join(r.imageDir, language, "Dockerfile")

	// Read Dockerfile
	buildCtx, err := createBuildContext(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to create build context: %w", err)
	}

	resp, err := r.client.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Tags:       []string{langConfig.Image},
		Dockerfile: "Dockerfile",
		Remove:     true,
	})
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	defer resp.Body.Close()

	// Consume the build output
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}

// imageExists checks if a Docker image exists locally
func (r *DockerRunner) imageExists(ctx context.Context, imageName string) (bool, error) {
	images, err := r.client.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == imageName {
				return true, nil
			}
		}
	}
	return false, nil
}

// buildCommand replaces placeholders in the command template
func (r *DockerRunner) buildCommand(cmdTemplate []string, sourcePath string) []string {
	cmd := make([]string, len(cmdTemplate))
	for i, part := range cmdTemplate {
		cmd[i] = strings.ReplaceAll(part, "{source}", sourcePath)
	}
	return cmd
}

// createBuildContext creates a tar archive for Docker build
func createBuildContext(dockerfilePath string) (io.Reader, error) {
	// For simplicity, we'll use the Docker CLI approach
	// In production, you'd create a proper tar archive
	// For now, we'll build from the directory
	return nil, fmt.Errorf("use 'docker build' command directly for now")
}

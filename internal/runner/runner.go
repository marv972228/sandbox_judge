// Package runner provides container-based code execution.
package runner

import (
	"context"
	"errors"
	"time"
)

// Common errors returned by runners
var (
	ErrTimeLimitExceeded   = errors.New("time limit exceeded")
	ErrMemoryLimitExceeded = errors.New("memory limit exceeded")
	ErrRuntimeError        = errors.New("runtime error")
	ErrCompilationError    = errors.New("compilation error")
	ErrImageNotFound       = errors.New("runner image not found")
)

// Verdict represents the result of a test case execution
type Verdict string

const (
	VerdictAccepted            Verdict = "AC"  // Correct answer
	VerdictWrongAnswer         Verdict = "WA"  // Incorrect output
	VerdictTimeLimitExceeded   Verdict = "TLE" // Exceeded time limit
	VerdictMemoryLimitExceeded Verdict = "MLE" // Exceeded memory limit
	VerdictRuntimeError        Verdict = "RE"  // Crashed or non-zero exit
	VerdictCompilationError    Verdict = "CE"  // Failed to compile
	VerdictSystemError         Verdict = "SE"  // Internal judge error
)

// RunConfig specifies execution parameters
type RunConfig struct {
	// Language identifier (e.g., "python", "go", "cpp")
	Language string

	// Path to the source file on the host
	SourcePath string

	// Stdin input to provide to the program
	Stdin string

	// TimeLimit is the maximum execution time
	TimeLimit time.Duration

	// MemoryLimit is the maximum memory in bytes (0 = no limit)
	MemoryLimit int64

	// WorkDir is an optional working directory inside the container
	WorkDir string
}

// RunResult contains the outcome of a code execution
type RunResult struct {
	// Stdout from the program
	Stdout string

	// Stderr from the program
	Stderr string

	// ExitCode of the program (0 = success)
	ExitCode int

	// Duration is the wall clock time taken
	Duration time.Duration

	// MemoryUsed is peak memory usage in bytes (if measurable)
	MemoryUsed int64

	// Verdict is the high-level result
	Verdict Verdict

	// Error contains any execution error
	Error error
}

// Runner defines the interface for code execution backends
type Runner interface {
	// Run executes code with the given configuration
	Run(ctx context.Context, config RunConfig) (*RunResult, error)

	// Supported returns the list of supported language identifiers
	Supported() []string

	// Cleanup releases any resources held by the runner
	Cleanup() error
}

// LanguageConfig defines how to run a specific language
type LanguageConfig struct {
	// Image is the Docker image to use
	Image string

	// CompileCmd is the command to compile (empty for interpreted languages)
	CompileCmd []string

	// RunCmd is the command template to run the code
	// Use {source} as placeholder for the source file path
	RunCmd []string

	// FileExtension is the expected source file extension
	FileExtension string
}

// DefaultLanguageConfigs provides default configurations for common languages
var DefaultLanguageConfigs = map[string]LanguageConfig{
	"python": {
		Image:         "sandbox-judge-python:latest",
		CompileCmd:    nil, // Interpreted
		RunCmd:        []string{"python3", "{source}"},
		FileExtension: ".py",
	},
	"python3": {
		Image:         "sandbox-judge-python:latest",
		CompileCmd:    nil,
		RunCmd:        []string{"python3", "{source}"},
		FileExtension: ".py",
	},
}

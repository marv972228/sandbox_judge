package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/marv972228/sandbox_judge/internal/runner"
)

// Simple test to verify DockerRunner works
func main() {
	// Get the project root (assuming we run from project root)
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// Create runner
	imageDir := filepath.Join(projectRoot, "docker")
	r, err := runner.NewDockerRunner(imageDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating runner: %v\n", err)
		os.Exit(1)
	}
	defer r.Cleanup()

	// Test 1: Correct solution
	fmt.Println("=== Test 1: Correct Solution ===")
	result := runTest(r, projectRoot, "correct.py", "2 7 11 15\n9\n", 5*time.Second)
	if result.Verdict == runner.VerdictAccepted && result.Stdout == "0 1\n" {
		fmt.Println("✅ Test 1 PASSED!")
	} else {
		fmt.Printf("❌ Test 1 FAILED! Verdict: %s, Output: %q\n", result.Verdict, result.Stdout)
	}

	// Test 2: Second test case
	fmt.Println("\n=== Test 2: Second Test Case ===")
	result = runTest(r, projectRoot, "correct.py", "3 2 4\n6\n", 5*time.Second)
	if result.Verdict == runner.VerdictAccepted && result.Stdout == "1 2\n" {
		fmt.Println("✅ Test 2 PASSED!")
	} else {
		fmt.Printf("❌ Test 2 FAILED! Verdict: %s, Output: %q\n", result.Verdict, result.Stdout)
	}

	// Test 3: TLE solution
	fmt.Println("\n=== Test 3: Time Limit Exceeded ===")
	result = runTest(r, projectRoot, "tle.py", "test\n", 2*time.Second)
	if result.Verdict == runner.VerdictTimeLimitExceeded {
		fmt.Printf("✅ Test 3 PASSED! (TLE detected in %v)\n", result.Duration)
	} else {
		fmt.Printf("❌ Test 3 FAILED! Expected TLE, got: %s\n", result.Verdict)
	}

	fmt.Println("\n=== All tests completed ===")
}

func runTest(r *runner.DockerRunner, projectRoot, solution, stdin string, timeout time.Duration) *runner.RunResult {
	config := runner.RunConfig{
		Language:    "python",
		SourcePath:  filepath.Join(projectRoot, "solutions", "two-sum", solution),
		Stdin:       stdin,
		TimeLimit:   timeout,
		MemoryLimit: 256 * 1024 * 1024, // 256MB
	}

	fmt.Printf("Running %s with timeout %v...\n", solution, timeout)

	ctx := context.Background()
	result, err := r.Run(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return &runner.RunResult{Verdict: runner.VerdictSystemError}
	}

	fmt.Printf("Verdict: %s, Duration: %v, ExitCode: %d\n", result.Verdict, result.Duration, result.ExitCode)
	if result.Stdout != "" {
		fmt.Printf("Stdout: %q\n", result.Stdout)
	}
	if result.Stderr != "" {
		fmt.Printf("Stderr: %q\n", result.Stderr)
	}

	return result
}

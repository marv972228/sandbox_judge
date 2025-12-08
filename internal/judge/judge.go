// Package judge provides the core orchestration for evaluating submissions.
package judge

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/marv972228/sandbox_judge/internal/compare"
	"github.com/marv972228/sandbox_judge/internal/problem"
	"github.com/marv972228/sandbox_judge/internal/runner"
)

// TestResult holds the outcome of a single test case
type TestResult struct {
	// TestCase is the test case that was run
	TestCase problem.TestCase

	// Verdict is the result (AC, WA, TLE, RE, etc.)
	Verdict runner.Verdict

	// Duration is how long execution took
	Duration time.Duration

	// Expected output (for display on WA)
	Expected string

	// Actual output from the submission
	Actual string

	// Error message if any
	Error string
}

// Result holds the overall outcome of judging a submission
type Result struct {
	// ProblemID is the problem that was judged
	ProblemID string

	// TestResults contains results for each test case
	TestResults []TestResult

	// FinalVerdict is the overall verdict
	FinalVerdict runner.Verdict

	// TotalDuration is the sum of all test case durations
	TotalDuration time.Duration

	// Passed is the count of AC test cases
	Passed int

	// Total is the total number of test cases
	Total int
}

// Judge orchestrates the evaluation of submissions
type Judge struct {
	problemLoader *problem.Loader
	runner        runner.Runner
	comparator    compare.Comparator
}

// Config holds configuration for the Judge
type Config struct {
	ProblemsDir string
	DockerDir   string
}

// New creates a new Judge instance
func New(cfg Config) (*Judge, error) {
	// Create problem loader
	loader := problem.NewLoader(cfg.ProblemsDir)

	// Create Docker runner
	r, err := runner.NewDockerRunner(cfg.DockerDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create runner: %w", err)
	}

	// Use default comparator
	comp := compare.NewDefaultComparator()

	return &Judge{
		problemLoader: loader,
		runner:        r,
		comparator:    comp,
	}, nil
}

// Run evaluates a submission against a problem
func (j *Judge) Run(ctx context.Context, problemID, solutionPath string) (*Result, error) {
	// Load the problem
	prob, err := j.problemLoader.Load(problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to load problem %s: %w", problemID, err)
	}

	// Load test cases
	testCases, err := j.problemLoader.LoadTestCases(problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to load test cases: %w", err)
	}

	if len(testCases) == 0 {
		return nil, fmt.Errorf("no test cases found for problem %s", problemID)
	}

	// Get absolute path for solution
	absSolutionPath, err := filepath.Abs(solutionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve solution path: %w", err)
	}

	// Detect language from file extension
	ext := filepath.Ext(solutionPath)
	language := extensionToLanguage(ext)
	if language == "" {
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	// Prepare result
	result := &Result{
		ProblemID:    problemID,
		TestResults:  make([]TestResult, 0, len(testCases)),
		FinalVerdict: runner.VerdictAccepted,
		Total:        len(testCases),
	}

	// Run each test case
	for _, tc := range testCases {
		testResult := j.runTestCase(ctx, prob, tc, absSolutionPath, language)
		result.TestResults = append(result.TestResults, testResult)
		result.TotalDuration += testResult.Duration

		if testResult.Verdict == runner.VerdictAccepted {
			result.Passed++
		} else if result.FinalVerdict == runner.VerdictAccepted {
			// First non-AC verdict becomes the final verdict
			result.FinalVerdict = testResult.Verdict
		}
	}

	return result, nil
}

// runTestCase runs a single test case and returns the result
func (j *Judge) runTestCase(ctx context.Context, prob *problem.Problem, tc problem.TestCase, solutionPath, language string) TestResult {
	// Configure the run
	cfg := runner.RunConfig{
		Language:    language,
		SourcePath:  solutionPath,
		Stdin:       tc.Input,
		TimeLimit:   time.Duration(prob.TimeLimitMS) * time.Millisecond,
		MemoryLimit: int64(prob.MemoryLimitMB) * 1024 * 1024,
	}

	// Run the solution
	runResult, err := j.runner.Run(ctx, cfg)
	if err != nil {
		return TestResult{
			TestCase: tc,
			Verdict:  runner.VerdictSystemError,
			Error:    err.Error(),
		}
	}

	// Check for non-AC verdicts from runner (TLE, MLE, RE)
	if runResult.Verdict != runner.VerdictAccepted {
		return TestResult{
			TestCase: tc,
			Verdict:  runResult.Verdict,
			Duration: runResult.Duration,
			Expected: tc.Expected,
			Actual:   runResult.Stdout,
			Error:    runResult.Stderr,
		}
	}

	// Compare output
	comparison := j.comparator.Compare(tc.Expected, runResult.Stdout)

	verdict := runner.VerdictAccepted
	if !comparison.Match {
		verdict = runner.VerdictWrongAnswer
	}

	return TestResult{
		TestCase: tc,
		Verdict:  verdict,
		Duration: runResult.Duration,
		Expected: comparison.Expected,
		Actual:   comparison.Actual,
	}
}

// RunSingleTest runs only a specific test case by number (1-indexed)
func (j *Judge) RunSingleTest(ctx context.Context, problemID, solutionPath string, testNum int) (*Result, error) {
	// Load the problem
	prob, err := j.problemLoader.Load(problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to load problem %s: %w", problemID, err)
	}

	// Load test cases
	testCases, err := j.problemLoader.LoadTestCases(problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to load test cases: %w", err)
	}

	if testNum < 1 || testNum > len(testCases) {
		return nil, fmt.Errorf("test %d does not exist (problem has %d tests)", testNum, len(testCases))
	}

	// Get absolute path for solution
	absSolutionPath, err := filepath.Abs(solutionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve solution path: %w", err)
	}

	// Detect language
	ext := filepath.Ext(solutionPath)
	language := extensionToLanguage(ext)
	if language == "" {
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	// Run single test
	tc := testCases[testNum-1]
	testResult := j.runTestCase(ctx, prob, tc, absSolutionPath, language)

	result := &Result{
		ProblemID:     problemID,
		TestResults:   []TestResult{testResult},
		FinalVerdict:  testResult.Verdict,
		TotalDuration: testResult.Duration,
		Total:         1,
	}

	if testResult.Verdict == runner.VerdictAccepted {
		result.Passed = 1
	}

	return result, nil
}

// Close releases resources held by the Judge
func (j *Judge) Close() error {
	return j.runner.Cleanup()
}

// extensionToLanguage maps file extensions to language identifiers
func extensionToLanguage(ext string) string {
	switch ext {
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".go":
		return "go"
	case ".cpp", ".cc", ".cxx":
		return "cpp"
	case ".c":
		return "c"
	case ".java":
		return "java"
	case ".rs":
		return "rust"
	default:
		return ""
	}
}

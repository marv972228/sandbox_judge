package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/marv972228/sandbox_judge/internal/judge"
	"github.com/marv972228/sandbox_judge/internal/problem"
	"github.com/marv972228/sandbox_judge/internal/runner"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"

	// Global flags
	cfgFile     string
	problemsDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "judge",
	Short: "Local code evaluation and benchmarking system",
	Long: `Sandbox Judge is a local code evaluation system inspired by LeetCode and Codeforces.

Practice coding problems locally with automated judging, multiple language support,
and performance benchmarking.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.judge.yaml)")
	rootCmd.PersistentFlags().StringVar(&problemsDir, "problems", "./problems", "path to problems directory")

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(runCmd)
}

// getLoader returns a problem loader for the configured problems directory.
func getLoader() *problem.Loader {
	return problem.NewLoader(problemsDir)
}

// versionCmd prints version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sandbox-judge %s\n", Version)
	},
}

// listCmd lists all available problems
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available problems",
	Long:  `List all problems available in the problem store, showing ID, title, difficulty, and tags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		loader := getLoader()
		problems, err := loader.ListProblems()
		if err != nil {
			return err
		}

		if len(problems) == 0 {
			fmt.Println("No problems found.")
			fmt.Printf("Add problems to: %s\n", problemsDir)
			return nil
		}

		// Print header
		fmt.Printf("%-20s %-30s %-10s %s\n", "ID", "TITLE", "DIFFICULTY", "TAGS")
		fmt.Println(strings.Repeat("-", 80))

		// Print each problem
		for _, p := range problems {
			tags := strings.Join(p.Tags, ", ")
			fmt.Printf("%-20s %-30s %-10s %s\n", p.ID, truncate(p.Title, 28), p.Difficulty, tags)
		}

		fmt.Printf("\n%d problem(s) found\n", len(problems))
		return nil
	},
}

// showCmd displays a problem description
var showCmd = &cobra.Command{
	Use:   "show <problem-id>",
	Short: "Show problem description",
	Long:  `Display the full problem description including examples, constraints, and input/output format.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		problemID := args[0]
		loader := getLoader()

		p, err := loader.Load(problemID)
		if err != nil {
			return err
		}

		// Print problem header
		fmt.Printf("# %s\n", p.Title)
		fmt.Printf("Difficulty: %s | Tags: %s\n", p.Difficulty, strings.Join(p.Tags, ", "))
		fmt.Printf("Time Limit: %dms | Memory Limit: %dMB\n", p.TimeLimitMS, p.MemoryLimitMB)
		fmt.Println()

		// Print description
		fmt.Println("## Description")
		fmt.Println()
		fmt.Println(strings.TrimSpace(p.Description))
		fmt.Println()

		// Print input/output format
		if p.InputFormat != "" {
			fmt.Println("## Input Format")
			fmt.Println()
			fmt.Println(strings.TrimSpace(p.InputFormat))
			fmt.Println()
		}

		if p.OutputFormat != "" {
			fmt.Println("## Output Format")
			fmt.Println()
			fmt.Println(strings.TrimSpace(p.OutputFormat))
			fmt.Println()
		}

		// Print constraints
		if len(p.Constraints) > 0 {
			fmt.Println("## Constraints")
			fmt.Println()
			for _, c := range p.Constraints {
				fmt.Printf("- %s\n", c)
			}
			fmt.Println()
		}

		// Print examples
		if len(p.Examples) > 0 {
			fmt.Println("## Examples")
			fmt.Println()
			for i, ex := range p.Examples {
				fmt.Printf("### Example %d\n", i+1)
				fmt.Println()
				fmt.Println("**Input:**")
				fmt.Println("```")
				fmt.Print(strings.TrimSpace(ex.Input))
				fmt.Println()
				fmt.Println("```")
				fmt.Println()
				fmt.Println("**Output:**")
				fmt.Println("```")
				fmt.Print(strings.TrimSpace(ex.Output))
				fmt.Println()
				fmt.Println("```")
				if ex.Explanation != "" {
					fmt.Println()
					fmt.Printf("**Explanation:** %s\n", ex.Explanation)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

// runCmd runs a solution against a problem
var runCmd = &cobra.Command{
	Use:   "run <problem-id> <solution-file>",
	Short: "Run solution against problem",
	Long: `Execute your solution against all test cases for the specified problem.

The solution will be run in a sandboxed container with resource limits.
Results show verdict (AC/WA/TLE/RE) and timing for each test case.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		problemID := args[0]
		solutionFile := args[1]
		verbose, _ := cmd.Flags().GetBool("verbose")
		testNum, _ := cmd.Flags().GetInt("test")

		// Verify solution file exists
		if _, err := os.Stat(solutionFile); os.IsNotExist(err) {
			return fmt.Errorf("solution file not found: %s", solutionFile)
		}

		// Get absolute path for problems dir
		absProblemDir, err := filepath.Abs(problemsDir)
		if err != nil {
			return fmt.Errorf("failed to resolve problems directory: %w", err)
		}

		// Get docker directory (relative to problems dir's parent)
		dockerDir := filepath.Join(filepath.Dir(absProblemDir), "docker")

		// Create judge
		j, err := judge.New(judge.Config{
			ProblemsDir: absProblemDir,
			DockerDir:   dockerDir,
		})
		if err != nil {
			return fmt.Errorf("failed to create judge: %w", err)
		}
		defer j.Close()

		// Create context with timeout (5 minutes max for entire run)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		fmt.Printf("Running %s...\n", problemID)

		var result *judge.Result
		if testNum > 0 {
			result, err = j.RunSingleTest(ctx, problemID, solutionFile, testNum)
		} else {
			result, err = j.Run(ctx, problemID, solutionFile)
		}
		if err != nil {
			return err
		}

		// Print results
		for i, tr := range result.TestResults {
			testName := tr.TestCase.Name
			if testName == "" {
				testName = fmt.Sprintf("Test %d", i+1)
			}

			verdictStr := colorVerdict(tr.Verdict)
			fmt.Printf("  %s: %s (%v)\n", testName, verdictStr, tr.Duration.Round(time.Millisecond))

			// Show diff on WA if verbose
			if verbose && tr.Verdict == runner.VerdictWrongAnswer {
				fmt.Println("    Expected:")
				for _, line := range strings.Split(tr.Expected, "\n") {
					fmt.Printf("      %s\n", line)
				}
				fmt.Println("    Actual:")
				for _, line := range strings.Split(tr.Actual, "\n") {
					fmt.Printf("      %s\n", line)
				}
			}

			// Show error on RE/TLE/SE if verbose
			if verbose && tr.Error != "" {
				fmt.Printf("    Error: %s\n", tr.Error)
			}
		}

		// Print summary
		fmt.Println()
		summaryVerdict := colorVerdict(result.FinalVerdict)
		if result.FinalVerdict == runner.VerdictAccepted {
			fmt.Printf("Result: %s (%d/%d tests passed)\n", summaryVerdict, result.Passed, result.Total)
		} else {
			fmt.Printf("Result: %s (%d/%d tests passed)\n", summaryVerdict, result.Passed, result.Total)
		}
		fmt.Printf("Total time: %v\n", result.TotalDuration.Round(time.Millisecond))

		return nil
	},
}

// colorVerdict returns a colored verdict string
func colorVerdict(v runner.Verdict) string {
	// ANSI color codes
	const (
		green  = "\033[32m"
		red    = "\033[31m"
		yellow = "\033[33m"
		reset  = "\033[0m"
	)

	switch v {
	case runner.VerdictAccepted:
		return green + "AC" + reset
	case runner.VerdictWrongAnswer:
		return red + "WA" + reset
	case runner.VerdictTimeLimitExceeded:
		return yellow + "TLE" + reset
	case runner.VerdictMemoryLimitExceeded:
		return yellow + "MLE" + reset
	case runner.VerdictRuntimeError:
		return red + "RE" + reset
	case runner.VerdictCompilationError:
		return red + "CE" + reset
	case runner.VerdictSystemError:
		return red + "SE" + reset
	default:
		return string(v)
	}
}

func init() {
	// Flags for run command
	runCmd.Flags().BoolP("verbose", "v", false, "Show detailed output including input/output diff on failure")
	runCmd.Flags().IntP("test", "t", 0, "Run only a specific test case (0 = all)")
	runCmd.Flags().Duration("timeout", 0, "Override the problem's time limit")
}

// truncate shortens a string to maxLen, adding "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// init runs before main
func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		// TODO: Load config from file
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", cfgFile)
	}
}

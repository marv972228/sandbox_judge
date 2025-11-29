package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"

	// Global flags
	cfgFile string
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

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(runCmd)
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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: List problems")
	},
}

// showCmd displays a problem description
var showCmd = &cobra.Command{
	Use:   "show <problem-id>",
	Short: "Show problem description",
	Long:  `Display the full problem description including examples, constraints, and input/output format.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		problemID := args[0]
		fmt.Printf("TODO: Show problem %s\n", problemID)
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
	Run: func(cmd *cobra.Command, args []string) {
		problemID := args[0]
		solutionFile := args[1]
		verbose, _ := cmd.Flags().GetBool("verbose")
		testNum, _ := cmd.Flags().GetInt("test")

		if verbose {
			fmt.Println("Verbose mode enabled")
		}
		if testNum > 0 {
			fmt.Printf("Running only test case %d\n", testNum)
		}

		fmt.Printf("TODO: Run %s against %s\n", solutionFile, problemID)
	},
}

func init() {
	// Flags for run command
	runCmd.Flags().BoolP("verbose", "v", false, "Show detailed output including input/output diff on failure")
	runCmd.Flags().IntP("test", "t", 0, "Run only a specific test case (0 = all)")
	runCmd.Flags().Duration("timeout", 0, "Override the problem's time limit")
}

// SetVersion allows setting version from main (via ldflags)
func SetVersion(v string) {
	Version = v
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

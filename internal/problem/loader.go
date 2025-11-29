// Package problem handles loading and parsing problem definitions.
package problem

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader handles loading problems from a directory.
type Loader struct {
	problemsDir string
}

// NewLoader creates a new problem loader for the given directory.
func NewLoader(problemsDir string) *Loader {
	return &Loader{problemsDir: problemsDir}
}

// Load reads a problem by ID from the problems directory.
func (l *Loader) Load(id string) (*Problem, error) {
	problemDir := filepath.Join(l.problemsDir, id)

	// Check if problem directory exists
	if _, err := os.Stat(problemDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("problem not found: %s", id)
	}

	// Load problem.yaml
	yamlPath := filepath.Join(problemDir, "problem.yaml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read problem.yaml: %w", err)
	}

	var problem Problem
	if err := yaml.Unmarshal(data, &problem); err != nil {
		return nil, fmt.Errorf("failed to parse problem.yaml: %w", err)
	}

	// Set defaults
	problem.Defaults()

	// Ensure ID matches directory name
	if problem.ID == "" {
		problem.ID = id
	}

	return &problem, nil
}

// LoadTestCases loads all test cases for a problem.
func (l *Loader) LoadTestCases(id string) ([]TestCase, error) {
	problemDir := filepath.Join(l.problemsDir, id)
	testsDir := filepath.Join(problemDir, "tests")

	var testCases []TestCase

	// Load sample tests
	sampleTests, err := l.loadTestsFromDir(filepath.Join(testsDir, "sample"), "sample")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	testCases = append(testCases, sampleTests...)

	// Load hidden tests
	hiddenTests, err := l.loadTestsFromDir(filepath.Join(testsDir, "hidden"), "hidden")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	testCases = append(testCases, hiddenTests...)

	if len(testCases) == 0 {
		return nil, fmt.Errorf("no test cases found for problem: %s", id)
	}

	return testCases, nil
}

// loadTestsFromDir loads test cases from a directory (sample or hidden).
func (l *Loader) loadTestsFromDir(dir, prefix string) ([]TestCase, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Find all .in files
	var testCases []TestCase
	inFiles := make(map[string]bool)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".in") {
			baseName := strings.TrimSuffix(name, ".in")
			inFiles[baseName] = true
		}
	}

	// Sort test names for consistent ordering
	var names []string
	for name := range inFiles {
		names = append(names, name)
	}
	sort.Strings(names)

	// Load each test case
	for _, name := range names {
		inPath := filepath.Join(dir, name+".in")
		outPath := filepath.Join(dir, name+".out")

		input, err := os.ReadFile(inPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", inPath, err)
		}

		expected, err := os.ReadFile(outPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", outPath, err)
		}

		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("%s/%s", prefix, name),
			Input:    string(input),
			Expected: string(expected),
		})
	}

	return testCases, nil
}

// List returns all available problem IDs.
func (l *Loader) List() ([]string, error) {
	entries, err := os.ReadDir(l.problemsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read problems directory: %w", err)
	}

	var ids []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Check if it has a problem.yaml
		yamlPath := filepath.Join(l.problemsDir, entry.Name(), "problem.yaml")
		if _, err := os.Stat(yamlPath); err == nil {
			ids = append(ids, entry.Name())
		}
	}

	sort.Strings(ids)
	return ids, nil
}

// ListProblems returns all available problems with metadata.
func (l *Loader) ListProblems() ([]*Problem, error) {
	ids, err := l.List()
	if err != nil {
		return nil, err
	}

	var problems []*Problem
	for _, id := range ids {
		p, err := l.Load(id)
		if err != nil {
			// Log warning but continue
			fmt.Fprintf(os.Stderr, "Warning: failed to load problem %s: %v\n", id, err)
			continue
		}
		problems = append(problems, p)
	}

	return problems, nil
}

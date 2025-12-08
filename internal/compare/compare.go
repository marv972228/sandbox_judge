// Package compare provides output comparison strategies.
package compare

import (
	"strings"
)

// Result represents the outcome of a comparison
type Result struct {
	// Match indicates if the outputs are considered equal
	Match bool

	// Expected is the expected output (possibly normalized)
	Expected string

	// Actual is the actual output (possibly normalized)
	Actual string

	// DiffLine is the first line number where difference was found (1-indexed, 0 if match)
	DiffLine int

	// DiffExpected is the expected content at the diff line
	DiffExpected string

	// DiffActual is the actual content at the diff line
	DiffActual string
}

// Comparator defines the interface for output comparison strategies
type Comparator interface {
	// Compare compares expected and actual outputs
	Compare(expected, actual string) Result

	// Name returns the comparator's identifier
	Name() string
}

// DefaultComparator performs whitespace-tolerant comparison
// - Trims leading/trailing whitespace from each line
// - Normalizes line endings (CRLF -> LF)
// - Ignores trailing blank lines
type DefaultComparator struct{}

// NewDefaultComparator creates a new DefaultComparator
func NewDefaultComparator() *DefaultComparator {
	return &DefaultComparator{}
}

// Name returns the comparator's identifier
func (c *DefaultComparator) Name() string {
	return "default"
}

// Compare performs whitespace-tolerant comparison
func (c *DefaultComparator) Compare(expected, actual string) Result {
	// Normalize both strings
	expectedLines := c.normalizeLines(expected)
	actualLines := c.normalizeLines(actual)

	// Compare line by line
	maxLines := len(expectedLines)
	if len(actualLines) > maxLines {
		maxLines = len(actualLines)
	}

	for i := 0; i < maxLines; i++ {
		var expLine, actLine string

		if i < len(expectedLines) {
			expLine = expectedLines[i]
		}
		if i < len(actualLines) {
			actLine = actualLines[i]
		}

		if expLine != actLine {
			return Result{
				Match:        false,
				Expected:     strings.Join(expectedLines, "\n"),
				Actual:       strings.Join(actualLines, "\n"),
				DiffLine:     i + 1,
				DiffExpected: expLine,
				DiffActual:   actLine,
			}
		}
	}

	return Result{
		Match:    true,
		Expected: strings.Join(expectedLines, "\n"),
		Actual:   strings.Join(actualLines, "\n"),
	}
}

// normalizeLines splits input into lines, trims each line, and removes trailing empty lines
func (c *DefaultComparator) normalizeLines(s string) []string {
	// Normalize line endings (CRLF -> LF)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	// Split into lines
	lines := strings.Split(s, "\n")

	// Trim each line
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}

	// Remove trailing empty lines
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}

// StrictComparator performs exact byte-for-byte comparison
type StrictComparator struct{}

// NewStrictComparator creates a new StrictComparator
func NewStrictComparator() *StrictComparator {
	return &StrictComparator{}
}

// Name returns the comparator's identifier
func (c *StrictComparator) Name() string {
	return "strict"
}

// Compare performs exact comparison
func (c *StrictComparator) Compare(expected, actual string) Result {
	if expected == actual {
		return Result{
			Match:    true,
			Expected: expected,
			Actual:   actual,
		}
	}

	// Find first difference
	expLines := strings.Split(expected, "\n")
	actLines := strings.Split(actual, "\n")

	maxLines := len(expLines)
	if len(actLines) > maxLines {
		maxLines = len(actLines)
	}

	for i := 0; i < maxLines; i++ {
		var expLine, actLine string
		if i < len(expLines) {
			expLine = expLines[i]
		}
		if i < len(actLines) {
			actLine = actLines[i]
		}

		if expLine != actLine {
			return Result{
				Match:        false,
				Expected:     expected,
				Actual:       actual,
				DiffLine:     i + 1,
				DiffExpected: expLine,
				DiffActual:   actLine,
			}
		}
	}

	// Shouldn't reach here if strings differ, but just in case
	return Result{
		Match:    false,
		Expected: expected,
		Actual:   actual,
		DiffLine: 1,
	}
}

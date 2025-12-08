package compare

import (
	"testing"
)

func TestDefaultComparator_ExactMatch(t *testing.T) {
	c := NewDefaultComparator()

	result := c.Compare("hello\nworld", "hello\nworld")
	if !result.Match {
		t.Errorf("Expected match, got mismatch at line %d", result.DiffLine)
	}
}

func TestDefaultComparator_TrailingNewline(t *testing.T) {
	c := NewDefaultComparator()

	// Should match even with different trailing newlines
	result := c.Compare("hello\nworld\n", "hello\nworld")
	if !result.Match {
		t.Errorf("Expected match (trailing newline), got mismatch at line %d", result.DiffLine)
	}

	result = c.Compare("hello\nworld", "hello\nworld\n")
	if !result.Match {
		t.Errorf("Expected match (trailing newline reversed), got mismatch at line %d", result.DiffLine)
	}

	result = c.Compare("hello\nworld\n\n\n", "hello\nworld")
	if !result.Match {
		t.Errorf("Expected match (multiple trailing newlines), got mismatch at line %d", result.DiffLine)
	}
}

func TestDefaultComparator_LeadingTrailingWhitespace(t *testing.T) {
	c := NewDefaultComparator()

	// Should match with whitespace differences
	result := c.Compare("  hello  \n  world  ", "hello\nworld")
	if !result.Match {
		t.Errorf("Expected match (whitespace trimmed), got mismatch at line %d", result.DiffLine)
	}

	result = c.Compare("hello\nworld", "  hello  \n  world  ")
	if !result.Match {
		t.Errorf("Expected match (whitespace trimmed reversed), got mismatch at line %d", result.DiffLine)
	}
}

func TestDefaultComparator_CRLFNormalization(t *testing.T) {
	c := NewDefaultComparator()

	// Should match with different line endings
	result := c.Compare("hello\r\nworld\r\n", "hello\nworld\n")
	if !result.Match {
		t.Errorf("Expected match (CRLF normalized), got mismatch at line %d", result.DiffLine)
	}

	result = c.Compare("hello\rworld\r", "hello\nworld")
	if !result.Match {
		t.Errorf("Expected match (CR normalized), got mismatch at line %d", result.DiffLine)
	}
}

func TestDefaultComparator_Mismatch(t *testing.T) {
	c := NewDefaultComparator()

	result := c.Compare("hello\nworld", "hello\nearth")
	if result.Match {
		t.Error("Expected mismatch, got match")
	}
	if result.DiffLine != 2 {
		t.Errorf("Expected diff at line 2, got line %d", result.DiffLine)
	}
	if result.DiffExpected != "world" {
		t.Errorf("Expected DiffExpected='world', got '%s'", result.DiffExpected)
	}
	if result.DiffActual != "earth" {
		t.Errorf("Expected DiffActual='earth', got '%s'", result.DiffActual)
	}
}

func TestDefaultComparator_DifferentLineCount(t *testing.T) {
	c := NewDefaultComparator()

	// Actual has fewer lines
	result := c.Compare("hello\nworld\nfoo", "hello\nworld")
	if result.Match {
		t.Error("Expected mismatch (fewer lines), got match")
	}
	if result.DiffLine != 3 {
		t.Errorf("Expected diff at line 3, got line %d", result.DiffLine)
	}

	// Actual has more lines
	result = c.Compare("hello\nworld", "hello\nworld\nfoo")
	if result.Match {
		t.Error("Expected mismatch (more lines), got match")
	}
	if result.DiffLine != 3 {
		t.Errorf("Expected diff at line 3, got line %d", result.DiffLine)
	}
}

func TestDefaultComparator_EmptyStrings(t *testing.T) {
	c := NewDefaultComparator()

	result := c.Compare("", "")
	if !result.Match {
		t.Error("Expected match for empty strings")
	}

	result = c.Compare("\n\n\n", "")
	if !result.Match {
		t.Error("Expected match (blank lines vs empty)")
	}

	result = c.Compare("", "\n\n\n")
	if !result.Match {
		t.Error("Expected match (empty vs blank lines)")
	}
}

func TestDefaultComparator_SingleLine(t *testing.T) {
	c := NewDefaultComparator()

	result := c.Compare("42", "42")
	if !result.Match {
		t.Error("Expected match for single line")
	}

	result = c.Compare("42\n", "42")
	if !result.Match {
		t.Error("Expected match for single line with trailing newline")
	}

	result = c.Compare("  42  ", "42")
	if !result.Match {
		t.Error("Expected match for single line with whitespace")
	}
}

func TestDefaultComparator_Numbers(t *testing.T) {
	c := NewDefaultComparator()

	// Typical judge output
	result := c.Compare("0 1\n", "0 1")
	if !result.Match {
		t.Error("Expected match for number output")
	}

	result = c.Compare("0 1", "0  1") // Different spacing
	if result.Match {
		t.Error("Expected mismatch for different internal spacing")
	}
}

func TestStrictComparator_ExactMatch(t *testing.T) {
	c := NewStrictComparator()

	result := c.Compare("hello\nworld", "hello\nworld")
	if !result.Match {
		t.Error("Expected match for exact strings")
	}
}

func TestStrictComparator_TrailingNewlineMismatch(t *testing.T) {
	c := NewStrictComparator()

	result := c.Compare("hello\nworld\n", "hello\nworld")
	if result.Match {
		t.Error("Expected mismatch for different trailing newline (strict)")
	}
}

func TestStrictComparator_WhitespaceMismatch(t *testing.T) {
	c := NewStrictComparator()

	result := c.Compare("hello ", "hello")
	if result.Match {
		t.Error("Expected mismatch for trailing space (strict)")
	}
}

func TestComparator_Name(t *testing.T) {
	def := NewDefaultComparator()
	if def.Name() != "default" {
		t.Errorf("Expected name 'default', got '%s'", def.Name())
	}

	strict := NewStrictComparator()
	if strict.Name() != "strict" {
		t.Errorf("Expected name 'strict', got '%s'", strict.Name())
	}
}

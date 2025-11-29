package problem

// Problem represents a coding problem with metadata and test cases.
type Problem struct {
	ID          string     `yaml:"id"`
	Title       string     `yaml:"title"`
	Difficulty  Difficulty `yaml:"difficulty"`
	Tags        []string   `yaml:"tags"`
	Description string     `yaml:"description"`

	// I/O format documentation
	InputFormat  string `yaml:"input_format"`
	OutputFormat string `yaml:"output_format"`

	// Constraints and limits
	Constraints   []string `yaml:"constraints"`
	TimeLimitMS   int      `yaml:"time_limit_ms"`
	MemoryLimitMB int      `yaml:"memory_limit_mb"`

	// Comparison settings
	Comparison     ComparisonMode `yaml:"comparison"`
	FloatTolerance float64        `yaml:"float_tolerance,omitempty"`
	Comparator     string         `yaml:"comparator,omitempty"`

	// Examples shown in problem description
	Examples []Example `yaml:"examples"`
}

// Difficulty represents problem difficulty level.
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

// ComparisonMode defines how output is compared against expected.
type ComparisonMode string

const (
	CompareDefault   ComparisonMode = "default"   // Whitespace-tolerant
	CompareStrict    ComparisonMode = "strict"    // Exact match
	CompareFloat     ComparisonMode = "float"     // Floating point tolerance
	CompareUnordered ComparisonMode = "unordered" // Order-independent lines
	CompareCustom    ComparisonMode = "custom"    // Custom comparator script
)

// Example represents a sample input/output pair shown in the problem description.
type Example struct {
	Input       string `yaml:"input"`
	Output      string `yaml:"output"`
	Explanation string `yaml:"explanation,omitempty"`
}

// TestCase represents a single test case with input and expected output.
type TestCase struct {
	Name     string // e.g., "sample/1" or "hidden/edge_case"
	Input    string
	Expected string
}

// Defaults sets default values for optional fields.
func (p *Problem) Defaults() {
	if p.TimeLimitMS == 0 {
		p.TimeLimitMS = 1000 // 1 second default
	}
	if p.MemoryLimitMB == 0 {
		p.MemoryLimitMB = 256 // 256 MB default
	}
	if p.Comparison == "" {
		p.Comparison = CompareDefault
	}
}

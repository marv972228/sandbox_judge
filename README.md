# Sandbox Judge

A self-hosted code evaluation and benchmarking system inspired by LeetCode and Codeforces. Practice coding problems locally with automated judging, multiple language support, and performance measurement.

> 🚧 **Work in Progress** - This project is under active development. See [TASKS.md](TASKS.md) for current progress.

## Features (Planned)

- 📝 **Problem Library** - Curated coding challenges organized by topic and difficulty
- 🐳 **Containerized Execution** - Safe, isolated code execution via Docker
- 🌐 **Multi-Language Support** - Python, JavaScript, Go, C++, and more
- ⏱️ **Performance Metrics** - Execution time and memory usage tracking
- 🖥️ **CLI & Web UI** - Use from terminal or browser
- 📊 **Progress Tracking** - Track your solving history and stats

## Requirements

- **Go 1.21+** - For building the CLI
- **Docker** - For sandboxed code execution (coming soon)
- **Make** - For build commands

## Quick Start

```bash
# Clone the repository
git clone https://github.com/marv972228/sandbox_judge.git
cd sandbox_judge

# Build the CLI
make build

# Verify installation
./bin/judge version
./bin/judge --help
```

## Usage

```bash
# List available problems
judge list

# View a problem description
judge show <problem-id>

# Run your solution against a problem
judge run <problem-id> <solution-file>

# Run with options
judge run two-sum solution.py --verbose --test 1
```

### Available Commands

| Command | Description |
|---------|-------------|
| `judge list` | List all available problems |
| `judge show <id>` | Display problem description |
| `judge run <id> <file>` | Run solution against test cases |
| `judge version` | Print version information |
| `judge completion` | Generate shell autocompletion |

### Run Command Flags

| Flag | Description |
|------|-------------|
| `-v, --verbose` | Show detailed output including diffs |
| `-t, --test <n>` | Run only a specific test case |
| `--timeout <duration>` | Override time limit |

## Development

```bash
# Build
make build

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean

# See all targets
make help
```

## Project Structure

```
sandbox_judge/
├── cmd/judge/          # CLI application
├── internal/           # Private packages
│   ├── judge/          # Core judging logic
│   ├── runner/         # Container execution
│   ├── problem/        # Problem loading
│   ├── compare/        # Output comparison
│   └── storage/        # Data persistence
├── pkg/api/            # Public API types
├── problems/           # Problem definitions
├── docker/             # Language runner images
└── web/                # Web UI (future)
```

## Documentation

- [DESIGN.md](DESIGN.md) - Architecture and design decisions
- [TASKS.md](TASKS.md) - Development roadmap and progress
- [IDEAS.md](IDEAS.md) - Future ideas and vision

## License

MIT

## Contributing

Contributions welcome! Please read the design docs first to understand the architecture.
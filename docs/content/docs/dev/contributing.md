# Contributing

Thank you for your interest in contributing to Sandbox Judge!

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/sandbox_judge.git
   cd sandbox_judge
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build and test:
   ```bash
   make build
   make test
   make docker-build
   ```

## Development Workflow

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test -v ./internal/compare/...
```

### Building

```bash
# Build the CLI binary
make build

# Build Docker images
make docker-build
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `go vet` before committing

```bash
# Format code
gofmt -w .

# Check for issues
go vet ./...
```

## Project Structure

```
sandbox_judge/
├── cmd/judge/       # CLI application
├── internal/        # Internal packages
│   ├── judge/       # Core orchestration
│   ├── problem/     # Problem loading
│   ├── runner/      # Code execution
│   └── compare/     # Output comparison
├── docker/          # Dockerfiles
├── problems/        # Problem library
└── docs/            # Documentation
```

## Making Changes

### Adding a New Feature

1. Create a branch: `git checkout -b feature/my-feature`
2. Make your changes
3. Add tests
4. Update documentation if needed
5. Run tests: `make test`
6. Commit with a clear message
7. Push and create a Pull Request

### Adding a New Problem

See [Creating Problems](../problems/creating.md) for the problem format.

1. Create directory: `problems/your-problem/`
2. Add `problem.yaml` with metadata
3. Add test cases in `tests/sample/` and `tests/hidden/`
4. Test with a reference solution
5. Submit a PR

### Adding a New Language

1. Create Dockerfile: `docker/language/Dockerfile`
2. Add language config in `internal/runner/runner.go`
3. Update `Makefile` docker-build target
4. Add tests
5. Update documentation

## Pull Request Guidelines

- Keep PRs focused on a single change
- Write clear commit messages
- Include tests for new functionality
- Update documentation as needed
- Ensure all tests pass

## Reporting Issues

When reporting bugs, please include:

- Go version (`go version`)
- Docker version (`docker --version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback

## Questions?

Open an issue or start a discussion on GitHub.

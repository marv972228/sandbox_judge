# Architecture Overview

This document describes the high-level architecture of Sandbox Judge.

## System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      Sandbox Judge                          │
├─────────────────────────────────────────────────────────────┤
│  CLI (cmd/judge)                                            │
│    └── Cobra commands: run, list, show                      │
├─────────────────────────────────────────────────────────────┤
│  Judge Engine (internal/judge)                              │
│    └── Orchestrates: load problem → run tests → compare     │
├─────────────────────────────────────────────────────────────┤
│  Components                                                 │
│    ├── Problem Loader (internal/problem)                    │
│    │     └── Parses YAML, loads test cases                  │
│    ├── Docker Runner (internal/runner)                      │
│    │     └── Executes code in containers                    │
│    └── Comparator (internal/compare)                        │
│          └── Compares expected vs actual output             │
├─────────────────────────────────────────────────────────────┤
│  Storage                                                    │
│    ├── problems/           # Problem definitions            │
│    └── data/               # Submissions, progress          │
└─────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Docker Containers                         │
│    sandbox-judge-python:latest                              │
│    (sandbox-judge-go:latest, etc. - planned)                │
└─────────────────────────────────────────────────────────────┘
```

## Package Structure

```
sandbox_judge/
├── cmd/
│   └── judge/          # CLI entry point
│       ├── main.go
│       └── cmd/        # Cobra commands
│           └── root.go
├── internal/
│   ├── judge/          # Core orchestration
│   │   └── judge.go
│   ├── problem/        # Problem loading
│   │   ├── types.go
│   │   └── loader.go
│   ├── runner/         # Code execution
│   │   ├── runner.go   # Interface
│   │   └── docker.go   # Docker implementation
│   └── compare/        # Output comparison
│       └── compare.go
├── docker/             # Dockerfiles for each language
│   └── python/
│       └── Dockerfile
├── problems/           # Problem library
│   └── two-sum/
│       ├── problem.yaml
│       └── tests/
└── docs/               # Documentation (MkDocs)
```

## Core Components

### CLI (cmd/judge)

The command-line interface built with [Cobra](https://github.com/spf13/cobra).

- **Responsibility:** Parse commands, invoke judge engine, display results
- **Key Files:** `cmd/judge/cmd/root.go`

### Judge Engine (internal/judge)

The central orchestrator that ties everything together.

- **Responsibility:** Load problem, run each test case, aggregate results
- **Key Interface:**
  ```go
  type Judge struct {
      problemLoader *problem.Loader
      runner        runner.Runner
      comparator    compare.Comparator
  }

  func (j *Judge) Run(ctx context.Context, problemID, solutionPath string) (*Result, error)
  ```

### Problem Loader (internal/problem)

Loads problem definitions from YAML files.

- **Responsibility:** Parse `problem.yaml`, load test cases from files
- **Key Types:**
  ```go
  type Problem struct {
      ID            string
      Title         string
      Difficulty    string
      Tags          []string
      TimeLimitMS   int
      MemoryLimitMB int
      Description   string
      // ...
  }

  type TestCase struct {
      Name     string
      Input    string
      Expected string
  }
  ```

### Docker Runner (internal/runner)

Executes user code safely in Docker containers.

- **Responsibility:** Create container, bind-mount source, run with limits
- **Key Interface:**
  ```go
  type Runner interface {
      Run(ctx context.Context, config RunConfig) (*RunResult, error)
      Supported() []string
      Cleanup() error
  }
  ```

### Comparator (internal/compare)

Compares expected output against actual output.

- **Responsibility:** Determine if output is correct (handling whitespace, etc.)
- **Key Types:**
  ```go
  type Comparator interface {
      Compare(expected, actual string) *CompareResult
  }

  // Implementations:
  // - DefaultComparator: whitespace-tolerant
  // - StrictComparator: exact match
  ```

## Execution Flow

When you run `judge run two-sum solution.py`:

```
1. CLI parses arguments
   └── problem-id: "two-sum", solution: "solution.py"

2. Judge.Run() is called
   │
   ├── 3. Problem Loader loads "two-sum"
   │      └── Reads problem.yaml, loads test cases
   │
   ├── 4. For each test case:
   │      │
   │      ├── 5. Runner executes solution
   │      │      ├── Create Docker container
   │      │      ├── Mount solution.py → /sandbox/solution.py
   │      │      ├── Run: python3 /sandbox/solution.py
   │      │      ├── Pipe test input to stdin
   │      │      ├── Capture stdout/stderr
   │      │      └── Enforce time/memory limits
   │      │
   │      └── 6. Comparator checks output
   │             └── Compare expected vs actual
   │
   └── 7. Aggregate results (AC/WA/TLE/RE)

8. CLI displays results with colors
```

## Security Model

Code execution is sandboxed using Docker:

| Protection | Implementation |
|------------|----------------|
| Network isolation | `NetworkMode: "none"` |
| Memory limits | `--memory` flag |
| Time limits | Context timeout + container stop |
| Process limits | `--pids-limit` |
| Read-only source | Mount with `ReadOnly: true` |
| Non-root execution | Container runs as `runner` user |
| Auto-cleanup | `AutoRemove: true` |

## Future Architecture (Web UI)

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Browser    │────▶│  API Server  │────▶│    Judge     │
│   (React)    │◀────│   (Go HTTP)  │◀────│   Engine     │
└──────────────┘     └──────────────┘     └──────────────┘
```

The web UI will add:
- REST API layer (`web/server/`)
- React frontend (`web/frontend/`)
- WebSocket for live updates (future)

# Sandbox Judge - Design Document

## Introduction

**Sandbox Judge** is a local code evaluation and benchmarking system inspired by online competitive programming platforms like LeetCode, HackerRank, and Codeforces. The goal is to provide a self-hosted environment where you can:

1. **Write and test solutions** to programming problems in multiple languages
2. **Evaluate correctness** by comparing your output against expected answers
3. **Benchmark performance** by measuring execution time and resource usage
4. **Practice locally** without relying on external services

### How Code Evaluation Works

The core concept is straightforward:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Your Solution  │───▶│  Test Inputs    │───▶│  Your Output    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                      │
                                                      ▼
                                              ┌───────────────┐
                                              │   Compare     │
                                              └───────────────┘
                                                      │
                                                      ▼
┌─────────────────┐                           ┌───────────────┐
│ Expected Output │◀──────────────────────────│    Verdict    │
└─────────────────┘                           └───────────────┘
```

1. **Compile** (if needed) your solution code
2. **Execute** your code with predefined test inputs
3. **Capture** stdout/stderr and resource metrics
4. **Compare** your output against expected output
5. **Report** verdict (Accepted, Wrong Answer, Time Limit Exceeded, etc.)

### Why Containers?

To support multiple programming languages safely and consistently, we use **containerized execution**:

- **Isolation**: User code runs in a sandboxed environment, preventing malicious code from affecting the host system
- **Reproducibility**: Same environment across different machines
- **Resource Control**: Containers allow us to enforce CPU time limits, memory limits, and prevent infinite loops
- **Multi-language Support**: Each language can have its own container image with the appropriate compiler/interpreter

---

## Terminology

| Term | Description |
|------|-------------|
| **Problem** | A programming challenge with a description, constraints, and test cases |
| **Test Case** | A single input/output pair used to verify correctness |
| **Submission** | User's code submitted for evaluation |
| **Verdict** | The result of evaluating a submission (see below) |
| **Judge** | The system that evaluates submissions |
| **Sandbox** | The isolated execution environment (container) |
| **Time Limit (TL)** | Maximum allowed execution time per test case |
| **Memory Limit (ML)** | Maximum allowed memory usage |

### Verdict Types

| Verdict | Abbreviation | Description |
|---------|--------------|-------------|
| **Accepted** | AC | All test cases passed |
| **Wrong Answer** | WA | Output doesn't match expected |
| **Time Limit Exceeded** | TLE | Execution exceeded time limit |
| **Memory Limit Exceeded** | MLE | Memory usage exceeded limit |
| **Runtime Error** | RE | Program crashed (segfault, exception, etc.) |
| **Compilation Error** | CE | Code failed to compile |
| **Presentation Error** | PE | Output format incorrect (whitespace issues) |

---

## System Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         Sandbox Judge                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐ │
│  │     CLI     │    │   Web UI    │    │   Problem Store     │ │
│  │  Interface  │    │  (React?)   │    │   (YAML/JSON)       │ │
│  └──────┬──────┘    └──────┬──────┘    └──────────┬──────────┘ │
│         │                  │                      │             │
│         └────────┬─────────┴──────────────────────┘             │
│                  │                                              │
│                  ▼                                              │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                      Judge Engine                           ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ ││
│  │  │  Compiler   │  │  Executor   │  │  Result Comparator  │ ││
│  │  │  Service    │  │  Service    │  │                     │ ││
│  │  └─────────────┘  └─────────────┘  └─────────────────────┘ ││
│  └─────────────────────────────────────────────────────────────┘│
│                  │                                              │
│                  ▼                                              │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                   Container Runtime                         ││
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐       ││
│  │  │ Python  │  │  C/C++  │  │  Java   │  │   Go    │  ...  ││
│  │  │ Runner  │  │  Runner │  │  Runner │  │  Runner │       ││
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘       ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

### Components

#### 1. Problem Store
Stores problem definitions in a structured format (YAML or JSON):

```yaml
# problems/two-sum/problem.yaml
id: two-sum
title: Two Sum
difficulty: easy
description: |
  Given an array of integers `nums` and an integer `target`, 
  return indices of the two numbers that add up to `target`.
  
constraints:
  - 2 <= nums.length <= 10^4
  - -10^9 <= nums[i] <= 10^9
  
time_limit_ms: 1000
memory_limit_mb: 256

test_cases:
  - input: |
      [2,7,11,15]
      9
    output: |
      [0,1]
  - input: |
      [3,2,4]
      6
    output: |
      [1,2]

# Hidden test cases for full evaluation
hidden_test_cases_file: hidden_tests.yaml
```

#### 2. Judge Engine
The core orchestrator that:
- Receives submissions
- Dispatches compilation jobs
- Runs test cases
- Collects and compares results
- Reports verdicts with metrics

#### 3. Container Runtime
Uses Docker (or Podman) to:
- Spin up isolated containers per submission
- Mount user code into the container
- Execute with resource limits (`--memory`, `--cpus`, `--timeout`)
- Capture stdout, stderr, and exit codes

#### 4. CLI Interface
```bash
# Submit a solution
sandbox-judge submit problems/two-sum solution.py

# Run with specific test case
sandbox-judge run problems/two-sum solution.py --test 1

# List available problems
sandbox-judge list

# Show problem description
sandbox-judge show two-sum
```

#### 5. Web UI
A browser-based interface providing:
- Problem browser with search/filter
- Code editor (Monaco editor like VS Code)
- Real-time submission status
- Submission history and statistics
- Side-by-side problem description and editor (like LeetCode)

---

## Performance Measurement

### Metrics to Capture

| Metric | How to Measure |
|--------|----------------|
| **Wall Clock Time** | Time from process start to exit |
| **CPU Time** | User + system CPU time (from `rusage`) |
| **Memory (Peak)** | Maximum resident set size |
| **Memory (Average)** | Sampled memory usage over execution |

### Implementation Approaches

1. **Container Resource Limits**
   ```bash
   docker run --memory=256m --cpus=1 --timeout=5s ...
   ```

2. **Process Monitoring**
   - Use `cgroups` (Linux) for precise resource tracking
   - `/usr/bin/time -v` for quick measurements
   - Custom wrapper that polls `/proc/<pid>/status`

3. **Seccomp Profiles**
   - Restrict system calls for additional security
   - Prevent network access, file system writes outside workspace

---

## Supported Languages (Initial)

| Language | Compiler/Runtime | Container Base |
|----------|------------------|----------------|
| Python 3 | python3.11 | `python:3.11-slim` |
| JavaScript | Node.js 20 | `node:20-slim` |
| TypeScript | ts-node / tsc | `node:20-slim` |
| C | gcc 12 | `gcc:12` |
| C++ | g++ 12 | `gcc:12` |
| Java | OpenJDK 17 | `openjdk:17-slim` |
| Go | go 1.21 | `golang:1.21-alpine` |
| Rust | rustc 1.73 | `rust:1.73-slim` |

---

## Project Structure

```
sandbox_judge/
├── cmd/
│   └── judge/
│       └── main.go              # CLI entry point
├── internal/
│   ├── judge/
│   │   └── judge.go             # Core orchestration
│   ├── runner/
│   │   ├── runner.go            # Container execution interface
│   │   └── docker.go            # Docker implementation
│   ├── problem/
│   │   ├── loader.go            # YAML parsing
│   │   └── types.go             # Problem structs
│   ├── compare/
│   │   ├── compare.go           # Comparison interface
│   │   ├── default.go           # Whitespace-tolerant
│   │   ├── strict.go            # Exact match
│   │   ├── float.go             # Floating point
│   │   └── unordered.go         # Set comparison
│   └── config/
│       └── languages.go         # Language configurations
├── pkg/
│   └── api/                     # Shared types for CLI/Web
│       └── types.go
├── web/
│   ├── frontend/                # React + TypeScript
│   │   ├── src/
│   │   │   ├── components/
│   │   │   ├── pages/
│   │   │   └── App.tsx
│   │   ├── package.json
│   │   └── tsconfig.json
│   └── server/                  # Go HTTP server (API)
│       └── server.go
├── problems/                    # Problem definitions
│   ├── two-sum/
│   │   ├── problem.yaml
│   │   ├── tests/
│   │   │   ├── sample/
│   │   │   │   ├── 1.in
│   │   │   │   └── 1.out
│   │   │   └── hidden/
│   │   │       ├── 1.in
│   │   │       └── 1.out
│   │   └── solutions/
│   │       └── optimal.py
│   └── ...
├── docker/                      # Language runner Dockerfiles
│   ├── base/
│   │   └── Dockerfile           # Common base image
│   ├── python/
│   │   └── Dockerfile
│   ├── javascript/
│   │   └── Dockerfile
│   ├── go/
│   │   └── Dockerfile
│   └── cpp/
│       └── Dockerfile
├── configs/
│   └── languages.yaml           # Language compile/run commands
├── go.mod
├── go.sum
├── Makefile
├── DESIGN.md
└── README.md
```

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Implementation Language** | Go | Fast, excellent Docker SDK, great for CLI tools |
| **Web UI** | React + TypeScript | Industry standard, rich ecosystem, good DX |
| **Problem Format** | YAML | Human-readable, good for multi-line strings |
| **I/O Style** | Stdin/Stdout (MVP) | Language-agnostic, simpler to implement |

---

## Problem I/O Styles

### Stdin/Stdout (MVP - Primary)

The traditional competitive programming approach. Users write complete programs:

```python
# User's solution reads from stdin, writes to stdout
nums = list(map(int, input().split()))
target = int(input())

# ... solve ...

print(result[0], result[1])
```

**Judge flow:**
```
input.txt → stdin → [user program] → stdout → compare with expected.txt
```

### Function Signature (Future Enhancement)

LeetCode-style where users implement just a function:

```python
class Solution:
    def twoSum(self, nums: List[int], target: int) -> List[int]:
        # Just the algorithm, no I/O handling
        ...
```

This requires **per-language harness templates** that:
1. Deserialize test input → function arguments
2. Call user's function
3. Serialize return value → stdout

We'll add this in a future phase with a template system.

---

## Comparison Strategy

Use a **tiered comparison approach** with configurable modes per problem:

### Default: Whitespace-Tolerant Comparison
```yaml
# problem.yaml
comparison: default  # or omit, this is the default
```

Rules:
- Trim leading/trailing whitespace from each line
- Normalize line endings (CRLF → LF)
- Ignore trailing empty lines
- Exact content match otherwise

### Strict: Exact Match
```yaml
comparison: strict
```
Byte-for-byte comparison. Useful for problems where formatting matters.

### Float: Floating Point Tolerance
```yaml
comparison: float
float_tolerance: 1e-6  # absolute or relative difference
```
For numerical problems where floating point precision varies.

### Unordered: Order-Independent
```yaml
comparison: unordered
```
Each line is compared as a set (order doesn't matter). Useful for "print all solutions" problems.

### Custom: User-Defined Comparator
```yaml
comparison: custom
comparator: checker.py  # Script that returns 0 for correct, 1 for wrong
```
For complex validation (e.g., "any valid path" in graph problems).

---

## Test Case Management

### Structure

```
problems/
└── two-sum/
    ├── problem.yaml       # Metadata, description, constraints
    ├── tests/
    │   ├── sample/        # Visible to user (shown in problem description)
    │   │   ├── 1.in
    │   │   ├── 1.out
    │   │   ├── 2.in
    │   │   └── 2.out
    │   └── hidden/        # Used for full evaluation, not shown
    │       ├── 1.in
    │       ├── 1.out
    │       ├── edge_empty.in
    │       ├── edge_empty.out
    │       ├── large_random.in
    │       └── large_random.out
    └── solutions/         # Reference solutions (for validation)
        ├── optimal.py
        └── brute_force.py
```

### Why Separate Files?

1. **Large test cases** - Some inputs are megabytes; embedding in YAML is unwieldy
2. **Binary data** - Easier to handle as separate files
3. **Git-friendly** - Can see diffs for individual test cases
4. **Generator scripts** - Can generate `.in`/`.out` files programmatically

### problem.yaml (Simplified)

```yaml
id: two-sum
title: Two Sum
difficulty: easy
tags: [array, hash-table]

time_limit_ms: 1000
memory_limit_mb: 256

comparison: default

description: |
  Given an array of integers `nums` and an integer `target`,
  return indices of the two numbers that add up to `target`.
  
  You may assume each input has exactly one solution,
  and you may not use the same element twice.

input_format: |
  Line 1: Space-separated integers (the array)
  Line 2: A single integer (the target)

output_format: |
  Two space-separated integers (the indices)

constraints:
  - 2 <= nums.length <= 10^4
  - -10^9 <= nums[i] <= 10^9
  - Only one valid answer exists

examples:
  - input: |
      2 7 11 15
      9
    output: |
      0 1
    explanation: nums[0] + nums[1] = 2 + 7 = 9

# Test cases loaded from tests/ directory automatically
```

### Test Case Generators (Future)

For stress testing, support generator scripts:

```yaml
# tests/generators/random_large.yaml
generator: random_large.py
count: 10  # Generate 10 test cases
seed: 42   # Reproducible
```

```python
# random_large.py - generates input, reference solution generates output
import random
import sys

n = random.randint(1000, 10000)
nums = [random.randint(-10**9, 10**9) for _ in range(n)]
# ... ensure valid solution exists ...
print(' '.join(map(str, nums)))
print(target)
```

---

## Data Persistence & Multi-User Considerations

### MVP: File-Based Storage (Single User)

For a single-user local setup, **no database is needed**. We can use the filesystem:

```
sandbox_judge/
├── problems/                    # Problem definitions (read-only)
├── data/                        # User data (gitignored)
│   ├── submissions/             # Submission history
│   │   ├── two-sum/
│   │   │   ├── 2024-01-15_143022_abc123.json
│   │   │   └── 2024-01-15_150045_def456.json
│   │   └── ...
│   ├── progress.json            # Problem completion status
│   └── settings.json            # User preferences
```

**Submission record example:**
```json
{
  "id": "abc123",
  "problem_id": "two-sum",
  "language": "python",
  "submitted_at": "2024-01-15T14:30:22Z",
  "verdict": "accepted",
  "runtime_ms": 45,
  "memory_kb": 12400,
  "test_results": [
    {"test": "sample/1", "verdict": "AC", "time_ms": 12},
    {"test": "sample/2", "verdict": "AC", "time_ms": 8},
    {"test": "hidden/1", "verdict": "AC", "time_ms": 45}
  ],
  "source_file": "solution.py"
}
```

**Why this works for MVP:**
- Zero setup, no database to install
- Human-readable, easy to debug
- Git-friendly (can version control your progress)
- Fast enough for single user

---

### Future: Database for Multi-User

When scaling to multiple users behind a proxy with JWT auth:

| Approach | Pros | Cons |
|----------|------|------|
| **SQLite** | Zero config, single file, embedded | Limited concurrency, not ideal for distributed |
| **PostgreSQL** | Robust, great concurrency, JSON support | Requires setup/hosting |
| **SQLite → PostgreSQL** | Start simple, migrate later | Need abstraction layer |

**Recommended: SQLite for now, PostgreSQL-ready**

Use a **repository pattern** so storage is swappable:

```go
// internal/storage/storage.go
type Store interface {
    // Submissions
    SaveSubmission(ctx context.Context, sub *Submission) error
    GetSubmission(ctx context.Context, id string) (*Submission, error)
    ListSubmissions(ctx context.Context, userID, problemID string) ([]*Submission, error)
    
    // Progress
    GetProgress(ctx context.Context, userID string) (*Progress, error)
    UpdateProgress(ctx context.Context, userID, problemID string, status ProblemStatus) error
    
    // Users (future)
    GetUser(ctx context.Context, userID string) (*User, error)
    CreateUser(ctx context.Context, user *User) error
}

// internal/storage/file.go
type FileStore struct { ... }      // MVP: filesystem-based

// internal/storage/sqlite.go  
type SQLiteStore struct { ... }    // Phase 2: embedded DB

// internal/storage/postgres.go
type PostgresStore struct { ... }  // Phase 3: multi-user
```

---

### Multi-User Architecture (Future)

```
                    ┌─────────────────┐
                    │  Reverse Proxy  │
                    │  (nginx/traefik)│
                    └────────┬────────┘
                             │ JWT validation
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                     Sandbox Judge API                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ Auth Middle │  │   Router    │  │   Judge Engine      │ │
│  │   ware      │──│             │──│                     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
│                            │                                │
│                            ▼                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Storage Layer (Interface)              │   │
│  │  FileStore │ SQLiteStore │ PostgresStore            │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │    Database     │
                    │  (PostgreSQL)   │
                    └─────────────────┘
```

**Key considerations for multi-user:**

1. **User isolation** - Each user's submissions are separate
2. **Concurrent execution** - Queue system for running submissions (Redis, or DB-backed queue)
3. **Rate limiting** - Prevent abuse
4. **Resource quotas** - Limit containers per user

---

### Schema Preview (Future PostgreSQL)

```sql
-- Users (managed externally via JWT, we just store ID)
CREATE TABLE users (
    id UUID PRIMARY KEY,
    external_id VARCHAR(255) UNIQUE,  -- From JWT sub claim
    created_at TIMESTAMP DEFAULT NOW()
);

-- Submissions
CREATE TABLE submissions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    problem_id VARCHAR(100) NOT NULL,
    language VARCHAR(50) NOT NULL,
    source_code TEXT NOT NULL,
    verdict VARCHAR(20),  -- AC, WA, TLE, etc.
    runtime_ms INT,
    memory_kb INT,
    submitted_at TIMESTAMP DEFAULT NOW(),
    judged_at TIMESTAMP
);

-- Test Results (per submission)
CREATE TABLE test_results (
    id UUID PRIMARY KEY,
    submission_id UUID REFERENCES submissions(id),
    test_name VARCHAR(100),
    verdict VARCHAR(20),
    runtime_ms INT,
    memory_kb INT,
    error_output TEXT
);

-- User Progress
CREATE TABLE progress (
    user_id UUID REFERENCES users(id),
    problem_id VARCHAR(100),
    status VARCHAR(20),  -- not_attempted, attempted, solved
    best_runtime_ms INT,
    best_memory_kb INT,
    solved_at TIMESTAMP,
    PRIMARY KEY (user_id, problem_id)
);
```

---

### Recommendation Summary

| Phase | Storage | Auth | Use Case |
|-------|---------|------|----------|
| **MVP** | Filesystem (JSON) | None | Single user, local |
| **Phase 2** | SQLite | Optional JWT | Single user, persistent |
| **Phase 3** | PostgreSQL | JWT via proxy | Multi-user, hosted |

**For now:** Build with the storage interface pattern, implement `FileStore` first. This keeps the door open for databases later without over-engineering upfront.

---

## Multi-User Readiness (MVP Considerations)

While MVP is single-user, these small decisions now prevent painful refactors later:

### 1. Thread User Context Through Code

Even for single-user, pass a `userID` parameter:

```go
// ❌ Don't do this - hardcoded assumption
func (j *Judge) Submit(problemID string, code []byte) (*Result, error)

// ✅ Do this - user context from the start
func (j *Judge) Submit(ctx context.Context, userID, problemID string, code []byte) (*Result, error)
```

For MVP, just hardcode the caller:
```go
result, err := judge.Submit(ctx, "local", "two-sum", code)
```

Later, extract `userID` from JWT:
```go
userID := auth.UserIDFromContext(ctx)  // From JWT middleware
result, err := judge.Submit(ctx, userID, "two-sum", code)
```

### 2. Config-Driven Paths

```go
// ❌ Hardcoded paths
dataDir := "./data/submissions"

// ✅ Config-based, can become user-scoped
type Config struct {
    DataDir     string  // "./data" for MVP
    ProblemsDir string  // "./problems"
}

func (c *Config) UserDataDir(userID string) string {
    return filepath.Join(c.DataDir, "users", userID)
}
```

### 3. Container Naming

Include user context in container names for easier cleanup:

```go
// Container name format: judge-{userID}-{submissionID}
containerName := fmt.Sprintf("judge-%s-%s", userID, submissionID)
```

### 4. API Design

Design REST endpoints that work for both single and multi-user:

```
# MVP: userID is implicit (always "local")
POST /api/problems/{id}/submit
GET  /api/problems/{id}/submissions
GET  /api/progress

# Multi-user: Same endpoints, userID from JWT
# No URL changes needed - auth middleware extracts user
```

### 5. What NOT to Build Yet

Don't over-engineer these until needed:

| Feature | Skip for MVP | Add When |
|---------|--------------|----------|
| User registration/login | ✓ | Multi-user phase |
| Rate limiting | ✓ | Multi-user phase |
| Submission queue | ✓ | Concurrent users |
| Database | ✓ | Phase 2+ |
| WebSocket for live updates | ✓ | Web UI polish |

### Summary

The key principle: **Make the single-user case a special case of multi-user, not a separate code path.**

By passing `userID` through your code and using config-driven paths, the transition from single-user to multi-user becomes:
1. Add JWT middleware
2. Swap `FileStore` for `PostgresStore`  
3. Deploy behind a proxy

No core logic changes required.

---

## MVP Scope

For the minimum viable product:

- [ ] CLI tool to run submissions
- [ ] Docker-based execution for Python and JavaScript
- [ ] Basic problem format (YAML)
- [ ] Time limit enforcement
- [ ] Simple output comparison (exact match, whitespace tolerant)
- [ ] 5-10 sample problems

### Future Enhancements

- [ ] Web UI with code editor
- [ ] More languages
- [ ] Custom comparators (for floating point, unordered output)
- [ ] Stress testing / random test generation
- [ ] Solution explanations and hints
- [ ] Progress tracking and statistics
- [ ] Import problems from LeetCode/Codeforces

---

## Documentation Site

Use **MkDocs** with the **Material** theme to generate a polished documentation site.

### Why MkDocs?

- Write docs in Markdown (already doing this)
- Beautiful, responsive output with Material theme
- Built-in search
- Easy to host on GitHub Pages, Netlify, or self-host
- Versioning support

### Proposed Structure

```
docs/
├── index.md                 # Home / overview
├── getting-started/
│   ├── installation.md      # How to install
│   ├── quickstart.md        # First problem in 5 minutes
│   └── configuration.md     # Config options
├── user-guide/
│   ├── cli-reference.md     # All CLI commands
│   ├── web-ui.md            # Using the web interface
│   ├── problems.md          # Browsing and solving
│   └── languages.md         # Supported languages
├── problem-authoring/
│   ├── format.md            # problem.yaml schema
│   ├── test-cases.md        # Writing test cases
│   ├── comparators.md       # Comparison modes
│   └── examples.md          # Sample problems walkthrough
├── deployment/
│   ├── docker.md            # Running in Docker
│   ├── docker-compose.md    # Full stack setup
│   └── multi-user.md        # Scaling to multiple users
├── development/
│   ├── architecture.md      # System design (from DESIGN.md)
│   ├── contributing.md      # How to contribute
│   ├── adding-languages.md  # Adding new language support
│   └── api-reference.md     # REST API docs
└── changelog.md             # Version history
```

### mkdocs.yml (Preview)

```yaml
site_name: Sandbox Judge
site_description: Local code evaluation and benchmarking system
repo_url: https://github.com/marv972228/sandbox_judge

theme:
  name: material
  palette:
    primary: indigo
    accent: amber
  features:
    - navigation.tabs
    - navigation.sections
    - search.highlight
    - content.code.copy

markdown_extensions:
  - pymdownx.highlight
  - pymdownx.superfences
  - pymdownx.tabbed
  - admonitions
  - toc:
      permalink: true

nav:
  - Home: index.md
  - Getting Started:
    - Installation: getting-started/installation.md
    - Quickstart: getting-started/quickstart.md
    - Configuration: getting-started/configuration.md
  - User Guide:
    - CLI Reference: user-guide/cli-reference.md
    - Web UI: user-guide/web-ui.md
  - Problem Authoring:
    - Format: problem-authoring/format.md
    - Test Cases: problem-authoring/test-cases.md
  - Deployment:
    - Docker: deployment/docker.md
  - Development:
    - Architecture: development/architecture.md
    - Contributing: development/contributing.md
```

---

## References & Inspiration

- [LeetCode](https://leetcode.com) - Popular interview prep platform
- [Codeforces](https://codeforces.com) - Competitive programming
- [Judge0](https://github.com/judge0/judge0) - Open source online judge API
- [DMOJ](https://github.com/DMOJ/judge-server) - Open source judge server
- [Piston](https://github.com/engineer-man/piston) - Code execution engine

---

## Next Steps

1. Validate this design - does it cover your needs?
2. Choose implementation language
3. Set up basic project structure
4. Implement MVP CLI + Docker runner
5. Add sample problems
6. Build Web UI


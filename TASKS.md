# Sandbox Judge - Development Tasks

This document tracks implementation progress with testable milestones.

---

## Phase 1: Foundation (CLI + Single Language)

**Goal:** Run a Python solution against a problem and get a verdict.

### 1.1 Project Setup
- [ ] Initialize Go module (`go mod init`)
- [ ] Create directory structure
- [ ] Add basic `Makefile` with `build`, `test`, `run` targets
- [ ] Create `.gitignore` (binaries, data/, vendor/, etc.)

**✅ Checkpoint:** `make build` produces a binary

---

### 1.2 Problem Loader
- [ ] Define `Problem` struct in `internal/problem/types.go`
- [ ] Implement YAML parser in `internal/problem/loader.go`
- [ ] Create first sample problem: `problems/two-sum/`
  - [ ] `problem.yaml` with metadata
  - [ ] `tests/sample/1.in`, `tests/sample/1.out`
  - [ ] `tests/sample/2.in`, `tests/sample/2.out`
- [ ] Add CLI command: `judge show <problem-id>`

**✅ Checkpoint:** `./judge show two-sum` prints problem description

---

### 1.3 Docker Runner (Python Only)
- [ ] Create Python runner Dockerfile (`docker/python/Dockerfile`)
- [ ] Define `Runner` interface in `internal/runner/runner.go`
- [ ] Implement `DockerRunner` in `internal/runner/docker.go`
  - [ ] Build/pull image
  - [ ] Mount source code
  - [ ] Execute with stdin from test input
  - [ ] Capture stdout, stderr, exit code
  - [ ] Enforce time limit (`--stop-timeout`)
  - [ ] Enforce memory limit (`--memory`)
- [ ] Add basic resource measurement (wall clock time)

**✅ Checkpoint:** Can manually run Python code in container and see output

---

### 1.4 Output Comparator
- [ ] Define `Comparator` interface in `internal/compare/compare.go`
- [ ] Implement `DefaultComparator` (whitespace-tolerant)
  - [ ] Trim lines
  - [ ] Normalize line endings
  - [ ] Ignore trailing blank lines

**✅ Checkpoint:** Unit tests pass for comparator edge cases

---

### 1.5 Judge Engine (Core Loop)
- [ ] Implement `Judge` in `internal/judge/judge.go`
  - [ ] Load problem
  - [ ] For each test case:
    - [ ] Run submission in container
    - [ ] Compare output
    - [ ] Record verdict (AC/WA/TLE/RE)
  - [ ] Aggregate results
- [ ] Add CLI command: `judge run <problem-id> <solution-file>`

**✅ Checkpoint:** `./judge run two-sum solution.py` returns verdict

```bash
# Example output
$ ./judge run two-sum solutions/correct.py
Running two-sum...
  Test 1: AC (12ms)
  Test 2: AC (8ms)
Result: ACCEPTED (2/2 tests passed)
```

---

### 1.6 Basic CLI Polish
- [ ] Add `judge list` - list all problems
- [ ] Add `judge run --test N` - run specific test only
- [ ] Add `judge run --verbose` - show input/output diff on failure
- [ ] Colorized output (green AC, red WA, yellow TLE)

**✅ Checkpoint:** CLI feels usable for daily practice

---

## Phase 2: Multi-Language Support

**Goal:** Support Python, JavaScript, Go, and C++.

### 2.1 Language Configuration
- [ ] Create `configs/languages.yaml`
  ```yaml
  python:
    image: sandbox-judge-python
    compile: null
    run: ["python3", "{file}"]
    extension: .py
  ```
- [ ] Load language config at startup
- [ ] Auto-detect language from file extension

**✅ Checkpoint:** Config loads without errors

---

### 2.2 Additional Runners
- [ ] JavaScript/Node.js
  - [ ] `docker/javascript/Dockerfile`
  - [ ] Test with sample problem
- [ ] Go
  - [ ] `docker/go/Dockerfile`
  - [ ] Handle compilation step
- [ ] C++
  - [ ] `docker/cpp/Dockerfile`
  - [ ] Compile then run

**✅ Checkpoint:** Same problem solved in 4 languages, all get AC

---

### 2.3 Compilation Support
- [ ] Add compile step to runner for compiled languages
- [ ] Capture compilation errors → CE verdict
- [ ] Cache compiled binaries (optional optimization)

**✅ Checkpoint:** C++ syntax error returns "Compilation Error"

---

## Phase 3: Storage & Progress Tracking

**Goal:** Track submissions and see your history.

### 3.1 Storage Interface
- [ ] Define `Store` interface in `internal/storage/storage.go`
- [ ] Implement `FileStore` (JSON files)
  - [ ] `SaveSubmission()`
  - [ ] `GetSubmission()`
  - [ ] `ListSubmissions()`

**✅ Checkpoint:** Submissions persist to `data/submissions/`

---

### 3.2 Progress Tracking
- [ ] Track problem status: not_attempted → attempted → solved
- [ ] Add CLI command: `judge status` - show overall progress
- [ ] Add CLI command: `judge history <problem-id>` - show past submissions

**✅ Checkpoint:** `./judge status` shows progress summary

```bash
$ ./judge status
Problems: 3/10 solved
  ✓ two-sum (AC, 45ms)
  ✓ reverse-string (AC, 12ms)
  ✗ valid-parentheses (WA, 2 attempts)
  - binary-search (not attempted)
  ...
```

---

## Phase 4: Problem Library

**Goal:** Curated problem set covering common interview topics at varying difficulties.

### 4.1 Problem Categories & Coverage

Target: **30+ problems** across these categories:

| Category | Easy | Medium | Hard | Total |
|----------|------|--------|------|-------|
| Arrays | 3 | 2 | 1 | 6 |
| Strings | 3 | 2 | 1 | 6 |
| Hash Tables | 2 | 2 | 1 | 5 |
| Linked Lists | 2 | 2 | 1 | 5 |
| Stacks & Queues | 2 | 1 | 1 | 4 |
| Trees & Graphs | 2 | 2 | 1 | 5 |
| Dynamic Programming | 1 | 2 | 2 | 5 |
| Sorting & Searching | 2 | 2 | 1 | 5 |
| **Total** | **17** | **15** | **9** | **41** |

---

### 4.2 Arrays (6 problems)
- [ ] **Easy:** Two Sum
- [ ] **Easy:** Best Time to Buy and Sell Stock
- [ ] **Easy:** Contains Duplicate
- [ ] **Medium:** Product of Array Except Self
- [ ] **Medium:** Maximum Subarray (Kadane's)
- [ ] **Hard:** Trapping Rain Water

**✅ Checkpoint:** All array problems have 3+ hidden test cases

---

### 4.3 Strings (6 problems)
- [ ] **Easy:** Reverse String
- [ ] **Easy:** Valid Palindrome
- [ ] **Easy:** Valid Anagram
- [ ] **Medium:** Longest Substring Without Repeating Characters
- [ ] **Medium:** Group Anagrams
- [ ] **Hard:** Minimum Window Substring

**✅ Checkpoint:** All string problems complete

---

### 4.4 Hash Tables (5 problems)
- [ ] **Easy:** Two Sum (already in Arrays, can cross-tag)
- [ ] **Easy:** First Unique Character in a String
- [ ] **Medium:** LRU Cache
- [ ] **Medium:** Subarray Sum Equals K
- [ ] **Hard:** Longest Consecutive Sequence

**✅ Checkpoint:** All hash table problems complete

---

### 4.5 Linked Lists (5 problems)
- [ ] **Easy:** Reverse Linked List
- [ ] **Easy:** Merge Two Sorted Lists
- [ ] **Medium:** Remove Nth Node From End
- [ ] **Medium:** Linked List Cycle (detect cycle)
- [ ] **Hard:** Merge K Sorted Lists

**✅ Checkpoint:** All linked list problems complete

---

### 4.6 Stacks & Queues (4 problems)
- [ ] **Easy:** Valid Parentheses
- [ ] **Easy:** Implement Queue using Stacks
- [ ] **Medium:** Min Stack
- [ ] **Hard:** Largest Rectangle in Histogram

**✅ Checkpoint:** All stack/queue problems complete

---

### 4.7 Trees & Graphs (5 problems)
- [ ] **Easy:** Maximum Depth of Binary Tree
- [ ] **Easy:** Invert Binary Tree
- [ ] **Medium:** Binary Tree Level Order Traversal
- [ ] **Medium:** Number of Islands
- [ ] **Hard:** Serialize and Deserialize Binary Tree

**✅ Checkpoint:** All tree/graph problems complete

---

### 4.8 Dynamic Programming (5 problems)
- [ ] **Easy:** Climbing Stairs
- [ ] **Medium:** Coin Change
- [ ] **Medium:** Longest Increasing Subsequence
- [ ] **Hard:** Edit Distance
- [ ] **Hard:** Word Break II

**✅ Checkpoint:** All DP problems complete

---

### 4.9 Sorting & Searching (5 problems)
- [ ] **Easy:** Binary Search
- [ ] **Easy:** First Bad Version
- [ ] **Medium:** Search in Rotated Sorted Array
- [ ] **Medium:** Find Peak Element
- [ ] **Hard:** Median of Two Sorted Arrays

**✅ Checkpoint:** All sorting/searching problems complete

---

### 4.10 Problem Quality Checklist

For each problem, ensure:
- [ ] Clear description with examples
- [ ] Input/output format documented
- [ ] Constraints listed
- [ ] 2+ sample test cases (visible)
- [ ] 3+ hidden test cases (edge cases, large inputs)
- [ ] At least one reference solution
- [ ] Appropriate time/memory limits set

**✅ Checkpoint:** All 30+ problems pass quality checklist

---

### 4.11 Additional Comparators
- [ ] Implement `StrictComparator` (exact match)
- [ ] Implement `FloatComparator` (tolerance-based)
- [ ] Implement `UnorderedComparator` (set comparison)
- [ ] Add problem that uses float comparison

**✅ Checkpoint:** Float problem works correctly

---

## Phase 5: Web UI (Basic)

**Goal:** Browser-based interface for viewing problems and submitting.

### 5.1 API Server
- [ ] Create HTTP server in `web/server/`
- [ ] `GET /api/problems` - list problems
- [ ] `GET /api/problems/{id}` - get problem details
- [ ] `POST /api/problems/{id}/submit` - submit solution
- [ ] `GET /api/problems/{id}/submissions` - list submissions
- [ ] `GET /api/status` - get progress

**✅ Checkpoint:** API responds correctly via curl

---

### 5.2 React Frontend Setup
- [ ] Initialize React + TypeScript project in `web/frontend/`
- [ ] Set up Vite for development
- [ ] Add Tailwind CSS for styling
- [ ] Create basic layout (sidebar + main content)

**✅ Checkpoint:** `npm run dev` shows "Hello World"

---

### 5.3 Problem List Page
- [ ] Fetch and display problems
- [ ] Show difficulty, tags, solved status
- [ ] Filter by difficulty/tag
- [ ] Search by title

**✅ Checkpoint:** Can browse all problems in browser

---

### 5.4 Problem Detail Page
- [ ] Display problem description (markdown rendered)
- [ ] Show examples with input/output
- [ ] Integrate Monaco editor for code
- [ ] Language selector dropdown
- [ ] Submit button

**✅ Checkpoint:** Can view problem and write code

---

### 5.5 Submission & Results
- [ ] Submit code from editor
- [ ] Show "Judging..." state
- [ ] Display verdict with test case breakdown
- [ ] Show runtime/memory stats

**✅ Checkpoint:** Full submit → judge → result flow works

---

## Phase 6: Polish & Quality of Life

### 6.1 Better Error Handling
- [ ] Friendly error messages for common issues
- [ ] Docker not running detection
- [ ] Invalid problem ID handling

### 6.2 Performance
- [ ] Pre-pull Docker images on first run
- [ ] Container reuse (optional)
- [ ] Parallel test execution (optional)

### 6.3 Documentation
- [ ] README with setup instructions
- [ ] How to add new problems guide
- [ ] How to add new languages guide

---

## Phase 7: Documentation Site

**Goal:** Professional docs site for users, contributors, and deployers.

### 7.1 MkDocs Setup
- [ ] Install MkDocs and Material theme
- [ ] Create `mkdocs.yml` configuration
- [ ] Create `docs/` directory structure
- [ ] Add `make docs` and `make docs-serve` targets

**✅ Checkpoint:** `make docs-serve` shows docs at localhost:8000

---

### 7.2 User Documentation
- [ ] Installation guide (prerequisites, setup)
- [ ] Quickstart tutorial (solve first problem)
- [ ] CLI reference (all commands documented)
- [ ] Web UI guide (screenshots, workflows)
- [ ] Supported languages list

**✅ Checkpoint:** New user can get started from docs alone

---

### 7.3 Problem Author Documentation
- [ ] `problem.yaml` schema reference
- [ ] Test case format and best practices
- [ ] Comparator modes explained
- [ ] Walkthrough: creating a new problem

**✅ Checkpoint:** Can create a new problem following only the docs

---

### 7.4 Deployment Documentation
- [ ] Docker single-container setup
- [ ] Docker Compose full stack
- [ ] Environment variables reference
- [ ] Multi-user deployment guide (future)

**✅ Checkpoint:** Can deploy from docs without reading source

---

### 7.5 Developer Documentation
- [ ] Architecture overview (adapted from DESIGN.md)
- [ ] Contributing guide (code style, PR process)
- [ ] Adding a new language guide
- [ ] API reference (REST endpoints)

**✅ Checkpoint:** New contributor can understand codebase

---

### 7.6 Publishing
- [ ] GitHub Actions to build docs on push
- [ ] Deploy to GitHub Pages (or alternative)
- [ ] Add version dropdown (if needed)

**✅ Checkpoint:** Docs live at `https://marv972228.github.io/sandbox_judge`

---

## Future Ideas (Backlog)

Not planned for MVP, but captured for later:

- [ ] Function signature mode (LeetCode style)
- [ ] Test case generators
- [ ] Import problems from LeetCode/Codeforces
- [ ] Solution explanations and hints
- [ ] Timed contest mode
- [ ] WebSocket for live judging updates
- [ ] Multi-user support with JWT auth
- [ ] Leaderboards
- [ ] Spaced repetition for problem review

---

## Current Focus

> **Now working on:** Phase 1.1 - Project Setup

---

## Notes

- Update this file as tasks complete
- Add new tasks as discovered
- Each checkpoint should be demonstrable

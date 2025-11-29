- [Sandbox Judge - Development Tasks](#sandbox-judge---development-tasks)
  - [Phase 1: Foundation (CLI + Single Language)](#phase-1-foundation-cli--single-language)
    - [1.1 Project Setup](#11-project-setup)
    - [1.2 Problem Loader](#12-problem-loader)
    - [1.3 Docker Runner (Python Only)](#13-docker-runner-python-only)
    - [1.4 Output Comparator](#14-output-comparator)
    - [1.5 Judge Engine (Core Loop)](#15-judge-engine-core-loop)
    - [1.6 Basic CLI Polish](#16-basic-cli-polish)
  - [Phase 2: Multi-Language Support](#phase-2-multi-language-support)
    - [2.1 Language Configuration](#21-language-configuration)
    - [2.2 Additional Runners](#22-additional-runners)
    - [2.3 Compilation Support](#23-compilation-support)
  - [Phase 3: Storage \& Progress Tracking](#phase-3-storage--progress-tracking)
    - [3.1 Storage Interface](#31-storage-interface)
    - [3.2 Progress Tracking](#32-progress-tracking)
  - [Phase 4: Problem Library](#phase-4-problem-library)
    - [4.1 Problem Categories \& Coverage](#41-problem-categories--coverage)
    - [4.2 Arrays (6 problems)](#42-arrays-6-problems)
    - [4.3 Strings (6 problems)](#43-strings-6-problems)
    - [4.4 Hash Tables (5 problems)](#44-hash-tables-5-problems)
    - [4.5 Linked Lists (5 problems)](#45-linked-lists-5-problems)
    - [4.6 Stacks \& Queues (4 problems)](#46-stacks--queues-4-problems)
    - [4.7 Trees \& Graphs (5 problems)](#47-trees--graphs-5-problems)
    - [4.8 Dynamic Programming (5 problems)](#48-dynamic-programming-5-problems)
    - [4.9 Sorting \& Searching (5 problems)](#49-sorting--searching-5-problems)
    - [4.10 Problem Quality Checklist](#410-problem-quality-checklist)
    - [4.11 Additional Comparators](#411-additional-comparators)
  - [Phase 5: Web UI (Basic)](#phase-5-web-ui-basic)
    - [5.1 API Server](#51-api-server)
    - [5.2 React Frontend Setup](#52-react-frontend-setup)
    - [5.3 Problem List Page](#53-problem-list-page)
    - [5.4 Problem Detail Page](#54-problem-detail-page)
    - [5.5 Submission \& Results](#55-submission--results)
  - [Phase 6: Polish \& Quality of Life](#phase-6-polish--quality-of-life)
    - [6.1 Better Error Handling](#61-better-error-handling)
    - [6.2 Performance](#62-performance)
    - [6.3 Documentation](#63-documentation)
  - [Phase 7: Documentation Site](#phase-7-documentation-site)
    - [7.1 MkDocs Setup](#71-mkdocs-setup)
    - [7.2 User Documentation](#72-user-documentation)
    - [7.3 Problem Author Documentation](#73-problem-author-documentation)
    - [7.4 Deployment Documentation](#74-deployment-documentation)
    - [7.5 Developer Documentation](#75-developer-documentation)
    - [7.6 Publishing](#76-publishing)
  - [Future Ideas (Backlog)](#future-ideas-backlog)
  - [Current Focus](#current-focus)
  - [Notes](#notes)


# Sandbox Judge - Development Tasks

This document tracks implementation progress with testable milestones.

---

## Phase 1: Foundation (CLI + Single Language)

**Goal:** Run a Python solution against a problem and get a verdict.

### 1.1 Project Setup
- [x] Initialize Go module (`go mod init`)
- [x] Create directory structure
- [x] Add basic `Makefile` with `build`, `test`, `run` targets
- [x] Create `.gitignore` (binaries, data/, vendor/, etc.)
- [x] Add Cobra for CLI framework (discovered: need flags, subcommands, help generation)
- [x] Refactor `main.go` to use Cobra commands

**✅ Checkpoint:** `make build` produces a binary, `judge --help` shows auto-generated help

---

### 1.2 Problem Loader
- [x] Define `Problem` struct in `internal/problem/types.go`
- [x] Implement YAML parser in `internal/problem/loader.go`
- [x] Create first sample problem: `problems/two-sum/`
  - [x] `problem.yaml` with metadata
  - [x] `tests/sample/1.in`, `tests/sample/1.out`
  - [x] `tests/sample/2.in`, `tests/sample/2.out`
- [x] Add CLI commands: `judge list` and `judge show <problem-id>`

**✅ Checkpoint:** `./judge show two-sum` prints problem description ✓

---

### 1.3 Docker Runner (Python Only)
- [x] Create Python runner Dockerfile (`docker/python/Dockerfile`)
- [x] Define `Runner` interface in `internal/runner/runner.go`
- [x] Implement `DockerRunner` in `internal/runner/docker.go`
  - [x] Build/pull image
  - [x] Mount source code
  - [x] Execute with stdin from test input
  - [x] Capture stdout, stderr, exit code
  - [x] Enforce time limit (`--stop-timeout`)
  - [x] Enforce memory limit (`--memory`)
- [x] Add basic resource measurement (wall clock time)
- [x] Verify Go code works (cmd/testrunner - AC, TLE verdicts tested)

**✅ Checkpoint:** Can manually run Python code in container and see output ✓

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

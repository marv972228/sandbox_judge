# Problem Format Overview

Sandbox Judge uses a simple file-based format for problems. Each problem is a directory containing configuration and test cases.

## Directory Structure

```
problems/
└── two-sum/
    ├── problem.yaml       # Problem configuration
    └── tests/
        ├── sample/        # Sample test cases (visible to users)
        │   ├── 1.in       # Input for test 1
        │   ├── 1.out      # Expected output for test 1
        │   ├── 2.in
        │   └── 2.out
        └── hidden/        # Hidden test cases (optional)
            ├── 1.in
            ├── 1.out
            └── ...
```

## problem.yaml

The `problem.yaml` file contains all metadata for a problem:

```yaml
id: two-sum
title: Two Sum
difficulty: easy
tags:
  - arrays
  - hash-table

time_limit_ms: 1000
memory_limit_mb: 256

description: |
  Given an array of integers nums and an integer target, return
  indices of the two numbers such that they add up to target.
  
  You may assume that each input would have exactly one solution,
  and you may not use the same element twice.

input_format: |
  - Line 1: Space-separated integers (the array)
  - Line 2: The target sum

output_format: |
  - Two space-separated indices

examples:
  - input: |
      2 7 11 15
      9
    output: |
      0 1
    explanation: "nums[0] + nums[1] = 2 + 7 = 9"
```

## Test Cases

Test cases are pairs of `.in` and `.out` files:

- **`.in` files** - Input to be passed to the solution via stdin
- **`.out` files** - Expected output to compare against stdout

### Sample vs Hidden Tests

- **Sample tests** (`tests/sample/`) - Shown to users, used for debugging
- **Hidden tests** (`tests/hidden/`) - Not shown, used for final grading

Both are executed when running `judge run`.

## Field Reference

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier (matches directory name) |
| `title` | Yes | Human-readable title |
| `difficulty` | Yes | `easy`, `medium`, or `hard` |
| `tags` | No | List of category tags |
| `time_limit_ms` | Yes | Time limit in milliseconds |
| `memory_limit_mb` | Yes | Memory limit in megabytes |
| `description` | Yes | Problem description (markdown) |
| `input_format` | No | Description of input format |
| `output_format` | No | Description of output format |
| `examples` | No | Example inputs/outputs with explanations |
| `comparator` | No | Output comparison mode (default: `default`) |

## Comparators

| Mode | Description |
|------|-------------|
| `default` | Whitespace-tolerant (trims lines, ignores trailing blanks) |
| `strict` | Exact byte-for-byte match |
| `float` | Floating-point with tolerance (coming soon) |
| `unordered` | Order-independent comparison (coming soon) |

## Next Steps

- [Creating Problems](creating.md) - Step-by-step guide to creating problems

# judge show

Show a problem's description and details.

## Synopsis

```bash
judge show <problem-id> [flags]
```

## Description

The `show` command displays detailed information about a problem, including:

- Title and difficulty
- Tags/categories
- Time and memory limits
- Problem description
- Input/output format
- Example test cases

## Arguments

| Argument | Description |
|----------|-------------|
| `problem-id` | The ID of the problem to display |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--help` | `-h` | Help for show |

## Examples

### Basic Usage

```bash
judge show two-sum
```

Output:
```
two-sum - Two Sum
=================
Difficulty: easy
Tags: arrays, hash-table
Time Limit: 1000ms | Memory: 256MB

Description:
Given an array of integers nums and an integer target, return
indices of the two numbers such that they add up to target.

You may assume that each input would have exactly one solution,
and you may not use the same element twice.

Input Format:
- Line 1: Space-separated integers (the array)
- Line 2: The target sum

Output Format:
- Two space-separated indices

Example:
Input:
2 7 11 15
9

Output:
0 1
```

## See Also

- [judge list](list.md) - List available problems
- [judge run](run.md) - Run a solution

# Quick Start

Get up and running with Sandbox Judge in 5 minutes.

## Your First Problem

Sandbox Judge comes with a sample problem called `two-sum`. Let's solve it!

### 1. View the Problem

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
```

### 2. Understand the Format

**Input:**
```
2 7 11 15
9
```
- Line 1: Space-separated integers (the array)
- Line 2: The target sum

**Output:**
```
0 1
```
- Two space-separated indices

### 3. Write Your Solution

Create a file `my_solution.py`:

```python
#!/usr/bin/env python3

def two_sum(nums, target):
    seen = {}
    for i, num in enumerate(nums):
        complement = target - num
        if complement in seen:
            return [seen[complement], i]
        seen[num] = i
    return []

if __name__ == "__main__":
    nums = list(map(int, input().split()))
    target = int(input())
    
    result = two_sum(nums, target)
    print(result[0], result[1])
```

### 4. Test Your Solution

```bash
judge run two-sum my_solution.py
```

If correct, you'll see:
```
Running two-sum...
  sample/1: AC (45ms)
  sample/2: AC (38ms)

Result: AC (2/2 tests passed)
Total time: 83ms
```

🎉 **Congratulations!** You've solved your first problem!

## Understanding Verdicts

| Verdict | Meaning |
|---------|---------|
| **AC** | Accepted - Your solution is correct |
| **WA** | Wrong Answer - Output doesn't match expected |
| **TLE** | Time Limit Exceeded - Solution too slow |
| **RE** | Runtime Error - Code crashed |
| **SE** | System Error - Judge issue (not your fault) |

## Using Verbose Mode

To see more details when debugging:

```bash
judge run two-sum my_solution.py --verbose
```

On WA, you'll see the expected vs actual output:
```
Running two-sum...
  sample/1: WA (45ms)
    Expected:
      0 1
    Actual:
      1 0

Result: WA (0/2 tests passed)
```

## Running Specific Tests

To run only one test case:

```bash
# Run only test case 1
judge run two-sum my_solution.py --test 1
```

## Listing Problems

See all available problems:

```bash
judge list
```

Output:
```
Available Problems:
  two-sum      [easy]   arrays, hash-table
```

## Next Steps

- [CLI Reference](../cli/overview.md) - Learn all commands
- [Problem Format](../problems/overview.md) - Create your own problems
- [Creating Problems](../problems/creating.md) - Add custom problems

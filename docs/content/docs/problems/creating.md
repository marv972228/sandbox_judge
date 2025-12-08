# Creating Problems

This guide walks you through creating a new problem for Sandbox Judge.

## Quick Start

1. Create a directory in `problems/`
2. Add `problem.yaml`
3. Add test cases in `tests/sample/`
4. (Optional) Add hidden tests in `tests/hidden/`

## Step-by-Step Example

Let's create a "Reverse String" problem.

### 1. Create Directory Structure

```bash
mkdir -p problems/reverse-string/tests/sample
mkdir -p problems/reverse-string/tests/hidden
```

### 2. Create problem.yaml

```bash
cat > problems/reverse-string/problem.yaml << 'EOF'
id: reverse-string
title: Reverse String
difficulty: easy
tags:
  - strings
  - two-pointers

time_limit_ms: 1000
memory_limit_mb: 256

description: |
  Write a function that reverses a string.
  
  The input string is given as a single line.

input_format: |
  A single line containing the string to reverse.

output_format: |
  The reversed string.

examples:
  - input: |
      hello
    output: |
      olleh
    explanation: "The characters are reversed"
  - input: |
      Hannah
    output: |
      hannaH
    explanation: "Note that case is preserved"
EOF
```

### 3. Create Sample Test Cases

```bash
# Test 1
echo "hello" > problems/reverse-string/tests/sample/1.in
echo "olleh" > problems/reverse-string/tests/sample/1.out

# Test 2
echo "Hannah" > problems/reverse-string/tests/sample/2.in
echo "hannaH" > problems/reverse-string/tests/sample/2.out
```

### 4. Create Hidden Test Cases

```bash
# Edge case: empty string
echo "" > problems/reverse-string/tests/hidden/1.in
echo "" > problems/reverse-string/tests/hidden/1.out

# Edge case: single character
echo "a" > problems/reverse-string/tests/hidden/2.in
echo "a" > problems/reverse-string/tests/hidden/2.out

# Large input
python3 -c "print('a' * 10000)" > problems/reverse-string/tests/hidden/3.in
python3 -c "print('a' * 10000)" > problems/reverse-string/tests/hidden/3.out
```

### 5. Verify the Problem

```bash
# List problems to see your new problem
judge list

# Show problem details
judge show reverse-string
```

### 6. Test with a Solution

Create a solution file:

```python
#!/usr/bin/env python3
s = input()
print(s[::-1])
```

Run it:

```bash
judge run reverse-string solution.py
```

## Best Practices

### Test Case Design

1. **Start simple** - Basic cases that demonstrate the problem
2. **Edge cases** - Empty input, single element, boundaries
3. **Large inputs** - Test performance (close to limits)
4. **Corner cases** - Negative numbers, duplicates, etc.

### Time/Memory Limits

| Difficulty | Time Limit | Memory |
|------------|------------|--------|
| Easy | 1000-2000ms | 256MB |
| Medium | 1000-3000ms | 256MB |
| Hard | 2000-5000ms | 512MB |

!!! tip "Tip"
    Run your reference solution and set the limit to ~10x the actual runtime.

### Problem Description

Write clear descriptions that include:

- What the problem is asking
- Constraints and guarantees
- Input/output format
- At least one worked example

### Naming Conventions

- Use lowercase with hyphens: `two-sum`, `reverse-string`
- Match the `id` field to the directory name
- Use descriptive, searchable names

## Validating Problems

Before sharing a problem, verify:

- [ ] `judge show <problem-id>` displays correctly
- [ ] At least one correct solution gets AC
- [ ] Time limits are reasonable (not too tight, not too loose)
- [ ] Edge cases are covered
- [ ] Description is clear and complete

## Troubleshooting

### Problem Not Found

```
Error: problem not found: my-problem
```

Check that:
1. Directory exists: `problems/my-problem/`
2. `problem.yaml` exists and is valid YAML
3. `id` in YAML matches directory name

### Test Cases Not Found

```
Error: no test cases found for problem my-problem
```

Check that:
1. `tests/sample/` directory exists
2. Files are named `N.in` and `N.out` (where N is a number)
3. Both `.in` and `.out` files exist for each test

### YAML Parse Errors

Common YAML issues:
- Incorrect indentation (use spaces, not tabs)
- Missing quotes around strings with special characters
- Forgetting the `|` for multi-line strings

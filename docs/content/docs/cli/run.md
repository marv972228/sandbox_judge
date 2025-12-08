# judge run

Run a solution against a problem's test cases.

## Synopsis

```bash
judge run <problem-id> <solution-file> [flags]
```

## Description

The `run` command executes your solution against all test cases for a given problem. It:

1. Loads the problem configuration
2. Runs your solution in a Docker container for each test case
3. Compares output against expected results
4. Reports verdicts (AC, WA, TLE, RE)

## Arguments

| Argument | Description |
|----------|-------------|
| `problem-id` | The ID of the problem (e.g., `two-sum`) |
| `solution-file` | Path to your solution file |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Show detailed output including input/output diff on failure |
| `--test int` | `-t` | Run only a specific test case (0 = all) |
| `--timeout duration` | | Override the problem's time limit |
| `--help` | `-h` | Help for run |

## Examples

### Basic Usage

```bash
# Run all test cases
judge run two-sum solution.py
```

Output:
```
Running two-sum...
  sample/1: AC (45ms)
  sample/2: AC (38ms)

Result: AC (2/2 tests passed)
Total time: 83ms
```

### Verbose Mode

See detailed output on failures:

```bash
judge run two-sum solution.py --verbose
```

On Wrong Answer:
```
Running two-sum...
  sample/1: WA (45ms)
    Expected:
      0 1
    Actual:
      1 0
  sample/2: AC (38ms)

Result: WA (1/2 tests passed)
Total time: 83ms
```

On Runtime Error:
```
Running two-sum...
  sample/1: RE (45ms)
    Error: Traceback (most recent call last):
      File "/sandbox/solution.py", line 5, in <module>
        result = nums[10]  # IndexError
    IndexError: list index out of range
```

### Run Specific Test

Run only test case 1:

```bash
judge run two-sum solution.py --test 1
```

### Override Time Limit

Set a custom timeout:

```bash
judge run two-sum solution.py --timeout 5s
```

## Verdicts

Verdicts are colorized in the terminal for quick visual feedback:

| Verdict | Color | Description |
|---------|-------|-------------|
| **AC** (Accepted) | 🟢 Green | Output matches expected exactly |
| **WA** (Wrong Answer) | 🔴 Red | Output doesn't match expected |
| **TLE** (Time Limit Exceeded) | 🟡 Yellow | Execution exceeded time limit |
| **MLE** (Memory Limit Exceeded) | 🟡 Yellow | Execution exceeded memory limit |
| **RE** (Runtime Error) | 🔴 Red | Program crashed or non-zero exit |
| **CE** (Compilation Error) | 🔴 Red | Failed to compile (compiled languages) |
| **SE** (System Error) | 🔴 Red | Internal judge error |

## Output Comparison

By default, output comparison is **whitespace-tolerant**:

- Leading/trailing whitespace on lines is trimmed
- Trailing blank lines are ignored
- Line endings are normalized

For problems requiring exact matching, the problem can specify `comparator: strict`.

## Supported Languages

Currently detected by file extension:

| Extension | Language |
|-----------|----------|
| `.py` | Python 3 |

More languages coming soon!

## See Also

- [judge list](list.md) - List available problems
- [judge show](show.md) - Show problem details

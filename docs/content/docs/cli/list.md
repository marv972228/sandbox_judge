# judge list

List all available problems.

## Synopsis

```bash
judge list [flags]
```

## Description

The `list` command displays all problems available in the problems directory. It shows:

- Problem ID
- Title
- Difficulty level
- Tags/categories

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--help` | `-h` | Help for list |

## Examples

### Basic Usage

```bash
judge list
```

Output:
```
ID                   TITLE                          DIFFICULTY TAGS
--------------------------------------------------------------------------------
two-sum              Two Sum                        easy       array, hash-table

1 problem(s) found
```

## Global Flags

You can specify a custom problems directory:

```bash
judge list --problems /path/to/my/problems
```

## See Also

- [judge show](show.md) - Show problem details
- [judge run](run.md) - Run a solution

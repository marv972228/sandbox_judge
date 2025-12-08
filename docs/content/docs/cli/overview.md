# CLI Overview

Sandbox Judge provides a command-line interface for running and managing coding problems.

## Basic Usage

```bash
judge [command] [flags]
```

## Available Commands

| Command | Description |
|---------|-------------|
| `run` | Run a solution against a problem |
| `list` | List all available problems |
| `show` | Show problem description |
| `help` | Help about any command |

## Global Flags

These flags work with any command:

| Flag | Description |
|------|-------------|
| `--config string` | Config file (default: `$HOME/.judge.yaml`) |
| `--problems string` | Path to problems directory (default: `./problems`) |
| `-h, --help` | Help for judge |
| `-v, --version` | Version for judge |

## Command Details

### judge run

Run a solution file against a problem's test cases.

```bash
judge run <problem-id> <solution-file> [flags]
```

**Example:**
```bash
judge run two-sum solution.py
judge run two-sum solution.py --verbose
judge run two-sum solution.py --test 1
```

See [judge run](run.md) for full details.

---

### judge list

List all available problems.

```bash
judge list [flags]
```

**Example:**
```bash
judge list
```

See [judge list](list.md) for full details.

---

### judge show

Show a problem's description and details.

```bash
judge show <problem-id> [flags]
```

**Example:**
```bash
judge show two-sum
```

See [judge show](show.md) for full details.

---

## Configuration

Sandbox Judge can be configured via:

1. **Command-line flags** (highest priority)
2. **Config file** (`~/.judge.yaml`)
3. **Environment variables** (prefix: `JUDGE_`)

### Config File Example

```yaml
# ~/.judge.yaml
problems: /home/user/my-problems
```

### Environment Variables

```bash
export JUDGE_PROBLEMS=/home/user/my-problems
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success (AC for run command) |
| 1 | Error or non-AC verdict |

## Getting Help

For help on any command:

```bash
judge --help
judge run --help
judge list --help
```

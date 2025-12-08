# Sandbox Judge

**A self-hosted, local code evaluation system for practicing coding problems.**

Sandbox Judge lets you run LeetCode-style problems on your own machine with instant feedback. No internet required, no rate limits, complete privacy.

## Features

- 🐳 **Docker-based sandboxing** - Safe, isolated code execution
- ⚡ **Instant feedback** - AC, WA, TLE, RE verdicts in seconds
- 🐍 **Multi-language support** - Python (more coming soon)
- 📝 **YAML problem format** - Easy to create and share problems
- 🎯 **CLI-first design** - Fast, keyboard-driven workflow
- 📊 **Progress tracking** - Track your solved problems (coming soon)

## Quick Example

```bash
# Run a solution against a problem
$ judge run two-sum solution.py
Running two-sum...
  sample/1: AC (45ms)
  sample/2: AC (38ms)

Result: AC (2/2 tests passed)
Total time: 83ms
```

## Why Sandbox Judge?

| Feature | LeetCode | Sandbox Judge |
|---------|----------|---------------|
| Works offline | ❌ | ✅ |
| No rate limits | ❌ | ✅ |
| Custom problems | Limited | ✅ |
| Privacy | ❌ | ✅ |
| Free forever | Freemium | ✅ |

## Getting Started

Ready to start practicing?

1. [Install Sandbox Judge](getting-started/installation.md)
2. [Run your first problem](getting-started/quickstart.md)
3. [Explore CLI commands](cli/overview.md)

## Project Status

Sandbox Judge is under active development. Current capabilities:

- ✅ Python code execution
- ✅ Basic verdicts (AC, WA, TLE, RE)
- ✅ CLI interface
- 🚧 Web UI (planned)
- 🚧 More languages (planned)
- 🚧 Progress tracking (planned)

## License

MIT License - see [LICENSE](https://github.com/marv972228/sandbox_judge/blob/main/LICENSE) for details.

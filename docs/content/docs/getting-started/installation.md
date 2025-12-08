# Installation

This guide will help you install Sandbox Judge on your system.

## Prerequisites

Before installing Sandbox Judge, ensure you have:

- **Go 1.21+** - For building from source
- **Docker** - For sandboxed code execution
- **Make** - For build automation (optional but recommended)

### Verify Prerequisites

```bash
# Check Go version
go version
# Expected: go version go1.21.x or higher

# Check Docker is running
docker info
# Should show Docker daemon information

# Check Make (optional)
make --version
```

## Installation Methods

### From Source (Recommended)

```bash
# Clone the repository
git clone https://github.com/marv972228/sandbox_judge.git
cd sandbox_judge

# Build the binary
make build

# Build Docker images for code execution
make docker-build

# Verify installation
./bin/judge --version
./bin/judge --help
```

### Add to PATH (Optional)

To run `judge` from anywhere:

```bash
# Option 1: Copy to a directory in your PATH
sudo cp ./bin/judge /usr/local/bin/

# Option 2: Add the bin directory to your PATH
echo 'export PATH="$PATH:/path/to/sandbox_judge/bin"' >> ~/.bashrc
source ~/.bashrc
```

## Verify Installation

After installation, verify everything works:

```bash
# Show help
judge --help

# List available problems
judge list

# Run the sample problem
judge run two-sum solutions/two-sum/correct.py
```

You should see output like:

```
Running two-sum...
  sample/1: AC (45ms)
  sample/2: AC (38ms)

Result: AC (2/2 tests passed)
Total time: 83ms
```

## Troubleshooting

### Docker Permission Denied

If you see "permission denied" errors with Docker:

```bash
# Add your user to the docker group
sudo usermod -aG docker $USER

# Log out and back in, or run:
newgrp docker
```

### Docker Image Not Found

If you see "image not found" errors:

```bash
# Build the Docker images
make docker-build

# Verify images exist
docker images | grep sandbox-judge
```

### Go Module Errors

If you see Go module errors:

```bash
# Download dependencies
go mod download

# Tidy modules
go mod tidy
```

## Next Steps

Now that you have Sandbox Judge installed:

- [Quick Start Guide](quickstart.md) - Run your first problem
- [CLI Reference](../cli/overview.md) - Learn all available commands

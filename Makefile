.PHONY: build test run clean help docker-build docker-build-python docs docs-serve docs-build

# Binary name
BINARY=judge
# Build directory
BUILD_DIR=bin

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/marv972228/sandbox_judge/cmd/judge/cmd.Version=$(VERSION)"

## build: Build the binary
build:
	@echo "Building $(BINARY)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/judge

## run: Run the application
run: build
	./$(BUILD_DIR)/$(BINARY) $(ARGS)

## test: Run tests
test:
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.txt coverage.html

## fmt: Format code
fmt:
	go fmt ./...

## lint: Run linter
lint:
	@which golangci-lint > /dev/null || (echo "Install golangci-lint: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

## tidy: Tidy dependencies
tidy:
	go mod tidy

## docker-build: Build all Docker runner images
docker-build: docker-build-python
	@echo "All Docker images built"

## docker-build-python: Build Python runner image
docker-build-python:
	@echo "Building Python runner image..."
	docker build -t sandbox-judge-python:latest ./docker/python

## docs: Build documentation (alias for docs-build)
docs: docs-build

## docs-serve: Start documentation dev server with live reload
docs-serve:
	@echo "Starting docs server at http://localhost:8000..."
	@docker build -t sandbox-judge-docs ./docs
	@docker run --rm -it -p 8000:8000 -v $(PWD)/docs/content:/docs sandbox-judge-docs

## docs-build: Build static documentation site
docs-build:
	@echo "Building documentation..."
	@docker build -t sandbox-judge-docs ./docs
	@docker run --rm -v $(PWD)/docs/content:/docs sandbox-judge-docs mkdocs build
	@echo "Documentation built to docs/content/site/"

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

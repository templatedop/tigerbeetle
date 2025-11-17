.PHONY: help test test-coverage fmt lint vet build clean examples

# Default target
help:
	@echo "Available targets:"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo "  vet            - Run go vet"
	@echo "  build          - Build all packages"
	@echo "  clean          - Clean build artifacts"
	@echo "  examples       - Build all examples"
	@echo "  all            - Run fmt, vet, lint, test"

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

# Run go vet
vet:
	go vet ./...

# Build all packages
build:
	go build ./...

# Build examples
examples:
	@echo "Building examples..."
	@for example in examples/*.go; do \
		echo "Building $$example..."; \
		go build -o /dev/null $$example || exit 1; \
	done
	@echo "All examples built successfully"

# Clean build artifacts
clean:
	go clean ./...
	rm -f coverage.txt coverage.html

# Run all checks
all: fmt vet test
	@echo "All checks passed!"

# Install dependencies
deps:
	go mod download
	go mod tidy

# Update dependencies
update-deps:
	go get -u ./...
	go mod tidy

# Contributing to TigerBeetle Go Library

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Getting Started

1. **Fork the repository**
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-tigerbeetle-lib.git
   cd go-tigerbeetle-lib
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```

## Development Workflow

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**

3. **Run tests**:
   ```bash
   go test ./...
   ```

4. **Run formatting**:
   ```bash
   go fmt ./...
   ```

5. **Run linting** (if you have golangci-lint installed):
   ```bash
   golangci-lint run
   ```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write clear, descriptive variable and function names
- Add comments for exported functions and types
- Keep functions focused and reasonably sized

### Testing

- Write tests for new features
- Ensure all tests pass before submitting PR
- Aim for good test coverage
- Use table-driven tests where appropriate

Example test structure:

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected int
    }{
        {"case 1", 1, 2},
        {"case 2", 2, 4},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := YourFunction(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Documentation

- Update README.md if adding new features
- Add GoDoc comments for all exported types and functions
- Include examples for complex features
- Update CHANGELOG.md (if present)

### Commit Messages

Write clear, descriptive commit messages:

```
Add escrow pattern for secure transactions

- Implement Escrow, ReleaseEscrow, and CancelEscrow methods
- Add comprehensive tests
- Update documentation with examples
```

### Pull Requests

1. **Update your fork**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push your changes**:
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request** on GitHub

4. **PR Description should include**:
   - What changes were made
   - Why these changes are needed
   - Any breaking changes
   - Related issues (if any)

## Areas for Contribution

We welcome contributions in these areas:

### New Features
- Additional transaction patterns
- Enhanced query capabilities
- Performance optimizations
- Additional utility functions

### Documentation
- Improved examples
- Tutorial content
- API documentation
- Architecture diagrams

### Testing
- Additional test cases
- Integration tests
- Benchmarks
- Test utilities

### Bug Fixes
- Bug reports
- Bug fixes
- Edge case handling

## Code Review Process

1. Maintainers will review your PR
2. Address any feedback
3. Once approved, your PR will be merged

## Questions?

- Open an issue for discussion
- Check existing issues and PRs
- Review TigerBeetle documentation

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

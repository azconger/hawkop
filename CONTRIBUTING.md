# Contributing to HawkOp CLI

Thank you for your interest in contributing to HawkOp! This document provides guidelines for contributing to the project.

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git
- Make (for build automation)

### Getting Started

1. Fork and clone the repository:
```bash
git clone https://github.com/your-username/hawkop.git
cd hawkop
```

2. Install development dependencies:
```bash
make dev
```

3. Verify setup by running tests:
```bash
make check
```

## Development Workflow

### Before Making Changes

1. Create a feature branch:
```bash
git checkout -b feature/your-feature-name
```

2. Ensure you have a StackHawk API key for testing:
```bash
./hawkop init
```

### During Development

1. **Write tests first** - Add tests for new functionality
2. **Follow existing patterns** - Look at existing commands for consistency
3. **Run checks frequently**:
```bash
make quick-test  # Fast feedback during development
make check       # Full validation before committing
```

### Code Standards

#### Go Code Style
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add package and function documentation comments
- Handle errors appropriately with user-friendly messages
- Use consistent error message format with ❌ prefix

#### CLI Command Patterns
- Follow existing flag patterns (`--format`, `--org`, `--limit`)
- Implement both table and JSON output formats
- Use pageSize=1000 for all API requests (efficiency)
- Apply filters after pagination to get latest data
- Handle `interface{}` types for flexible API responses

#### Testing Requirements
- Unit tests for all new functionality
- Use testify framework with test suites
- Mock external dependencies (API calls)
- Test both success and error cases
- Maintain or improve code coverage

### API Integration Guidelines

#### Rate Limiting
- Respect StackHawk's 360 requests/minute limit
- Use 167ms intervals between requests
- Implement proper retry logic for 429 responses

#### Authentication
- Handle JWT token refresh on 401 responses
- Never log or expose API keys or tokens
- Use X-ApiKey header format for authentication

#### Error Handling
- Handle all HTTP status codes appropriately
- Provide meaningful error messages to users
- Distinguish between client errors (4xx) and server errors (5xx)

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test suite
go test ./internal/api/
go test ./cmd/

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
make test  # Generates coverage.html
```

### Writing Tests

1. **API Client Tests**: Use mock HTTP servers
2. **Command Tests**: Test command structure and flags
3. **Integration Tests**: Test end-to-end workflows (future)

Example test structure:
```go
type CommandTestSuite struct {
    suite.Suite
    mockClient *api.MockClient
}

func (suite *CommandTestSuite) TestCommand_Success() {
    // Test implementation
}
```

## Code Quality

### Required Checks

Before submitting a PR, ensure:

- [ ] All tests pass: `make test`
- [ ] Code is properly formatted: `go fmt ./...`
- [ ] Static analysis passes: `go vet ./...`
- [ ] Linting passes: `golangci-lint run`
- [ ] No race conditions: `go test -race ./...`

### Automation

The project includes:
- **GitHub Actions CI**: Runs tests, linting, and builds on every PR
- **Release Automation**: GoReleaser builds multi-platform binaries on tags
- **Make Commands**: Consistent development workflow

## Submitting Changes

### Pull Request Process

1. Ensure your code follows all standards above
2. Run `make check` to validate everything locally
3. Push your feature branch and create a PR
4. Provide clear description of changes and rationale
5. Update documentation if needed

### PR Requirements

- All CI checks must pass
- Code coverage should not decrease significantly
- At least one review from maintainer
- Clear commit messages following conventional commits

### Commit Message Format

Use clear, descriptive commit messages:
```
feat: add scan alerts filtering by severity
fix: handle JWT refresh on 401 responses
docs: update installation instructions
test: add integration tests for scan commands
```

## Project Architecture

### Directory Structure
```
hawkop/
├── cmd/                    # CLI command implementations
├── internal/
│   ├── api/               # StackHawk API client
│   ├── config/            # Configuration management
│   └── format/            # Output formatting
├── .github/workflows/     # CI/CD automation
├── examples/              # Example configurations
└── docs/                  # Additional documentation
```

### Adding New Commands

1. Create command file in `cmd/`
2. Add corresponding test file
3. Register command in appropriate parent command
4. Update documentation and help text
5. Add API client methods if needed

### Adding New API Endpoints

1. Add types to `internal/api/types.go`
2. Add client methods to `internal/api/client.go`
3. Add tests to `internal/api/client_test.go`
4. Update mock handlers in `internal/api/mock.go`

## Getting Help

- Check existing issues and documentation
- Ask questions in issue discussions
- Reference [StackHawk API Documentation](https://apidocs.stackhawk.com/)
- Review [CLAUDE.md](CLAUDE.md) for development context

## Security

- Never commit API keys or tokens
- Follow secure coding practices
- Report security issues privately to maintainers
- This is a defensive security tool - no malicious functionality

Thank you for contributing to HawkOp!
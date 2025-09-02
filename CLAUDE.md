# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HawkOp is a professional-grade Go CLI companion utility for the StackHawk scanner and platform. It provides developers and security teams with streamlined access to StackHawk's dynamic application security testing (DAST) capabilities directly from the terminal.

This project follows enterprise Go development standards and includes comprehensive testing, documentation, and CI/CD workflows.

The CLI follows GitHub's `gh` CLI design patterns and supports core StackHawk operations including:
- Authentication and API key management (`init`, `status`)
- Organization and resource management (`org list`, `org set`, `org get`, `org clear`)
- User management (`user list` with role filtering)
- Team management (`team list`)
- Application management (`app list` with status filtering)
- **Scan management (`scan list`, `scan get`, `scan alerts`)**
- Version information (`version`)

## Project Structure

- `main.go` - The main application entry point
- `go.mod` - Go module definition specifying Go 1.24
- `LICENSE` - MIT license for open source compliance
- `README.md` - Comprehensive project documentation with installation and usage
- `CONTRIBUTING.md` - Development guidelines and contribution process
- `.gitignore` - Professional Go project gitignore
- `Makefile` - Build automation with proper ldflags and version handling
- `.goreleaser.yaml` - Multi-platform release automation
- `cmd/` - Cobra CLI command definitions
  - `root.go` - Root command and CLI setup
  - `init.go` - API key initialization command
  - `status.go` - Configuration status command
  - `version.go` - Version information command (updated for new version package)
  - `org.go` - Organization management commands
  - `user.go` - User listing and management commands
  - `team.go` - Team listing and management commands
  - `app.go` - Application listing and management commands
  - `scan.go` - Scan listing and analysis commands
- `internal/config/` - Configuration management with YAML support
  - `config.go` - YAML config file handling, JWT management
  - `config_test.go` - Comprehensive test suite with 35.1% coverage
- `internal/api/` - StackHawk API client
  - `client.go` - HTTP client with automatic JWT refresh
  - `client_test.go` - Complete test suite with mock server (34.4% coverage)
  - `types.go` - API response data structures
- `internal/format/` - Output formatting utilities
  - `table.go` - Table formatting for CLI output
  - `table_test.go` - Complete test suite with 100% coverage
- `internal/version/` - Build-time version information
  - `version.go` - Version, build time, and git commit handling
- `examples/` - Example configuration files
- `HAWKOPDOC.md` - Comprehensive CLI documentation and command specifications

## Common Development Commands

### Professional Build System (Preferred)
- `make build` - Build with proper version ldflags
- `make test` - Run full test suite with coverage
- `make lint` - Run formatter, vet, and linter
- `make clean` - Clean build artifacts
- `make install` - Install to GOPATH/bin
- `make release-snapshot` - Test release build locally

### Direct Go Commands (Alternative)
- `go run main.go` - Compile and run the application directly
- `go build` - Compile the application into an executable binary
- `go test ./...` - Run all tests including API client and command tests
- `go test ./internal/api/` - Run API client tests with mock server
- `go test ./cmd/` - Run command structure and flag tests
- `go vet ./...` - Run Go's built-in static analyzer
- `go fmt ./...` - Format Go source code
- `go mod tidy` - Clean up module dependencies

## Architecture Notes

HawkOp is now a fully functional CLI application with comprehensive StackHawk platform integration. The architecture includes:

### Command Structure
- Root command with subcommands following `hawkop <command> <subcommand>` pattern
- **Implemented Commands**: `init`, `version`, `status`, `org`, `app`, `user`, `team`, `scan`
- **Future Commands**: `policy`, `finding`, `report`
- **Subcommands**: `list`, `get`, `set`, `clear`, `alerts`
- **Consistent flag patterns**: `--format`, `--org`, `--limit`, `--app`, `--env`, `--status`, `--severity`

### Configuration Management  
- API key storage in `~/.config/hawkop/` (encrypted)
- Default organization and CLI preferences
- Connection status tracking

### Current Dependencies
- `github.com/spf13/cobra` - CLI framework for command structure
- `golang.org/x/term` - Terminal utilities for secure password input
- `github.com/stretchr/testify` - Testing framework with assertions, mocks, and test suites
- `gopkg.in/yaml.v3` - YAML parsing and marshaling for configuration
- Standard library HTTP client for StackHawk API

### Implemented Features ✅
1. ✅ **CLI Framework** - Complete Cobra-based CLI with hierarchical commands
2. ✅ **Security & Auth** - Secure config storage, JWT management, API key protection
3. ✅ **Configuration** - `~/.config/hawkop/config.yaml` with 600 permissions
4. ✅ **Core Commands** - `init`, `status`, `version`, `org`, `user`, `team`, `app`, `scan`
5. ✅ **Resource Management** - List and manage orgs, users, teams, applications
6. ✅ **Scan Analysis** - List scans, view details, analyze security alerts
7. ✅ **Smart Filtering** - App/env/status/role/severity filters across commands
8. ✅ **Output Formats** - Professional table formatting + JSON for automation
9. ✅ **Enterprise Ready** - Pagination, organization awareness, role-based access
10. ✅ **Real Security Data** - Live integration showing actual vulnerability findings
11. ✅ **Extensible Architecture** - Clean patterns for adding new commands and reports
12. ✅ **Production Quality** - Error handling, validation, user-friendly messaging
13. ✅ **Testing Infrastructure** - Comprehensive test suites with testify framework
14. ✅ **CI/CD Pipeline** - GitHub Actions for automated testing and releases
15. ✅ **Release Management** - GoReleaser for multi-platform binary distribution
16. ✅ **Professional Standards** - MIT license, comprehensive docs, contribution guidelines
17. ✅ **Build Automation** - Makefile with proper version injection and workflows
18. ✅ **YAML Configuration** - Human-readable YAML format for configuration files

### API Endpoints Integrated
- `GET /api/v1/auth/login` - JWT authentication with X-ApiKey header
- `GET /api/v1/user` - Get current user info and organizations
- `GET /api/v1/org/{orgId}/members` - List organization members/users (max pageSize=1000)
- `GET /api/v1/org/{orgId}/teams` - List organization teams (max pageSize=1000)
- `GET /api/v2/org/{orgId}/apps` - List organization applications (max pageSize=1000)
- `GET /api/v1/scan/{orgId}` - List organization scans with metadata (max pageSize=1000)
- `GET /api/v1/scan/{scanId}/alerts` - Get security alerts for specific scan

### API Standards Compliance ✅
- **Rate Limiting**: 360 requests/minute compliance with 167ms intervals
- **Pagination**: Default pageSize=1000 (maximum) to minimize API requests
- **Error Handling**: Comprehensive HTTP status code handling (400, 401, 403, 404, 409, 422, 429)
- **Retry Logic**: Automatic JWT refresh on 401, rate limit retry on 429
- **Query Parameters**: Proper URL encoding and parameter validation

### Future Enhancement Opportunities
1. **Advanced Scan Features**
   - `scan finding` - Individual finding details with request/response data
   - `scan message` - Raw HTTP request/response analysis
   - Scan filtering by date ranges and criticality

2. **Enterprise Reporting** 
   - `app summary` - Cross-application security posture dashboard
   - `app report` - MTTR analysis, scan coverage, policy compliance
   - Historical trending and ROI metrics

3. **Application Deep Dive**
   - `app get` - Application metadata, configuration, policy assignment
   - `app scans` - Scan history and trends for specific applications
   - Attack Surface repository mappings

4. **Policy & Configuration Management**
   - `policy list` - Available scan policies
   - Policy assignment and configuration
   - Environment and configuration management

5. **Advanced Features**
   - Export capabilities (CSV, PDF reports)
   - Interactive mode for guided workflows
   - Scan result comparison and diff analysis

### Configuration File Format
The config file is stored at `~/.config/hawkop/config.yaml` with 600 permissions:
```yaml
api_key: your-api-key
org_id: optional-default-org-id
jwt:
  token: jwt-token
  expires_at: 2024-12-25T15:30:45Z
```

## Development Standards Checklist

### Before Submitting Code
- [ ] All checks pass: `make lint` (includes fmt, vet, golangci-lint)
- [ ] All tests pass: `make test` (includes race detection and coverage)
- [ ] Build succeeds: `make build` (with proper version injection)
- [ ] API rate limiting respected (167ms intervals)
- [ ] Error handling implemented for all API calls
- [ ] Proper JWT token refresh on 401 responses
- [ ] Consistent flag patterns across commands
- [ ] Table output formatted correctly
- [ ] JSON output properly structured
- [ ] Security: No secrets in code or logs
- [ ] YAML configuration properly handled

### Testing Requirements
- Unit tests for all API client methods
- Command structure tests for all CLI commands
- Flag validation tests
- Mock server tests for API integration
- Error case testing (401, 404, rate limits)
- Input validation testing

### Code Style Guidelines
- Follow existing patterns in `cmd/` directory
- Use consistent error messages with ❌ prefix
- Implement both table and JSON output formats
- Use pageSize=1000 for all API requests
- Apply filters after pagination for latest data
- Handle interface{} types for flexible API responses

## Reference Materials

### StackHawk Resources

- StackHawk API OpenAPI Spec: https://download.stackhawk.com/openapi/stackhawk-openapi.json
- StackHawk API Authentication Reference: https://apidocs.stackhawk.com/reference/login
- StackHawk API Documentation: https://apidocs.stackhawk.com/
- StackHawk Documentation: https://docs.stackhawk.com/
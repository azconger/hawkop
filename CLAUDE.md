# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HawkOp is a Go CLI companion utility for the StackHawk scanner and platform. It provides developers and security teams with streamlined access to StackHawk's dynamic application security testing (DAST) capabilities directly from the terminal.

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
- `cmd/` - Cobra CLI command definitions
  - `root.go` - Root command and CLI setup
  - `init.go` - API key initialization command
  - `status.go` - Configuration status command
  - `version.go` - Version information command
  - `org.go` - Organization management commands
  - `user.go` - User listing and management commands
  - `team.go` - Team listing and management commands
  - `app.go` - Application listing and management commands
  - `scan.go` - Scan listing and analysis commands
- `internal/config/` - Configuration management
  - `config.go` - Config file handling, JWT management
- `internal/api/` - StackHawk API client
  - `client.go` - HTTP client with automatic JWT refresh
  - `types.go` - API response data structures
- `internal/format/` - Output formatting utilities
  - `table.go` - Table formatting for CLI output
- `examples/` - Example configuration files
- `HAWKOPDOC.md` - Comprehensive CLI documentation and command specifications

## Common Development Commands

### Building and Running
- `go run main.go` - Compile and run the application directly
- `go build` - Compile the application into an executable binary
- `go build -o hawkop` - Build with a specific output name

### Testing and Quality
- `go test` - Run tests (though no test files exist currently)
- `go vet` - Run Go's built-in static analyzer
- `go fmt` - Format Go source code

### Module Management
- `go mod tidy` - Clean up module dependencies
- `go mod download` - Download module dependencies

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
- Standard library HTTP client for StackHawk API

### Implemented Features ✅
1. ✅ **CLI Framework** - Complete Cobra-based CLI with hierarchical commands
2. ✅ **Security & Auth** - Secure config storage, JWT management, API key protection
3. ✅ **Configuration** - `~/.config/hawkop/config.json` with 600 permissions
4. ✅ **Core Commands** - `init`, `status`, `version`, `org`, `user`, `team`, `app`, `scan`
5. ✅ **Resource Management** - List and manage orgs, users, teams, applications
6. ✅ **Scan Analysis** - List scans, view details, analyze security alerts
7. ✅ **Smart Filtering** - App/env/status/role/severity filters across commands
8. ✅ **Output Formats** - Professional table formatting + JSON for automation
9. ✅ **Enterprise Ready** - Pagination, organization awareness, role-based access
10. ✅ **Real Security Data** - Live integration showing actual vulnerability findings
11. ✅ **Extensible Architecture** - Clean patterns for adding new commands and reports
12. ✅ **Production Quality** - Error handling, validation, user-friendly messaging

### API Endpoints Integrated
- `GET /api/v1/auth/login` - JWT authentication with X-ApiKey header
- `GET /api/v1/user` - Get current user info and organizations
- `GET /api/v1/org/{orgId}/members` - List organization members/users  
- `GET /api/v1/org/{orgId}/teams` - List organization teams
- `GET /api/v2/org/{orgId}/apps` - List organization applications
- `GET /api/v1/scan/{orgId}` - List organization scans with metadata
- `GET /api/v1/scan/{scanId}/alerts` - Get security alerts for specific scan

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
The config file is stored at `~/.config/hawkop/config.json` with 600 permissions:
```json
{
  "api_key": "your-api-key",
  "org_id": "optional-default-org-id", 
  "jwt": {
    "token": "jwt-token",
    "expires_at": "2024-12-25T15:30:45Z"
  }
}
```

## Reference Materials

### StackHawk Resources

- StackHawk API OpenAPI Spec: https://download.stackhawk.com/openapi/stackhawk-openapi.json
- StackHawk API Authentication Reference: https://apidocs.stackhawk.com/reference/login
- StackHawk API Documentation: https://apidocs.stackhawk.com/
- StackHawk Documentation: https://docs.stackhawk.com/
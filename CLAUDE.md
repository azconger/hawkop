# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HawkOp is a Go CLI companion utility for the StackHawk scanner and platform. It provides developers and security teams with streamlined access to StackHawk's dynamic application security testing (DAST) capabilities directly from the terminal.

The CLI follows GitHub's `gh` CLI design patterns and supports core StackHawk operations including:
- Authentication and API key management (`init`, `status`)
- Organization and resource management (`org list`, `app list`, `user list`, `team list`)  
- Security policy management (`policy list`)
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
- `internal/config/` - Configuration management
  - `config.go` - Config file handling, JWT management
- `internal/api/` - StackHawk API client
  - `client.go` - HTTP client with automatic JWT refresh
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

This is currently a minimal Go application with placeholder code that needs to be replaced with the actual CLI implementation. The planned architecture includes:

### Command Structure
- Root command with subcommands following `hawkop <command> <subcommand>` pattern
- Commands: `init`, `version`, `status`, `org`, `app`, `user`, `team`, `policy`
- Subcommands: `list` for resource commands
- Consistent flag patterns: `--format`, `--org`, `--limit`, `--type`

### Configuration Management  
- API key storage in `~/.config/hawkop/` (encrypted)
- Default organization and CLI preferences
- Connection status tracking

### Current Dependencies
- `github.com/spf13/cobra` - CLI framework for command structure
- `golang.org/x/term` - Terminal utilities for secure password input
- Standard library HTTP client for StackHawk API
- StackHawk API OpenAPI Spec: https://download.stackhawk.com/openapi/stackhawk-openapi.json

### Implemented Features ✅
1. ✅ Basic CLI structure with Cobra framework
2. ✅ Configuration management (`~/.config/hawkop/config.json`)
3. ✅ Credential storage with file permissions (600)
4. ✅ JWT token management with expiration checking
5. ✅ Commands: `init`, `status`, `version`, `org set/get/clear`
6. ✅ API client with automatic JWT refresh
7. ✅ Secure API key input (hidden from terminal)

### Next Implementation Steps
1. Resource listing commands (`org list`, `app list`, `user list`, `team list`, `policy list`)
2. API endpoints integration using StackHawk OpenAPI spec
3. Table/JSON output formatting
4. Error handling and user-friendly messages
5. Additional configuration options

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
# HawkOp Documentation

HawkOp is a command-line interface (CLI) companion utility for the StackHawk scanner and platform. It provides developers and security teams with streamlined access to StackHawk's dynamic application security testing (DAST) capabilities directly from the terminal.

## Overview

HawkOp follows the design patterns established by GitHub's `gh` CLI, offering an intuitive and familiar experience for developers already comfortable with modern CLI tools. The tool enables efficient management of StackHawk resources including organizations, applications, users, teams, and scan policies.

## Installation

```bash
# Download and install hawkop (installation method TBD)
```

## Getting Started

### Initial Setup

Before using HawkOp, you'll need to authenticate with your StackHawk API key:

```bash
hawkop init
```

This command will prompt you to enter your StackHawk API key, which will be securely cached for future use.

## Commands

### Authentication & Setup

#### `hawkop init`
Cache your StackHawk API key for authenticated requests.

```bash
hawkop init
```

**Interactive prompts:**
- StackHawk API Key
- Default organization (optional)

### General

#### `hawkop version`
Display the current version of HawkOp CLI.

```bash
hawkop version
```

#### `hawkop status`
Display CLI connection status and current context information.

```bash
hawkop status
```

**Output includes:**
- Connection status (connected/disconnected)
- Current organization name
- Authenticated user name
- API endpoint
- CLI version

### Organization Management

#### `hawkop org list`
List all organizations you belong to.

```bash
hawkop org list
```

**Options:**
- `--format table|json` - Output format (default: table)
- `--limit N` - Limit number of results

### Application Management

#### `hawkop app list`
List applications in the current organization.

```bash
hawkop app list
```

**Options:**
- `--org ORG_NAME` - Specify organization (uses default if not provided)
- `--format table|json` - Output format (default: table)
- `--limit N` - Limit number of results

### User Management

#### `hawkop user list`
List users in the current organization.

```bash
hawkop user list
```

**Options:**
- `--org ORG_NAME` - Specify organization (uses default if not provided)
- `--format table|json` - Output format (default: table)
- `--role admin|member|viewer` - Filter by user role
- `--limit N` - Limit number of results

### Team Management

#### `hawkop team list`
List teams in the current organization.

```bash
hawkop team list
```

**Options:**
- `--org ORG_NAME` - Specify organization (uses default if not provided)
- `--format table|json` - Output format (default: table)
- `--limit N` - Limit number of results

### Policy Management

#### `hawkop policy list`
List all available scan policies (both built-in and custom).

```bash
hawkop policy list
```

**Options:**
- `--org ORG_NAME` - Specify organization for custom policies
- `--format table|json` - Output format (default: table)
- `--type builtin|custom|all` - Filter by policy type (default: all)
- `--limit N` - Limit number of results

## Configuration

HawkOp stores configuration in `~/.config/hawkop/` including:
- API key (encrypted)
- Default organization
- CLI preferences

## Common Workflows

### Quick Status Check
```bash
hawkop status
hawkop org list
hawkop app list
```

### Setting Up a New Environment
```bash
hawkop init
hawkop org list
hawkop app list --org "My Organization"
```

### Policy Review
```bash
hawkop policy list --type builtin
hawkop policy list --type custom --org "My Organization"
```

## Output Formats

Most list commands support multiple output formats:

- **table** (default): Human-readable tabular format
- **json**: Machine-readable JSON format for scripting

Example:
```bash
hawkop app list --format json | jq '.[] | select(.name | contains("prod"))'
```

## Error Handling

HawkOp provides clear error messages and suggestions for common issues:
- Authentication failures
- Network connectivity problems
- Invalid organization or resource names
- Permission denied scenarios

## Support

For issues, feature requests, or questions:
- GitHub Issues: [Repository URL]
- StackHawk Documentation: [StackHawk Docs URL]
- Community Support: [Community URL]
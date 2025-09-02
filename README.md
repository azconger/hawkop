# HawkOp CLI

A powerful command-line companion for the StackHawk dynamic application security testing (DAST) platform. HawkOp provides developers and security teams with streamlined access to StackHawk's capabilities directly from the terminal.

## Features

- 🔐 **Secure Authentication** - API key management with automatic JWT token refresh
- 🏢 **Organization Management** - List and manage organizations, set defaults
- 👥 **User Management** - List organization members with role filtering
- 🏗️ **Team Management** - View and manage teams within organizations
- 📱 **Application Management** - List applications with status filtering
- 🔍 **Scan Management** - List scans, view details, analyze security alerts
- 📊 **Flexible Output** - Table or JSON format support for automation
- 🔧 **Configuration** - Persistent settings and credential storage
- 🚀 **Enterprise Ready** - Rate limiting, pagination, role-based access

## Installation

### Download Release Binary

Download the latest release for your platform from [GitHub Releases](https://github.com/azconger/hawkop/releases).

### Build from Source

```bash
git clone https://github.com/azconger/hawkop.git
cd hawkop
make build
```

### Install with Go

```bash
go install github.com/azconger/hawkop@latest
```

## Quick Start

First, initialize HawkOp with your StackHawk API key:

```bash
./hawkop init
```

## Commands

### Authentication & Configuration

```bash
# Initialize with API key
hawkop init

# Check configuration status
hawkop status

# Show version information
hawkop version
```

### Organization Management

```bash
# List all organizations
hawkop org list

# Set default organization
hawkop org set <org-id>

# Get current default organization
hawkop org get

# Clear default organization
hawkop org clear
```

### User Management

```bash
# List all users in organization
hawkop user list

# List users with specific role
hawkop user list --role admin

# Limit results and use JSON format
hawkop user list --limit 5 --format json

# Use specific organization
hawkop user list --org <org-id>
```

### Team Management

```bash
# List all teams in organization
hawkop team list

# Limit results
hawkop team list --limit 10

# JSON output
hawkop team list --format json
```

### Application Management

```bash
# List all applications
hawkop app list

# Filter by application status
hawkop app list --status ACTIVE

# Limit and format results
hawkop app list --limit 5 --format json

# Use specific organization
hawkop app list --org <org-id>
```

### Scan Management

```bash
# List recent scans (default: 100 latest)
hawkop scan list

# Filter by application and environment
hawkop scan list --app "My App" --env production

# Filter by scan status
hawkop scan list --status COMPLETED

# Get detailed scan information
hawkop scan get <scan-id>

# View scan statistics
hawkop scan get <scan-id> --view stats

# List security alerts for a scan
hawkop scan alerts <scan-id>

# Filter alerts by severity
hawkop scan alerts <scan-id> --severity High
```

## Configuration

HawkOp stores configuration in `~/.config/hawkop/config.json` with secure file permissions (600). The configuration includes:

- API key (encrypted storage)
- Default organization ID
- JWT tokens with automatic refresh

## Output Formats

### Table Format (Default)
```
ID                                    NAME              STATUS  TYPE
------------------------------------  ----------------  ------  --------
058b994a-b95e-4562-ad0a-de8175164c60  api1              ACTIVE  STANDARD
acbf5d2d-e3e3-4e06-808a-e5085ba525db  Broken Crystals   ACTIVE  STANDARD
```

### JSON Format
```json
[
  {
    "applicationId": "058b994a-b95e-4562-ad0a-de8175164c60",
    "name": "api1",
    "applicationStatus": "ACTIVE",
    "applicationType": "STANDARD"
  }
]
```

## Common Flags

- `--format, -f` - Output format (table|json)
- `--limit, -l` - Limit number of results (0 = no limit)
- `--org, -o` - Override default organization
- `--role, -r` - Filter by user role (admin|member|owner)
- `--status, -s` - Filter by application status (ACTIVE|ENV_INCOMPLETE)

## API Integration

HawkOp integrates with the StackHawk API using the following endpoints:

- **Authentication**: `GET /api/v1/auth/login`
- **User Info**: `GET /api/v1/user`
- **Organization Members**: `GET /api/v1/org/{orgId}/members`
- **Organization Teams**: `GET /api/v1/org/{orgId}/teams`
- **Organization Applications**: `GET /api/v2/org/{orgId}/apps`

## Development

### Requirements

- Go 1.24+
- StackHawk API key

### Build

```bash
go build -o hawkop .
```

### Development Commands

```bash
# Build the binary
make build

# Run all tests with coverage
make test

# Format code and run all checks
make check

# Install development tools
make dev

# Create a new release
make release
```

### Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/stretchr/testify` - Testing framework
- `golang.org/x/term` - Secure terminal input

## Security

- API keys are stored securely with file permissions 600
- JWT tokens are automatically refreshed as needed
- No sensitive data is logged or exposed in output
- Rate limiting respects StackHawk's 360 requests/minute limit

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues related to HawkOp CLI, please create an issue in this repository.
For StackHawk platform support, visit [StackHawk Documentation](https://docs.stackhawk.com/).
# cfcli - Cloudflare DNS CLI

A fast, powerful command-line interface for managing Cloudflare DNS records written in Go. This is a Go implementation inspired by [danielpigott/cloudflare-cli](https://github.com/danielpigott/cloudflare-cli).

## Features

- Full CRUD operations on DNS records (Create, Read, Update, Delete)
- Support for all common DNS record types (A, AAAA, CNAME, MX, TXT, NS, SRV)
- Multiple output formats (table, JSON, CSV)
- Configuration file support for managing multiple accounts
- Advanced filtering and querying capabilities
- Support for Cloudflare API Tokens and API Keys (legacy)

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/skyline/cfcli/releases).

#### Linux/macOS

```bash
# Download the appropriate binary for your platform
# Example for Linux x86_64:
wget https://github.com/skyline/cfcli/releases/download/v1.0.0/cfcli_1.0.0_Linux_x86_64.tar.gz
tar -xzf cfcli_1.0.0_Linux_x86_64.tar.gz
chmod +x cfcli
sudo mv cfcli /usr/local/bin/
```

#### Windows

Download the `.zip` file for your architecture and extract `cfcli.exe` to a directory in your PATH.

### Homebrew (macOS/Linux)

```bash
brew install skyline/tap/cfcli
```

### Scoop (Windows)

```bash
scoop bucket add skyline https://github.com/skyline/scoop-bucket
scoop install cfcli
```

### Using Go Install

```bash
go install github.com/skyline/cfcli@latest
```

### Docker

```bash
# Pull the image
docker pull ghcr.io/skyline/cfcli:latest

# Run with environment variables
docker run --rm -e CF_API_KEY=your-token -e CF_API_DOMAIN=example.com ghcr.io/skyline/cfcli:latest ls

# Or with config file
docker run --rm -v ~/.cfcli.yml:/root/.cfcli.yml ghcr.io/skyline/cfcli:latest ls
```

### From Source

```bash
git clone https://github.com/skyline/cfcli.git
cd cfcli
make build
sudo make install
```

## Configuration

You can configure `cfcli` using a YAML configuration file located at `$HOME/.cfcli.yml`.

### Single Account Setup

```yaml
defaults:
    token: <your-cloudflare-api-token>
    domain: example.com
```

### Multiple Accounts Setup

```yaml
defaults:
    account: work
accounts:
    work:
        token: <cloudflare-token>
        domain: example.com
    personal:
        token: <cloudflare-token>
        domain: mysite.com
```

Then use `-u personal` to switch between accounts.

### API Tokens vs API Keys

- **API Tokens** (recommended): Use only the token without email
  ```yaml
  token: <cloudflare-token>
  ```

- **API Keys** (legacy): Require both token and email
  ```yaml
  token: <cloudflare-api-key>
  email: you@example.com
  ```

### Environment Variables

You can also use environment variables:

```bash
export CF_API_KEY=your-token
export CF_API_EMAIL=you@example.com  # Only for API Keys
export CF_API_DOMAIN=example.com
```

## Usage

### List Zones

```bash
cfcli -k <token> zones
```

### List DNS Records

```bash
# List all records for a domain
cfcli -d example.com -k <token> ls

# List records with filtering
cfcli -d example.com -k <token> -q name:test ls

# Output as JSON
cfcli -d example.com -k <token> -f json ls

# Output as CSV
cfcli -d example.com -k <token> -f csv ls
```

### Add DNS Record

```bash
# Add an A record
cfcli -d example.com -k <token> -t A add mail 1.2.3.4

# Add a CNAME record
cfcli -d example.com -k <token> -t CNAME add www example.com

# Add an A record with proxy enabled
cfcli -d example.com -k <token> -t A -a add test 1.1.1.1

# Add an MX record with priority
cfcli -d example.com -k <token> -t MX -p 10 add @ mail.example.com

# Add a record with custom TTL
cfcli -d example.com -k <token> -t A --ttl 300 add test 5.6.7.8
```

### Edit DNS Record

```bash
# Edit an A record
cfcli -d example.com -k <token> -t A edit mail 5.6.7.8

# Change record type from A to CNAME
cfcli -d example.com -k <token> -t A -n CNAME edit test example.com

# Edit record with custom TTL
cfcli -d example.com -k <token> -t A --ttl 300 edit mail 1.2.3.4
```

### Find DNS Record

```bash
# Find records by name
cfcli -d example.com -k <token> find test

# Find records by type
cfcli -d example.com -k <token> -t A find test

# Find records with filters
cfcli -d example.com -k <token> find -q content:1.1.1.1
```

### Delete DNS Record

```bash
# Delete all records with a name
cfcli -d example.com -k <token> rm test

# Delete specific record by type
cfcli -d example.com -k <token> -t A rm test

# Delete with additional filters
cfcli -d example.com -k <token> -t A rm test -q content:1.1.1.1
```

## Command Line Options

```
Flags:
  -u, --account string   Named account from config file
  -a, --activate         Activate cloudflare (enable proxy) after creating record
  -c, --config string    config file (default is $HOME/.cfcli.yml)
  -d, --domain string    Domain to operate on
  -e, --email string     Email of your cloudflare account
  -f, --format string    Output format: table, json, csv (default "table")
  -h, --help             help for cfcli
  -n, --newtype string   New type when editing a record
  -p, --priority int     Priority for MX or SRV records
  -q, --query string     Comma-separated filters (e.g., content:1.1.1.1,type:A)
  -k, --token string     API token for your cloudflare account
  -l, --ttl int          TTL in seconds (1 for auto, 120-86400) (default 1)
  -t, --type string      Type of DNS record (A, AAAA, CNAME, MX, TXT, NS, SRV)
```

## Examples

### Complete Workflow

```bash
# List all zones
cfcli -k <token> zones

# List all DNS records for a domain
cfcli -d example.com -k <token> ls

# Add a new A record
cfcli -d example.com -k <token> -t A add test 192.0.2.1

# Verify it was created
cfcli -d example.com -k <token> -q name:test ls

# Update the record
cfcli -d example.com -k <token> -t A edit test.example.com 192.0.2.2

# Change it to a CNAME
cfcli -d example.com -k <token> -t A -n CNAME edit test.example.com target.com

# Delete the record
cfcli -d example.com -k <token> -t CNAME rm test.example.com
```

### Using Config File

Create `~/.cfcli.yml`:
```yaml
defaults:
    token: tJzwm7RiIrLz9iQiHMAtjFrkFg1Z1cA7ap_1yGTf
    domain: example.com
```

Then you can omit the token and domain from commands:
```bash
cfcli ls
cfcli -t A add test 192.0.2.1
cfcli -t A edit test 192.0.2.2
cfcli -t A rm test
```

## Supported DNS Record Types

- **A**: IPv4 address record
- **AAAA**: IPv6 address record
- **CNAME**: Canonical name record
- **MX**: Mail exchange record (supports priority)
- **TXT**: Text record
- **NS**: Name server record
- **SRV**: Service record (supports priority)

## Comparison with Node.js Version

This Go implementation provides:
- ✅ Faster execution (compiled binary vs interpreted JavaScript)
- ✅ Single binary distribution (no npm/node dependencies)
- ✅ Lower memory footprint
- ✅ Cross-platform builds with no runtime dependencies
- ✅ Compatible command-line interface
- ✅ All features from the original cloudflare-cli

## Building from Source

```bash
# Clone the repository
git clone https://github.com/skyline/cfcli.git
cd cfcli

# Install dependencies
go mod download

# Build using Make
make build

# Or build manually
go build -o cfcli

# Install system-wide
sudo make install
# Or manually
sudo mv cfcli /usr/local/bin/
```

## Development

### Using Make

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean

# Install to ~/bin
make dev-install

# Create release snapshot (requires goreleaser)
make release-snapshot
```

### Manual Cross-Platform Builds

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o cfcli-linux-amd64

# Linux ARM64 (Raspberry Pi, ARM servers)
GOOS=linux GOARCH=arm64 go build -o cfcli-linux-arm64

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o cfcli-darwin-amd64

# macOS Apple Silicon (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o cfcli-darwin-arm64

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o cfcli-windows-amd64.exe

# FreeBSD
GOOS=freebsd GOARCH=amd64 go build -o cfcli-freebsd-amd64
```

## Automated Releases

Releases are automated using GitHub Actions and GoReleaser. Every time a version tag (e.g., `v1.0.0`) is pushed to the repository:

1. Binaries are built for all supported platforms
2. Docker images are created and pushed to GitHub Container Registry
3. Checksums are generated
4. Release notes are automatically created
5. Artifacts are uploaded to GitHub Releases

Supported platforms include:
- Linux (amd64, arm64, arm, 386)
- macOS/Darwin (amd64, arm64)
- Windows (amd64, 386)
- FreeBSD, OpenBSD, NetBSD (amd64, arm64, arm, 386)

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Credits

Inspired by [danielpigott/cloudflare-cli](https://github.com/danielpigott/cloudflare-cli)

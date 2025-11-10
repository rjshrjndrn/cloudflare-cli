# cfcli Cheat Sheet

Quick reference for common cfcli commands.

## Installation

```bash
# Homebrew
brew install skyline/tap/cfcli

# Download binary
curl -L https://github.com/skyline/cfcli/releases/latest/download/cfcli_Linux_x86_64.tar.gz | tar xz

# Docker
docker pull ghcr.io/skyline/cfcli:latest
```

## Configuration

```bash
# Create config file
cat > ~/.cfcli.yml << 'EOF'
defaults:
    token: YOUR_API_TOKEN
    domain: example.com
EOF
```

## Basic Commands

```bash
# List all zones
cfcli -k TOKEN zones

# List all DNS records
cfcli -k TOKEN -d example.com ls

# List with filters
cfcli -d example.com -q name:www ls
cfcli -d example.com -t A ls
```

## Create Records

```bash
# A record
cfcli -d example.com -t A add test 192.0.2.1

# A record with proxy enabled
cfcli -d example.com -t A -a add www 192.0.2.1

# CNAME record
cfcli -d example.com -t CNAME add blog target.com

# MX record with priority
cfcli -d example.com -t MX -p 10 add @ mail.example.com

# TXT record
cfcli -d example.com -t TXT add @ "v=spf1 include:_spf.google.com ~all"

# Custom TTL
cfcli -d example.com -t A --ttl 300 add test 192.0.2.1
```

## Update Records

```bash
# Update IP
cfcli -d example.com -t A edit test.example.com 192.0.2.2

# Change type A to CNAME
cfcli -d example.com -t A -n CNAME edit test.example.com target.com

# Update TTL
cfcli -d example.com -t A --ttl 600 edit test.example.com 192.0.2.1
```

## Delete Records

```bash
# Delete by name and type
cfcli -d example.com -t A rm test.example.com

# Delete with filter
cfcli -d example.com rm test -q content:192.0.2.1
```

## Search/Find

```bash
# Find by name
cfcli -d example.com find test

# Find by content
cfcli -d example.com find -q content:192.0.2.1

# Find by type
cfcli -d example.com -t A find test
```

## Output Formats

```bash
# Table (default)
cfcli -d example.com ls

# JSON
cfcli -d example.com -f json ls

# CSV
cfcli -d example.com -f csv ls > records.csv
```

## Multiple Accounts

```yaml
# ~/.cfcli.yml
defaults:
    account: work
accounts:
    work:
        token: WORK_TOKEN
        domain: work.com
    personal:
        token: PERSONAL_TOKEN
        domain: personal.com
```

```bash
# Use default account
cfcli ls

# Use specific account
cfcli -u personal ls
```

## Environment Variables

```bash
export CF_API_KEY=your-token
export CF_API_DOMAIN=example.com

# Now run without flags
cfcli ls
```

## Common Patterns

```bash
# Add www record pointing to root
cfcli -d example.com -t CNAME add www example.com

# Add mail server
cfcli -d example.com -t A add mail 192.0.2.10
cfcli -d example.com -t MX -p 10 add @ mail.example.com

# Add SPF record
cfcli -d example.com -t TXT add @ "v=spf1 include:_spf.example.com ~all"

# Wildcard record
cfcli -d example.com -t A add "*" 192.0.2.1

# Root domain A record
cfcli -d example.com -t A add @ 192.0.2.1
```

## Filters

```bash
# By content
-q content:192.0.2.1

# By type
-q type:A

# By name
-q name:www

# Multiple filters
-q content:192.0.2.1,type:A
```

## Tips

- Use `-a` flag to enable Cloudflare proxy (orange cloud)
- TTL of `1` means "auto"
- Use `@` for root domain in record name
- Combine `-q` with `-t` for precise filtering
- Export to CSV: `-f csv > file.csv`
- Chain commands with `&&` for batch operations

## Help

```bash
# General help
cfcli --help

# Command specific help
cfcli add --help
cfcli edit --help

# Version
cfcli --version
```

## Docker Usage

```bash
# With environment variables
docker run --rm \
  -e CF_API_KEY=token \
  -e CF_API_DOMAIN=example.com \
  ghcr.io/skyline/cfcli:latest ls

# With config file
docker run --rm \
  -v ~/.cfcli.yml:/root/.cfcli.yml \
  ghcr.io/skyline/cfcli:latest ls
```

## Quick Examples

```bash
# Complete workflow
cfcli zones                                              # List zones
cfcli -d example.com ls                                  # List records
cfcli -d example.com -t A add test 192.0.2.1            # Add record
cfcli -d example.com -q name:test ls                     # Verify
cfcli -d example.com -t A edit test.example.com 192.0.2.2  # Update
cfcli -d example.com -t A rm test.example.com           # Delete
```

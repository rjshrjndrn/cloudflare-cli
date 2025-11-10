# Quick Start Guide

Get started with cfcli in 5 minutes!

## 1. Installation

### macOS/Linux

```bash
# Download the latest release
curl -L https://github.com/skyline/cfcli/releases/latest/download/cfcli_Linux_x86_64.tar.gz | tar xz

# Make it executable and move to PATH
chmod +x cfcli
sudo mv cfcli /usr/local/bin/
```

### Windows (PowerShell)

```powershell
# Download from GitHub Releases
# https://github.com/skyline/cfcli/releases/latest
# Extract and add to PATH
```

### Using Homebrew

```bash
brew install skyline/tap/cfcli
```

## 2. Get Your Cloudflare API Token

1. Log in to the [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. Go to **My Profile** → **API Tokens**
3. Click **Create Token**
4. Use the **Edit zone DNS** template
5. Select your zone and create the token
6. Copy the token (you won't see it again!)

## 3. Test the Connection

```bash
# List your zones
cfcli -k YOUR_TOKEN zones
```

You should see a table with your Cloudflare zones!

## 4. Create a Configuration File (Optional but Recommended)

```bash
# Create the config file
cat > ~/.cfcli.yml << EOF
defaults:
    token: YOUR_TOKEN
    domain: example.com
EOF
```

Now you can run commands without specifying the token:

```bash
cfcli ls
```

## 5. Common Tasks

### List DNS Records

```bash
# List all records
cfcli -d example.com ls

# Filter by name
cfcli -d example.com -q name:www ls

# Filter by type
cfcli -d example.com -t A ls
```

### Add a DNS Record

```bash
# Add an A record
cfcli -d example.com -t A add test 192.0.2.1

# Add a CNAME record
cfcli -d example.com -t CNAME add blog myblog.com

# Add an MX record with priority
cfcli -d example.com -t MX -p 10 add @ mail.example.com
```

### Update a DNS Record

```bash
# Update an A record
cfcli -d example.com -t A edit test.example.com 192.0.2.2

# Change type from A to CNAME
cfcli -d example.com -t A -n CNAME edit test.example.com target.com
```

### Delete a DNS Record

```bash
# Delete a specific record
cfcli -d example.com -t A rm test.example.com

# Delete all records with a name
cfcli -d example.com rm test.example.com
```

## 6. Advanced Usage

### Use Different Output Formats

```bash
# JSON output
cfcli -d example.com -f json ls

# CSV output
cfcli -d example.com -f csv ls > records.csv
```

### Manage Multiple Accounts

```yaml
# ~/.cfcli.yml
defaults:
    account: work
accounts:
    work:
        token: WORK_TOKEN
        domain: work-site.com
    personal:
        token: PERSONAL_TOKEN
        domain: my-site.com
```

```bash
# Use work account (default)
cfcli ls

# Use personal account
cfcli -u personal ls
```

### Enable Cloudflare Proxy

```bash
# Add an A record with proxy enabled
cfcli -d example.com -t A -a add www 192.0.2.1

# The -a flag enables the orange cloud (proxy)
```

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check out [CONTRIBUTING.md](CONTRIBUTING.md) if you want to contribute
- Report issues on [GitHub](https://github.com/skyline/cfcli/issues)

## Troubleshooting

### "Zone not found"
Make sure you're using the correct domain name and that it exists in your Cloudflare account.

### "Authentication failed"
Check that your API token is correct and has the necessary permissions.

### "Permission denied" when installing
Use `sudo` when moving the binary to `/usr/local/bin/` or install to `~/bin/` instead.

## Tips

1. **Use config file**: Save time by storing your credentials in `~/.cfcli.yml`
2. **Tab completion**: Generate shell completions with `cfcli completion bash|zsh|fish`
3. **Short names**: Use `ls` instead of `listrecords`, `rm` instead of `remove`
4. **Filters**: Use `-q` for powerful filtering: `-q content:1.1.1.1,type:A`
5. **Help anytime**: Run `cfcli --help` or `cfcli COMMAND --help` for assistance

Happy DNS managing! 🚀

# Contributing to cfcli

Thank you for your interest in contributing to cfcli! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, but recommended)

### Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/cfcli.git
   cd cfcli
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Build the project:
   ```bash
   make build
   # or
   go build -o cfcli .
   ```

5. Run tests:
   ```bash
   make test
   # or
   go test -v ./...
   ```

## Project Structure

```
cfcli/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and global flags
│   ├── add.go             # Add DNS record command
│   ├── edit.go            # Edit DNS record command
│   ├── remove.go          # Remove DNS record command
│   ├── list.go            # List DNS records command
│   ├── find.go            # Find DNS records command
│   └── zones.go           # List zones command
├── internal/
│   ├── cloudflare/        # Cloudflare API client wrapper
│   │   └── client.go
│   └── config/            # Configuration file handling
│       └── config.go
├── main.go                # Application entry point
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── Makefile               # Build automation
├── Dockerfile             # Docker image definition
├── .goreleaser.yml        # GoReleaser configuration
└── .github/
    └── workflows/
        └── release.yml    # GitHub Actions release workflow
```

## Making Changes

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions and types

### Commit Messages

Use conventional commit format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(add): support for CAA records

Add support for creating and managing CAA records

Closes #123

fix(edit): handle nil pointer when updating proxied status

Properly check for nil proxied status before dereferencing
```

### Adding New Commands

1. Create a new file in `cmd/` directory (e.g., `cmd/newcommand.go`)
2. Define the cobra command
3. Register it in `init()` function with `rootCmd.AddCommand()`
4. Update README.md with usage examples

### Adding New DNS Record Types

1. Update `internal/cloudflare/client.go`:
   - Add support in `AddDNSRecord()`
   - Add support in `UpdateDNSRecord()`
2. Update documentation in README.md

## Testing

### Manual Testing

You can use the provided test token (read-only) or set up your own:

```bash
# List zones
./cfcli -k YOUR_TOKEN zones

# Test CRUD operations on a test domain
./cfcli -k YOUR_TOKEN -d test.example.com -t A add test 192.0.2.1
./cfcli -k YOUR_TOKEN -d test.example.com -t A edit test.example.com 192.0.2.2
./cfcli -k YOUR_TOKEN -d test.example.com -t A rm test.example.com
```

### Running Tests

```bash
make test
```

## Building and Releasing

### Local Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build with GoReleaser (snapshot)
make release-snapshot
```

### Creating a Release

Releases are automated via GitHub Actions and GoReleaser.

1. Update version in code if needed
2. Create and push a tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Create Docker images
   - Generate release notes
   - Upload artifacts to GitHub Releases

### Supported Platforms

The automated release builds for:
- Linux (amd64, arm64, arm, 386)
- macOS (amd64, arm64)
- Windows (amd64, 386)
- FreeBSD (amd64, arm64, arm, 386)
- OpenBSD (amd64, arm64, arm, 386)
- NetBSD (amd64, arm64, arm, 386)

## Pull Request Process

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feat/your-feature
   ```

2. Make your changes and commit:
   ```bash
   git add .
   git commit -m "feat: your feature description"
   ```

3. Push to your fork:
   ```bash
   git push origin feat/your-feature
   ```

4. Create a Pull Request on GitHub

5. Ensure CI checks pass

6. Address review comments if any

7. Once approved, a maintainer will merge your PR

## Code Review Guidelines

When reviewing PRs, consider:

- **Functionality**: Does it work as intended?
- **Code quality**: Is it readable and maintainable?
- **Tests**: Are there adequate tests?
- **Documentation**: Is it properly documented?
- **Performance**: Are there any performance concerns?
- **Security**: Are there any security implications?

## Questions?

If you have questions, feel free to:
- Open an issue
- Start a discussion on GitHub
- Reach out to maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

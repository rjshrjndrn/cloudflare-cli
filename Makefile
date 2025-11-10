.PHONY: build clean test install dev-install release-snapshot

BINARY_NAME=cfcli
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf dist/

test:
	@echo "Running tests..."
	go test -v ./...

install: build
	@echo "Installing to /usr/local/bin..."
	sudo mv $(BINARY_NAME) /usr/local/bin/

dev-install: build
	@echo "Installing to ~/bin..."
	mkdir -p ~/bin
	mv $(BINARY_NAME) ~/bin/

release-snapshot:
	@echo "Building release snapshot..."
	goreleaser release --snapshot --clean

# Cross-platform builds
build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 .

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 .

build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 .

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 .

build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe .

build-all: build-linux build-linux-arm64 build-darwin build-darwin-arm64 build-windows
	@echo "Built binaries for all platforms"

help:
	@echo "Available targets:"
	@echo "  build              - Build the binary for current platform"
	@echo "  clean              - Remove build artifacts"
	@echo "  test               - Run tests"
	@echo "  install            - Build and install to /usr/local/bin"
	@echo "  dev-install        - Build and install to ~/bin"
	@echo "  build-all          - Build for all platforms"
	@echo "  release-snapshot   - Build release snapshot with goreleaser"
	@echo "  help               - Show this help message"

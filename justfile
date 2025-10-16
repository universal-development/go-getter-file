# justfile for go-getter-file

# Default recipe - show available commands
default:
    @just --list

# Get version from git tags
version:
    @git describe --tags --abbrev=12 --dirty --broken 2>/dev/null || echo "dev"

# Build the application with version information
build:
    #!/usr/bin/env bash
    VERSION=$(git describe --tags --abbrev=12 --dirty --broken 2>/dev/null || echo "dev")
    echo "Building go-getter-file version: $VERSION"
    go build -ldflags "-X main.version=$VERSION" -o go-getter-file

# Build for multiple platforms
build-all:
    #!/usr/bin/env bash
    VERSION=$(git describe --tags --abbrev=12 --dirty --broken 2>/dev/null || echo "dev")
    echo "Building go-getter-file version: $VERSION for multiple platforms"

    # Linux
    GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o dist/go-getter-file-linux-amd64
    GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$VERSION" -o dist/go-getter-file-linux-arm64

    # macOS
    GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o dist/go-getter-file-darwin-amd64
    GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$VERSION" -o dist/go-getter-file-darwin-arm64

    # Windows
    GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o dist/go-getter-file-windows-amd64.exe

    echo "Build complete. Binaries in dist/"

# Run the application with example config
run:
    @just build
    ./go-getter-file example.go.getter.yaml

# Run unit tests
test:
    go test -v ./internal/...

# Run integration tests
test-integration:
    cd test && go test -v -timeout 5m

# Run all tests (unit + integration)
test-all:
    go test -v ./internal/...
    cd test && go test -v -timeout 5m

# Run tests with coverage
test-coverage:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
    rm -f go-getter-file
    rm -f go-getter-file.exe
    rm -rf dist/
    rm -f coverage.out coverage.html
    rm -rf downloaded-*.md github-*.md
    rm -f test/go-getter-file

# Install the application
install:
    #!/usr/bin/env bash
    VERSION=$(git describe --tags --abbrev=12 --dirty --broken 2>/dev/null || echo "dev")
    go install -ldflags "-X main.version=$VERSION"

# Format code
fmt:
    go fmt ./...

# Lint code
lint:
    @command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Install: https://golangci-lint.run/usage/install/"; exit 1; }
    golangci-lint run

# Tidy dependencies
tidy:
    go mod tidy

# Verify dependencies
verify:
    go mod verify

# Show current version
show-version:
    @just version

# Development workflow - format, test, build
dev: fmt test build

# Development workflow with all tests
dev-full: fmt test-all build

# Prepare for release - format, lint, test-all, build-all
release: fmt lint test-all build-all
    @echo "Release build complete!"

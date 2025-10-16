# go-getter-file

[![CI/CD](https://github.com/universal-development/go-getter-file/actions/workflows/ci.yml/badge.svg)](https://github.com/universal-development/go-getter-file/actions/workflows/ci.yml)

CLI application which do configuration files for go-getter v2.

Features:
* download files through configuration file
* download multiple files in parallel, by default files in format `*.go.getter.yaml`
* scan directories for configuration files
* usage of embedded go-getter library or external go-getter executable

Configuration files are in YAML format, see example below, `*.go.getter.yaml`
Example usage:

Process single configuration file:
```bash
go-getter-file file.go.getter.yaml
```
Process multiple configuration files:
```bash
go-getter-file file1.go.getter.yaml file2.go.getter.yaml
```
Process all configuration files in a directory:
```bash
go-getter-file configs-v1 configs-v2
```

Example configuration file

```yaml
# project1.go.getter.yaml
version: 1
name: "project1"

# Global configuration for all sources
config:
  # Optional: number of parallel fetches (default: 4)
  #parallelism: 4
  # Optional: number of retries for fetching each source (default: 3)
  #retries: 3
  # Optional: timeout for each fetch operation (default: 30s)
  #timeout: 30s
  # Optional: specify a custom path for go-getter operations, if not set use internal go-getter
  #go-getter-path: "/opt/go-getter"

sources:
  - url: "https://example.com/file1.txt"
    dest: "local-file1.txt"
    # Optional: override global timeout for this source
    timeout: 60s
  - url: "https://example.com/file2.txt"
    dest: "local-file2.txt"
  
  - url: "https://example.com/config/"
    dest: "local-config/"
    recursive: true
```

## Development

### Building

The project uses [just](https://github.com/casey/just) as a command runner for common tasks.

```bash
# Build the application
just build

# Build for all platforms
just build-all

# Run tests
just test

# Run tests with coverage
just test-coverage

# Clean build artifacts
just clean

# See all available commands
just --list
```

### Testing

```bash
# Run unit tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### CI/CD

The project uses GitHub Actions for continuous integration and deployment:

- **On Push/PR**: Runs tests and builds binaries for all platforms
- **On Tag (v*)**: Creates a GitHub release with:
  - Multi-platform binaries (Linux, macOS, Windows)
  - Compressed archives (.tar.gz, .zip)
  - SHA256 checksums
  - Auto-generated release notes
- **Manual Trigger**: Supports workflow_dispatch for custom releases

#### Creating Releases

**Standard Release:**
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

**Pre-Release (auto-detected):**
Tags containing `alpha`, `beta`, `rc`, `pre`, `preview`, or `dev` are automatically marked as pre-releases:
```bash
git tag -a v1.0.0-rc1 -m "Release candidate 1"
git push origin v1.0.0-rc1
```

**Manual Release (via GitHub Actions UI):**
1. Go to Actions → CI/CD → Run workflow
2. Enter the tag name (e.g., `v1.0.0`)
3. Optionally check "Create as draft release"
4. Optionally check "Mark as pre-release"
5. Click "Run workflow"

This allows you to:
- Create releases without pushing tags
- Override automatic pre-release detection
- Create draft releases for review before publishing

The version is automatically derived from git tags using `git describe --tags --abbrev=12 --dirty --broken`.

## TODO

- [ ] Add support for configurable configuration file patterns (e.g. `*.getter.yaml`, `*.config.yaml`) through CLI flag/env variable
- [ ] Add support for excluding certain files or directories through CLI flag/env variable
- [ ] Add support for dry-run mode to preview actions without making changes
- [ ] Add support for logging levels (info, debug, error) through CLI flag/env
- [ ] Add support for outputting results to a log file through CLI flag/env
- [ ] Add support for validating configuration files before processing
- [ ] Add support for more advanced go-getter options (e.g. authentication, proxies) through configuration file
- [ ] Add opentelemetry tracing and metrics
- [ ] Add unit and integration tests
- [ ] Add Dockerfile for containerized usage
- [ ] Add just file for building project
- [ ] Add GitHub Actions workflow for CI/CD
- [ ] Add custom headers to requests

## License

This code is released under the MIT License. See [LICENSE](LICENSE).
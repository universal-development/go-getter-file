# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) whenever releases are tagged.

## [Unreleased]
### Changed
- `just cleanup` now runs `go mod tidy` to ensure module metadata stays in sync.

## [0.0.2] - 2025-10-19
### Added
- `internal/app` package that exposes the CLI workflow for programmatic use.
- `just fmt-imports` task to run gofmt and goimports across the codebase.
- `just cleanup` aggregate task to execute all formatting helpers.
- Repository-managed `pre-commit` hook that runs the cleanup routine automatically when enabled.
- `just install-hooks` helper to wire up the repository-managed git hooks.
- Dependabot configuration covering GitHub Actions and grouped Go module updates.

### Changed
- Integration tests execute the CLI in-process, automatically skipping network-dependent cases when connectivity is unavailable.
- Processor now aggregates detailed error messages across configuration files and sources to simplify troubleshooting.

## [0.0.1] - 2025-10-16
### Added
- Command-line interface for processing go-getter configuration files with graceful interrupt handling.
- YAML configuration loader that applies sensible defaults, validates input, and supports per-source overrides.
- Fetcher implementation that can use the embedded go-getter library or an external binary with retry and timeout controls.
- Processor that expands file and directory arguments, enforces per-config parallelism, and streams progress output.
- Unit tests for configuration loading, fetcher behavior, and path expansion plus a comprehensive integration test suite with GitHub-backed fixtures.
- Project automation helpers including the `justfile` task runner definitions.

### Changed
- Updated go-getter dependency stack to the current v2 APIs for downloads.

### Documentation
- Authored README covering usage examples, development workflow, and release process guidance.

### CI/CD
- Introduced GitHub Actions workflows for testing, multi-platform builds, and release automation, along with subsequent pipeline refinements.

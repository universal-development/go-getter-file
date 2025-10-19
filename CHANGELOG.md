# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) whenever releases are tagged.

## [Unreleased]
### Added
- `internal/app` package that exposes the CLI workflow for programmatic use.

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

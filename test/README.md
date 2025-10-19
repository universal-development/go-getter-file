# Integration Tests

This directory contains integration tests for go-getter-file.

## Structure

```
test/
├── integration_test.go    # Integration test suite
├── fixtures/              # Test fixture files
│   ├── github-single-file.go.getter.yaml
│   ├── github-directory.go.getter.yaml
│   ├── multiple-files.go.getter.yaml
│   ├── config1.go.getter.yaml
│   ├── config2.go.getter.yaml
│   ├── invalid.go.getter.yaml
│   ├── simple.go.getter.yaml
│   ├── multiple.go.getter.yaml
│   └── batch-configs/
│       ├── batch1.go.getter.yaml
│       ├── batch2.go.getter.yaml
│       └── batch3.go.getter.yaml
└── README.md
```

## Running Tests

### Run all integration tests
```bash
just test-integration
```

### Run from test directory
```bash
cd test
go test -v
```

### Run specific test
```bash
cd test
go test -v -run TestFetchFromGitHub
```

### Run all tests (unit + integration)
```bash
just test-all
```

## Test Coverage

The integration tests cover:

1. **TestFetchFromGitHub** - Fetch a single file from GitHub using HTTPS URL
2. **TestFetchFromGitHubDirectory** - Fetch and extract a ZIP archive from GitHub
3. **TestProcessMultipleFiles** - Process a config with multiple sources in parallel
4. **TestProcessMultipleConfigFiles** - Process multiple config files at once
5. **TestProcessDirectory** - Scan and process all configs in a directory
6. **TestVersionFlag** - Test --version flag output
7. **TestHelpFlag** - Test --help flag output
8. **TestInvalidConfig** - Test error handling for invalid configs
9. **TestNonExistentFile** - Test error handling for missing files

## How Integration Tests Work

1. **TestMain** builds the binary once into a temporary directory before running any tests
2. Each test runs the binary with different configurations and working directories
3. Tests use fixture files from `fixtures/` directory instead of generating configs
4. The binary is cleaned up automatically after all tests complete
5. Each test uses `t.TempDir()` for isolated working directories

## Test Requirements

- Internet connection (tests fetch from GitHub)
- Go 1.21 or later
- Approximately 60 seconds timeout for network operations

## Test Fixtures

The `fixtures/` directory contains sample configuration files used by tests. These fixtures serve both as test inputs and as examples for users:

### Single File Tests
- `github-single-file.go.getter.yaml` - Fetch a single file via HTTPS
- `simple.go.getter.yaml` - Basic single-file fetch

### Multiple Files Tests
- `multiple-files.go.getter.yaml` - Multiple parallel downloads
- `multiple.go.getter.yaml` - Alternative multiple file config

### Multi-Config Tests
- `config1.go.getter.yaml` - First config for multi-config test
- `config2.go.getter.yaml` - Second config for multi-config test

### Directory Tests
- `batch-configs/batch*.go.getter.yaml` - Multiple configs in a directory

### Error Handling Tests
- `invalid.go.getter.yaml` - Invalid config for error testing
- `github-directory.go.getter.yaml` - ZIP archive extraction

## Writing New Tests

Integration tests should:
1. Use the `runCLI(workDir, args...)` helper to execute the binary
2. Use `t.TempDir()` for temporary directories
3. Reference fixture files from `fixturesPath` variable
4. Verify both success cases and error handling
5. Log useful information for debugging

Example:
```go
func TestNewFeature(t *testing.T) {
    testDir := t.TempDir()

    // Use fixture file
    fixtureFile := filepath.Join(fixturesPath, "my-fixture.go.getter.yaml")

    // Run CLI from test directory
    output, err := runCLI(testDir, fixtureFile)
    if err != nil {
        t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
    }

    // Verify results
    outputFile := filepath.Join(testDir, "output", "file.txt")
    if _, err := os.Stat(outputFile); os.IsNotExist(err) {
        t.Errorf("Expected output file %s does not exist", outputFile)
    }

    t.Logf("Test completed successfully")
}
```

## Troubleshooting

### Build Failures in IDE

If you see "Failed to build binary" when running tests in your IDE, make sure:
1. The test working directory is set to `test/`
2. Go modules are properly initialized
3. All dependencies are downloaded (`go mod download`)

The build command explicitly sets the working directory to the project root, so it should work from any location.

### Network Failures

Integration tests require internet access to fetch files from GitHub. If tests fail with network errors:
1. Check your internet connection
2. Verify GitHub is accessible
3. Check if you're behind a proxy (may need to configure Git/HTTP proxy settings)

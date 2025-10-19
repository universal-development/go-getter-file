package test

import (
	"bytes"
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/universal-development/go-getter-file/internal/app"
)

const (
	testCLIVersion      = "test"
	networkProbeTimeout = 5 * time.Second
)

var fixturesPath string

// TestMain sets up test environment
func TestMain(m *testing.M) {
	var err error

	// Get absolute path to fixtures directory
	fixturesPath, err = filepath.Abs("fixtures")
	if err != nil {
		panic("Failed to get fixtures path: " + err.Error())
	}

	// Run tests
	os.Exit(m.Run())
}

// runCLI executes the go-getter-file CLI directly via the app package
func runCLI(t *testing.T, workDir string, args ...string) (string, error) {
	t.Helper()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	if workDir != "" {
		if err := os.Chdir(workDir); err != nil {
			t.Fatalf("Failed to change to workDir %s: %v", workDir, err)
		}
		defer func() {
			if err := os.Chdir(origDir); err != nil {
				t.Fatalf("Failed to restore working directory: %v", err)
			}
		}()
	}

	var stdout bytes.Buffer

	err = app.Run(context.Background(), testCLIVersion, args, &stdout)
	output := stdout.String()
	return output, err
}

func ensureNetwork(t *testing.T) {
	t.Helper()

	conn, err := net.DialTimeout("tcp", "raw.githubusercontent.com:443", networkProbeTimeout)
	if err != nil {
		t.Skipf("skipping test: network unavailable (%v)", err)
		return
	}
	_ = conn.Close()
}

// TestFetchFromGitHub tests fetching a single file from GitHub
func TestFetchFromGitHub(t *testing.T) {
	ensureNetwork(t)
	testDir := t.TempDir()

	// Use fixture file
	fixtureFile := filepath.Join(fixturesPath, "github-single-file.go.getter.yaml")

	// Run the application from test directory (so relative paths in fixture work)
	output, err := runCLI(t, testDir, fixtureFile)
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	// Verify output directory exists and contains README.md
	// When dest is "output", go-getter places the file at output/README.md
	outputFile := filepath.Join(testDir, "output", "README.md")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// List what files actually exist for debugging
		if files, _ := filepath.Glob(filepath.Join(testDir, "output", "*")); len(files) > 0 {
			t.Logf("Files found in output/: %v", files)
		}
		t.Errorf("Expected output file %s does not exist", outputFile)
	}

	// Verify file has content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if len(content) == 0 {
		t.Error("Output file is empty")
	}

	t.Logf("Successfully fetched file (%d bytes)", len(content))
}

// TestFetchFromGitHubDirectory tests fetching a directory from GitHub
func TestFetchFromGitHubDirectory(t *testing.T) {
	ensureNetwork(t)
	testDir := t.TempDir()

	// Use fixture file
	fixtureFile := filepath.Join(fixturesPath, "github-directory.go.getter.yaml")

	// Run the application from test directory
	output, err := runCLI(t, testDir, fixtureFile)
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	// Verify output directory exists and contains files
	outputDir := filepath.Join(testDir, "output")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("Expected output directory %s does not exist", outputDir)
	}

	// Check that some files were downloaded
	files, err := filepath.Glob(filepath.Join(outputDir, "*"))
	if err != nil {
		t.Fatalf("Failed to glob for files: %v", err)
	}
	if len(files) == 0 {
		t.Error("No files found in downloaded directory")
	}

	t.Logf("Successfully fetched directory with %d files", len(files))
}

// TestProcessMultipleFiles tests processing a config with multiple sources
func TestProcessMultipleFiles(t *testing.T) {
	ensureNetwork(t)
	testDir := t.TempDir()

	// Use fixture file
	fixtureFile := filepath.Join(fixturesPath, "multiple-files.go.getter.yaml")

	// Run the application from test directory
	output, err := runCLI(t, testDir, fixtureFile)
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	// Verify both output directories exist and contain the expected files
	// Fixture uses dest: "output1" and dest: "output2"
	outputFiles := []string{
		filepath.Join(testDir, "output1", "README.md"),
		filepath.Join(testDir, "output2", "CHANGELOG.md"),
	}
	for _, file := range outputFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Expected output file %s does not exist", file)
		}

		// Verify file has content
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Failed to read output file %s: %v", file, err)
		}
		if len(content) == 0 {
			t.Errorf("Output file %s is empty", file)
		}
	}

	t.Logf("Successfully fetched %d files", len(outputFiles))
}

// TestProcessMultipleConfigFiles tests processing multiple config files at once
func TestProcessMultipleConfigFiles(t *testing.T) {
	ensureNetwork(t)
	testDir := t.TempDir()

	// Use fixtures directly
	config1 := filepath.Join(fixturesPath, "config1.go.getter.yaml")
	config2 := filepath.Join(fixturesPath, "config2.go.getter.yaml")

	// Run with both config files from test directory
	output, err := runCLI(t, testDir, config1, config2)
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	// Verify both outputs exist
	// Fixtures use dest: "out1" and dest: "out2"
	files := []string{
		filepath.Join(testDir, "out1", "README.md"),
		filepath.Join(testDir, "out2", "CHANGELOG.md"),
	}
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Expected output file %s does not exist", file)
		}

		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Failed to read output file %s: %v", file, err)
		}
		if len(content) == 0 {
			t.Errorf("Output file %s is empty", file)
		}
	}

	t.Logf("Successfully processed %d config files", 2)
}

// TestProcessDirectory tests processing all config files in a directory
func TestProcessDirectory(t *testing.T) {
	ensureNetwork(t)
	testDir := t.TempDir()

	// Use fixtures batch-configs directory
	configsDir := filepath.Join(fixturesPath, "batch-configs")

	// Run with directory from test directory
	output, err := runCLI(t, testDir, configsDir)
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	// Verify all output files exist
	// Fixtures use dest: "batch_out1", "batch_out2", "batch_out3"
	expectedFiles := []struct {
		dir  string
		file string
	}{
		{"batch_out1", "README.md"},
		{"batch_out2", "CHANGELOG.md"},
		{"batch_out3", "LICENSE"},
	}

	for _, expected := range expectedFiles {
		outputFile := filepath.Join(testDir, expected.dir, expected.file)
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Errorf("Expected output file %s does not exist", outputFile)
		}

		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output file %s: %v", outputFile, err)
		}
		if len(content) == 0 {
			t.Errorf("Output file %s is empty", outputFile)
		}
	}

	t.Logf("Successfully processed directory with %d config files", len(expectedFiles))
}

// TestVersionFlag tests the --version flag
func TestVersionFlag(t *testing.T) {
	output, err := runCLI(t, "", "--version")
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	outputStr := output
	if !strings.Contains(outputStr, "go-getter-file version") {
		t.Errorf("Expected version string, got: %s", outputStr)
	}

	t.Logf("Version output: %s", strings.TrimSpace(outputStr))
}

// TestHelpFlag tests the --help flag
func TestHelpFlag(t *testing.T) {
	output, err := runCLI(t, "", "--help")
	if err != nil {
		// Help might exit with non-zero, but should still produce output
		if len(output) == 0 {
			t.Fatalf("Command failed with no output: %v", err)
		}
	}

	outputStr := output
	expectedStrings := []string{"Usage:", "Options:", "Examples:"}
	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Expected help output to contain '%s', got: %s", expected, outputStr)
		}
	}

	t.Log("Help output verified")
}

// TestInvalidConfig tests handling of invalid configuration
func TestInvalidConfig(t *testing.T) {
	// Use invalid fixture
	configFile := filepath.Join(fixturesPath, "invalid.go.getter.yaml")

	// Run the application - should fail
	_, err := runCLI(t, "", configFile)
	if err == nil {
		t.Error("Expected command to fail with invalid config, but it succeeded")
	}

	if err != nil && !strings.Contains(err.Error(), "invalid config file") {
		t.Errorf("Expected error to mention invalid config, got: %v", err)
	}

	t.Log("Invalid config correctly rejected")
}

// TestNonExistentFile tests handling of non-existent config file
func TestNonExistentFile(t *testing.T) {
	output, err := runCLI(t, "", "/nonexistent/path/to/config.yaml")
	if err == nil {
		t.Error("Expected command to fail with non-existent file, but it succeeded")
	}

	if err != nil {
		if !strings.Contains(err.Error(), "failed to expand path") {
			t.Errorf("Unexpected error message: %v", err)
		}
	}

	if len(output) == 0 {
		t.Errorf("Expected output describing missing file, got empty output")
	}

	t.Log("Non-existent file correctly handled")
}

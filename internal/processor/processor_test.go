package processor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create test files
	testFiles := []string{
		"test1.go.getter.yaml",
		"test2.go.getter.yaml",
		"other.yaml",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tmpDir, file)
		if err := os.WriteFile(filePath, []byte("test: data"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	tests := []struct {
		name      string
		path      string
		wantCount int
		wantError bool
	}{
		{
			name:      "single file",
			path:      filepath.Join(tmpDir, "test1.go.getter.yaml"),
			wantCount: 1,
			wantError: false,
		},
		{
			name:      "directory with multiple matching files",
			path:      tmpDir,
			wantCount: 2, // test1 and test2, not other.yaml
			wantError: false,
		},
		{
			name:      "nonexistent file",
			path:      filepath.Join(tmpDir, "nonexistent.yaml"),
			wantCount: 0,
			wantError: true,
		},
		{
			name:      "nonexistent directory",
			path:      filepath.Join(tmpDir, "nonexistent-dir"),
			wantCount: 0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := expandPath(tt.path)

			if tt.wantError {
				if err == nil {
					t.Errorf("expandPath() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expandPath() unexpected error: %v", err)
				}
				if len(files) != tt.wantCount {
					t.Errorf("expandPath() returned %d files, want %d", len(files), tt.wantCount)
				}
			}
		})
	}
}

func TestExpandPathEmptyDirectory(t *testing.T) {
	// Create an empty directory
	tmpDir := t.TempDir()

	_, err := expandPath(tmpDir)
	if err == nil {
		t.Error("expandPath() expected error for empty directory, got nil")
	}
}

func TestNew(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create test files
	testFile1 := filepath.Join(tmpDir, "test1.go.getter.yaml")
	testFile2 := filepath.Join(tmpDir, "test2.go.getter.yaml")

	for _, file := range []string{testFile1, testFile2} {
		if err := os.WriteFile(file, []byte("test: data"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	tests := []struct {
		name      string
		paths     []string
		wantCount int
		wantError bool
	}{
		{
			name:      "single file",
			paths:     []string{testFile1},
			wantCount: 1,
			wantError: false,
		},
		{
			name:      "multiple files",
			paths:     []string{testFile1, testFile2},
			wantCount: 2,
			wantError: false,
		},
		{
			name:      "directory",
			paths:     []string{tmpDir},
			wantCount: 2,
			wantError: false,
		},
		{
			name:      "mix of files and directory",
			paths:     []string{testFile1, tmpDir},
			wantCount: 3, // testFile1 appears twice (once directly, once from dir scan)
			wantError: false,
		},
		{
			name:      "empty paths",
			paths:     []string{},
			wantCount: 0,
			wantError: true,
		},
		{
			name:      "nonexistent path",
			paths:     []string{filepath.Join(tmpDir, "nonexistent.yaml")},
			wantCount: 0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc, err := New(tt.paths)

			if tt.wantError {
				if err == nil {
					t.Errorf("New() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("New() unexpected error: %v", err)
				}
				if proc == nil {
					t.Fatal("New() returned nil processor")
				}
				if len(proc.configFiles) != tt.wantCount {
					t.Errorf("New() processor has %d config files, want %d", len(proc.configFiles), tt.wantCount)
				}
			}
		})
	}
}

func TestNewWithNestedDirectories(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create test files in both directories
	mainFile := filepath.Join(tmpDir, "main.go.getter.yaml")
	subFile := filepath.Join(subDir, "sub.go.getter.yaml")

	for _, file := range []string{mainFile, subFile} {
		if err := os.WriteFile(file, []byte("test: data"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Test that scanning main directory doesn't include subdirectory files
	proc, err := New([]string{tmpDir})
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// Should only find files in the top-level directory
	if len(proc.configFiles) != 1 {
		t.Errorf("New() found %d files in main directory, want 1 (subdirectories should not be scanned)", len(proc.configFiles))
	}

	// Test scanning both directories
	proc2, err := New([]string{tmpDir, subDir})
	if err != nil {
		t.Fatalf("New() with both directories unexpected error: %v", err)
	}

	if len(proc2.configFiles) != 2 {
		t.Errorf("New() with both directories found %d files, want 2", len(proc2.configFiles))
	}
}

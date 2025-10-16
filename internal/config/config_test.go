package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSetDefaults(t *testing.T) {
	tests := []struct {
		name     string
		config   FileConfig
		expected Config
	}{
		{
			name: "all defaults",
			config: FileConfig{
				Version: 1,
				Name:    "test",
				Config:  Config{},
			},
			expected: Config{
				Parallelism: 4,
				Retries:     3,
				Timeout:     30 * time.Second,
			},
		},
		{
			name: "partial defaults",
			config: FileConfig{
				Version: 1,
				Name:    "test",
				Config: Config{
					Parallelism: 10,
					Retries:     5,
				},
			},
			expected: Config{
				Parallelism: 10,
				Retries:     5,
				Timeout:     30 * time.Second,
			},
		},
		{
			name: "no defaults needed",
			config: FileConfig{
				Version: 1,
				Name:    "test",
				Config: Config{
					Parallelism:  8,
					Retries:      2,
					Timeout:      60 * time.Second,
					GoGetterPath: "/custom/path",
				},
			},
			expected: Config{
				Parallelism:  8,
				Retries:      2,
				Timeout:      60 * time.Second,
				GoGetterPath: "/custom/path",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.SetDefaults()
			if tt.config.Config.Parallelism != tt.expected.Parallelism {
				t.Errorf("Parallelism = %d, want %d", tt.config.Config.Parallelism, tt.expected.Parallelism)
			}
			if tt.config.Config.Retries != tt.expected.Retries {
				t.Errorf("Retries = %d, want %d", tt.config.Config.Retries, tt.expected.Retries)
			}
			if tt.config.Config.Timeout != tt.expected.Timeout {
				t.Errorf("Timeout = %v, want %v", tt.config.Config.Timeout, tt.expected.Timeout)
			}
			if tt.config.Config.GoGetterPath != tt.expected.GoGetterPath {
				t.Errorf("GoGetterPath = %s, want %s", tt.config.Config.GoGetterPath, tt.expected.GoGetterPath)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    FileConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid config",
			config: FileConfig{
				Version: 1,
				Name:    "test-project",
				Sources: []Source{
					{URL: "https://example.com/file.txt", Dest: "local.txt"},
				},
			},
			wantError: false,
		},
		{
			name: "missing version",
			config: FileConfig{
				Name: "test-project",
				Sources: []Source{
					{URL: "https://example.com/file.txt", Dest: "local.txt"},
				},
			},
			wantError: true,
			errorMsg:  "version is required",
		},
		{
			name: "missing name",
			config: FileConfig{
				Version: 1,
				Sources: []Source{
					{URL: "https://example.com/file.txt", Dest: "local.txt"},
				},
			},
			wantError: true,
			errorMsg:  "name is required",
		},
		{
			name: "no sources",
			config: FileConfig{
				Version: 1,
				Name:    "test-project",
				Sources: []Source{},
			},
			wantError: true,
			errorMsg:  "at least one source is required",
		},
		{
			name: "source missing url",
			config: FileConfig{
				Version: 1,
				Name:    "test-project",
				Sources: []Source{
					{Dest: "local.txt"},
				},
			},
			wantError: true,
			errorMsg:  "source 0: url is required",
		},
		{
			name: "source missing dest",
			config: FileConfig{
				Version: 1,
				Name:    "test-project",
				Sources: []Source{
					{URL: "https://example.com/file.txt"},
				},
			},
			wantError: true,
			errorMsg:  "source 0: dest is required",
		},
		{
			name: "multiple sources with one invalid",
			config: FileConfig{
				Version: 1,
				Name:    "test-project",
				Sources: []Source{
					{URL: "https://example.com/file1.txt", Dest: "local1.txt"},
					{URL: "https://example.com/file2.txt"}, // missing dest
				},
			},
			wantError: true,
			errorMsg:  "source 1: dest is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantError {
				if err == nil {
					t.Errorf("Validate() expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		content   string
		wantError bool
		validate  func(*testing.T, *FileConfig)
	}{
		{
			name: "valid minimal config",
			content: `version: 1
name: "test-project"
sources:
  - url: "https://example.com/file.txt"
    dest: "local.txt"
`,
			wantError: false,
			validate: func(t *testing.T, cfg *FileConfig) {
				if cfg.Version != 1 {
					t.Errorf("Version = %d, want 1", cfg.Version)
				}
				if cfg.Name != "test-project" {
					t.Errorf("Name = %s, want test-project", cfg.Name)
				}
				if len(cfg.Sources) != 1 {
					t.Errorf("Sources count = %d, want 1", len(cfg.Sources))
				}
				// Check defaults were set
				if cfg.Config.Parallelism != 4 {
					t.Errorf("Parallelism = %d, want 4", cfg.Config.Parallelism)
				}
				if cfg.Config.Retries != 3 {
					t.Errorf("Retries = %d, want 3", cfg.Config.Retries)
				}
				if cfg.Config.Timeout != 30*time.Second {
					t.Errorf("Timeout = %v, want 30s", cfg.Config.Timeout)
				}
			},
		},
		{
			name: "valid config with custom settings",
			content: `version: 1
name: "custom-project"
config:
  parallelism: 8
  retries: 5
  timeout: 60s
  go-getter-path: "/usr/bin/go-getter"
sources:
  - url: "https://example.com/file1.txt"
    dest: "local1.txt"
  - url: "https://example.com/file2.txt"
    dest: "local2.txt"
    timeout: 120s
`,
			wantError: false,
			validate: func(t *testing.T, cfg *FileConfig) {
				if cfg.Config.Parallelism != 8 {
					t.Errorf("Parallelism = %d, want 8", cfg.Config.Parallelism)
				}
				if cfg.Config.Retries != 5 {
					t.Errorf("Retries = %d, want 5", cfg.Config.Retries)
				}
				if cfg.Config.Timeout != 60*time.Second {
					t.Errorf("Timeout = %v, want 60s", cfg.Config.Timeout)
				}
				if cfg.Config.GoGetterPath != "/usr/bin/go-getter" {
					t.Errorf("GoGetterPath = %s, want /usr/bin/go-getter", cfg.Config.GoGetterPath)
				}
				if len(cfg.Sources) != 2 {
					t.Errorf("Sources count = %d, want 2", len(cfg.Sources))
				}
				if cfg.Sources[1].Timeout != 120*time.Second {
					t.Errorf("Source timeout = %v, want 120s", cfg.Sources[1].Timeout)
				}
			},
		},
		{
			name: "invalid yaml",
			content: `version: 1
name: "test"
sources:
  - url: "test
    invalid yaml here
`,
			wantError: true,
		},
		{
			name: "invalid config - missing sources",
			content: `version: 1
name: "test-project"
sources: []
`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			filePath := filepath.Join(tmpDir, tt.name+".yaml")
			err := os.WriteFile(filePath, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Load config
			cfg, err := LoadConfig(filePath)

			if tt.wantError {
				if err == nil {
					t.Errorf("LoadConfig() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("LoadConfig() unexpected error: %v", err)
				}
				if tt.validate != nil {
					tt.validate(t, cfg)
				}
			}
		})
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/to/config.yaml")
	if err == nil {
		t.Error("LoadConfig() expected error for nonexistent file, got nil")
	}
}

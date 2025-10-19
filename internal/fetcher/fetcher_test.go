package fetcher

import (
	"testing"
	"time"

	"github.com/universal-development/go-getter-file/internal/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name            string
		config          config.Config
		wantExternalBin bool
	}{
		{
			name: "embedded go-getter (no path specified)",
			config: config.Config{
				Parallelism: 4,
				Retries:     3,
				Timeout:     30 * time.Second,
			},
			wantExternalBin: false,
		},
		{
			name: "external go-getter binary",
			config: config.Config{
				Parallelism:  4,
				Retries:      3,
				Timeout:      30 * time.Second,
				GoGetterPath: "/usr/bin/go-getter",
			},
			wantExternalBin: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(tt.config)

			if f == nil {
				t.Fatal("New() returned nil fetcher")
			}

			if f.config.Parallelism != tt.config.Parallelism {
				t.Errorf("Parallelism = %d, want %d", f.config.Parallelism, tt.config.Parallelism)
			}

			if f.config.Retries != tt.config.Retries {
				t.Errorf("Retries = %d, want %d", f.config.Retries, tt.config.Retries)
			}

			if f.config.Timeout != tt.config.Timeout {
				t.Errorf("Timeout = %v, want %v", f.config.Timeout, tt.config.Timeout)
			}

			if f.useExternalBin != tt.wantExternalBin {
				t.Errorf("useExternalBin = %v, want %v", f.useExternalBin, tt.wantExternalBin)
			}

			if f.config.GoGetterPath != tt.config.GoGetterPath {
				t.Errorf("GoGetterPath = %s, want %s", f.config.GoGetterPath, tt.config.GoGetterPath)
			}
		})
	}
}

func TestFetcherConfiguration(t *testing.T) {
	cfg := config.Config{
		Parallelism:  10,
		Retries:      5,
		Timeout:      60 * time.Second,
		GoGetterPath: "/custom/go-getter",
	}

	f := New(cfg)

	if f.config.Parallelism != 10 {
		t.Errorf("Expected parallelism 10, got %d", f.config.Parallelism)
	}

	if f.config.Retries != 5 {
		t.Errorf("Expected retries 5, got %d", f.config.Retries)
	}

	if f.config.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", f.config.Timeout)
	}

	if !f.useExternalBin {
		t.Error("Expected useExternalBin to be true")
	}
}

func TestFetcherDefaults(t *testing.T) {
	cfg := config.Config{}
	f := New(cfg)

	if f.useExternalBin {
		t.Error("Expected useExternalBin to be false when no path is set")
	}

	if f.config.GoGetterPath != "" {
		t.Errorf("Expected empty GoGetterPath, got %s", f.config.GoGetterPath)
	}
}

// Note: Testing actual fetch operations would require mocking or integration tests
// The following are placeholder tests that could be expanded with proper mocking

func TestFetchSourceTimeout(t *testing.T) {
	// This is a basic structure test
	// In a real scenario, you'd want to:
	// 1. Mock the go-getter client
	// 2. Test timeout behavior
	// 3. Test retry logic

	cfg := config.Config{
		Retries: 3,
		Timeout: 1 * time.Second,
	}

	source := config.Source{
		URL:     "https://example.com/file.txt",
		Dest:    "/tmp/test.txt",
		Timeout: 2 * time.Second, // Override global timeout
	}

	f := New(cfg)

	// Verify that the fetcher was created with correct config
	if f.config.Timeout != 1*time.Second {
		t.Errorf("Global timeout = %v, want 1s", f.config.Timeout)
	}

	if source.Timeout != 2*time.Second {
		t.Errorf("Source timeout = %v, want 2s", source.Timeout)
	}

	// The actual fetch would fail without network, but we've verified the structure
	// Integration tests would cover the actual fetching behavior
}

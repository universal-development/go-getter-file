package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the global configuration for all sources
type Config struct {
	Parallelism  int           `yaml:"parallelism,omitempty"`
	Retries      int           `yaml:"retries,omitempty"`
	Timeout      time.Duration `yaml:"timeout,omitempty"`
	GoGetterPath string        `yaml:"go-getter-path,omitempty"`
}

// Source represents a single source to fetch
type Source struct {
	URL       string        `yaml:"url"`
	Dest      string        `yaml:"dest"`
	Timeout   time.Duration `yaml:"timeout,omitempty"`
	Recursive bool          `yaml:"recursive,omitempty"`
}

// FileConfig represents the complete configuration file structure
type FileConfig struct {
	Version int      `yaml:"version"`
	Name    string   `yaml:"name"`
	Config  Config   `yaml:"config"`
	Sources []Source `yaml:"sources"`
}

// SetDefaults sets default values for the configuration
func (c *FileConfig) SetDefaults() {
	if c.Config.Parallelism == 0 {
		c.Config.Parallelism = 4
	}
	if c.Config.Retries == 0 {
		c.Config.Retries = 3
	}
	if c.Config.Timeout == 0 {
		c.Config.Timeout = 30 * time.Second
	}
}

// LoadConfig loads a configuration file from the given path
func LoadConfig(path string) (*FileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config FileConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	config.SetDefaults()

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config file %s: %w", path, err)
	}

	return &config, nil
}

// Validate validates the configuration
func (c *FileConfig) Validate() error {
	if c.Version == 0 {
		return fmt.Errorf("version is required")
	}
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source is required")
	}

	for i, source := range c.Sources {
		if source.URL == "" {
			return fmt.Errorf("source %d: url is required", i)
		}
		if source.Dest == "" {
			return fmt.Errorf("source %d: dest is required", i)
		}
	}

	return nil
}

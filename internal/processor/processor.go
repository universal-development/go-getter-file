package processor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/universal-development/go-getter-file/internal/config"
	"github.com/universal-development/go-getter-file/internal/fetcher"
)

// Processor handles processing configuration files
type Processor struct {
	configFiles []string
}

// New creates a new Processor
func New(paths []string) (*Processor, error) {
	var configFiles []string

	for _, path := range paths {
		files, err := expandPath(path)
		if err != nil {
			return nil, fmt.Errorf("failed to expand path %s: %w", path, err)
		}
		configFiles = append(configFiles, files...)
	}

	if len(configFiles) == 0 {
		return nil, fmt.Errorf("no configuration files found")
	}

	return &Processor{
		configFiles: configFiles,
	}, nil
}

// expandPath expands a path to a list of configuration files
func expandPath(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// If it's a file, return it directly
	if !info.IsDir() {
		return []string{path}, nil
	}

	// If it's a directory, scan for *.go.getter.yaml files
	pattern := filepath.Join(path, "*.go.getter.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no *.go.getter.yaml files found in directory %s", path)
	}

	return matches, nil
}

// Process processes all configuration files
func (p *Processor) Process(ctx context.Context) error {
	fmt.Printf("Processing %d configuration file(s)\n", len(p.configFiles))

	var wg sync.WaitGroup
	errors := make([]error, len(p.configFiles))

	for i, configFile := range p.configFiles {
		wg.Add(1)
		go func(idx int, file string) {
			defer wg.Done()
			errors[idx] = p.processConfigFile(ctx, file)
		}(i, configFile)
	}

	wg.Wait()

	// Check for errors
	var hasError bool
	var details []string
	for i, err := range errors {
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", p.configFiles[i], err)
			hasError = true
			details = append(details, fmt.Sprintf("%s: %v", p.configFiles[i], err))
		}
	}

	if hasError {
		return fmt.Errorf("some configuration files failed to process: %s", strings.Join(details, "; "))
	}

	return nil
}

// processConfigFile processes a single configuration file
func (p *Processor) processConfigFile(ctx context.Context, path string) error {
	fmt.Printf("\n==> Processing config: %s\n", path)

	cfg, err := config.LoadConfig(path)
	if err != nil {
		return err
	}

	fmt.Printf("Config: %s (version: %d)\n", cfg.Name, cfg.Version)
	fmt.Printf("Sources: %d, Parallelism: %d, Retries: %d\n",
		len(cfg.Sources), cfg.Config.Parallelism, cfg.Config.Retries)

	f := fetcher.New(cfg.Config)

	// Process sources with parallelism
	return p.processSources(ctx, f, cfg.Sources, cfg.Config.Parallelism)
}

// processSources processes all sources with the specified parallelism
func (p *Processor) processSources(ctx context.Context, f *fetcher.Fetcher, sources []config.Source, parallelism int) error {
	semaphore := make(chan struct{}, parallelism)
	var wg sync.WaitGroup
	errors := make([]error, len(sources))

	for i, source := range sources {
		wg.Add(1)
		go func(idx int, src config.Source) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fmt.Printf("  [%d/%d] Fetching %s -> %s\n", idx+1, len(sources), src.URL, src.Dest)
			err := f.FetchSource(ctx, src)
			if err != nil {
				errors[idx] = err
				fmt.Printf("  [%d/%d] Failed: %v\n", idx+1, len(sources), err)
			} else {
				fmt.Printf("  [%d/%d] Success: %s\n", idx+1, len(sources), src.Dest)
			}
		}(i, source)
	}

	wg.Wait()

	// Check for errors
	var hasError bool
	var details []string
	for idx, err := range errors {
		if err != nil {
			hasError = true
			details = append(details, fmt.Sprintf("%s: %v", sources[idx].URL, err))
		}
	}

	if hasError {
		return fmt.Errorf("some sources failed to download: %s", strings.Join(details, "; "))
	}

	return nil
}

package fetcher

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/hashicorp/go-getter/v2"
	"github.com/universal-development/go-getter-file/internal/config"
)

// Fetcher handles downloading files using go-getter
type Fetcher struct {
	config         config.Config
	useExternalBin bool
}

// New creates a new Fetcher instance
func New(cfg config.Config) *Fetcher {
	return &Fetcher{
		config:         cfg,
		useExternalBin: cfg.GoGetterPath != "",
	}
}

// FetchSource downloads a single source with retries
func (f *Fetcher) FetchSource(ctx context.Context, source config.Source) error {
	timeout := source.Timeout
	if timeout == 0 {
		timeout = f.config.Timeout
	}

	retries := f.config.Retries
	var lastErr error

	for attempt := 0; attempt <= retries; attempt++ {
		if attempt > 0 {
			fmt.Printf("  Retry %d/%d for %s\n", attempt, retries, source.URL)
		}

		err := f.fetch(ctx, source, timeout)
		if err == nil {
			return nil
		}

		lastErr = err
		if attempt < retries {
			// Wait a bit before retrying
			time.Sleep(time.Second * time.Duration(attempt+1))
		}
	}

	return fmt.Errorf("failed after %d retries: %w", retries, lastErr)
}

// fetch performs the actual download
func (f *Fetcher) fetch(ctx context.Context, source config.Source, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if f.useExternalBin {
		return f.fetchExternal(ctx, source)
	}

	return f.fetchEmbedded(ctx, source)
}

// fetchEmbedded uses the embedded go-getter library
func (f *Fetcher) fetchEmbedded(ctx context.Context, source config.Source) error {
	client := &getter.Client{}

	req := &getter.Request{
		Src: source.URL,
		Dst: source.Dest,
	}

	_, err := client.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("go-getter failed for %s: %w", source.URL, err)
	}

	return nil
}

// fetchExternal uses an external go-getter binary
func (f *Fetcher) fetchExternal(ctx context.Context, source config.Source) error {
	args := []string{source.URL, source.Dest}

	cmd := exec.CommandContext(ctx, f.config.GoGetterPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("external go-getter failed for %s: %w\nOutput: %s",
			source.URL, err, string(output))
	}

	return nil
}

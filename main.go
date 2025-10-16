package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/universal-development/go-getter-file/internal/processor"
)

// version is set via ldflags during build
var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		printUsage()
		return fmt.Errorf("no configuration files or directories specified")
	}

	// Check for help or version flags
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage()
		return nil
	}

	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		fmt.Printf("go-getter-file version %s\n", version)
		return nil
	}

	// Print version on startup
	fmt.Printf("go-getter-file version %s\n", version)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, cancelling...")
		cancel()
	}()

	// Get paths from arguments
	paths := os.Args[1:]

	// Create processor
	proc, err := processor.New(paths)
	if err != nil {
		return err
	}

	// Process all configuration files
	if err := proc.Process(ctx); err != nil {
		return err
	}

	fmt.Println("\nAll configuration files processed successfully!")
	return nil
}

func printUsage() {
	fmt.Printf(`go-getter-file version %s

CLI application for fetching files using go-getter with YAML configuration.

Usage:
  go-getter-file [options] <config-file-or-directory>...

Options:
  -h, --help     Show this help message
  -v, --version  Show version information

Arguments:
  One or more configuration files (*.go.getter.yaml) or directories.
  Directories will be scanned for *.go.getter.yaml files.

Examples:
  # Process a single configuration file
  go-getter-file project1.go.getter.yaml

  # Process multiple configuration files
  go-getter-file project1.go.getter.yaml project2.go.getter.yaml

  # Process all configuration files in directories
  go-getter-file configs-v1/ configs-v2/

  # Mix files and directories
  go-getter-file project1.go.getter.yaml configs/

Configuration:
  Configuration files are in YAML format (*.go.getter.yaml)
  See README.md for configuration file format and examples.

`, version)
}

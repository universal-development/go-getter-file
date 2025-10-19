package app

import (
	"context"
	"fmt"
	"io"

	"github.com/universal-development/go-getter-file/internal/processor"
)

// Run executes the CLI application logic using the provided context and arguments.
func Run(ctx context.Context, version string, args []string, stdout io.Writer) error {
	if len(args) == 0 {
		printUsage(stdout, version)
		return fmt.Errorf("no configuration files or directories specified")
	}

	if args[0] == "-h" || args[0] == "--help" {
		printUsage(stdout, version)
		return nil
	}

	if args[0] == "-v" || args[0] == "--version" {
		fmt.Fprintf(stdout, "go-getter-file version %s\n", version)
		return nil
	}

	fmt.Fprintf(stdout, "go-getter-file version %s\n", version)

	proc, err := processor.New(args)
	if err != nil {
		return err
	}

	if err := proc.Process(ctx); err != nil {
		return err
	}

	fmt.Fprintln(stdout, "\nAll configuration files processed successfully!")
	return nil
}

func printUsage(w io.Writer, version string) {
	fmt.Fprintf(w, `go-getter-file version %s

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

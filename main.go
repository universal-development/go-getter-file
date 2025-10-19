package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/universal-development/go-getter-file/internal/app"
)

// version is set via ldflags during build
var version = "dev"

func main() {
	if err := execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Fprintln(os.Stdout, "\nReceived interrupt signal, cancelling...")
		cancel()
	}()

	return app.Run(ctx, version, os.Args[1:], os.Stdout)
}

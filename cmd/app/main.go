package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Glitch-guy0/authService/modules/logger"
)

// Version will be set during build
type VersionInfo struct {
	Version string
	Commit  string
	Date    string
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Create a context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize version info
	versionInfo := VersionInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}

	// Start the application
	if err := run(ctx, versionInfo); err != nil {
		os.Stderr.WriteString("error: " + err.Error() + "\n")
		os.Exit(1)
	}
}

func run(ctx context.Context, version VersionInfo) error {
	// Initialize logger
	log := logger.New()

	// Log startup information
	log.Info("Starting auth-service version %s (commit: %s, built: %s)",
		version.Version, version.Commit, version.Date)

	// TODO: Initialize configuration

	// TODO: Initialize application context

	// TODO: Initialize HTTP server

	// Wait for interrupt signal to gracefully shut down the server
	<-ctx.Done()
	log.Info("Shutting down...")

	// TODO: Add graceful shutdown logic

	return nil
}

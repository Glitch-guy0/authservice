package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Glitch-guy0/authService/src/modules/core"
	"github.com/Glitch-guy0/authService/src/modules/logger"
	"github.com/Glitch-guy0/authService/src/modules/server"
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

	// Initialize application context
	appCtx := core.NewAppContextWithDefaults(log)

	// Initialize HTTP server
	server := server.NewServerWithDefaults(appCtx)
	server.Initialize()

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Info("Server starting on %s", server.GetAddress())
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Log server started successfully
	log.Info("HTTP server started successfully",
		"address", server.GetAddress(),
		"port", server.GetConfig().Port,
		"mode", server.GetConfig().Mode,
	)

	// Wait for interrupt signal to gracefully shut down the server
	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Info("Shutting down...")

		// Graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("Server shutdown failed", "error", err)
			return err
		}

		log.Info("Server shutdown successfully")
	}

	return nil
}

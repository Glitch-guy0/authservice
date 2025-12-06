package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Glitch-guy0/authService/modules/config"
	"github.com/Glitch-guy0/authService/modules/core"
	"github.com/Glitch-guy0/authService/modules/logger"
	"github.com/Glitch-guy0/authService/modules/server"
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
		stop() // Call stop before exiting
		os.Exit(1)
	}
}

func run(ctx context.Context, version VersionInfo) error {
	// Initialize configuration first
	if err := config.Init("./configs"); err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}

	// Initialize logger
	log := logger.New()

	// Log startup information
	log.Info("Starting auth-service version %s (commit: %s, built: %s)",
		version.Version, version.Commit, version.Date)

	// Initialize application context with loaded configuration
	appCtx := core.NewAppContext(log, config.AllSettings())

	// Initialize HTTP server
	server := server.NewServerFromConfig(appCtx)
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
	log.WithFields(map[string]interface{}{
		"address": server.GetAddress(),
		"port":    server.GetConfig().Port,
		"mode":    server.GetConfig().Mode,
	}).Info("HTTP server started successfully")

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
			log.WithField("error", err).Error("Server shutdown failed")
			return err
		}

		log.Info("Server shutdown successfully")
	}

	return nil
}

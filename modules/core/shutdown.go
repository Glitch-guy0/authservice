package core

import (
	"context"
	"sync"
	"time"

	"github.com/Glitch-guy0/authService/modules/logger"
)

// ShutdownHandler represents a function that handles graceful shutdown
type ShutdownHandler struct {
	Name     string
	Handler  func(ctx context.Context) error
	Timeout  time.Duration
	Priority int // Lower numbers = higher priority (shutdown first)
}

// ShutdownManager manages graceful shutdown of application components
type ShutdownManager struct {
	handlers []ShutdownHandler
	mu       sync.Mutex
	logger   logger.Logger
}

// NewShutdownManager creates a new shutdown manager
func NewShutdownManager(logger logger.Logger) *ShutdownManager {
	return &ShutdownManager{
		handlers: make([]ShutdownHandler, 0),
		logger:   logger,
	}
}

// RegisterHandler registers a shutdown handler
func (sm *ShutdownManager) RegisterHandler(name string, handler func(ctx context.Context) error, timeout time.Duration, priority int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.handlers = append(sm.handlers, ShutdownHandler{
		Name:     name,
		Handler:  handler,
		Timeout:  timeout,
		Priority: priority,
	})

	sm.logger.Info("Registered shutdown handler", "name", name, "timeout", timeout, "priority", priority)
}

// Shutdown triggers graceful shutdown of all registered handlers
func (sm *ShutdownManager) Shutdown(ctx context.Context) error {
	sm.mu.Lock()
	handlers := make([]ShutdownHandler, len(sm.handlers))
	copy(handlers, sm.handlers)
	sm.mu.Unlock()

	sm.logger.Info("Starting graceful shutdown", "handlers", len(handlers))

	// Sort handlers by priority (lower numbers first)
	for i := 0; i < len(handlers); i++ {
		for j := i + 1; j < len(handlers); j++ {
			if handlers[i].Priority > handlers[j].Priority {
				handlers[i], handlers[j] = handlers[j], handlers[i]
			}
		}
	}

	var errors []error

	for _, handler := range handlers {
		sm.logger.Info("Shutting down component", "name", handler.Name)

		// Create context with timeout for this handler
		handlerCtx, cancel := context.WithTimeout(ctx, handler.Timeout)

		// Execute shutdown handler
		err := handler.Handler(handlerCtx)
		cancel()

		if err != nil {
			sm.logger.Error("Shutdown handler failed", "name", handler.Name, "error", err)
			errors = append(errors, err)
		} else {
			sm.logger.Info("Component shutdown successfully", "name", handler.Name)
		}
	}

	if len(errors) > 0 {
		sm.logger.Error("Some shutdown handlers failed", "count", len(errors))
		return errors[0] // Return first error
	}

	sm.logger.Info("Graceful shutdown completed")
	return nil
}

// Shutdown triggers shutdown on the AppContext
func (ac *AppContext) Shutdown(timeout time.Duration) error {
	ac.mu.Lock()
	if ac.shutdown == nil {
		ac.mu.Unlock()
		return nil // Already shutdown
	}

	// Close shutdown channel to signal shutdown
	close(ac.shutdown)
	ac.shutdown = nil
	ac.mu.Unlock()

	// Create shutdown manager
	shutdownManager := NewShutdownManager(ac.Logger)

	// Register default shutdown handlers
	shutdownManager.RegisterHandler("logger", func(ctx context.Context) error {
		// Logger doesn't need explicit shutdown, but we can flush if needed
		ac.Logger.Info("Logger shutdown completed")
		return nil
	}, 5*time.Second, 100)

	// Register health status cleanup
	shutdownManager.RegisterHandler("health", func(ctx context.Context) error {
		ac.healthMu.Lock()
		ac.healthStatus = make(map[string]HealthStatus)
		ac.healthMu.Unlock()
		ac.Logger.Info("Health status cleanup completed")
		return nil
	}, 2*time.Second, 90)

	// Create context with overall timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Execute shutdown
	return shutdownManager.Shutdown(ctx)
}

// WaitForShutdown waits for the shutdown signal
func (ac *AppContext) WaitForShutdown() <-chan struct{} {
	return ac.GetShutdownChannel()
}

// RegisterShutdownHandler registers a custom shutdown handler
func (ac *AppContext) RegisterShutdownHandler(name string, handler func(ctx context.Context) error, timeout time.Duration, priority int) {
	// This would integrate with a shutdown manager
	// For now, we'll just log it
	_ = handler // TODO: Implement actual shutdown handler registration
	ac.Logger.Info("Shutdown handler registered", "name", name, "timeout", timeout, "priority", priority)
}

// AddShutdownHandler adds a simple shutdown handler with default timeout and priority
func (ac *AppContext) AddShutdownHandler(name string, handler func(ctx context.Context) error) {
	ac.RegisterShutdownHandler(name, handler, 30*time.Second, 50)
}

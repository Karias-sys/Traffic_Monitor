package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Karias-sys/Traffic_Monitor/internal/config"
	"github.com/Karias-sys/Traffic_Monitor/pkg/logger"
)

// Version information (set by build-time ldflags)
var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Create a context that's cancelled on interrupt
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Run the application
	if err := run(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run(ctx context.Context) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	logger, err := logger.New(logger.Config{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Log startup information
	logger.WithComponent("main").Info("Starting Netwatch Traffic Monitor")

	// Log version information
	logger.WithComponent("main").Info("Version: " + version)
	logger.WithComponent("main").Info("Build Time: " + buildTime)
	logger.WithComponent("main").Info("Git Commit: " + gitCommit)

	// Log configuration (excluding sensitive data)
	logger.WithComponent("main").Info(fmt.Sprintf("Host: %s", cfg.Host))
	logger.WithComponent("main").Info(fmt.Sprintf("Port: %d", cfg.Port))
	logger.WithComponent("main").Info(fmt.Sprintf("Interface: %s", cfg.Interface))
	logger.WithComponent("main").Info(fmt.Sprintf("Log Level: %s", cfg.LogLevel))
	logger.WithComponent("main").Info(fmt.Sprintf("Development Mode: %t", cfg.DevMode))

	// TODO: Initialize and start application components
	// This will be implemented in subsequent stories:
	// - Packet capture engine
	// - Flow processing
	// - Web server and API
	// - WebSocket handler
	// - Metrics endpoint

	logger.WithComponent("main").Info("Application initialized successfully")

	// Wait for context cancellation (interrupt signal)
	<-ctx.Done()

	logger.WithComponent("main").Info("Shutting down application")

	// TODO: Graceful shutdown of components
	// This will be implemented in subsequent stories when components are added

	logger.WithComponent("main").Info("Application shutdown complete")

	return nil
}
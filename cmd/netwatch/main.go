package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Karias-sys/Traffic_Monitor/internal/capture"
	"github.com/Karias-sys/Traffic_Monitor/internal/config"
	"github.com/Karias-sys/Traffic_Monitor/internal/metrics"
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

	// Initialize interface manager
	interfaceManager := capture.NewInterfaceManager(logger.WithComponent("interface").Logger)

	// Validate and resolve interface
	var interfaceName string
	if cfg.Interface == "any" {
		defaultIface, err := interfaceManager.GetDefaultInterface()
		if err != nil {
			return fmt.Errorf("failed to get default interface: %w", err)
		}
		interfaceName = defaultIface.Name
		logger.WithComponent("interface").Info(fmt.Sprintf("Using default interface: %s", interfaceName))
	} else {
		if err := interfaceManager.ValidateInterface(cfg.Interface); err != nil {
			return fmt.Errorf("interface validation failed: %w", err)
		}
		interfaceName = cfg.Interface
	}

	// Initialize metrics collector
	metricsCollector := metrics.NewSystemMetricsCollector(logger.WithComponent("metrics").Logger)

	// Initialize packet capture engine
	captureEngine := capture.NewPacketCaptureEngine(logger.WithComponent("capture").Logger)
	captureEngine.SetMetricsCollector(metricsCollector)

	// Start packet capture
	if err := captureEngine.StartCapture(interfaceName); err != nil {
		return fmt.Errorf("failed to start packet capture: %w", err)
	}

	logger.WithComponent("main").Info("Application initialized successfully")

	// TODO: In future stories, add:
	// - Flow processing
	// - Web server and API
	// - WebSocket handler
	// - Metrics HTTP endpoint

	// Wait for context cancellation (interrupt signal)
	<-ctx.Done()

	logger.WithComponent("main").Info("Shutting down application")

	// Graceful shutdown
	if captureEngine.IsRunning() {
		if err := captureEngine.StopCapture(); err != nil {
			logger.WithComponent("main").Error(fmt.Sprintf("Error stopping capture engine: %v", err))
		}
	}

	logger.WithComponent("main").Info("Application shutdown complete")

	return nil
}

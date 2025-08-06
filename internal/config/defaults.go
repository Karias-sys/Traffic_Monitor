package config

import "time"

func getDefaults() *Config {
	return &Config{
		// Server configuration - localhost-first design
		Host: "localhost",
		Port: 8080,

		// Capture configuration - optimized for performance
		Interface:       "any", // Capture on all interfaces by default
		SnapLength:      1600,  // Sufficient for most packets including headers
		Promiscuous:     false, // Start non-promiscuous for security
		Timeout:         100 * time.Millisecond, // Balance between responsiveness and CPU usage
		BufferSize:      32 * 1024 * 1024, // 32MB buffer for high throughput
		FlowTimeout:     5 * time.Minute,   // Flow idle timeout
		MaxFlows:        100000, // Maximum flows to track (memory limit consideration)
		CleanupInterval: 30 * time.Second, // Regular cleanup to maintain <5% CPU target

		// Logging configuration
		LogLevel:  "info",  // Default to info level
		LogFormat: "json",  // Structured logging for monitoring

		// Security configuration - optional by default
		EnableAuth: false, // Localhost-first, optional security
		AuthToken:  "",    // Empty by default

		// Performance configuration
		CPUProfile:    "", // No profiling by default
		MemProfile:    "", // No profiling by default
		MetricsPort:   9090, // Standard metrics port
		EnableMetrics: true,  // Enable metrics for monitoring

		// Development configuration
		DevMode: false, // Production mode by default
	}
}
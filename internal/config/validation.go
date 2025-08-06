package config

import (
	"fmt"
	"net"
	"strings"
)

func Validate(cfg *Config) error {
	if err := validateNetwork(cfg); err != nil {
		return fmt.Errorf("network validation failed: %w", err)
	}

	if err := validateCapture(cfg); err != nil {
		return fmt.Errorf("capture validation failed: %w", err)
	}

	if err := validateLogging(cfg); err != nil {
		return fmt.Errorf("logging validation failed: %w", err)
	}

	if err := validateSecurity(cfg); err != nil {
		return fmt.Errorf("security validation failed: %w", err)
	}

	if err := validatePerformance(cfg); err != nil {
		return fmt.Errorf("performance validation failed: %w", err)
	}

	return nil
}

func validateNetwork(cfg *Config) error {
	// Validate host
	if cfg.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	// Validate IP address if not localhost
	if cfg.Host != "localhost" {
		if ip := net.ParseIP(cfg.Host); ip == nil {
			return fmt.Errorf("invalid host IP address: %s", cfg.Host)
		}
	}

	// Validate port range
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got: %d", cfg.Port)
	}

	// Validate metrics port
	if cfg.EnableMetrics {
		if cfg.MetricsPort < 1 || cfg.MetricsPort > 65535 {
			return fmt.Errorf("metrics port must be between 1 and 65535, got: %d", cfg.MetricsPort)
		}

		// Ensure metrics port is different from main port
		if cfg.MetricsPort == cfg.Port {
			return fmt.Errorf("metrics port cannot be the same as main port: %d", cfg.Port)
		}
	}

	return nil
}

func validateCapture(cfg *Config) error {
	// Validate interface (allow "any" as special case)
	if cfg.Interface == "" {
		return fmt.Errorf("interface cannot be empty")
	}

	// Validate snap length
	if cfg.SnapLength < 64 || cfg.SnapLength > 65535 {
		return fmt.Errorf("snap length must be between 64 and 65535, got: %d", cfg.SnapLength)
	}

	// Validate timeout
	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive, got: %v", cfg.Timeout)
	}

	// Validate buffer size (minimum 1KB, maximum 1GB)
	if cfg.BufferSize < 1024 {
		return fmt.Errorf("buffer size must be at least 1KB, got: %d", cfg.BufferSize)
	}
	if cfg.BufferSize > 1024*1024*1024 {
		return fmt.Errorf("buffer size must not exceed 1GB for memory management, got: %d", cfg.BufferSize)
	}

	// Validate flow timeout
	if cfg.FlowTimeout <= 0 {
		return fmt.Errorf("flow timeout must be positive, got: %v", cfg.FlowTimeout)
	}

	// Validate max flows (memory management constraint)
	if cfg.MaxFlows < 1000 {
		return fmt.Errorf("max flows must be at least 1000, got: %d", cfg.MaxFlows)
	}
	if cfg.MaxFlows > 1000000 {
		return fmt.Errorf("max flows must not exceed 1M for memory management, got: %d", cfg.MaxFlows)
	}

	// Validate cleanup interval
	if cfg.CleanupInterval <= 0 {
		return fmt.Errorf("cleanup interval must be positive, got: %v", cfg.CleanupInterval)
	}

	return nil
}

func validateLogging(cfg *Config) error {
	// Validate log level
	validLevels := []string{"debug", "info", "warn", "error"}
	validLevel := false
	for _, level := range validLevels {
		if strings.ToLower(cfg.LogLevel) == level {
			validLevel = true
			break
		}
	}
	if !validLevel {
		return fmt.Errorf("invalid log level: %s, must be one of: %v", cfg.LogLevel, validLevels)
	}

	// Validate log format
	validFormats := []string{"json", "text"}
	validFormat := false
	for _, format := range validFormats {
		if strings.ToLower(cfg.LogFormat) == format {
			validFormat = true
			break
		}
	}
	if !validFormat {
		return fmt.Errorf("invalid log format: %s, must be one of: %v", cfg.LogFormat, validFormats)
	}

	return nil
}

func validateSecurity(cfg *Config) error {
	// If auth is enabled, token must be provided
	if cfg.EnableAuth && cfg.AuthToken == "" {
		return fmt.Errorf("authentication token required when auth is enabled")
	}

	// Token length validation for security
	if cfg.EnableAuth && len(cfg.AuthToken) < 16 {
		return fmt.Errorf("authentication token must be at least 16 characters long")
	}

	return nil
}

func validatePerformance(cfg *Config) error {
	// Validate that flow storage won't exceed memory limits
	// Rough estimate: each flow ~1KB, max 1GB total
	estimatedMemory := cfg.MaxFlows * 1024 // bytes
	maxMemory := 1024 * 1024 * 1024        // 1GB

	if estimatedMemory > maxMemory {
		return fmt.Errorf("estimated flow memory usage (%d MB) exceeds 1GB limit", estimatedMemory/(1024*1024))
	}

	// Validate cleanup interval doesn't conflict with flow timeout
	if cfg.CleanupInterval > cfg.FlowTimeout {
		return fmt.Errorf("cleanup interval (%v) should not exceed flow timeout (%v)", cfg.CleanupInterval, cfg.FlowTimeout)
	}

	return nil
}

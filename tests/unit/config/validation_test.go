package config_test

import (
	"testing"
	"time"

	"github.com/Karias-sys/Traffic_Monitor/internal/config"
	"github.com/stretchr/testify/assert"
)

// getValidConfig returns a minimal valid configuration for testing
func getValidConfig(host string, port int, metricsPort int) *config.Config {
	return &config.Config{
		// Network settings
		Host:          host,
		Port:          port,
		MetricsPort:   metricsPort,
		EnableMetrics: true,
		
		// Capture settings
		Interface:       "eth0",
		SnapLength:      1600,
		Timeout:         100 * time.Millisecond,
		BufferSize:      32 * 1024 * 1024,
		FlowTimeout:     5 * time.Minute,
		MaxFlows:        100000,
		CleanupInterval: 30 * time.Second,
		
		// Logging settings
		LogLevel:  "info",
		LogFormat: "json",
		
		// Security settings
		EnableAuth: false,
		AuthToken:  "",
	}
}

func TestValidateNetwork(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *config.Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid localhost configuration",
			cfg: getValidConfig("localhost", 8080, 9090),
			wantError: false,
		},
		{
			name: "valid IP address configuration",
			cfg: getValidConfig("192.168.1.100", 8080, 9090),
			wantError: false,
		},
		{
			name: "empty host",
			cfg: func() *config.Config {
				cfg := getValidConfig("", 8080, 9090)
				cfg.Host = ""
				return cfg
			}(),
			wantError: true,
			errorMsg:  "host cannot be empty",
		},
		{
			name: "invalid IP address",
			cfg: func() *config.Config {
				cfg := getValidConfig("999.999.999.999", 8080, 9090)
				return cfg
			}(),
			wantError: true,
			errorMsg:  "invalid host IP address",
		},
		{
			name: "port too low",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 0, 9090)
				return cfg
			}(),
			wantError: true,
			errorMsg:  "port must be between 1 and 65535",
		},
		{
			name: "port too high",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 65536, 9090)
				return cfg
			}(),
			wantError: true,
			errorMsg:  "port must be between 1 and 65535",
		},
		{
			name: "metrics port same as main port",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 8080)
				return cfg
			}(),
			wantError: true,
			errorMsg:  "metrics port cannot be the same as main port",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.Validate(tt.cfg)
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCapture(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *config.Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid capture configuration",
			cfg: getValidConfig("localhost", 8080, 9090),
			wantError: false,
		},
		{
			name: "empty interface",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.Interface = ""
				return cfg
			}(),
			wantError: true,
			errorMsg:  "interface cannot be empty",
		},
		{
			name: "snap length too small",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.SnapLength = 32
				return cfg
			}(),
			wantError: true,
			errorMsg:  "snap length must be between 64 and 65535",
		},
		{
			name: "snap length too large",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.SnapLength = 70000
				return cfg
			}(),
			wantError: true,
			errorMsg:  "snap length must be between 64 and 65535",
		},
		{
			name: "negative timeout",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.Timeout = -1 * time.Second
				return cfg
			}(),
			wantError: true,
			errorMsg:  "timeout must be positive",
		},
		{
			name: "buffer size too small",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.BufferSize = 512
				return cfg
			}(),
			wantError: true,
			errorMsg:  "buffer size must be at least 1KB",
		},
		{
			name: "buffer size too large",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.BufferSize = 2 * 1024 * 1024 * 1024 // 2GB
				return cfg
			}(),
			wantError: true,
			errorMsg:  "buffer size must not exceed 1GB",
		},
		{
			name: "max flows too small",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.MaxFlows = 500
				return cfg
			}(),
			wantError: true,
			errorMsg:  "max flows must be at least 1000",
		},
		{
			name: "max flows too large",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.MaxFlows = 2000000
				return cfg
			}(),
			wantError: true,
			errorMsg:  "max flows must not exceed 1M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.Validate(tt.cfg)
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLogging(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *config.Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid logging configuration",
			cfg: getValidConfig("localhost", 8080, 9090),
			wantError: false,
		},
		{
			name: "valid debug level",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.LogLevel = "debug"
				cfg.LogFormat = "text"
				return cfg
			}(),
			wantError: false,
		},
		{
			name: "invalid log level",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.LogLevel = "trace"
				return cfg
			}(),
			wantError: true,
			errorMsg:  "invalid log level: trace",
		},
		{
			name: "invalid log format",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.LogFormat = "xml"
				return cfg
			}(),
			wantError: true,
			errorMsg:  "invalid log format: xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.Validate(tt.cfg)
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSecurity(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *config.Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "auth disabled",
			cfg: getValidConfig("localhost", 8080, 9090),
			wantError: false,
		},
		{
			name: "auth enabled with valid token",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.EnableAuth = true
				cfg.AuthToken = "valid-token-123456"
				return cfg
			}(),
			wantError: false,
		},
		{
			name: "auth enabled without token",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.EnableAuth = true
				cfg.AuthToken = ""
				return cfg
			}(),
			wantError: true,
			errorMsg:  "authentication token required when auth is enabled",
		},
		{
			name: "auth enabled with short token",
			cfg: func() *config.Config {
				cfg := getValidConfig("localhost", 8080, 9090)
				cfg.EnableAuth = true
				cfg.AuthToken = "short"
				return cfg
			}(),
			wantError: true,
			errorMsg:  "authentication token must be at least 16 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.Validate(tt.cfg)
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
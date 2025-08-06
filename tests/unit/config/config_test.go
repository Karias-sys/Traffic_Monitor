package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/Karias-sys/Traffic_Monitor/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDefault(t *testing.T) {
	// Clear any existing environment variables
	clearEnvVars()

	cfg, err := config.Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify default values
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "any", cfg.Interface)
	assert.Equal(t, int32(1600), cfg.SnapLength)
	assert.Equal(t, false, cfg.Promiscuous)
	assert.Equal(t, 100*time.Millisecond, cfg.Timeout)
	assert.Equal(t, 32*1024*1024, cfg.BufferSize)
	assert.Equal(t, 5*time.Minute, cfg.FlowTimeout)
	assert.Equal(t, 100000, cfg.MaxFlows)
	assert.Equal(t, 30*time.Second, cfg.CleanupInterval)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "json", cfg.LogFormat)
	assert.Equal(t, false, cfg.EnableAuth)
	assert.Equal(t, "", cfg.AuthToken)
	assert.Equal(t, 9090, cfg.MetricsPort)
	assert.Equal(t, true, cfg.EnableMetrics)
	assert.Equal(t, false, cfg.DevMode)
}

func TestLoadFromEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		validate func(t *testing.T, cfg *config.Config)
	}{
		{
			name: "host and port",
			envVars: map[string]string{
				"NETWATCH_HOST":         "0.0.0.0",
				"NETWATCH_PORT":         "8888",
				"NETWATCH_METRICS_PORT": "9999",
			},
			validate: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "0.0.0.0", cfg.Host)
				assert.Equal(t, 8888, cfg.Port)
				assert.Equal(t, 9999, cfg.MetricsPort)
			},
		},
		{
			name: "capture settings",
			envVars: map[string]string{
				"NETWATCH_INTERFACE":        "eth0",
				"NETWATCH_SNAP_LENGTH":      "2048",
				"NETWATCH_PROMISCUOUS":      "true",
				"NETWATCH_TIMEOUT":          "200ms",
				"NETWATCH_BUFFER_SIZE":      "65536",
				"NETWATCH_FLOW_TIMEOUT":     "10m",
				"NETWATCH_MAX_FLOWS":        "50000",
				"NETWATCH_CLEANUP_INTERVAL": "60s",
			},
			validate: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "eth0", cfg.Interface)
				assert.Equal(t, int32(2048), cfg.SnapLength)
				assert.Equal(t, true, cfg.Promiscuous)
				assert.Equal(t, 200*time.Millisecond, cfg.Timeout)
				assert.Equal(t, 65536, cfg.BufferSize)
				assert.Equal(t, 10*time.Minute, cfg.FlowTimeout)
				assert.Equal(t, 50000, cfg.MaxFlows)
				assert.Equal(t, 60*time.Second, cfg.CleanupInterval)
			},
		},
		{
			name: "logging settings",
			envVars: map[string]string{
				"NETWATCH_LOG_LEVEL":  "debug",
				"NETWATCH_LOG_FORMAT": "text",
			},
			validate: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "debug", cfg.LogLevel)
				assert.Equal(t, "text", cfg.LogFormat)
			},
		},
		{
			name: "security settings",
			envVars: map[string]string{
				"NETWATCH_ENABLE_AUTH": "true",
				"NETWATCH_AUTH_TOKEN":  "test-token-12345678",
			},
			validate: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, true, cfg.EnableAuth)
				assert.Equal(t, "test-token-12345678", cfg.AuthToken)
			},
		},
		{
			name: "performance settings",
			envVars: map[string]string{
				"NETWATCH_METRICS_PORT":   "8090",
				"NETWATCH_ENABLE_METRICS": "false",
				"NETWATCH_DEV_MODE":       "true",
			},
			validate: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, 8090, cfg.MetricsPort)
				assert.Equal(t, false, cfg.EnableMetrics)
				assert.Equal(t, true, cfg.DevMode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			clearEnvVars()

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Load configuration
			cfg, err := config.Load()
			require.NoError(t, err)
			require.NotNil(t, cfg)

			// Run validation
			tt.validate(t, cfg)

			// Clean up
			clearEnvVars()
		})
	}
}

func clearEnvVars() {
	envVars := []string{
		"NETWATCH_HOST",
		"NETWATCH_PORT",
		"NETWATCH_INTERFACE",
		"NETWATCH_SNAP_LENGTH",
		"NETWATCH_PROMISCUOUS",
		"NETWATCH_TIMEOUT",
		"NETWATCH_BUFFER_SIZE",
		"NETWATCH_FLOW_TIMEOUT",
		"NETWATCH_MAX_FLOWS",
		"NETWATCH_CLEANUP_INTERVAL",
		"NETWATCH_LOG_LEVEL",
		"NETWATCH_LOG_FORMAT",
		"NETWATCH_ENABLE_AUTH",
		"NETWATCH_AUTH_TOKEN",
		"NETWATCH_CPU_PROFILE",
		"NETWATCH_MEM_PROFILE",
		"NETWATCH_METRICS_PORT",
		"NETWATCH_ENABLE_METRICS",
		"NETWATCH_DEV_MODE",
	}

	for _, env := range envVars {
		os.Unsetenv(env)
	}
}

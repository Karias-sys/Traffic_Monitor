package config

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

const (
	// Security constants for safe integer conversion
	MaxInt32 = math.MaxInt32 // 2147483647
	MinInt32 = math.MinInt32 // -2147483648
)

type Config struct {
	// Server configuration
	Host string `json:"host"`
	Port int    `json:"port"`

	// Capture configuration
	Interface         string        `json:"interface"`
	SnapLength        int32         `json:"snap_length"`
	Promiscuous       bool          `json:"promiscuous"`
	Timeout           time.Duration `json:"timeout"`
	BufferSize        int           `json:"buffer_size"`
	RingBlockSize     uint32        `json:"ring_block_size"`
	RingBlockCount    uint32        `json:"ring_block_count"`
	ChannelBufferSize int           `json:"channel_buffer_size"`
	FlowTimeout       time.Duration `json:"flow_timeout"`
	MaxFlows          int           `json:"max_flows"`
	CleanupInterval   time.Duration `json:"cleanup_interval"`

	// Logging configuration
	LogLevel  string `json:"log_level"`
	LogFormat string `json:"log_format"`

	// Security configuration
	EnableAuth bool   `json:"enable_auth"`
	AuthToken  string `json:"auth_token"`

	// Performance configuration
	CPUProfile    string `json:"cpu_profile"`
	MemProfile    string `json:"mem_profile"`
	MetricsPort   int    `json:"metrics_port"`
	EnableMetrics bool   `json:"enable_metrics"`

	// Development configuration
	DevMode bool `json:"dev_mode"`
}

func Load() (*Config, error) {
	cfg := getDefaults()

	if err := loadFromEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	if err := loadFromFlags(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse command line flags: %w", err)
	}

	if err := Validate(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func LoadWithInterfaceValidation(validator InterfaceValidator) (*Config, error) {
	SetInterfaceValidator(validator)
	return Load()
}

func loadFromEnv(cfg *Config) error {
	if host := os.Getenv("NETWATCH_HOST"); host != "" {
		cfg.Host = host
	}

	if port := os.Getenv("NETWATCH_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}

	if iface := os.Getenv("NETWATCH_INTERFACE"); iface != "" {
		cfg.Interface = iface
	}

	if snapLen := os.Getenv("NETWATCH_SNAP_LENGTH"); snapLen != "" {
		if s, err := strconv.ParseInt(snapLen, 10, 32); err == nil {
			// Additional validation for security (though ParseInt already limits to int32 range)
			if s < 0 {
				cfg.SnapLength = 0 // Default to safe value for negative inputs
			} else {
				cfg.SnapLength = int32(s)
			}
		}
	}

	if promiscuous := os.Getenv("NETWATCH_PROMISCUOUS"); promiscuous != "" {
		if p, err := strconv.ParseBool(promiscuous); err == nil {
			cfg.Promiscuous = p
		}
	}

	if timeout := os.Getenv("NETWATCH_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			cfg.Timeout = t
		}
	}

	if bufferSize := os.Getenv("NETWATCH_BUFFER_SIZE"); bufferSize != "" {
		if b, err := strconv.Atoi(bufferSize); err == nil {
			cfg.BufferSize = b
		}
	}

	if ringBlockSize := os.Getenv("NETWATCH_RING_BLOCK_SIZE"); ringBlockSize != "" {
		if r, err := strconv.ParseUint(ringBlockSize, 10, 32); err == nil {
			cfg.RingBlockSize = uint32(r)
		}
	}

	if ringBlockCount := os.Getenv("NETWATCH_RING_BLOCK_COUNT"); ringBlockCount != "" {
		if r, err := strconv.ParseUint(ringBlockCount, 10, 32); err == nil {
			cfg.RingBlockCount = uint32(r)
		}
	}

	if channelBufferSize := os.Getenv("NETWATCH_CHANNEL_BUFFER_SIZE"); channelBufferSize != "" {
		if c, err := strconv.Atoi(channelBufferSize); err == nil {
			cfg.ChannelBufferSize = c
		}
	}

	if flowTimeout := os.Getenv("NETWATCH_FLOW_TIMEOUT"); flowTimeout != "" {
		if f, err := time.ParseDuration(flowTimeout); err == nil {
			cfg.FlowTimeout = f
		}
	}

	if maxFlows := os.Getenv("NETWATCH_MAX_FLOWS"); maxFlows != "" {
		if m, err := strconv.Atoi(maxFlows); err == nil {
			cfg.MaxFlows = m
		}
	}

	if cleanupInterval := os.Getenv("NETWATCH_CLEANUP_INTERVAL"); cleanupInterval != "" {
		if c, err := time.ParseDuration(cleanupInterval); err == nil {
			cfg.CleanupInterval = c
		}
	}

	if logLevel := os.Getenv("NETWATCH_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	if logFormat := os.Getenv("NETWATCH_LOG_FORMAT"); logFormat != "" {
		cfg.LogFormat = logFormat
	}

	if enableAuth := os.Getenv("NETWATCH_ENABLE_AUTH"); enableAuth != "" {
		if e, err := strconv.ParseBool(enableAuth); err == nil {
			cfg.EnableAuth = e
		}
	}

	if authToken := os.Getenv("NETWATCH_AUTH_TOKEN"); authToken != "" {
		cfg.AuthToken = authToken
	}

	if cpuProfile := os.Getenv("NETWATCH_CPU_PROFILE"); cpuProfile != "" {
		cfg.CPUProfile = cpuProfile
	}

	if memProfile := os.Getenv("NETWATCH_MEM_PROFILE"); memProfile != "" {
		cfg.MemProfile = memProfile
	}

	if metricsPort := os.Getenv("NETWATCH_METRICS_PORT"); metricsPort != "" {
		if m, err := strconv.Atoi(metricsPort); err == nil {
			cfg.MetricsPort = m
		}
	}

	if enableMetrics := os.Getenv("NETWATCH_ENABLE_METRICS"); enableMetrics != "" {
		if e, err := strconv.ParseBool(enableMetrics); err == nil {
			cfg.EnableMetrics = e
		}
	}

	if devMode := os.Getenv("NETWATCH_DEV_MODE"); devMode != "" {
		if d, err := strconv.ParseBool(devMode); err == nil {
			cfg.DevMode = d
		}
	}

	return nil
}

func loadFromFlags(cfg *Config) error {
	// Skip flag parsing if flags have already been parsed
	// This prevents issues in testing environments
	if flag.Parsed() {
		return nil
	}

	host := flag.String("host", cfg.Host, "Host to bind to")
	port := flag.Int("port", cfg.Port, "Port to listen on")
	iface := flag.String("interface", cfg.Interface, "Network interface to capture on")
	snapLength := flag.Int("snap-length", int(cfg.SnapLength), "Maximum packet capture length")
	promiscuous := flag.Bool("promiscuous", cfg.Promiscuous, "Enable promiscuous mode")
	timeout := flag.Duration("timeout", cfg.Timeout, "Packet capture timeout")
	bufferSize := flag.Int("buffer-size", cfg.BufferSize, "Packet capture buffer size")
	ringBlockSize := flag.Uint("ring-block-size", uint(cfg.RingBlockSize), "Ring buffer block size")
	ringBlockCount := flag.Uint("ring-block-count", uint(cfg.RingBlockCount), "Ring buffer block count")
	channelBufferSize := flag.Int("channel-buffer-size", cfg.ChannelBufferSize, "Packet channel buffer size")
	flowTimeout := flag.Duration("flow-timeout", cfg.FlowTimeout, "Flow timeout duration")
	maxFlows := flag.Int("max-flows", cfg.MaxFlows, "Maximum number of flows to track")
	cleanupInterval := flag.Duration("cleanup-interval", cfg.CleanupInterval, "Flow cleanup interval")
	logLevel := flag.String("log-level", cfg.LogLevel, "Logging level (debug, info, warn, error)")
	logFormat := flag.String("log-format", cfg.LogFormat, "Log format (json, text)")
	enableAuth := flag.Bool("enable-auth", cfg.EnableAuth, "Enable authentication")
	authToken := flag.String("auth-token", cfg.AuthToken, "Authentication token")
	cpuProfile := flag.String("cpu-profile", cfg.CPUProfile, "Write CPU profile to file")
	memProfile := flag.String("mem-profile", cfg.MemProfile, "Write memory profile to file")
	metricsPort := flag.Int("metrics-port", cfg.MetricsPort, "Metrics endpoint port")
	enableMetrics := flag.Bool("enable-metrics", cfg.EnableMetrics, "Enable metrics endpoint")
	devMode := flag.Bool("dev-mode", cfg.DevMode, "Enable development mode")

	flag.Parse()

	cfg.Host = *host
	cfg.Port = *port
	cfg.Interface = *iface
	// Secure conversion with bounds checking to prevent integer overflow
	if *snapLength < 0 || *snapLength > MaxInt32 {
		return fmt.Errorf("snap-length must be between 0 and %d, got: %d", MaxInt32, *snapLength)
	}
	cfg.SnapLength = int32(*snapLength)
	cfg.Promiscuous = *promiscuous
	cfg.Timeout = *timeout
	cfg.BufferSize = *bufferSize
	cfg.RingBlockSize = uint32(*ringBlockSize)
	cfg.RingBlockCount = uint32(*ringBlockCount)
	cfg.ChannelBufferSize = *channelBufferSize
	cfg.FlowTimeout = *flowTimeout
	cfg.MaxFlows = *maxFlows
	cfg.CleanupInterval = *cleanupInterval
	cfg.LogLevel = *logLevel
	cfg.LogFormat = *logFormat
	cfg.EnableAuth = *enableAuth
	cfg.AuthToken = *authToken
	cfg.CPUProfile = *cpuProfile
	cfg.MemProfile = *memProfile
	cfg.MetricsPort = *metricsPort
	cfg.EnableMetrics = *enableMetrics
	cfg.DevMode = *devMode

	return nil
}

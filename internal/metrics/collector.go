package metrics

import (
	"log/slog"
	"sync"
	"time"
)

type SystemMetricsCollector struct {
	mu                *sync.RWMutex
	logger            *slog.Logger
	captureStatistics CaptureMetrics
	systemStatistics  SystemMetrics
	enabled           bool
}

type CaptureMetrics struct {
	PacketsReceived    uint64    `json:"packets_received"`
	PacketsDropped     uint64    `json:"packets_dropped"`
	BytesReceived      uint64    `json:"bytes_received"`
	RingUtilization    float64   `json:"ring_utilization"`
	ErrorCount         uint64    `json:"error_count"`
	LastPacketTime     time.Time `json:"last_packet_time"`
	CaptureStartTime   time.Time `json:"capture_start_time"`
	UptimeSeconds      float64   `json:"uptime_seconds"`
}

type SystemMetrics struct {
	CPUUsagePercent    float64   `json:"cpu_usage_percent"`
	MemoryUsageMB      float64   `json:"memory_usage_mb"`
	MemoryUsagePercent float64   `json:"memory_usage_percent"`
	GoroutineCount     int       `json:"goroutine_count"`
	LastUpdateTime     time.Time `json:"last_update_time"`
}

type AllMetrics struct {
	Capture CaptureMetrics `json:"capture"`
	System  SystemMetrics  `json:"system"`
	Updated time.Time      `json:"updated"`
}

func NewSystemMetricsCollector(logger *slog.Logger) *SystemMetricsCollector {
	return &SystemMetricsCollector{
		mu:      &sync.RWMutex{},
		logger:  logger,
		enabled: true,
	}
}

func (smc *SystemMetricsCollector) UpdateCaptureMetrics(
	packetsReceived, packetsDropped, bytesReceived, errorCount uint64,
	ringUtilization float64,
	lastPacketTime time.Time,
	captureStartTime time.Time,
) {
	smc.mu.Lock()
	defer smc.mu.Unlock()

	if !smc.enabled {
		return
	}

	smc.captureStatistics = CaptureMetrics{
		PacketsReceived:  packetsReceived,
		PacketsDropped:   packetsDropped,
		BytesReceived:    bytesReceived,
		RingUtilization:  ringUtilization,
		ErrorCount:       errorCount,
		LastPacketTime:   lastPacketTime,
		CaptureStartTime: captureStartTime,
		UptimeSeconds:    time.Since(captureStartTime).Seconds(),
	}

	smc.logger.Debug("updated capture metrics",
		slog.Uint64("packets_received", packetsReceived),
		slog.Uint64("packets_dropped", packetsDropped),
		slog.Uint64("bytes_received", bytesReceived),
		slog.Float64("ring_utilization", ringUtilization),
		slog.Uint64("error_count", errorCount))
}

func (smc *SystemMetricsCollector) UpdateSystemMetrics(
	cpuPercent, memoryMB, memoryPercent float64,
	goroutineCount int,
) {
	smc.mu.Lock()
	defer smc.mu.Unlock()

	if !smc.enabled {
		return
	}

	smc.systemStatistics = SystemMetrics{
		CPUUsagePercent:    cpuPercent,
		MemoryUsageMB:      memoryMB,
		MemoryUsagePercent: memoryPercent,
		GoroutineCount:     goroutineCount,
		LastUpdateTime:     time.Now(),
	}

	smc.logger.Debug("updated system metrics",
		slog.Float64("cpu_percent", cpuPercent),
		slog.Float64("memory_mb", memoryMB),
		slog.Float64("memory_percent", memoryPercent),
		slog.Int("goroutines", goroutineCount))
}

func (smc *SystemMetricsCollector) GetCaptureMetrics() CaptureMetrics {
	smc.mu.RLock()
	defer smc.mu.RUnlock()
	return smc.captureStatistics
}

func (smc *SystemMetricsCollector) GetSystemMetrics() SystemMetrics {
	smc.mu.RLock()
	defer smc.mu.RUnlock()
	return smc.systemStatistics
}

func (smc *SystemMetricsCollector) GetAllMetrics() AllMetrics {
	smc.mu.RLock()
	defer smc.mu.RUnlock()

	return AllMetrics{
		Capture: smc.captureStatistics,
		System:  smc.systemStatistics,
		Updated: time.Now(),
	}
}

func (smc *SystemMetricsCollector) Enable() {
	smc.mu.Lock()
	defer smc.mu.Unlock()
	smc.enabled = true
	smc.logger.Info("metrics collection enabled")
}

func (smc *SystemMetricsCollector) Disable() {
	smc.mu.Lock()
	defer smc.mu.Unlock()
	smc.enabled = false
	smc.logger.Info("metrics collection disabled")
}

func (smc *SystemMetricsCollector) IsEnabled() bool {
	smc.mu.RLock()
	defer smc.mu.RUnlock()
	return smc.enabled
}

func (smc *SystemMetricsCollector) Reset() {
	smc.mu.Lock()
	defer smc.mu.Unlock()

	smc.captureStatistics = CaptureMetrics{}
	smc.systemStatistics = SystemMetrics{}

	smc.logger.Info("metrics reset")
}
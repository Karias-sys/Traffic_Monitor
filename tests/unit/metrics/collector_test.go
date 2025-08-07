package metrics

import (
	"testing"
	"time"

	"github.com/Karias-sys/Traffic_Monitor/internal/metrics"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
)

func createTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func TestNewSystemMetricsCollector(t *testing.T) {
	logger := createTestLogger()
	collector := metrics.NewSystemMetricsCollector(logger)

	assert.NotNil(t, collector)
	assert.True(t, collector.IsEnabled())
}

func TestSystemMetricsCollector_UpdateCaptureMetrics(t *testing.T) {
	logger := createTestLogger()
	collector := metrics.NewSystemMetricsCollector(logger)

	now := time.Now()
	startTime := now.Add(-time.Minute)

	collector.UpdateCaptureMetrics(
		100,   // packetsReceived
		5,     // packetsDropped
		50000, // bytesReceived
		2,     // errorCount
		0.25,  // ringUtilization
		now,   // lastPacketTime
		startTime, // captureStartTime
	)

	metrics := collector.GetCaptureMetrics()
	assert.Equal(t, uint64(100), metrics.PacketsReceived)
	assert.Equal(t, uint64(5), metrics.PacketsDropped)
	assert.Equal(t, uint64(50000), metrics.BytesReceived)
	assert.Equal(t, uint64(2), metrics.ErrorCount)
	assert.Equal(t, 0.25, metrics.RingUtilization)
	assert.Equal(t, now, metrics.LastPacketTime)
	assert.Equal(t, startTime, metrics.CaptureStartTime)
	assert.Equal(t, 60.0, metrics.UptimeSeconds)
}

func TestSystemMetricsCollector_UpdateSystemMetrics(t *testing.T) {
	logger := createTestLogger()
	collector := metrics.NewSystemMetricsCollector(logger)

	collector.UpdateSystemMetrics(
		15.5,  // cpuPercent
		256.0, // memoryMB
		25.0,  // memoryPercent
		150,   // goroutineCount
	)

	systemMetrics := collector.GetSystemMetrics()
	assert.Equal(t, 15.5, systemMetrics.CPUUsagePercent)
	assert.Equal(t, 256.0, systemMetrics.MemoryUsageMB)
	assert.Equal(t, 25.0, systemMetrics.MemoryUsagePercent)
	assert.Equal(t, 150, systemMetrics.GoroutineCount)
	assert.False(t, systemMetrics.LastUpdateTime.IsZero())
}

func TestSystemMetricsCollector_GetAllMetrics(t *testing.T) {
	logger := createTestLogger()
	collector := metrics.NewSystemMetricsCollector(logger)

	now := time.Now()
	startTime := now.Add(-time.Minute)

	// Update both capture and system metrics
	collector.UpdateCaptureMetrics(200, 10, 100000, 1, 0.5, now, startTime)
	collector.UpdateSystemMetrics(20.0, 512.0, 50.0, 200)

	allMetrics := collector.GetAllMetrics()

	// Check capture metrics
	assert.Equal(t, uint64(200), allMetrics.Capture.PacketsReceived)
	assert.Equal(t, uint64(10), allMetrics.Capture.PacketsDropped)
	assert.Equal(t, uint64(100000), allMetrics.Capture.BytesReceived)
	assert.Equal(t, 0.5, allMetrics.Capture.RingUtilization)

	// Check system metrics
	assert.Equal(t, 20.0, allMetrics.System.CPUUsagePercent)
	assert.Equal(t, 512.0, allMetrics.System.MemoryUsageMB)
	assert.Equal(t, 50.0, allMetrics.System.MemoryUsagePercent)
	assert.Equal(t, 200, allMetrics.System.GoroutineCount)

	// Check timestamp
	assert.False(t, allMetrics.Updated.IsZero())
}

func TestSystemMetricsCollector_EnableDisable(t *testing.T) {
	logger := createTestLogger()
	collector := metrics.NewSystemMetricsCollector(logger)

	// Initially enabled
	assert.True(t, collector.IsEnabled())

	// Disable
	collector.Disable()
	assert.False(t, collector.IsEnabled())

	// Try to update metrics when disabled
	now := time.Now()
	collector.UpdateCaptureMetrics(100, 5, 50000, 2, 0.25, now, now)
	collector.UpdateSystemMetrics(15.5, 256.0, 25.0, 150)

	// Metrics should not be updated
	captureMetrics := collector.GetCaptureMetrics()
	systemMetrics := collector.GetSystemMetrics()

	assert.Equal(t, uint64(0), captureMetrics.PacketsReceived)
	assert.Equal(t, 0.0, systemMetrics.CPUUsagePercent)

	// Re-enable
	collector.Enable()
	assert.True(t, collector.IsEnabled())

	// Now updates should work
	collector.UpdateCaptureMetrics(100, 5, 50000, 2, 0.25, now, now)
	captureMetrics = collector.GetCaptureMetrics()
	assert.Equal(t, uint64(100), captureMetrics.PacketsReceived)
}

func TestSystemMetricsCollector_Reset(t *testing.T) {
	logger := createTestLogger()
	collector := metrics.NewSystemMetricsCollector(logger)

	now := time.Now()

	// Update metrics
	collector.UpdateCaptureMetrics(100, 5, 50000, 2, 0.25, now, now)
	collector.UpdateSystemMetrics(15.5, 256.0, 25.0, 150)

	// Verify metrics are set
	captureMetrics := collector.GetCaptureMetrics()
	systemMetrics := collector.GetSystemMetrics()
	assert.Equal(t, uint64(100), captureMetrics.PacketsReceived)
	assert.Equal(t, 15.5, systemMetrics.CPUUsagePercent)

	// Reset
	collector.Reset()

	// Verify metrics are cleared
	captureMetrics = collector.GetCaptureMetrics()
	systemMetrics = collector.GetSystemMetrics()
	assert.Equal(t, uint64(0), captureMetrics.PacketsReceived)
	assert.Equal(t, 0.0, systemMetrics.CPUUsagePercent)
	assert.True(t, captureMetrics.LastPacketTime.IsZero())
	assert.True(t, systemMetrics.LastUpdateTime.IsZero())
}

func TestSystemMetricsCollector_ConcurrentAccess(t *testing.T) {
	logger := createTestLogger()
	collector := metrics.NewSystemMetricsCollector(logger)

	// Test concurrent updates and reads
	done := make(chan bool, 2)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			now := time.Now()
			collector.UpdateCaptureMetrics(
				uint64(i), uint64(i/10), uint64(i*1000), uint64(i/50),
				float64(i)/100.0, now, now,
			)
			collector.UpdateSystemMetrics(
				float64(i), float64(i*10), float64(i), i*2,
			)
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			_ = collector.GetCaptureMetrics()
			_ = collector.GetSystemMetrics()
			_ = collector.GetAllMetrics()
		}
		done <- true
	}()

	// Wait for completion
	<-done
	<-done

	// Verify final state is consistent
	allMetrics := collector.GetAllMetrics()
	assert.GreaterOrEqual(t, allMetrics.Capture.PacketsReceived, uint64(0))
	assert.GreaterOrEqual(t, allMetrics.System.CPUUsagePercent, 0.0)
}
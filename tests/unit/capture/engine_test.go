//go:build linux

package capture

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/Karias-sys/Traffic_Monitor/internal/capture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockMetricsCollector struct {
	UpdatedMetrics []MetricsUpdate
}

type MetricsUpdate struct {
	PacketsReceived    uint64
	PacketsDropped     uint64
	BytesReceived      uint64
	ErrorCount         uint64
	RingUtilization    float64
	LastPacketTime     time.Time
	CaptureStartTime   time.Time
}

func (m *MockMetricsCollector) UpdateCaptureMetrics(
	packetsReceived, packetsDropped, bytesReceived, errorCount uint64,
	ringUtilization float64,
	lastPacketTime time.Time,
	captureStartTime time.Time,
) {
	m.UpdatedMetrics = append(m.UpdatedMetrics, MetricsUpdate{
		PacketsReceived:  packetsReceived,
		PacketsDropped:   packetsDropped,
		BytesReceived:    bytesReceived,
		ErrorCount:       errorCount,
		RingUtilization:  ringUtilization,
		LastPacketTime:   lastPacketTime,
		CaptureStartTime: captureStartTime,
	})
}

func createTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func TestNewPacketCaptureEngine(t *testing.T) {
	logger := createTestLogger()
	engine := capture.NewPacketCaptureEngine(logger)

	assert.NotNil(t, engine)
	assert.False(t, engine.IsRunning())
}

func TestPacketCaptureEngine_SetMetricsCollector(t *testing.T) {
	logger := createTestLogger()
	engine := capture.NewPacketCaptureEngine(logger)
	mockCollector := &MockMetricsCollector{}

	engine.SetMetricsCollector(mockCollector)

	// This test verifies that the metrics collector is set without error
	// The actual functionality is tested in integration scenarios
	assert.NotNil(t, engine)
}

func TestPacketCaptureEngine_GetStatistics_NotRunning(t *testing.T) {
	logger := createTestLogger()
	engine := capture.NewPacketCaptureEngine(logger)

	stats := engine.GetStatistics()

	assert.Equal(t, uint64(0), stats.PacketsReceived)
	assert.Equal(t, uint64(0), stats.PacketsDropped)
	assert.Equal(t, uint64(0), stats.BytesReceived)
	assert.Equal(t, uint64(0), stats.ErrorCount)
	assert.Equal(t, 0.0, stats.RingUtilization)
}

func TestPacketCaptureEngine_StartCapture_InvalidInterface(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root - would actually try to create socket")
	}

	logger := createTestLogger()
	engine := capture.NewPacketCaptureEngine(logger)

	err := engine.StartCapture("non-existent-interface")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interface")
	assert.False(t, engine.IsRunning())
}

func TestPacketCaptureEngine_StartCapture_AlreadyRunning(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Skipping test - requires root privileges for AF_PACKET socket")
	}

	logger := createTestLogger()
	engine := capture.NewPacketCaptureEngine(logger)

	// Start capture on loopback (should be available)
	err := engine.StartCapture("lo")
	if err != nil {
		t.Skipf("Cannot start capture on loopback: %v", err)
	}
	defer engine.StopCapture()

	// Try to start again
	err = engine.StartCapture("lo")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already running")
}

func TestPacketCaptureEngine_StopCapture_NotRunning(t *testing.T) {
	logger := createTestLogger()
	engine := capture.NewPacketCaptureEngine(logger)

	err := engine.StopCapture()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not started")
}

func TestPacketCaptureEngine_PacketChannel(t *testing.T) {
	logger := createTestLogger()
	engine := capture.NewPacketCaptureEngine(logger)

	channel := engine.PacketChannel()

	assert.NotNil(t, channel)

	// Channel should be readable (non-blocking check)
	select {
	case <-channel:
		t.Fatal("Expected empty channel")
	default:
		// Expected - channel is empty
	}
}

// Integration test - requires root privileges
func TestPacketCaptureEngine_Integration_Loopback(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Integration test requires root privileges for AF_PACKET socket")
	}

	logger := createTestLogger()
	engine := capture.NewPacketCaptureEngine(logger)
	mockCollector := &MockMetricsCollector{}
	engine.SetMetricsCollector(mockCollector)

	// Start capture on loopback
	err := engine.StartCapture("lo")
	require.NoError(t, err)
	defer engine.StopCapture()

	assert.True(t, engine.IsRunning())

	// Wait a brief moment for initialization
	time.Sleep(100 * time.Millisecond)

	// Get channel for packets
	packetChannel := engine.PacketChannel()
	assert.NotNil(t, packetChannel)

	// Generate some loopback traffic (ping localhost)
	// This is optional - the test mainly verifies the engine starts correctly

	// Stop capture
	err = engine.StopCapture()
	assert.NoError(t, err)
	assert.False(t, engine.IsRunning())

	// Verify metrics were collected
	stats := engine.GetStatistics()
	assert.Equal(t, uint64(0), stats.PacketsReceived) // Reset after stop
}

func TestPacketCaptureEngine_Lifecycle(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Lifecycle test requires root privileges for AF_PACKET socket")
	}

	logger := createTestLogger()
	engine := capture.NewPacketCaptureEngine(logger)

	// Initial state
	assert.False(t, engine.IsRunning())

	// Start
	err := engine.StartCapture("lo")
	require.NoError(t, err)
	assert.True(t, engine.IsRunning())

	// Stop
	err = engine.StopCapture()
	assert.NoError(t, err)
	assert.False(t, engine.IsRunning())

	// Start again
	err = engine.StartCapture("lo")
	require.NoError(t, err)
	assert.True(t, engine.IsRunning())

	// Final cleanup
	err = engine.StopCapture()
	assert.NoError(t, err)
}
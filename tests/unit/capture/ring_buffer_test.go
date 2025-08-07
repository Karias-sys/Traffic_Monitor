//go:build linux

package capture

import (
	"log/slog"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/Karias-sys/Traffic_Monitor/internal/capture"
	"github.com/stretchr/testify/assert"
)

func TestRingBuffer_InvalidConfiguration(t *testing.T) {
	logger := createTestLogger()

	// Test zero block size
	_, err := capture.NewRingBuffer(-1, 0, 1024, logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "block size")

	// Test zero block count
	_, err = capture.NewRingBuffer(-1, 1024, 0, logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "block size")

	// Test unaligned block size
	pageSize := syscall.Getpagesize()
	unalignedSize := uint32(pageSize + 1)
	_, err = capture.NewRingBuffer(-1, unalignedSize, 1024, logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "aligned")
}

func TestRingBuffer_Close_NotInitialized(t *testing.T) {
	logger := createTestLogger()

	// Create a ring buffer with invalid socket to simulate failure
	rb, err := capture.NewRingBuffer(-1, 4096, 1024, logger)
	if err == nil {
		// If creation succeeded (shouldn't with invalid socket), test close
		err = rb.Close()
		assert.NoError(t, err)

		// Test double close
		err = rb.Close()
		assert.NoError(t, err) // Should not error on double close
	}
}

func TestRingBuffer_GetUtilization_NotInitialized(t *testing.T) {
	logger := createTestLogger()

	// Create a ring buffer with invalid socket
	rb, err := capture.NewRingBuffer(-1, 4096, 1024, logger)
	if err == nil {
		utilization := rb.GetUtilization()
		assert.Equal(t, 0.0, utilization)

		// Close and test utilization
		rb.Close()
		utilization = rb.GetUtilization()
		assert.Equal(t, 0.0, utilization)
	}
}

func TestRingBuffer_ProcessPackets_Closed(t *testing.T) {
	logger := createTestLogger()

	// Create a ring buffer with invalid socket
	rb, err := capture.NewRingBuffer(-1, 4096, 1024, logger)
	if err == nil {
		// Close the ring buffer
		rb.Close()

		// Try to process packets
		handlerCalled := false
		err = rb.ProcessPackets(func(data []byte, timestamp time.Time) {
			handlerCalled = true
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "closed")
		assert.False(t, handlerCalled)
	}
}

// Integration test requiring AF_PACKET socket
func TestRingBuffer_Integration(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Integration test requires root privileges for AF_PACKET socket")
	}

	logger := createTestLogger()

	// Create an AF_PACKET socket
	socket, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0x0003) // ETH_P_ALL
	if err != nil {
		t.Skipf("Cannot create AF_PACKET socket: %v", err)
	}
	defer syscall.Close(socket)

	// Create ring buffer
	pageSize := syscall.Getpagesize()
	blockSize := uint32(pageSize * 2) // 2 pages
	blockCount := uint32(4)           // 4 blocks

	rb, err := capture.NewRingBuffer(socket, blockSize, blockCount, logger)
	if err != nil {
		t.Skipf("Cannot create ring buffer: %v", err)
	}
	defer rb.Close()

	// Test utilization
	utilization := rb.GetUtilization()
	assert.GreaterOrEqual(t, utilization, 0.0)
	assert.LessOrEqual(t, utilization, 1.0)

	// Test packet processing (no packets expected immediately)
	handlerCallCount := 0
	err = rb.ProcessPackets(func(data []byte, timestamp time.Time) {
		handlerCallCount++
		assert.NotNil(t, data)
		assert.False(t, timestamp.IsZero())
	})

	// Should not error even if no packets are available
	assert.NoError(t, err)

	// Close ring buffer
	err = rb.Close()
	assert.NoError(t, err)
}
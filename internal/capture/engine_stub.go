//go:build !linux

package capture

import (
	"errors"
	"log/slog"
	"time"
)

var (
	ErrPlatformNotSupported = errors.New("AF_PACKET capture is only supported on Linux")
)

type RawPacket struct {
	Timestamp time.Time
	Interface int
	Data      []byte
	Length    uint32
}

type CaptureStatistics struct {
	PacketsReceived    uint64
	PacketsDropped     uint64
	RingUtilization    float64
	LastPacketTime     time.Time
	BytesReceived      uint64
	ErrorCount         uint64
}

type MetricsCollector interface {
	UpdateCaptureMetrics(packetsReceived, packetsDropped, bytesReceived, errorCount uint64,
		ringUtilization float64, lastPacketTime time.Time, captureStartTime time.Time)
}

type PacketCaptureEngine struct {
	logger *slog.Logger
}

func NewPacketCaptureEngine(logger *slog.Logger) *PacketCaptureEngine {
	return &PacketCaptureEngine{
		logger: logger,
	}
}

func (e *PacketCaptureEngine) StartCapture(interfaceName string) error {
	e.logger.Error("AF_PACKET capture is not supported on this platform")
	return ErrPlatformNotSupported
}

func (e *PacketCaptureEngine) StopCapture() error {
	return ErrPlatformNotSupported
}

func (e *PacketCaptureEngine) PacketChannel() <-chan RawPacket {
	return make(<-chan RawPacket)
}

func (e *PacketCaptureEngine) GetStatistics() CaptureStatistics {
	return CaptureStatistics{}
}

func (e *PacketCaptureEngine) IsRunning() bool {
	return false
}

func (e *PacketCaptureEngine) SetMetricsCollector(collector MetricsCollector) {
	// No-op on unsupported platforms
}
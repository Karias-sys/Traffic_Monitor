//go:build linux

package capture

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

var (
	ErrEngineNotStarted = errors.New("packet capture engine not started")
	ErrEngineRunning    = errors.New("packet capture engine already running")
	ErrSocketCreation   = errors.New("failed to create AF_PACKET socket")
	ErrRingSetup        = errors.New("failed to setup ring buffer")
	ErrInterfaceBind    = errors.New("failed to bind socket to interface")
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
	mu                *sync.RWMutex
	logger            *slog.Logger
	socket            int
	interfaceIndex    int
	interfaceName     string
	running           bool
	ringBuffer        *RingBuffer
	packetChannel     chan RawPacket
	ctx               context.Context
	cancel            context.CancelFunc
	statistics        CaptureStatistics
	statisticsMu      *sync.RWMutex
	metricsCollector  MetricsCollector
	captureStartTime  time.Time
}

type EngineConfig struct {
	RingBlockSize    uint32
	RingFrameCount   uint32
	ChannelBufferSize int
}

func NewPacketCaptureEngine(logger *slog.Logger) *PacketCaptureEngine {
	return &PacketCaptureEngine{
		mu:             &sync.RWMutex{},
		logger:         logger,
		socket:         -1,
		packetChannel:  make(chan RawPacket, 1000),
		statisticsMu:   &sync.RWMutex{},
	}
}

func (e *PacketCaptureEngine) StartCapture(interfaceName string) error {
	return e.StartCaptureWithConfig(interfaceName, EngineConfig{
		RingBlockSize:     32 * 1024 * 1024,
		RingFrameCount:    1024,
		ChannelBufferSize: 1000,
	})
}

func (e *PacketCaptureEngine) StartCaptureWithConfig(interfaceName string, config EngineConfig) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return ErrEngineRunning
	}

	e.logger.Info("starting packet capture", slog.String("interface", interfaceName))

	interfaceIndex, err := e.getInterfaceIndex(interfaceName)
	if err != nil {
		e.logger.Error("failed to get interface index", 
			slog.String("interface", interfaceName),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to get interface index for %s: %w", interfaceName, err)
	}

	socket, err := e.createSocket()
	if err != nil {
		e.logger.Error("failed to create AF_PACKET socket", slog.String("error", err.Error()))
		return fmt.Errorf("%w: %v", ErrSocketCreation, err)
	}
	defer func() {
		if err != nil {
			unix.Close(socket)
		}
	}()

	if err = e.setupTPACKETv3(socket); err != nil {
		e.logger.Error("failed to setup TPACKETv3", slog.String("error", err.Error()))
		return fmt.Errorf("%w: %v", ErrRingSetup, err)
	}

	if err = e.bindSocket(socket, interfaceIndex); err != nil {
		e.logger.Error("failed to bind socket to interface", 
			slog.String("interface", interfaceName),
			slog.String("error", err.Error()))
		return fmt.Errorf("%w: interface %s: %v", ErrInterfaceBind, interfaceName, err)
	}

	ringBuffer, err := NewRingBuffer(socket, config.RingBlockSize, config.RingFrameCount, e.logger)
	if err != nil {
		e.logger.Error("failed to create ring buffer", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create ring buffer: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	e.socket = socket
	e.interfaceIndex = interfaceIndex
	e.interfaceName = interfaceName
	e.ringBuffer = ringBuffer
	e.ctx = ctx
	e.cancel = cancel
	e.running = true
	e.captureStartTime = time.Now()

	if config.ChannelBufferSize > 0 {
		e.packetChannel = make(chan RawPacket, config.ChannelBufferSize)
	}

	e.resetStatistics()

	go e.captureLoop()

	e.logger.Info("packet capture started successfully",
		slog.String("interface", interfaceName),
		slog.Int("interface_index", interfaceIndex))

	return nil
}

func (e *PacketCaptureEngine) StopCapture() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.running {
		return ErrEngineNotStarted
	}

	e.logger.Info("stopping packet capture", slog.String("interface", e.interfaceName))

	e.cancel()
	e.running = false

	if e.ringBuffer != nil {
		if err := e.ringBuffer.Close(); err != nil {
			e.logger.Error("failed to close ring buffer", slog.String("error", err.Error()))
		}
		e.ringBuffer = nil
	}

	if e.socket != -1 {
		if err := unix.Close(e.socket); err != nil {
			e.logger.Error("failed to close socket", slog.String("error", err.Error()))
		}
		e.socket = -1
	}

	close(e.packetChannel)
	e.packetChannel = make(chan RawPacket, 1000)

	e.logger.Info("packet capture stopped")
	return nil
}

func (e *PacketCaptureEngine) PacketChannel() <-chan RawPacket {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.packetChannel
}

func (e *PacketCaptureEngine) GetStatistics() CaptureStatistics {
	e.statisticsMu.RLock()
	stats := e.statistics
	e.statisticsMu.RUnlock()

	if e.ringBuffer != nil {
		stats.RingUtilization = e.ringBuffer.GetUtilization()
	}

	return stats
}

func (e *PacketCaptureEngine) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.running
}

func (e *PacketCaptureEngine) SetMetricsCollector(collector MetricsCollector) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metricsCollector = collector
}

func (e *PacketCaptureEngine) getInterfaceIndex(interfaceName string) (int, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return 0, fmt.Errorf("interface %s not found: %w", interfaceName, err)
	}
	return iface.Index, nil
}

func (e *PacketCaptureEngine) createSocket() (int, error) {
	socket, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(unix.ETH_P_ALL))
	if err != nil {
		return -1, fmt.Errorf("socket creation failed: %w", err)
	}

	return socket, nil
}

func (e *PacketCaptureEngine) setupTPACKETv3(socket int) error {
	version := int(3)
	if err := unix.SetsockoptInt(socket, unix.SOL_PACKET, unix.PACKET_VERSION, version); err != nil {
		return fmt.Errorf("failed to set PACKET_VERSION: %w", err)
	}

	req := &tpacketReq3{
		blockSize:       32 * 1024,
		blockNum:        1024,
		frameSize:       2048,
		frameNum:        1024 * 16,
		retireBlkTov:    100,
		sizeofPriv:      0,
		featureReqWord: 0,
	}

	reqBytes := (*[unsafe.Sizeof(tpacketReq3{})]byte)(unsafe.Pointer(req))[:]
	if err := unix.SetsockoptString(socket, unix.SOL_PACKET, unix.PACKET_RX_RING, string(reqBytes)); err != nil {
		return fmt.Errorf("failed to set PACKET_RX_RING: %w", err)
	}

	return nil
}

func (e *PacketCaptureEngine) bindSocket(socket int, interfaceIndex int) error {
	addr := &unix.SockaddrLinklayer{
		Protocol: unix.ETH_P_ALL,
		Ifindex:  interfaceIndex,
	}

	if err := unix.Bind(socket, addr); err != nil {
		return fmt.Errorf("bind failed: %w", err)
	}

	return nil
}

func (e *PacketCaptureEngine) captureLoop() {
	e.logger.Debug("starting capture loop")
	defer e.logger.Debug("capture loop ended")

	pollFds := []unix.PollFd{
		{
			Fd:     int32(e.socket),
			Events: unix.POLLIN,
		},
	}

	for {
		select {
		case <-e.ctx.Done():
			return
		default:
			ready, err := unix.Poll(pollFds, 100)
			if err != nil {
				if err == unix.EINTR {
					continue
				}
				e.updateErrorCount()
				e.logger.Error("poll failed", slog.String("error", err.Error()))
				time.Sleep(10 * time.Millisecond)
				continue
			}

			if ready > 0 && pollFds[0].Revents&unix.POLLIN != 0 {
				if err := e.processRingBuffer(); err != nil {
					e.updateErrorCount()
					e.logger.Debug("error processing ring buffer", slog.String("error", err.Error()))
				}
			}
		}
	}
}

func (e *PacketCaptureEngine) processRingBuffer() error {
	return e.ringBuffer.ProcessPackets(func(data []byte, timestamp time.Time) {
		if len(data) == 0 {
			return
		}
		
		packet := RawPacket{
			Timestamp: timestamp,
			Interface: e.interfaceIndex,
			Data:      make([]byte, len(data)),
			Length:    uint32(len(data)),
		}
		copy(packet.Data, data)

		select {
		case e.packetChannel <- packet:
			e.updatePacketStatistics(packet)
		default:
			e.updateDroppedCount()
			e.logger.Debug("packet channel full, dropping packet")
		}
	})
}

func (e *PacketCaptureEngine) updatePacketStatistics(packet RawPacket) {
	e.statisticsMu.Lock()
	e.statistics.PacketsReceived++
	e.statistics.BytesReceived += uint64(packet.Length)
	e.statistics.LastPacketTime = packet.Timestamp
	
	stats := e.statistics
	e.statisticsMu.Unlock()

	if e.metricsCollector != nil {
		ringUtilization := 0.0
		if e.ringBuffer != nil {
			ringUtilization = e.ringBuffer.GetUtilization()
		}
		e.metricsCollector.UpdateCaptureMetrics(
			stats.PacketsReceived,
			stats.PacketsDropped,
			stats.BytesReceived,
			stats.ErrorCount,
			ringUtilization,
			stats.LastPacketTime,
			e.captureStartTime,
		)
	}
}

func (e *PacketCaptureEngine) updateDroppedCount() {
	e.statisticsMu.Lock()
	defer e.statisticsMu.Unlock()
	e.statistics.PacketsDropped++
}

func (e *PacketCaptureEngine) updateErrorCount() {
	e.statisticsMu.Lock()
	defer e.statisticsMu.Unlock()
	e.statistics.ErrorCount++
}

func (e *PacketCaptureEngine) resetStatistics() {
	e.statisticsMu.Lock()
	defer e.statisticsMu.Unlock()
	e.statistics = CaptureStatistics{}
}



type tpacketReq3 struct {
	blockSize       uint32
	blockNum        uint32
	frameSize       uint32
	frameNum        uint32
	retireBlkTov    uint32
	sizeofPriv      uint32
	featureReqWord  uint32
}
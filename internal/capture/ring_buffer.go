//go:build linux

package capture

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

var (
	ErrRingBufferClosed = errors.New("ring buffer is closed")
	ErrMMapFailed       = errors.New("memory mapping failed")
	ErrInvalidBuffer    = errors.New("invalid ring buffer configuration")
)

type PacketHandler func(data []byte, timestamp time.Time)

type RingBuffer struct {
	mu          *sync.RWMutex
	logger      *slog.Logger
	socket      int
	buffer      []byte
	blockSize   uint32
	blockCount  uint32
	currentBlock uint32
	closed      bool
}

type tpacketHdrV1 struct {
	blockStatus uint32
	numPkts     uint32
	offsetToFirst uint32
	blockLen    uint32
	seqNum      uint64
	tsFirst     tpacketBlockTS
	tsLast      tpacketBlockTS
}

type tpacketBlockTS struct {
	sec  uint32
	nsec uint32
}

type tpacket3Hdr struct {
	nextOffset uint32
	sec        uint32
	nsec       uint32
	snaplen    uint32
	len        uint32
	status     uint32
	mac        uint16
	net        uint16
	hv1        tpacketHdrV1
}

func NewRingBuffer(socket int, blockSize, blockCount uint32, logger *slog.Logger) (*RingBuffer, error) {
	if blockSize == 0 || blockCount == 0 {
		return nil, fmt.Errorf("%w: block size and count must be > 0", ErrInvalidBuffer)
	}

	if blockSize%uint32(syscall.Getpagesize()) != 0 {
		return nil, fmt.Errorf("%w: block size must be aligned to page size", ErrInvalidBuffer)
	}

	totalSize := int(blockSize * blockCount)
	
	buffer, err := unix.Mmap(socket, 0, totalSize, 
		unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMMapFailed, err)
	}

	rb := &RingBuffer{
		mu:           &sync.RWMutex{},
		logger:       logger,
		socket:       socket,
		buffer:       buffer,
		blockSize:    blockSize,
		blockCount:   blockCount,
		currentBlock: 0,
		closed:       false,
	}

	logger.Info("ring buffer created successfully",
		slog.Int("block_size", int(blockSize)),
		slog.Int("block_count", int(blockCount)),
		slog.Int("total_size", totalSize))

	return rb, nil
}

func (rb *RingBuffer) ProcessPackets(handler PacketHandler) error {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.closed {
		return ErrRingBufferClosed
	}

	blockOffset := int(rb.currentBlock * rb.blockSize)
	if blockOffset >= len(rb.buffer) {
		return fmt.Errorf("invalid block offset: %d >= %d", blockOffset, len(rb.buffer))
	}

	blockData := rb.buffer[blockOffset:]
	if len(blockData) < int(unsafe.Sizeof(tpacketHdrV1{})) {
		return fmt.Errorf("insufficient data for block header")
	}

	blockHdr := (*tpacketHdrV1)(unsafe.Pointer(&blockData[0]))

	if blockHdr.blockStatus&unix.TP_STATUS_KERNEL != 0 {
		return nil
	}

	if blockHdr.blockStatus&unix.TP_STATUS_USER == 0 {
		return nil
	}

	processed, err := rb.processBlockPackets(blockData, blockHdr, handler)
	if err != nil {
		rb.logger.Debug("error processing block packets", slog.String("error", err.Error()))
	}

	blockHdr.blockStatus = unix.TP_STATUS_KERNEL
	rb.currentBlock = (rb.currentBlock + 1) % rb.blockCount

	rb.logger.Debug("processed block", 
		slog.Uint64("packets", uint64(processed)),
		slog.Uint64("block", uint64(rb.currentBlock-1)),
		slog.Uint64("total_packets", uint64(blockHdr.numPkts)))

	return nil
}

func (rb *RingBuffer) processBlockPackets(blockData []byte, blockHdr *tpacketHdrV1, handler PacketHandler) (uint32, error) {
	packetOffset := int(blockHdr.offsetToFirst)
	packetsProcessed := uint32(0)
	maxPackets := blockHdr.numPkts

	for packetsProcessed < maxPackets {
		if packetOffset >= int(rb.blockSize) {
			return packetsProcessed, fmt.Errorf("packet offset beyond block size: %d >= %d", packetOffset, rb.blockSize)
		}

		packetData := blockData[packetOffset:]
		if len(packetData) < int(unsafe.Sizeof(tpacket3Hdr{})) {
			return packetsProcessed, fmt.Errorf("insufficient data for packet header")
		}

		packetHdr := (*tpacket3Hdr)(unsafe.Pointer(&packetData[0]))
		
		if packetHdr.snaplen == 0 || packetHdr.snaplen > 65536 {
			return packetsProcessed, fmt.Errorf("invalid packet snaplen: %d", packetHdr.snaplen)
		}

		payloadOffset := int(packetHdr.mac)
		if payloadOffset >= len(packetData) || payloadOffset < 0 {
			return packetsProcessed, fmt.Errorf("invalid payload offset: %d", payloadOffset)
		}

		payloadEnd := payloadOffset + int(packetHdr.snaplen)
		if payloadEnd > len(packetData) {
			return packetsProcessed, fmt.Errorf("payload extends beyond packet data: %d > %d", payloadEnd, len(packetData))
		}

		payload := packetData[payloadOffset:payloadEnd]
		timestamp := time.Unix(int64(packetHdr.sec), int64(packetHdr.nsec))

		handler(payload, timestamp)
		packetsProcessed++
		
		if packetHdr.nextOffset == 0 {
			break
		}
		
		nextOffset := int(packetHdr.nextOffset)
		if nextOffset < int(unsafe.Sizeof(tpacket3Hdr{})) {
			return packetsProcessed, fmt.Errorf("invalid next offset: %d", nextOffset)
		}
		
		packetOffset += nextOffset
	}

	return packetsProcessed, nil
}

func (rb *RingBuffer) Close() error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.closed {
		return nil
	}

	rb.closed = true

	if rb.buffer != nil {
		if err := unix.Munmap(rb.buffer); err != nil {
			rb.logger.Error("failed to unmap ring buffer", slog.String("error", err.Error()))
			return fmt.Errorf("failed to unmap ring buffer: %w", err)
		}
		rb.buffer = nil
	}

	rb.logger.Info("ring buffer closed successfully")
	return nil
}

func (rb *RingBuffer) GetUtilization() float64 {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.closed || rb.buffer == nil {
		return 0.0
	}

	userBlocks := uint32(0)
	
	for i := uint32(0); i < rb.blockCount; i++ {
		blockOffset := int(i * rb.blockSize)
		if blockOffset >= len(rb.buffer) {
			continue
		}
		
		blockData := rb.buffer[blockOffset:]
		if len(blockData) < int(unsafe.Sizeof(tpacketHdrV1{})) {
			continue
		}
		
		blockHdr := (*tpacketHdrV1)(unsafe.Pointer(&blockData[0]))
		if blockHdr.blockStatus&unix.TP_STATUS_USER != 0 {
			userBlocks++
		}
	}
	
	return float64(userBlocks) / float64(rb.blockCount)
}
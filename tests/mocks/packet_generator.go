package mocks

import (
	"time"
)

// PacketGenerator provides utility functions for generating test packets
type PacketGenerator struct{}

// NewPacketGenerator creates a new packet generator
func NewPacketGenerator() *PacketGenerator {
	return &PacketGenerator{}
}

// GenerateEthernetFrame creates a basic Ethernet frame
func (pg *PacketGenerator) GenerateEthernetFrame(srcMAC, dstMAC []byte, etherType uint16, payload []byte) []byte {
	frame := make([]byte, 14+len(payload))
	
	// Destination MAC
	copy(frame[0:6], dstMAC)
	// Source MAC
	copy(frame[6:12], srcMAC)
	// EtherType
	frame[12] = byte(etherType >> 8)
	frame[13] = byte(etherType & 0xFF)
	// Payload
	copy(frame[14:], payload)
	
	return frame
}

// GenerateIPv4Packet creates an IPv4 packet with basic headers
func (pg *PacketGenerator) GenerateIPv4Packet(srcIP, dstIP []byte, protocol uint8, payload []byte) []byte {
	ipHeader := make([]byte, 20)
	
	ipHeader[0] = 0x45 // Version 4, IHL 5
	ipHeader[1] = 0x00 // ToS
	totalLen := 20 + len(payload)
	ipHeader[2] = byte(totalLen >> 8)
	ipHeader[3] = byte(totalLen & 0xFF)
	ipHeader[4], ipHeader[5] = 0x12, 0x34 // ID
	ipHeader[6], ipHeader[7] = 0x40, 0x00 // Flags, Fragment Offset
	ipHeader[8] = 0x40 // TTL
	ipHeader[9] = protocol
	ipHeader[10], ipHeader[11] = 0x00, 0x00 // Checksum (not calculated)
	copy(ipHeader[12:16], srcIP)
	copy(ipHeader[16:20], dstIP)
	
	packet := make([]byte, len(ipHeader)+len(payload))
	copy(packet, ipHeader)
	copy(packet[20:], payload)
	
	return packet
}

// GenerateTCPSegment creates a TCP segment
func (pg *PacketGenerator) GenerateTCPSegment(srcPort, dstPort uint16, flags uint8, payload []byte) []byte {
	tcpHeader := make([]byte, 20)
	
	tcpHeader[0] = byte(srcPort >> 8)
	tcpHeader[1] = byte(srcPort & 0xFF)
	tcpHeader[2] = byte(dstPort >> 8)
	tcpHeader[3] = byte(dstPort & 0xFF)
	tcpHeader[4], tcpHeader[5], tcpHeader[6], tcpHeader[7] = 0x12, 0x34, 0x56, 0x78 // Seq
	tcpHeader[8], tcpHeader[9], tcpHeader[10], tcpHeader[11] = 0x87, 0x65, 0x43, 0x21 // Ack
	tcpHeader[12] = 0x50 // Data offset (20 bytes)
	tcpHeader[13] = flags
	tcpHeader[14], tcpHeader[15] = 0x20, 0x00 // Window
	tcpHeader[16], tcpHeader[17] = 0x12, 0x34 // Checksum (not calculated)
	tcpHeader[18], tcpHeader[19] = 0x00, 0x00 // Urgent pointer
	
	segment := make([]byte, len(tcpHeader)+len(payload))
	copy(segment, tcpHeader)
	copy(segment[20:], payload)
	
	return segment
}

// GenerateUDPDatagram creates a UDP datagram
func (pg *PacketGenerator) GenerateUDPDatagram(srcPort, dstPort uint16, payload []byte) []byte {
	udpHeader := make([]byte, 8)
	
	udpHeader[0] = byte(srcPort >> 8)
	udpHeader[1] = byte(srcPort & 0xFF)
	udpHeader[2] = byte(dstPort >> 8)
	udpHeader[3] = byte(dstPort & 0xFF)
	
	length := 8 + len(payload)
	udpHeader[4] = byte(length >> 8)
	udpHeader[5] = byte(length & 0xFF)
	udpHeader[6], udpHeader[7] = 0x12, 0x34 // Checksum (not calculated)
	
	datagram := make([]byte, len(udpHeader)+len(payload))
	copy(datagram, udpHeader)
	copy(datagram[8:], payload)
	
	return datagram
}

// GenerateCompletePacket creates a complete Ethernet + IP + Transport packet
func (pg *PacketGenerator) GenerateCompletePacket(
	srcMAC, dstMAC []byte,
	srcIP, dstIP []byte,
	protocol uint8,
	srcPort, dstPort uint16,
	payload []byte,
) []byte {
	var transportSegment []byte
	
	switch protocol {
	case 6: // TCP
		transportSegment = pg.GenerateTCPSegment(srcPort, dstPort, 0x18, payload) // PSH, ACK
	case 17: // UDP
		transportSegment = pg.GenerateUDPDatagram(srcPort, dstPort, payload)
	default:
		transportSegment = payload
	}
	
	ipPacket := pg.GenerateIPv4Packet(srcIP, dstIP, protocol, transportSegment)
	ethernetFrame := pg.GenerateEthernetFrame(srcMAC, dstMAC, 0x0800, ipPacket)
	
	return ethernetFrame
}

// MockPacketHandler is a test helper for packet processing
type MockPacketHandler struct {
	ReceivedPackets []MockPacketData
}

type MockPacketData struct {
	Data      []byte
	Timestamp time.Time
}

func NewMockPacketHandler() *MockPacketHandler {
	return &MockPacketHandler{
		ReceivedPackets: make([]MockPacketData, 0),
	}
}

func (mph *MockPacketHandler) HandlePacket(data []byte, timestamp time.Time) {
	// Make a copy of the data since it might be reused
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	
	mph.ReceivedPackets = append(mph.ReceivedPackets, MockPacketData{
		Data:      dataCopy,
		Timestamp: timestamp,
	})
}

func (mph *MockPacketHandler) GetPacketCount() int {
	return len(mph.ReceivedPackets)
}

func (mph *MockPacketHandler) Reset() {
	mph.ReceivedPackets = mph.ReceivedPackets[:0]
}
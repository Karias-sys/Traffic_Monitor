//go:build linux

package capture

import (
	"encoding/binary"
	"fmt"
	"net"
)

type EthernetHeader struct {
	DstMAC    net.HardwareAddr
	SrcMAC    net.HardwareAddr
	EtherType uint16
}

type IPv4Header struct {
	Version    uint8
	IHL        uint8
	ToS        uint8
	Length     uint16
	ID         uint16
	Flags      uint8
	FragOffset uint16
	TTL        uint8
	Protocol   uint8
	Checksum   uint16
	SrcIP      net.IP
	DstIP      net.IP
}

type IPv6Header struct {
	Version      uint8
	TrafficClass uint8
	FlowLabel    uint32
	PayloadLen   uint16
	NextHeader   uint8
	HopLimit     uint8
	SrcIP        net.IP
	DstIP        net.IP
}

type TCPHeader struct {
	SrcPort    uint16
	DstPort    uint16
	SeqNum     uint32
	AckNum     uint32
	DataOffset uint8
	Flags      uint8
	Window     uint16
	Checksum   uint16
	UrgentPtr  uint16
}

type UDPHeader struct {
	SrcPort  uint16
	DstPort  uint16
	Length   uint16
	Checksum uint16
}

type ParsedPacket struct {
	Ethernet *EthernetHeader
	IPv4     *IPv4Header
	IPv6     *IPv6Header
	TCP      *TCPHeader
	UDP      *UDPHeader
	Payload  []byte
}

const (
	EtherTypeIPv4 = 0x0800
	EtherTypeIPv6 = 0x86DD
	EtherTypeARP  = 0x0806

	IPProtoTCP  = 6
	IPProtoUDP  = 17
	IPProtoICMP = 1

	EthernetHeaderSize = 14
	IPv4HeaderMinSize  = 20
	IPv6HeaderSize     = 40
	TCPHeaderMinSize   = 20
	UDPHeaderSize      = 8
)

func ParsePacket(data []byte) (*ParsedPacket, error) {
	if len(data) < EthernetHeaderSize {
		return nil, fmt.Errorf("packet too short for Ethernet header: %d bytes", len(data))
	}

	packet := &ParsedPacket{}
	offset := 0

	ethernet, err := parseEthernetHeader(data[offset:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ethernet header: %w", err)
	}
	packet.Ethernet = ethernet
	offset += EthernetHeaderSize

	if offset >= len(data) {
		return packet, nil
	}

	switch ethernet.EtherType {
	case EtherTypeIPv4:
		if err := parseIPv4Packet(packet, data, &offset); err != nil {
			return packet, nil // Return partial packet on parse error
		}
	case EtherTypeIPv6:
		if err := parseIPv6Packet(packet, data, &offset); err != nil {
			return packet, nil // Return partial packet on parse error
		}
	default:
		// Unknown ethernet type, add remaining data as payload
		if offset < len(data) {
			packet.Payload = data[offset:]
		}
	}

	if offset < len(data) {
		packet.Payload = data[offset:]
	}

	return packet, nil
}

func parseIPv4Packet(packet *ParsedPacket, data []byte, offset *int) error {
	ipv4, ipOffset, err := parseIPv4Header(data[*offset:])
	if err != nil {
		return err
	}
	packet.IPv4 = ipv4
	*offset += ipOffset

	if *offset >= len(data) {
		return nil
	}

	switch ipv4.Protocol {
	case IPProtoTCP:
		tcp, tcpOffset, err := parseTCPHeader(data[*offset:])
		if err != nil {
			return err
		}
		packet.TCP = tcp
		*offset += tcpOffset
	case IPProtoUDP:
		udp, udpOffset, err := parseUDPHeader(data[*offset:])
		if err != nil {
			return err
		}
		packet.UDP = udp
		*offset += udpOffset
	}
	return nil
}

func parseIPv6Packet(packet *ParsedPacket, data []byte, offset *int) error {
	ipv6, ipOffset, err := parseIPv6Header(data[*offset:])
	if err != nil {
		return err
	}
	packet.IPv6 = ipv6
	*offset += ipOffset

	if *offset >= len(data) {
		return nil
	}

	switch ipv6.NextHeader {
	case IPProtoTCP:
		tcp, tcpOffset, err := parseTCPHeader(data[*offset:])
		if err != nil {
			return err
		}
		packet.TCP = tcp
		*offset += tcpOffset
	case IPProtoUDP:
		udp, udpOffset, err := parseUDPHeader(data[*offset:])
		if err != nil {
			return err
		}
		packet.UDP = udp
		*offset += udpOffset
	}
	return nil
}

func parseEthernetHeader(data []byte) (*EthernetHeader, error) {
	if len(data) < EthernetHeaderSize {
		return nil, fmt.Errorf("insufficient data for Ethernet header: need %d, got %d", EthernetHeaderSize, len(data))
	}

	eth := &EthernetHeader{
		DstMAC:    make(net.HardwareAddr, 6),
		SrcMAC:    make(net.HardwareAddr, 6),
		EtherType: binary.BigEndian.Uint16(data[12:14]),
	}

	copy(eth.DstMAC, data[0:6])
	copy(eth.SrcMAC, data[6:12])

	return eth, nil
}

func parseIPv4Header(data []byte) (*IPv4Header, int, error) {
	if len(data) < IPv4HeaderMinSize {
		return nil, 0, fmt.Errorf("insufficient data for IPv4 header: need %d, got %d", IPv4HeaderMinSize, len(data))
	}

	ihl := data[0] & 0x0F
	headerLength := int(ihl) * 4

	if headerLength < IPv4HeaderMinSize {
		return nil, 0, fmt.Errorf("invalid IPv4 header length: %d", headerLength)
	}

	if len(data) < headerLength {
		return nil, 0, fmt.Errorf("insufficient data for IPv4 header: need %d, got %d", headerLength, len(data))
	}

	ipv4 := &IPv4Header{
		Version:    (data[0] & 0xF0) >> 4,
		IHL:        ihl,
		ToS:        data[1],
		Length:     binary.BigEndian.Uint16(data[2:4]),
		ID:         binary.BigEndian.Uint16(data[4:6]),
		Flags:      (data[6] & 0xE0) >> 5,
		FragOffset: binary.BigEndian.Uint16(data[6:8]) & 0x1FFF,
		TTL:        data[8],
		Protocol:   data[9],
		Checksum:   binary.BigEndian.Uint16(data[10:12]),
		SrcIP:      net.IP(data[12:16]),
		DstIP:      net.IP(data[16:20]),
	}

	return ipv4, headerLength, nil
}

func parseIPv6Header(data []byte) (*IPv6Header, int, error) {
	if len(data) < IPv6HeaderSize {
		return nil, 0, fmt.Errorf("insufficient data for IPv6 header: need %d, got %d", IPv6HeaderSize, len(data))
	}

	ipv6 := &IPv6Header{
		Version:      (data[0] & 0xF0) >> 4,
		TrafficClass: ((data[0] & 0x0F) << 4) | ((data[1] & 0xF0) >> 4),
		FlowLabel:    (uint32(data[1]&0x0F) << 16) | (uint32(data[2]) << 8) | uint32(data[3]),
		PayloadLen:   binary.BigEndian.Uint16(data[4:6]),
		NextHeader:   data[6],
		HopLimit:     data[7],
		SrcIP:        net.IP(data[8:24]),
		DstIP:        net.IP(data[24:40]),
	}

	return ipv6, IPv6HeaderSize, nil
}

func parseTCPHeader(data []byte) (*TCPHeader, int, error) {
	if len(data) < TCPHeaderMinSize {
		return nil, 0, fmt.Errorf("insufficient data for TCP header: need %d, got %d", TCPHeaderMinSize, len(data))
	}

	dataOffset := (data[12] & 0xF0) >> 4
	headerLength := int(dataOffset) * 4

	if headerLength < TCPHeaderMinSize {
		return nil, 0, fmt.Errorf("invalid TCP header length: %d", headerLength)
	}

	if len(data) < headerLength {
		return nil, 0, fmt.Errorf("insufficient data for TCP header: need %d, got %d", headerLength, len(data))
	}

	tcp := &TCPHeader{
		SrcPort:    binary.BigEndian.Uint16(data[0:2]),
		DstPort:    binary.BigEndian.Uint16(data[2:4]),
		SeqNum:     binary.BigEndian.Uint32(data[4:8]),
		AckNum:     binary.BigEndian.Uint32(data[8:12]),
		DataOffset: dataOffset,
		Flags:      data[13],
		Window:     binary.BigEndian.Uint16(data[14:16]),
		Checksum:   binary.BigEndian.Uint16(data[16:18]),
		UrgentPtr:  binary.BigEndian.Uint16(data[18:20]),
	}

	return tcp, headerLength, nil
}

func parseUDPHeader(data []byte) (*UDPHeader, int, error) {
	if len(data) < UDPHeaderSize {
		return nil, 0, fmt.Errorf("insufficient data for UDP header: need %d, got %d", UDPHeaderSize, len(data))
	}

	udp := &UDPHeader{
		SrcPort:  binary.BigEndian.Uint16(data[0:2]),
		DstPort:  binary.BigEndian.Uint16(data[2:4]),
		Length:   binary.BigEndian.Uint16(data[4:6]),
		Checksum: binary.BigEndian.Uint16(data[6:8]),
	}

	return udp, UDPHeaderSize, nil
}

func (p *ParsedPacket) String() string {
	result := "Packet:\n"
	
	if p.Ethernet != nil {
		result += fmt.Sprintf("  Ethernet: %s -> %s (Type: 0x%04x)\n",
			p.Ethernet.SrcMAC, p.Ethernet.DstMAC, p.Ethernet.EtherType)
	}
	
	if p.IPv4 != nil {
		result += fmt.Sprintf("  IPv4: %s -> %s (Proto: %d)\n",
			p.IPv4.SrcIP, p.IPv4.DstIP, p.IPv4.Protocol)
	}
	
	if p.IPv6 != nil {
		result += fmt.Sprintf("  IPv6: %s -> %s (Next: %d)\n",
			p.IPv6.SrcIP, p.IPv6.DstIP, p.IPv6.NextHeader)
	}
	
	if p.TCP != nil {
		result += fmt.Sprintf("  TCP: %d -> %d (Flags: 0x%02x)\n",
			p.TCP.SrcPort, p.TCP.DstPort, p.TCP.Flags)
	}
	
	if p.UDP != nil {
		result += fmt.Sprintf("  UDP: %d -> %d (Len: %d)\n",
			p.UDP.SrcPort, p.UDP.DstPort, p.UDP.Length)
	}
	
	if len(p.Payload) > 0 {
		result += fmt.Sprintf("  Payload: %d bytes\n", len(p.Payload))
	}
	
	return result
}
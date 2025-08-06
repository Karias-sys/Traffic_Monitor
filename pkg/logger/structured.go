package logger

import (
	"time"
)

const (
	FieldComponent   = "component"
	FieldOperation   = "operation"
	FieldDuration    = "duration_ms"
	FieldError       = "error"
	FieldPacketCount = "packet_count"
	FieldFlowCount   = "flow_count"
	FieldInterface   = "interface"
	FieldSourceIP    = "source_ip"
	FieldDestIP      = "dest_ip"
	FieldSourcePort  = "source_port"
	FieldDestPort    = "dest_port"
	FieldProtocol    = "protocol"
	FieldBytes       = "bytes"
	FieldPackets     = "packets"
	FieldTimestamp   = "timestamp"
	FieldMemoryUsage = "memory_usage_mb"
	FieldCPUUsage    = "cpu_usage_percent"
	FieldGoroutines  = "goroutines"
	FieldConnections = "connections"
	FieldHTTPMethod  = "http_method"
	FieldHTTPStatus  = "http_status"
	FieldHTTPPath    = "http_path"
	FieldUserAgent   = "user_agent"
	FieldRemoteAddr  = "remote_addr"
	FieldRequestID   = "request_id"
)

type StructuredFields struct {
	Component   string        `json:"component,omitempty"`
	Operation   string        `json:"operation,omitempty"`
	Duration    time.Duration `json:"duration_ms,omitempty"`
	Error       string        `json:"error,omitempty"`
	PacketCount int64         `json:"packet_count,omitempty"`
	FlowCount   int64         `json:"flow_count,omitempty"`
	Interface   string        `json:"interface,omitempty"`
	SourceIP    string        `json:"source_ip,omitempty"`
	DestIP      string        `json:"dest_ip,omitempty"`
	SourcePort  int           `json:"source_port,omitempty"`
	DestPort    int           `json:"dest_port,omitempty"`
	Protocol    string        `json:"protocol,omitempty"`
	Bytes       int64         `json:"bytes,omitempty"`
	Packets     int64         `json:"packets,omitempty"`
	Timestamp   time.Time     `json:"timestamp,omitempty"`
	MemoryUsage float64       `json:"memory_usage_mb,omitempty"`
	CPUUsage    float64       `json:"cpu_usage_percent,omitempty"`
	Goroutines  int           `json:"goroutines,omitempty"`
	Connections int           `json:"connections,omitempty"`
	HTTPMethod  string        `json:"http_method,omitempty"`
	HTTPStatus  int           `json:"http_status,omitempty"`
	HTTPPath    string        `json:"http_path,omitempty"`
	UserAgent   string        `json:"user_agent,omitempty"`
	RemoteAddr  string        `json:"remote_addr,omitempty"`
	RequestID   string        `json:"request_id,omitempty"`
}

func (s StructuredFields) ToMap() map[string]any {
	fields := make(map[string]any)

	if s.Component != "" {
		fields[FieldComponent] = s.Component
	}
	if s.Operation != "" {
		fields[FieldOperation] = s.Operation
	}
	if s.Duration > 0 {
		fields[FieldDuration] = s.Duration.Milliseconds()
	}
	if s.Error != "" {
		fields[FieldError] = s.Error
	}
	if s.PacketCount > 0 {
		fields[FieldPacketCount] = s.PacketCount
	}
	if s.FlowCount > 0 {
		fields[FieldFlowCount] = s.FlowCount
	}
	if s.Interface != "" {
		fields[FieldInterface] = s.Interface
	}
	if s.SourceIP != "" {
		fields[FieldSourceIP] = s.SourceIP
	}
	if s.DestIP != "" {
		fields[FieldDestIP] = s.DestIP
	}
	if s.SourcePort > 0 {
		fields[FieldSourcePort] = s.SourcePort
	}
	if s.DestPort > 0 {
		fields[FieldDestPort] = s.DestPort
	}
	if s.Protocol != "" {
		fields[FieldProtocol] = s.Protocol
	}
	if s.Bytes > 0 {
		fields[FieldBytes] = s.Bytes
	}
	if s.Packets > 0 {
		fields[FieldPackets] = s.Packets
	}
	if !s.Timestamp.IsZero() {
		fields[FieldTimestamp] = s.Timestamp
	}
	if s.MemoryUsage > 0 {
		fields[FieldMemoryUsage] = s.MemoryUsage
	}
	if s.CPUUsage > 0 {
		fields[FieldCPUUsage] = s.CPUUsage
	}
	if s.Goroutines > 0 {
		fields[FieldGoroutines] = s.Goroutines
	}
	if s.Connections > 0 {
		fields[FieldConnections] = s.Connections
	}
	if s.HTTPMethod != "" {
		fields[FieldHTTPMethod] = s.HTTPMethod
	}
	if s.HTTPStatus > 0 {
		fields[FieldHTTPStatus] = s.HTTPStatus
	}
	if s.HTTPPath != "" {
		fields[FieldHTTPPath] = s.HTTPPath
	}
	if s.UserAgent != "" {
		fields[FieldUserAgent] = s.UserAgent
	}
	if s.RemoteAddr != "" {
		fields[FieldRemoteAddr] = s.RemoteAddr
	}
	if s.RequestID != "" {
		fields[FieldRequestID] = s.RequestID
	}

	return fields
}

func (l *Logger) WithStructured(fields StructuredFields) *Logger {
	return l.WithFields(fields.ToMap())
}

func (l *Logger) InfoWithStructured(msg string, fields StructuredFields) {
	l.WithFields(fields.ToMap()).Info(msg)
}

func (l *Logger) ErrorWithStructured(msg string, fields StructuredFields) {
	l.WithFields(fields.ToMap()).Error(msg)
}

func (l *Logger) DebugWithStructured(msg string, fields StructuredFields) {
	l.WithFields(fields.ToMap()).Debug(msg)
}

func (l *Logger) WarnWithStructured(msg string, fields StructuredFields) {
	l.WithFields(fields.ToMap()).Warn(msg)
}

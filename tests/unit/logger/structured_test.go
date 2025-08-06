package logger_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/Karias-sys/Traffic_Monitor/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructuredFields_ToMap(t *testing.T) {
	now := time.Now()
	duration := 500 * time.Millisecond

	fields := logger.StructuredFields{
		Component:   "test-component",
		Operation:   "test-operation",
		Duration:    duration,
		Error:       "test error",
		PacketCount: 1000,
		FlowCount:   50,
		Interface:   "eth0",
		SourceIP:    "192.168.1.1",
		DestIP:      "192.168.1.2",
		SourcePort:  8080,
		DestPort:    443,
		Protocol:    "TCP",
		Bytes:       1024,
		Packets:     10,
		Timestamp:   now,
		MemoryUsage: 256.5,
		CPUUsage:    15.7,
		Goroutines:  100,
		Connections: 25,
		HTTPMethod:  "GET",
		HTTPStatus:  200,
		HTTPPath:    "/api/v1/flows",
		UserAgent:   "netwatch-client/1.0",
		RemoteAddr:  "192.168.1.100",
		RequestID:   "req-12345",
	}

	m := fields.ToMap()

	assert.Equal(t, "test-component", m[logger.FieldComponent])
	assert.Equal(t, "test-operation", m[logger.FieldOperation])
	assert.Equal(t, duration.Milliseconds(), m[logger.FieldDuration])
	assert.Equal(t, "test error", m[logger.FieldError])
	assert.Equal(t, int64(1000), m[logger.FieldPacketCount])
	assert.Equal(t, int64(50), m[logger.FieldFlowCount])
	assert.Equal(t, "eth0", m[logger.FieldInterface])
	assert.Equal(t, "192.168.1.1", m[logger.FieldSourceIP])
	assert.Equal(t, "192.168.1.2", m[logger.FieldDestIP])
	assert.Equal(t, 8080, m[logger.FieldSourcePort])
	assert.Equal(t, 443, m[logger.FieldDestPort])
	assert.Equal(t, "TCP", m[logger.FieldProtocol])
	assert.Equal(t, int64(1024), m[logger.FieldBytes])
	assert.Equal(t, int64(10), m[logger.FieldPackets])
	assert.Equal(t, now, m[logger.FieldTimestamp])
	assert.Equal(t, 256.5, m[logger.FieldMemoryUsage])
	assert.Equal(t, 15.7, m[logger.FieldCPUUsage])
	assert.Equal(t, 100, m[logger.FieldGoroutines])
	assert.Equal(t, 25, m[logger.FieldConnections])
	assert.Equal(t, "GET", m[logger.FieldHTTPMethod])
	assert.Equal(t, 200, m[logger.FieldHTTPStatus])
	assert.Equal(t, "/api/v1/flows", m[logger.FieldHTTPPath])
	assert.Equal(t, "netwatch-client/1.0", m[logger.FieldUserAgent])
	assert.Equal(t, "192.168.1.100", m[logger.FieldRemoteAddr])
	assert.Equal(t, "req-12345", m[logger.FieldRequestID])
}

func TestStructuredFields_ToMap_EmptyFields(t *testing.T) {
	fields := logger.StructuredFields{
		Component: "test-component",
		// All other fields are zero values
	}

	m := fields.ToMap()

	// Should only contain the non-zero field
	assert.Len(t, m, 1)
	assert.Equal(t, "test-component", m[logger.FieldComponent])

	// Verify zero values are not included
	_, exists := m[logger.FieldOperation]
	assert.False(t, exists)
	_, exists = m[logger.FieldDuration]
	assert.False(t, exists)
	_, exists = m[logger.FieldPacketCount]
	assert.False(t, exists)
}

func TestLoggerWithStructured(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "info",
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)

	fields := logger.StructuredFields{
		Component:   "packet-capture",
		Operation:   "process_packet",
		Duration:    250 * time.Millisecond,
		PacketCount: 500,
		Interface:   "eth0",
		Protocol:    "TCP",
		Bytes:       2048,
	}

	structuredLogger := l.WithStructured(fields)
	structuredLogger.Info("packet processed")

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "packet processed", logEntry["msg"])
	assert.Equal(t, "packet-capture", logEntry["component"])
	assert.Equal(t, "process_packet", logEntry["operation"])
	assert.Equal(t, float64(250), logEntry["duration_ms"]) // JSON numbers are float64
	assert.Equal(t, float64(500), logEntry["packet_count"])
	assert.Equal(t, "eth0", logEntry["interface"])
	assert.Equal(t, "TCP", logEntry["protocol"])
	assert.Equal(t, float64(2048), logEntry["bytes"])
}

func TestInfoWithStructured(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "info",
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)

	fields := logger.StructuredFields{
		Component:  "api-server",
		Operation:  "handle_request",
		HTTPMethod: "POST",
		HTTPPath:   "/api/v1/flows",
		HTTPStatus: 201,
		RemoteAddr: "192.168.1.50",
		RequestID:  "req-abc123",
	}

	l.InfoWithStructured("request completed", fields)

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "request completed", logEntry["msg"])
	assert.Equal(t, "INFO", logEntry["level"])
	assert.Equal(t, "api-server", logEntry["component"])
	assert.Equal(t, "handle_request", logEntry["operation"])
	assert.Equal(t, "POST", logEntry["http_method"])
	assert.Equal(t, "/api/v1/flows", logEntry["http_path"])
	assert.Equal(t, float64(201), logEntry["http_status"])
	assert.Equal(t, "192.168.1.50", logEntry["remote_addr"])
	assert.Equal(t, "req-abc123", logEntry["request_id"])
}

func TestErrorWithStructured(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "info",
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)

	fields := logger.StructuredFields{
		Component: "flow-processor",
		Operation: "update_flow",
		Error:     "connection timeout",
		Interface: "eth0",
		SourceIP:  "192.168.1.10",
		DestIP:    "192.168.1.20",
		Protocol:  "TCP",
	}

	l.ErrorWithStructured("flow update failed", fields)

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "flow update failed", logEntry["msg"])
	assert.Equal(t, "ERROR", logEntry["level"])
	assert.Equal(t, "flow-processor", logEntry["component"])
	assert.Equal(t, "update_flow", logEntry["operation"])
	assert.Equal(t, "connection timeout", logEntry["error"])
	assert.Equal(t, "eth0", logEntry["interface"])
	assert.Equal(t, "192.168.1.10", logEntry["source_ip"])
	assert.Equal(t, "192.168.1.20", logEntry["dest_ip"])
	assert.Equal(t, "TCP", logEntry["protocol"])
}

func TestDebugWithStructured(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "debug",
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)

	fields := logger.StructuredFields{
		Component:   "memory-monitor",
		Operation:   "cleanup",
		MemoryUsage: 512.7,
		FlowCount:   75000,
		Goroutines:  50,
	}

	l.DebugWithStructured("memory cleanup completed", fields)

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "memory cleanup completed", logEntry["msg"])
	assert.Equal(t, "DEBUG", logEntry["level"])
	assert.Equal(t, "memory-monitor", logEntry["component"])
	assert.Equal(t, "cleanup", logEntry["operation"])
	assert.Equal(t, 512.7, logEntry["memory_usage_mb"])
	assert.Equal(t, float64(75000), logEntry["flow_count"])
	assert.Equal(t, float64(50), logEntry["goroutines"])
}

func TestWarnWithStructured(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "warn",
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)

	fields := logger.StructuredFields{
		Component:   "performance-monitor",
		Operation:   "check_cpu",
		CPUUsage:    85.3,
		Connections: 1000,
	}

	l.WarnWithStructured("high CPU usage detected", fields)

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "high CPU usage detected", logEntry["msg"])
	assert.Equal(t, "WARN", logEntry["level"])
	assert.Equal(t, "performance-monitor", logEntry["component"])
	assert.Equal(t, "check_cpu", logEntry["operation"])
	assert.Equal(t, 85.3, logEntry["cpu_usage_percent"])
	assert.Equal(t, float64(1000), logEntry["connections"])
}
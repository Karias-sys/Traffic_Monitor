package mocks

import (
	"bytes"
	"io"

	"github.com/Karias-sys/Traffic_Monitor/pkg/logger"
)

// MockLogger provides a mock implementation of the logger for testing
type MockLogger struct {
	*logger.Logger
	Buffer *bytes.Buffer
}

// NewMockLogger creates a new mock logger that writes to a buffer
func NewMockLogger(level string) (*MockLogger, error) {
	buf := &bytes.Buffer{}
	
	l, err := logger.New(logger.Config{
		Level:  level,
		Format: "json",
		Writer: buf,
	})
	if err != nil {
		return nil, err
	}
	
	return &MockLogger{
		Logger: l,
		Buffer: buf,
	}, nil
}

// GetOutput returns the current buffer contents as a string
func (m *MockLogger) GetOutput() string {
	return m.Buffer.String()
}

// Clear clears the buffer
func (m *MockLogger) Clear() {
	m.Buffer.Reset()
}

// GetLines returns the buffer contents as separate lines
func (m *MockLogger) GetLines() []string {
	output := m.GetOutput()
	if output == "" {
		return []string{}
	}
	
	lines := []string{}
	for _, line := range bytes.Split([]byte(output), []byte("\n")) {
		if len(line) > 0 {
			lines = append(lines, string(line))
		}
	}
	return lines
}

// DiscardLogger creates a logger that discards all output
type DiscardLogger struct {
	*logger.Logger
}

// NewDiscardLogger creates a new logger that discards all output
func NewDiscardLogger(level string) (*DiscardLogger, error) {
	l, err := logger.New(logger.Config{
		Level:  level,
		Format: "json",
		Writer: io.Discard,
	})
	if err != nil {
		return nil, err
	}
	
	return &DiscardLogger{Logger: l}, nil
}
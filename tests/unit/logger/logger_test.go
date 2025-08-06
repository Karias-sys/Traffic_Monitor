package logger_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Karias-sys/Traffic_Monitor/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name      string
		config    logger.Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid json logger",
			config: logger.Config{
				Level:  "info",
				Format: "json",
				Writer: &bytes.Buffer{},
			},
			wantError: false,
		},
		{
			name: "valid text logger",
			config: logger.Config{
				Level:  "debug",
				Format: "text",
				Writer: &bytes.Buffer{},
			},
			wantError: false,
		},
		{
			name: "invalid log level",
			config: logger.Config{
				Level:  "invalid",
				Format: "json",
				Writer: &bytes.Buffer{},
			},
			wantError: true,
			errorMsg:  "invalid log level",
		},
		{
			name: "invalid format",
			config: logger.Config{
				Level:  "info",
				Format: "xml",
				Writer: &bytes.Buffer{},
			},
			wantError: true,
			errorMsg:  "unsupported log format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := logger.New(tt.config)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, l)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, l)
			}
		})
	}
}

func TestLoggerLevels(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "debug",
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)
	require.NotNil(t, l)

	// Test different log levels
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Len(t, lines, 4)

	// Verify each line is valid JSON and contains the expected message
	messages := []string{"debug message", "info message", "warn message", "error message"}
	levels := []string{"DEBUG", "INFO", "WARN", "ERROR"}

	for i, line := range lines {
		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(line), &logEntry)
		require.NoError(t, err, "Line %d should be valid JSON: %s", i, line)

		assert.Equal(t, messages[i], logEntry["msg"])
		assert.Equal(t, levels[i], logEntry["level"])
		assert.NotEmpty(t, logEntry["time"])
	}
}

func TestLoggerWithComponent(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "info",
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)

	componentLogger := l.WithComponent("test-component")
	componentLogger.Info("test message")

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "test message", logEntry["msg"])
	assert.Equal(t, "test-component", logEntry["component"])
}

func TestLoggerWithError(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "info",
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)

	testErr := assert.AnError
	errorLogger := l.WithError(testErr)
	errorLogger.Error("error occurred")

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "error occurred", logEntry["msg"])
	assert.Equal(t, testErr.Error(), logEntry["error"])
}

func TestLoggerWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "info",
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)

	fields := map[string]any{
		"user_id":   123,
		"operation": "test",
		"count":     42,
	}

	fieldsLogger := l.WithFields(fields)
	fieldsLogger.Info("operation completed")

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "operation completed", logEntry["msg"])
	assert.Equal(t, float64(123), logEntry["user_id"]) // JSON numbers are float64
	assert.Equal(t, "test", logEntry["operation"])
	assert.Equal(t, float64(42), logEntry["count"])
}

func TestLoggerFormattedMessages(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "info",
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)

	l.Infof("User %s logged in with ID %d", "john", 123)

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "User john logged in with ID 123", logEntry["msg"])
}

func TestLoggerLevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := logger.New(logger.Config{
		Level:  "warn", // Only warn and error should appear
		Format: "json",
		Writer: buf,
	})
	require.NoError(t, err)

	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should only have warn and error messages
	assert.Len(t, lines, 2)

	// Verify the messages
	var logEntry1, logEntry2 map[string]interface{}
	err = json.Unmarshal([]byte(lines[0]), &logEntry1)
	require.NoError(t, err)
	err = json.Unmarshal([]byte(lines[1]), &logEntry2)
	require.NoError(t, err)

	assert.Equal(t, "warn message", logEntry1["msg"])
	assert.Equal(t, "WARN", logEntry1["level"])
	assert.Equal(t, "error message", logEntry2["msg"])
	assert.Equal(t, "ERROR", logEntry2["level"])
}

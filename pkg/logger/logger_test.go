package logger

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteLog(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	tests := []struct {
		name          string
		level         string
		format        string
		args          []any
		expectPanic   bool
		expectedLevel string
		expectedMsg   string
	}{
		{
			name:          "Info level logging",
			level:         LogLevelInfo,
			format:        "test message %s",
			args:          []any{"data"},
			expectPanic:   false,
			expectedLevel: "info",
			expectedMsg:   "test message data",
		},
		{
			name:          "Error level logging",
			level:         LogLevelError,
			format:        "error occurred: %d",
			args:          []any{404},
			expectPanic:   false,
			expectedLevel: "error",
			expectedMsg:   "error occurred: 404",
		},
		{
			name:          "Fatal level logging triggers panic",
			level:         LogLevelFatal,
			format:        "fatal error: %s",
			args:          []any{"system crash"},
			expectPanic:   true,
			expectedLevel: "fatal",
			expectedMsg:   "fatal error: system crash",
		},
		{
			name:          "Empty format string",
			level:         LogLevelInfo,
			format:        "",
			args:          []any{},
			expectPanic:   false,
			expectedLevel: "info",
			expectedMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			if tt.expectPanic {
				assert.Panics(t, func() {
					WriteLog(tt.level, tt.format, tt.args...)
				})
			} else {
				WriteLog(tt.level, tt.format, tt.args...)
			}

			logOutput := buf.String()
			assert.Contains(t, logOutput, "[level:"+tt.expectedLevel+"]")
			assert.Contains(t, logOutput, tt.expectedMsg)
			assert.Contains(t, logOutput, "logger_test.go")

			// Additional check to verify log structure
			logParts := strings.SplitN(logOutput, "] ", 3)
			assert.Equal(t, 3, len(logParts), "Log should have three parts: timestamp, file location, and message")
		})
	}
}

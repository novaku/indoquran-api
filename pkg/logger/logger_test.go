package logger

import (
	"bytes"
	"log"
	"os"
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
		args          []interface{}
		expectedPanic bool
		expectedLog   string
	}{
		{
			name:          "Info level logging",
			level:         "INFO",
			format:        "test message %s",
			args:          []interface{}{"data"},
			expectedPanic: false,
			expectedLog:   "[level:INFO] test message data",
		},
		{
			name:          "Error level logging",
			level:         "ERROR",
			format:        "error occurred: %d",
			args:          []interface{}{404},
			expectedPanic: false,
			expectedLog:   "[level:ERROR] error occurred: 404",
		},
		{
			name:          "Fatal level logging triggers panic",
			level:         LogLevelFatal,
			format:        "fatal error: %s",
			args:          []interface{}{"system crash"},
			expectedPanic: true,
			expectedLog:   "[level:FATAL] fatal error: system crash",
		},
		{
			name:          "Empty format string",
			level:         "INFO",
			format:        "",
			args:          []interface{}{},
			expectedPanic: false,
			expectedLog:   "[level:INFO] ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			if tt.expectedPanic {
				assert.Panics(t, func() {
					WriteLog(tt.level, tt.format, tt.args...)
				})
			} else {
				WriteLog(tt.level, tt.format, tt.args...)
			}

			logOutput := buf.String()
			assert.Contains(t, logOutput, tt.expectedLog)
			assert.Contains(t, logOutput, "logger_test.go")
		})
	}
}

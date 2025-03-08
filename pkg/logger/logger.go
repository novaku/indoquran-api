package logger

import (
	"fmt"
	"log"
	"runtime"
)

const (
	LogLevelDebug   = "debug"
	LogLevelInfo    = "info"
	LogLevelWarning = "warning"
	LogLevelError   = "error"
	LogLevelFatal   = "fatal" // fatal log will throw panic
)

// WriteLog logs messages with filename and line number
func WriteLog(level, format string, args ...interface{}) {
	// Get the caller information
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	// Format the log message
	msg := fmt.Sprintf(format, args...)
	log.Printf("[%s:%d] [level:%s] %s", file, line, level, msg)

	// if level == LogLevelFatal then throw panic
	if level == LogLevelFatal {
		panic(msg)
	}
}

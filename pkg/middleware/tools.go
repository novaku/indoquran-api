package middleware

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Custom response writer to capture response body
type customResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Get the full URL of the request
func getFullURL(c *gin.Context) string {
	req := c.Request
	protocol := "http"
	if req.TLS != nil {
		protocol = "https"
	}
	// Construct the full URL
	return fmt.Sprintf("%s://%s%s", protocol, req.Host, req.RequestURI)
}

// Format duration for logging
func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%d nanoseconds", d.Nanoseconds())
	} else if d < time.Second {
		if d < time.Microsecond {
			return fmt.Sprintf("%d nanoseconds", d.Nanoseconds())
		}
		if d < time.Millisecond {
			return fmt.Sprintf("%d microseconds", d.Microseconds())
		}
		return fmt.Sprintf("%d milliseconds", d.Milliseconds())
	} else if d < time.Minute {
		return fmt.Sprintf("%d seconds", int64(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%d minutes", int64(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%d hours", int64(d.Hours()))
	} else {
		return fmt.Sprintf("%d days", int64(d.Hours()/24))
	}
}

// Implement Write method to capture response body
func (cw *customResponseWriter) Write(b []byte) (int, error) {
	cw.body.Write(b)
	return cw.ResponseWriter.Write(b)
}

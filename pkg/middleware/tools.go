package middleware

import (
	"bytes"
	"fmt"
	"strings"
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
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, c.Request.RequestURI)
}

// Format duration for logging
func formatDuration(d time.Duration) string {
	s := d.String()

	// Simplify duration format by removing trailing zero minutes and seconds
	s = strings.ReplaceAll(s, "h0m0s", "h")
	s = strings.ReplaceAll(s, "m0s", "m")
	s = strings.ReplaceAll(s, "µs0", "µs")
	s = strings.ReplaceAll(s, "ms0", "ms")

	return s
}

// Implement Write method to capture response body
func (cw *customResponseWriter) Write(b []byte) (int, error) {
	cw.body.Write(b)
	return cw.ResponseWriter.Write(b)
}

package middleware

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetFullURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		setup    func() (*gin.Context, *httptest.ResponseRecorder)
		expected string
	}{
		{
			name: "HTTP URL with no query parameters",
			setup: func() (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest(http.MethodGet, "/api/test", nil)
				c.Request.Host = "example.com"
				return c, w
			},
			expected: "http://example.com/api/test",
		},
		{
			name: "HTTP URL with query parameters",
			setup: func() (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest(http.MethodGet, "/api/test?param=value", nil)
				c.Request.Host = "example.com"
				return c, w
			},
			expected: "http://example.com/api/test?param=value",
		},
		{
			name: "HTTPS URL",
			setup: func() (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest(http.MethodGet, "/api/test", nil)
				c.Request.Host = "example.com"
				c.Request.TLS = &tls.ConnectionState{}
				return c, w
			},
			expected: "https://example.com/api/test",
		},
		{
			name: "URL with port number",
			setup: func() (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest(http.MethodGet, "/api/test", nil)
				c.Request.Host = "example.com:8080"
				return c, w
			},
			expected: "http://example.com:8080/api/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := tt.setup()
			result := getFullURL(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "500 nanoseconds",
			duration: 500 * time.Nanosecond,
			expected: "500 nanoseconds",
		},
		{
			name:     "2 microseconds",
			duration: 2 * time.Microsecond,
			expected: "2 microseconds",
		},
		{
			name:     "800 milliseconds",
			duration: 800 * time.Millisecond,
			expected: "800 milliseconds",
		},
		{
			name:     "30 seconds",
			duration: 30 * time.Second,
			expected: "30 seconds",
		},
		{
			name:     "45 minutes",
			duration: 45 * time.Minute,
			expected: "45 minutes",
		},
		{
			name:     "3 hours",
			duration: 3 * time.Hour,
			expected: "3 hours",
		},
		{
			name:     "5 days",
			duration: 5 * 24 * time.Hour,
			expected: "5 days",
		},
		{
			name:     "zero duration",
			duration: 0,
			expected: "0 nanoseconds",
		},
		{
			name:     "negative duration",
			duration: -5 * time.Second,
			expected: "-5 seconds",
		},
		{
			name:     "boundary between microseconds and milliseconds",
			duration: 999 * time.Microsecond,
			expected: "999 microseconds",
		},
		{
			name:     "boundary between hours and days",
			duration: 23 * time.Hour,
			expected: "23 hours",
		},
		{
			name:     "large number of days",
			duration: 365 * 24 * time.Hour,
			expected: "365 days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"indoquran-api/internal/database"
	"indoquran-api/internal/model"
	"indoquran-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Middleware to log request and response asynchronously
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture start time
		start := time.Now()

		// Capture request body
		var requestPayload string
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			bodyBytes, _ := c.GetRawData()
			requestPayload = string(bodyBytes)
			// Reset body to be available for further handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create a buffer to capture the response body
		responseBody := new(bytes.Buffer)
		c.Writer = &customResponseWriter{ResponseWriter: c.Writer, body: responseBody}

		// Proceed with request
		c.Next()

		// Capture end time and response details
		duration := time.Since(start)

		// Create log entry
		logEntry := model.ApiTrafficLog{
			IPAddress:      c.ClientIP(),
			Endpoint:       getFullURL(c),
			Duration:       formatDuration(duration),
			HTTPMethod:     c.Request.Method,
			RequestPayload: requestPayload,
			ResponseStatus: c.Writer.Status(),
			// ResponseBody:   responseBody.String(),
			ResponseBody: "",
			UserAgent:    c.Request.UserAgent(),
			Referer:      c.Request.Referer(),
		}

		// Log entry asynchronously
		go func() {
			if err := database.GetDB().Create(&logEntry).Error; err != nil {
				// Handle logging error (e.g., log to a file or monitoring system)
				// For simplicity, we're just printing it here
				logger.WriteLog(logger.LogLevelError, "Error logging API traffic: %#v", err)
			}
		}()
	}
}

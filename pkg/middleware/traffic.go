package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type APITrafficLog struct {
	gorm.Model
	IPAddress      string        `gorm:"column:ip_address"`
	Endpoint       string        `gorm:"column:endpoint"`
	Duration       time.Duration `gorm:"column:duration"`
	HTTPMethod     string        `gorm:"column:http_method"`
	RequestPayload string        `gorm:"column:request_payload"`
	ResponseStatus int           `gorm:"column:response_status"`
	ResponseBody   string        `gorm:"column:response_body"`
	UserAgent      string        `gorm:"column:user_agent"`
	Referer        string        `gorm:"column:referer"`
}

func (APITrafficLog) TableName() string {
	return "api_traffic_logs"
}

// LoggingMiddleware handles logging of API traffic
func LoggingMiddleware(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start)

		// Get request body if it exists
		var requestBody string
		if c.Request.Body != nil {
			if body, err := c.GetRawData(); err == nil {
				requestBody = string(body)
			}
		}

		// Create traffic log entry
		trafficLog := APITrafficLog{
			IPAddress:      c.ClientIP(),
			Endpoint:       c.Request.URL.String(),
			Duration:       duration,
			HTTPMethod:     c.Request.Method,
			RequestPayload: requestBody,
			ResponseStatus: c.Writer.Status(),
			ResponseBody:   "", // You might want to capture response body if needed
			UserAgent:      c.Request.UserAgent(),
			Referer:        c.Request.Referer(),
		}

		// Log to database asynchronously
		go func(log APITrafficLog) {
			if err := db.Create(&log).Error; err != nil {
				logger.Error("Error logging API traffic", zap.Error(err))
			}
		}(trafficLog)
	}
}

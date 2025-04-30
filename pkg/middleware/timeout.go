package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware applies a timeout to the request
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Create a new request with the context
		c.Request = c.Request.WithContext(ctx)

		// Create a channel to signal completion
		done := make(chan struct{})
		defer close(done)

		// Run the next handler in a goroutine
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// Handle panic
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "Internal server error",
					})
				}
				done <- struct{}{}
			}()
			c.Next()
		}()

		select {
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "Request timed out, please try again later.",
			})
		case <-done:
			// Request completed successfully
		}
	}
}

package middleware

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gopkg.in/yaml.v3"
)

// Config structure to hold endpoint-specific rate limits.
type Config struct {
	RateLimits map[string]RateLimitConfig `yaml:"rate_limits"`
}

// RateLimitConfig holds the rate and burst for each endpoint.
type RateLimitConfig struct {
	Rate  float64 `yaml:"rate"`
	Burst int     `yaml:"burst"`
}

// RateLimiter to track rate limits per client and endpoint.
type RateLimiter struct {
	limiters sync.Map
	config   Config
}

// LoadConfig reads the YAML config file and parses it.
func LoadConfig(path string) (Config, error) {
	var config Config
	file, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(file, &config)
	return config, err
}

// NewRateLimiter initializes the rate limiter using the loaded config.
func NewRateLimiter(config Config) *RateLimiter {
	return &RateLimiter{config: config}
}

// GetLimiter returns a rate limiter for a specific client IP and API path.
func (r *RateLimiter) GetLimiter(clientIP, path string) *rate.Limiter {
	key := fmt.Sprintf("%s:%s", clientIP, path)
	limiter, exists := r.limiters.Load(key)

	if !exists {
		// Fetch rate and burst from the configuration or use a default if not defined.
		rateLimitConfig, ok := r.config.RateLimits[path]
		if !ok {
			// Default rate limit (if not configured in config.yaml)
			rateLimitConfig = RateLimitConfig{Rate: 5, Burst: 10} // Default: 5 requests/second, burst of 10
		}

		limiter = rate.NewLimiter(rate.Limit(rateLimitConfig.Rate), rateLimitConfig.Burst)
		r.limiters.Store(key, limiter)
	}

	return limiter.(*rate.Limiter)
}

// RateLimitMiddleware applies rate limiting based on config for each endpoint.
func (r *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		path := c.FullPath() // Get the current API route path

		// Handle dynamic route for /api/v1/surat/:id and /api/v1/ayat/:id by treating it as /api/v1/surat/* and /api/v1/ayat/:id
		if c.Param("id") != "" && c.FullPath() == "/api/v1/surat/:id" {
			path = "/api/v1/surat/*"
		}
		if c.Param("id") != "" && c.FullPath() == "/api/v1/ayat/:id" {
			path = "/api/v1/ayat/*"
		}

		limiter := r.GetLimiter(clientIP, path)

		if !limiter.Allow() {
			// Too many requests for this endpoint
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Rate limit exceeded for endpoint %s. Please try again later.", path),
			})
			c.Abort()
			return
		}

		// Proceed with the request
		c.Next()
	}
}

package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type RateLimiter struct {
	redisClient *redis.Client
	mu          sync.RWMutex
	config      map[string]RateLimit
}

type RateLimit struct {
	Requests int
	Period   time.Duration
}

func NewRateLimiter(redisClient *redis.Client, config map[string]RateLimit) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		config:      config,
	}
}

// RateLimitMiddleware creates a middleware for rate limiting
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		endpoint := c.Request.URL.Path

		// Get rate limit config for this endpoint
		limit, exists := limiter.config[endpoint]
		if !exists {
			// Use default rate limit if no specific config exists
			limit = RateLimit{
				Requests: 100,       // Default to 100 requests
				Period:   time.Hour, // Default to per hour
			}
		}

		// Create a unique key for this IP and endpoint
		key := fmt.Sprintf("rate_limit:%s:%s", ip, endpoint)

		// Use Lua script for atomic operations
		script := `
			local current = redis.call('GET', KEYS[1])
			if current == false then
				redis.call('SETEX', KEYS[1], ARGV[2], 1)
				return 1
			end
			if tonumber(current) >= tonumber(ARGV[1]) then
				return -1
			end
			redis.call('INCR', KEYS[1])
			return tonumber(current) + 1
		`

		// Execute the Lua script
		result, err := limiter.redisClient.Eval(script, []string{key}, limit.Requests, limit.Period.Seconds()).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Error checking rate limit",
			})
			c.Abort()
			return
		}

		// Check if rate limit is exceeded
		if result == -1 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (r *RateLimiter) Status(ip, endpoint string) (int, error) {
	key := fmt.Sprintf("rate_limit:%s:%s", ip, endpoint)
	val, err := r.redisClient.Get(key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

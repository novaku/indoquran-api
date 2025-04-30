package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2" // Ensure this package is installed using `go get gopkg.in/yaml.v2`
)

// RateLimitConfig represents the rate limit configuration for a specific route.
type RateLimitConfig struct {
	Rate  float64 `yaml:"rate"`
	Burst int     `yaml:"burst"`
}

// Config represents the overall rate limit configuration.
type Config struct {
	RateLimits map[string]RateLimitConfig `yaml:"rate_limits"`
}

// LoadConfig loads the rate limit configuration from a YAML file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "rate_limit.yml")

	configContent := []byte(`
rate_limits:
  "/api/v1/search":
    rate: 2
    burst: 4
  "/api/v1/surat/*":
    rate: 5
    burst: 10
`)

	err := os.WriteFile(configPath, configContent, 0644)
	assert.NoError(t, err)

	// Test successful config loading
	config, err := LoadConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, 2.0, config.RateLimits["/api/v1/search"].Rate)
	assert.Equal(t, 4, config.RateLimits["/api/v1/search"].Burst)
	assert.Equal(t, 5.0, config.RateLimits["/api/v1/surat/*"].Rate)
	assert.Equal(t, 10, config.RateLimits["/api/v1/surat/*"].Burst)

	// Test loading non-existent config
	_, err = LoadConfig("non_existent.yml")
	assert.Error(t, err)

	// Test loading invalid config
	err = os.WriteFile(configPath, []byte("invalid:yaml:content"), 0644)
	assert.NoError(t, err)
	_, err = LoadConfig(configPath)
	assert.Error(t, err)
}

func TestRateLimiter(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	config := map[string]RateLimit{
		"/api/test": {
			Requests: 2,
			Period:   time.Second,
		},
	}

	rateLimiter := NewRateLimiter(rdb, config)

	// Test rate limit middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitMiddleware(rateLimiter))
	router.GET("/api/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	// Test successful requests
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/test", nil)
	req1.RemoteAddr = "127.0.0.1:1234"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/test", nil)
	req2.RemoteAddr = "127.0.0.1:1234"
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Test rate limit exceeded
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/test", nil)
	req3.RemoteAddr = "127.0.0.1:1234"
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusTooManyRequests, w3.Code)

	// Test status
	count, err := rateLimiter.Status("127.0.0.1", "/api/test")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := gin.New()

	server := NewServer(g)
	assert.NotNil(t, server)
	assert.Implements(t, (*IServer)(nil), server)
}

func TestRouteSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test rate limit config
	tmpDir := t.TempDir()
	configPath := tmpDir + "/rate_limit.yml"
	configContent := []byte(`
rate_limits:
  "/api/v1/search":
    rate: 10
    burst: 20
`)
	err := os.WriteFile(configPath, configContent, 0644)
	assert.NoError(t, err)

	// Set up viper configuration
	viper.Set("RATE_LIMIT_FILE", configPath)
	viper.Set("API_PORT", "8080")

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Root endpoint",
			path:           "/api/v1/",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"message": "Welcome to indoquran.web.id API v1.0",
			},
		},
		{
			name:           "CSP Report endpoint",
			path:           "/csp-report",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router for each test
			g := gin.New()
			server := NewServer(g)

			// Start the router in a goroutine
			go func() {
				server.RunRouter()
			}()

			// Give the server a moment to start
			time.Sleep(100 * time.Millisecond)

			// Create test request
			w := httptest.NewRecorder()
			var req *http.Request
			if tt.path == "/csp-report" {
				req = httptest.NewRequest(http.MethodPost, tt.path, nil)
			} else {
				req = httptest.NewRequest(http.MethodGet, tt.path, nil)
			}

			// Serve the request
			g.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// For endpoints that return JSON
			if tt.path == "/api/v1/" {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestGracefulShutdown(t *testing.T) {
	// Create a test server
	srv := &http.Server{
		Addr: ":8080",
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test graceful shutdown
	go func() {
		time.Sleep(100 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	}()

	// Wait for shutdown
	err := srv.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestMiddlewareConfiguration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := gin.New()
	server := NewServer(g).(*Server)

	// Setup test rate limit config
	tmpDir := t.TempDir()
	configPath := tmpDir + "/rate_limit.yml"
	configContent := []byte(`
rate_limits:
  "/test":
    rate: 1
    burst: 2
`)
	err := os.WriteFile(configPath, configContent, 0644)
	assert.NoError(t, err)
	viper.Set("RATE_LIMIT_FILE", configPath)

	// Test middleware setup
	server.RunRouter()

	// Verify CORS middleware
	engine := server.g
	handlers := engine.Handlers
	assert.Greater(t, len(handlers), 0)

	// Test rate limiting
	w := httptest.NewRecorder()
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		engine.ServeHTTP(w, req)
		if i < 2 {
			assert.NotEqual(t, http.StatusTooManyRequests, w.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code)
		}
	}

	// Test timeout middleware
	slowHandler := func(c *gin.Context) {
		time.Sleep(2 * time.Minute)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
	engine.GET("/slow", slowHandler)

	w = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusGatewayTimeout, w.Code)

	// Test CSP middleware
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	engine.ServeHTTP(w, req)
	assert.Contains(t, w.Header().Get("Content-Security-Policy"), "default-src 'self'")
}

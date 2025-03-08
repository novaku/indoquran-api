package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTimeoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		timeout        time.Duration
		sleepDuration  time.Duration
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Request completes within timeout",
			timeout:        100 * time.Millisecond,
			sleepDuration:  50 * time.Millisecond,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}`,
		},
		{
			name:           "Request exceeds timeout",
			timeout:        50 * time.Millisecond,
			sleepDuration:  100 * time.Millisecond,
			expectedStatus: http.StatusGatewayTimeout,
			expectedBody:   `{"error":"Request timed out, please try again later."}`,
		},
		{
			name:           "Zero timeout",
			timeout:        0,
			sleepDuration:  10 * time.Millisecond,
			expectedStatus: http.StatusGatewayTimeout,
			expectedBody:   `{"error":"Request timed out, please try again later."}`,
		},
		{
			name:           "Immediate response",
			timeout:        1 * time.Second,
			sleepDuration:  0,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(TimeoutMiddleware(tt.timeout))

			router.GET("/test", func(c *gin.Context) {
				time.Sleep(tt.sleepDuration)
				c.JSON(http.StatusOK, gin.H{"status": "success"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestTimeoutMiddleware_ConcurrentRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMiddleware(50 * time.Millisecond))

	router.GET("/test", func(c *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	var wg sync.WaitGroup
	requests := 5

	results := make([]int, requests)
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)
			results[index] = w.Code
		}(i)
	}

	wg.Wait()

	for _, code := range results {
		assert.Equal(t, http.StatusGatewayTimeout, code)
	}
}

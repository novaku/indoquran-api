package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestContentSecurityPolicy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedHeader string
	}{
		{
			name:           "GET request with CSP header",
			method:         "GET",
			path:           "/test",
			expectedHeader: "default-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; report-uri /csp-report;",
		},
		{
			name:           "POST request with CSP header",
			method:         "POST",
			path:           "/test",
			expectedHeader: "default-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; report-uri /csp-report;",
		},
		{
			name:           "OPTIONS request with CSP header",
			method:         "OPTIONS",
			path:           "/test",
			expectedHeader: "default-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; report-uri /csp-report;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(ContentSecurityPolicy())

			router.Any("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tt.expectedHeader, w.Header().Get("Content-Security-Policy"))
		})
	}
}

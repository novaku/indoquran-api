package middleware

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"indoquran-api/pkg/database"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Custom response writer that implements gin.ResponseWriter
type testResponseWriter struct {
	*httptest.ResponseRecorder
	closeNotifyCh chan bool
	size          int
	status        int
}

func (w *testResponseWriter) Size() int {
	return w.size
}

func (w *testResponseWriter) Status() int {
	return w.status
}

func (w *testResponseWriter) WriteString(s string) (int, error) {
	w.size += len(s)
	return w.ResponseRecorder.WriteString(s)
}

func (w *testResponseWriter) Written() bool {
	return w.ResponseRecorder.Code != 0
}

func (w *testResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseRecorder.WriteHeader(code)
}

func (w *testResponseWriter) WriteHeaderNow() {}

func (w *testResponseWriter) Pusher() http.Pusher {
	return nil
}

func (w *testResponseWriter) CloseNotify() <-chan bool {
	if w.closeNotifyCh == nil {
		w.closeNotifyCh = make(chan bool, 1)
	}
	return w.closeNotifyCh
}

func (w *testResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func TestLoggingMiddleware(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create a mock database
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	// Convert sql.DB to gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Save original DB and restore after test
	originalDB := database.GetDB()
	database.SetDB(gormDB)
	defer database.SetDB(originalDB)

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
		setupMock      func(sqlmock.Sqlmock)
	}{
		{
			name:           "GET request",
			method:         "GET",
			path:           "/api/test",
			body:           "",
			expectedStatus: http.StatusOK,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `api_traffic_logs`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:           "POST request with body",
			method:         "POST",
			path:           "/api/test",
			body:           `{"key":"value"}`,
			expectedStatus: http.StatusOK,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `api_traffic_logs`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router with middleware
			router := gin.New()
			logger, err := zap.NewProduction()
			assert.NoError(t, err)
			defer logger.Sync()
			router.Use(LoggingMiddleware(gormDB, logger))

			// Add test endpoint
			router.Any("/api/test", func(c *gin.Context) {
				c.JSON(tt.expectedStatus, gin.H{"status": "success"})
			})

			// Setup mock expectations
			tt.setupMock(mock)

			// Create request
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			// Add test headers
			req.Header.Set("User-Agent", "test-agent")
			req.Header.Set("Referer", "test-referer")
			req.Header.Set("X-Forwarded-For", "192.168.1.1")

			// Create custom response writer
			w := httptest.NewRecorder()
			gw := &testResponseWriter{
				ResponseRecorder: w,
				closeNotifyCh:    make(chan bool, 1),
			}

			// Perform request
			router.ServeHTTP(gw, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify response
			var response map[string]string
			err = json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, "success", response["status"])

			// Wait a bit for the async logging
			time.Sleep(100 * time.Millisecond)

			// Verify all mock expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCustomResponseWriter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		writeData      string
		expectedStatus int
	}{
		{
			name:           "Write success response",
			writeData:      "test response",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Write error response",
			writeData:      "error response",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new response recorder
			w := httptest.NewRecorder()

			// Create custom response writer
			body := new(bytes.Buffer)
			customWriter := &customResponseWriter{
				ResponseWriter: &testResponseWriter{
					ResponseRecorder: w,
					closeNotifyCh:    make(chan bool, 1),
				},
				body: body,
			}

			// Write status
			customWriter.WriteHeader(tt.expectedStatus)
			assert.Equal(t, tt.expectedStatus, customWriter.Status())

			// Write data
			n, err := customWriter.Write([]byte(tt.writeData))
			assert.NoError(t, err)
			assert.Equal(t, len(tt.writeData), n)

			// Verify written data
			assert.Equal(t, tt.writeData, body.String())
		})
	}
}

func TestGetFullURLWithQueryParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		path     string
		query    string
		expected string
	}{
		{
			name:     "URL with query parameters",
			path:     "/api/test",
			query:    "param1=value1&param2=value2",
			expected: "http://example.com/api/test?param1=value1&param2=value2",
		},
		{
			name:     "URL without query parameters",
			path:     "/api/test",
			query:    "",
			expected: "http://example.com/api/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			url := tt.path
			if tt.query != "" {
				url += "?" + tt.query
			}

			c.Request = httptest.NewRequest(http.MethodGet, url, nil)
			c.Request.Host = "example.com"

			result := getFullURL(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

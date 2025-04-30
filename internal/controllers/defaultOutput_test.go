package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockResponseFormatter implements ResponseFormatterInterface
type MockResponseFormatter struct {
	mock.Mock
}

func (m *MockResponseFormatter) FormatResponse(data interface{}, err error) (int, Response) {
	args := m.Called(data, err)
	return args.Get(0).(int), args.Get(1).(Response)
}

// MockResponseWriter implements ResponseWriterInterface
type MockResponseWriter struct {
	mock.Mock
}

func (m *MockResponseWriter) WriteResponse(c *gin.Context, data interface{}, err error) {
	m.Called(c, data, err)
}

func TestWriteResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		data          any
		err           error
		expectedCode  int
		expectedData  any
		expectedError string
	}{
		{
			name:          "success with data",
			data:          map[string]any{"key": "value"},
			err:           nil,
			expectedCode:  http.StatusOK,
			expectedData:  map[string]any{"key": "value"},
			expectedError: "",
		},
		{
			name:          "error case",
			data:          nil,
			err:           errors.New("test error"),
			expectedCode:  http.StatusBadRequest,
			expectedData:  nil,
			expectedError: "test error",
		},
		{
			name:          "success with nil data",
			data:          nil,
			err:           nil,
			expectedCode:  http.StatusOK,
			expectedData:  nil,
			expectedError: "",
		},
		{
			name:          "success with array data",
			data:          []any{"item1", "item2"},
			err:           nil,
			expectedCode:  http.StatusOK,
			expectedData:  []any{"item1", "item2"},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockFormatter := new(MockResponseFormatter)
			mockWriter := new(MockResponseWriter)

			// Setup test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Setup formatter mock
			response := Response{
				Version: "1.0",
				Data:    tt.expectedData,
				Error:   tt.expectedError,
			}
			mockFormatter.On("FormatResponse", tt.data, tt.err).Return(tt.expectedCode, response)

			// Setup writer mock
			mockWriter.On("WriteResponse", c, tt.data, tt.err).Return()

			// Create response writer with mock formatter
			writer := NewDefaultResponseWriter(mockFormatter)

			// Execute test
			writer.WriteResponse(c, tt.data, tt.err)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			var actualResponse Response
			err := json.NewDecoder(w.Body).Decode(&actualResponse)
			assert.NoError(t, err)

			assert.Equal(t, "1.0", actualResponse.Version)
			assert.Equal(t, tt.expectedData, actualResponse.Data)
			assert.Equal(t, tt.expectedError, actualResponse.Error)

			// Verify mock expectations
			mockFormatter.AssertExpectations(t)
		})
	}
}

func TestDefaultResponseWriter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		data          any
		err           error
		expectedCode  int
		expectedData  any
		expectedError string
	}{
		{
			name:          "success with data",
			data:          map[string]any{"key": "value"},
			err:           nil,
			expectedCode:  http.StatusOK,
			expectedData:  map[string]any{"key": "value"},
			expectedError: "",
		},
		{
			name:          "error case",
			data:          nil,
			err:           errors.New("test error"),
			expectedCode:  http.StatusBadRequest,
			expectedData:  nil,
			expectedError: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Get default writer
			writer := GetDefaultResponseWriter()

			// Execute test
			writer.WriteResponse(c, tt.data, tt.err)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			var response Response
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			assert.Equal(t, "1.0", response.Version)
			assert.Equal(t, tt.expectedData, response.Data)
			assert.Equal(t, tt.expectedError, response.Error)
		})
	}
}

package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		data          interface{}
		err           error
		expectedCode  int
		expectedData  interface{}
		expectedError string
	}{
		{
			name:          "success with data",
			data:          map[string]string{"key": "value"},
			err:           nil,
			expectedCode:  http.StatusOK,
			expectedData:  map[string]string{"key": "value"},
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
			data:          []string{"item1", "item2"},
			err:           nil,
			expectedCode:  http.StatusOK,
			expectedData:  []string{"item1", "item2"},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			WriteResponse(c, tt.data, tt.err)

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

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

type MockDetail struct {
	mock.Mock
}

func (m *MockDetail) GetAyat(id string) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func TestDetailAyat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		ayatID       string
		mockResponse interface{}
		mockError    error
		expectedCode int
	}{
		{
			name:         "success get ayat",
			ayatID:       "1",
			mockResponse: map[string]string{"text": "sample ayat"},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid ayat id",
			ayatID:       "invalid",
			mockResponse: nil,
			mockError:    errors.New("ayat not found"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "empty ayat id",
			ayatID:       "",
			mockResponse: nil,
			mockError:    errors.New("invalid ayat id"),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDetail := new(MockDetail)
			mockDetail.On("GetAyat", tt.ayatID).Return(tt.mockResponse, tt.mockError)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.ayatID}}

			DetailAyat(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response Response
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			if tt.mockError != nil {
				assert.Equal(t, tt.mockError.Error(), response.Error)
				assert.Nil(t, response.Data)
			} else {
				assert.Equal(t, tt.mockResponse, response.Data)
				assert.Empty(t, response.Error)
			}

			mockDetail.AssertExpectations(t)
		})
	}
}

package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSurat struct {
	mock.Mock
}

func (m *MockSurat) GetSuratList(id string) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func TestListSurat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		suratID      string
		mockResponse interface{}
		mockError    error
		expectedCode int
	}{
		{
			name:         "success get all surats",
			suratID:      "",
			mockResponse: []map[string]string{{"id": "1", "name": "Al-Fatihah"}},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "success get specific surat",
			suratID:      "1",
			mockResponse: map[string]string{"id": "1", "name": "Al-Fatihah"},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid surat id",
			suratID:      "invalid",
			mockResponse: nil,
			mockError:    errors.New("surat not found"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "server error",
			suratID:      "",
			mockResponse: nil,
			mockError:    errors.New("internal server error"),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSurat := new(MockSurat)
			mockSurat.On("GetSuratList", tt.suratID).Return(tt.mockResponse, tt.mockError)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.suratID != "" {
				c.Request, _ = http.NewRequest(http.MethodGet, "/?surat="+tt.suratID, nil)
			} else {
				c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
			}

			ListSurat(c)

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

			mockSurat.AssertExpectations(t)
		})
	}
}

type MockAyat struct {
	mock.Mock
}

func (m *MockAyat) GetAyatList(suratID string, page int, pageSize int) (interface{}, error) {
	args := m.Called(suratID, page, pageSize)
	return args.Get(0), args.Error(1)
}

func TestListAyatInSurat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		suratID      string
		page         string
		pageSize     string
		mockResponse interface{}
		mockError    error
		expectedCode int
	}{
		{
			name:         "success get ayat list with default pagination",
			suratID:      "1",
			page:         "",
			pageSize:     "",
			mockResponse: []map[string]string{{"number": "1", "text": "bismillah"}},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "success get ayat list with custom pagination",
			suratID:      "1",
			page:         "2",
			pageSize:     "5",
			mockResponse: []map[string]string{{"number": "6", "text": "sample"}},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid surat id",
			suratID:      "999",
			page:         "1",
			pageSize:     "10",
			mockResponse: nil,
			mockError:    errors.New("surat not found"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "negative page number",
			suratID:      "1",
			page:         "-1",
			pageSize:     "10",
			mockResponse: nil,
			mockError:    errors.New("invalid page number"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "zero page size",
			suratID:      "1",
			page:         "1",
			pageSize:     "0",
			mockResponse: nil,
			mockError:    errors.New("invalid page size"),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAyat := new(MockAyat)

			expectedPage := 1
			if tt.page != "" {
				expectedPage, _ = strconv.Atoi(tt.page)
			}

			expectedPageSize := 10
			if tt.pageSize != "" {
				expectedPageSize, _ = strconv.Atoi(tt.pageSize)
			}

			mockAyat.On("GetAyatList", tt.suratID, expectedPage, expectedPageSize).Return(tt.mockResponse, tt.mockError)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			url := "/surat/" + tt.suratID
			if tt.page != "" {
				url += "?p=" + tt.page
			}
			if tt.pageSize != "" {
				if tt.page != "" {
					url += "&"
				} else {
					url += "?"
				}
				url += "n=" + tt.pageSize
			}

			c.Request, _ = http.NewRequest(http.MethodGet, url, nil)
			c.Params = []gin.Param{{Key: "id", Value: tt.suratID}}

			ListAyatInSurat(c)

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

			mockAyat.AssertExpectations(t)
		})
	}
}

package controllers

import (
	"encoding/json"
	"fmt"
	"indoquran-api/internal/model"
	"indoquran-api/internal/services/list"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSuratService implements list.SuratService
type MockSuratService struct {
	mock.Mock
}

func (m *MockSuratService) GetSuratList(suratID string) ([]model.IdMuntakhab, error) {
	args := m.Called(suratID)
	return args.Get(0).([]model.IdMuntakhab), args.Error(1)
}

// MockAyatService implements list.AyatService
type MockAyatService struct {
	mock.Mock
}

func (m *MockAyatService) GetAyatList(suratID string, page, pageSize int) ([]*model.AyatDetail, error) {
	args := m.Called(suratID, page, pageSize)
	return args.Get(0).([]*model.AyatDetail), args.Error(1)
}

var (
	originalNewSurat = list.NewSurat
	originalNewAyat  = list.NewAyat
)

func TestList_GetSuratList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSurat := new(MockSuratService)
	listController := &ListController{suratService: mockSurat}

	tests := []struct {
		name           string
		suratID        string
		mockSetup      func(*MockSuratService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:    "Success - Get all surats",
			suratID: "",
			mockSetup: func(ms *MockSuratService) {
				ms.On("GetSuratList", "").Return([]model.IdMuntakhab{
					{Index: 1, Surat: 1, Ayat: 7, Text: "Al-Fatihah"},
					{Index: 2, Surat: 2, Ayat: 286, Text: "Al-Baqarah"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"status":  "success",
				"message": "Get surat list successfully",
				"data": []model.IdMuntakhab{
					{Index: 1, Surat: 1, Ayat: 7, Text: "Al-Fatihah"},
					{Index: 2, Surat: 2, Ayat: 286, Text: "Al-Baqarah"},
				},
			},
		},
		{
			name:    "Success - Get specific surat",
			suratID: "1",
			mockSetup: func(ms *MockSuratService) {
				ms.On("GetSuratList", "1").Return([]model.IdMuntakhab{
					{Index: 1, Surat: 1, Ayat: 7, Text: "Al-Fatihah"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"status":  "success",
				"message": "Get surat list successfully",
				"data": []model.IdMuntakhab{
					{Index: 1, Surat: 1, Ayat: 7, Text: "Al-Fatihah"},
				},
			},
		},
		{
			name:    "Error - Invalid surat ID",
			suratID: "invalid",
			mockSetup: func(ms *MockSuratService) {
				ms.On("GetSuratList", "invalid").Return([]model.IdMuntakhab{}, fmt.Errorf("invalid surat ID"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"status":  "error",
				"message": "invalid surat ID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mockSurat)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.suratID != "" {
				c.Params = append(c.Params, gin.Param{Key: "suratId", Value: tt.suratID})
			}

			listController.GetSuratList(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			expectedBody, _ := json.Marshal(tt.expectedBody)
			actualBody, _ := json.Marshal(response)
			assert.JSONEq(t, string(expectedBody), string(actualBody))

			mockSurat.AssertExpectations(t)
		})
	}
}

func TestList_GetAyatList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAyat := new(MockAyatService)
	listController := &ListController{ayatService: mockAyat}

	tests := []struct {
		name           string
		suratID        string
		page           string
		pageSize       string
		mockSetup      func(*MockAyatService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "Success - Get ayat list",
			suratID:  "1",
			page:     "1",
			pageSize: "10",
			mockSetup: func(ma *MockAyatService) {
				ma.On("GetAyatList", "1", 1, 10).Return([]model.QuranAyat{
					{
						AyatKey:    "0001001",
						AyatNumber: 1,
						Surat:      1,
						Ayat:       1,
						Text:       "بِسْمِ ٱللَّهِ",
						Simple:     stringPtr("Dengan nama Allah"),
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"status":  "success",
				"message": "Get ayat list successfully",
				"data": []model.QuranAyat{
					{
						AyatKey:    "0001001",
						AyatNumber: 1,
						Surat:      1,
						Ayat:       1,
						Text:       "بِسْمِ ٱللَّهِ",
						Simple:     stringPtr("Dengan nama Allah"),
					},
				},
			},
		},
		{
			name:     "Error - Invalid surat ID",
			suratID:  "invalid",
			page:     "1",
			pageSize: "10",
			mockSetup: func(ma *MockAyatService) {
				ma.On("GetAyatList", "invalid", 1, 10).Return([]model.QuranAyat{}, fmt.Errorf("invalid surat ID"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"status":  "error",
				"message": "invalid surat ID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mockAyat)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = append(c.Params, gin.Param{Key: "suratId", Value: tt.suratID})
			c.Request = httptest.NewRequest("GET", "/?page="+tt.page+"&pageSize="+tt.pageSize, nil)

			listController.GetAyatList(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			expectedBody, _ := json.Marshal(tt.expectedBody)
			actualBody, _ := json.Marshal(response)
			assert.JSONEq(t, string(expectedBody), string(actualBody))

			mockAyat.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

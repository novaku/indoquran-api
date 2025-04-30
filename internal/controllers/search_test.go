package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"indoquran-api/internal/model"
	"indoquran-api/internal/services/search"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearchService implements SearchServiceInterface
type MockSearchService struct {
	mock.Mock
}

func (m *MockSearchService) FullTextSearch(query string, juz, surat, page, rowsPerPage int) (*search.SearchResult, error) {
	args := m.Called(query, juz, surat, page, rowsPerPage)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*search.SearchResult), args.Error(1)
}

// MockSearchRequestValidator implements SearchRequestValidatorInterface
type MockSearchRequestValidator struct {
	mock.Mock
}

func (m *MockSearchRequestValidator) ValidateRequest(c *gin.Context) (*SearchRequest, error) {
	args := m.Called(c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SearchRequest), args.Error(1)
}

// MockSearchResponseBuilder implements SearchResponseBuilderInterface
type MockSearchResponseBuilder struct {
	mock.Mock
}

func (m *MockSearchResponseBuilder) BuildResponse(result *search.SearchResult, request *SearchRequest) *ResultJsonFormat {
	args := m.Called(result, request)
	return args.Get(0).(*ResultJsonFormat)
}

func TestSearchHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		query          string
		page           string
		juz            string
		surat          string
		rowsPerPage    string
		mockResults    *search.SearchResult
		mockError      error
		expectedCode   int
		expectedResult *ResultJsonFormat
	}{
		{
			name:        "successful search",
			query:       "test",
			page:        "1",
			juz:         "0",
			surat:       "0",
			rowsPerPage: "10",
			mockResults: &search.SearchResult{
				Results: []*model.AyatDetail{
					{
						ID:         1,
						Juz:        1,
						Surat:      1,
						Ayat:       1,
						TextIndo:   "test",
						TextArabic: "test",
					},
				},
				TotalCount: 20,
				Count: []*model.CountResult{
					{
						Type:       "surat",
						Identifier: 1,
						Count:      1,
					},
				},
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedResult: &ResultJsonFormat{
				Aggregate: []*model.CountResult{
					{
						Type:       "surat",
						Identifier: 1,
						Count:      1,
					},
				},
				Pagination: &Pagination{
					CurrentPage: 1,
					RowsPerPage: 10,
					TotalPages:  2,
					TotalRows:   20,
				},
				Results: []*model.AyatDetail{
					{
						ID:         1,
						Juz:        1,
						Surat:      1,
						Ayat:       1,
						TextIndo:   "test",
						TextArabic: "test",
					},
				},
			},
		},
		{
			name:           "empty query",
			query:          "",
			page:           "1",
			juz:            "0",
			surat:          "0",
			rowsPerPage:    "10",
			mockResults:    nil,
			mockError:      ErrInvalidRequest("search query is required"),
			expectedCode:   http.StatusBadRequest,
			expectedResult: nil,
		},
		{
			name:        "invalid page number",
			query:       "test",
			page:        "invalid",
			juz:         "0",
			surat:       "0",
			rowsPerPage: "10",
			mockResults: &search.SearchResult{
				Results:    []*model.AyatDetail{},
				TotalCount: 0,
				Count:      []*model.CountResult{},
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedResult: &ResultJsonFormat{
				Pagination: &Pagination{
					CurrentPage: 1,
					RowsPerPage: 10,
					TotalPages:  0,
					TotalRows:   0,
				},
				Results: []*model.AyatDetail{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockSearchService := new(MockSearchService)
			mockValidator := new(MockSearchRequestValidator)
			mockResponseBuilder := new(MockSearchResponseBuilder)

			// Setup request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest("GET", "/?q="+tt.query+"&p="+tt.page+"&juz="+tt.juz+"&surat="+tt.surat+"&n="+tt.rowsPerPage, nil)
			c.Request = req

			// Setup validator mock
			request := &SearchRequest{
				Query:       tt.query,
				Page:        1,
				Juz:         0,
				Surat:       0,
				RowsPerPage: 10,
			}
			mockValidator.On("ValidateRequest", c).Return(request, tt.mockError)

			// Setup search service mock
			if tt.mockError == nil {
				mockSearchService.On("FullTextSearch", tt.query, 0, 0, 1, 10).Return(tt.mockResults, nil)
				mockResponseBuilder.On("BuildResponse", tt.mockResults, request).Return(tt.expectedResult)
			}

			// Create controller with mocks
			controller := NewSearchController(mockSearchService, mockValidator, mockResponseBuilder)

			// Execute test
			controller.SearchHandler(c)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.mockError == nil {
				var response Response
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, response.Data)
			}

			// Verify mock expectations
			mockValidator.AssertExpectations(t)
			if tt.mockError == nil {
				mockSearchService.AssertExpectations(t)
				mockResponseBuilder.AssertExpectations(t)
			}
		})
	}
}

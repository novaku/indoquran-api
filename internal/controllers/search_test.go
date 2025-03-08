package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"indoquran-api/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSearch struct {
	mock.Mock
}

func (m *MockSearch) FullTextSearch(query string, juz, surat, page, rowsPerPage int) (interface{}, int64, interface{}, error) {
	args := m.Called(query, juz, surat, page, rowsPerPage)
	return args.Get(0), args.Get(1).(int64), args.Get(2), args.Error(3)
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
		mockResults    interface{}
		mockCount      int64
		mockAggregate  interface{}
		mockError      error
		expectedCode   int
		expectedResult ResultJsonFormat
	}{
		{
			name:          "successful search",
			query:         "test",
			page:          "1",
			juz:           "0",
			surat:         "0",
			rowsPerPage:   "10",
			mockResults:   []string{"result1", "result2"},
			mockCount:     20,
			mockAggregate: map[string]interface{}{"total": 20},
			mockError:     nil,
			expectedCode:  http.StatusOK,
			expectedResult: ResultJsonFormat{
				Aggregate: []*model.CountResult{
					{
						Type:       "",
						Identifier: 0,
						Count:      0,
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
						ID:         0,
						Juz:        0,
						Surat:      0,
						Ayat:       0,
						TextIndo:   "",
						TextArabic: "",
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
			mockCount:      0,
			mockAggregate:  nil,
			mockError:      errors.New("empty search query"),
			expectedCode:   http.StatusInternalServerError,
			expectedResult: ResultJsonFormat{},
		},
		{
			name:          "invalid page number",
			query:         "test",
			page:          "invalid",
			juz:           "0",
			surat:         "0",
			rowsPerPage:   "10",
			mockResults:   []string{},
			mockCount:     0,
			mockAggregate: nil,
			mockError:     nil,
			expectedCode:  http.StatusOK,
			expectedResult: ResultJsonFormat{
				Pagination: &Pagination{
					CurrentPage: 1,
					RowsPerPage: 10,
					TotalPages:  0,
					TotalRows:   0,
				},
				Results: []*model.AyatDetail{
					{
						ID:         0,
						Juz:        0,
						Surat:      0,
						Ayat:       0,
						TextIndo:   "",
						TextArabic: "",
					},
				},
			},
		},
		{
			name:          "custom rows per page",
			query:         "test",
			page:          "1",
			juz:           "0",
			surat:         "0",
			rowsPerPage:   "5",
			mockResults:   []string{"result1"},
			mockCount:     1,
			mockAggregate: map[string]interface{}{"total": 1},
			mockError:     nil,
			expectedCode:  http.StatusOK,
			expectedResult: ResultJsonFormat{
				Aggregate: []*model.CountResult{
					{
						Type:       "",
						Identifier: 0,
						Count:      0,
					},
				},
				Pagination: &Pagination{
					CurrentPage: 1,
					RowsPerPage: 5,
					TotalPages:  1,
					TotalRows:   1,
				},
				Results: []*model.AyatDetail{
					{
						ID:         0,
						Juz:        0,
						Surat:      0,
						Ayat:       0,
						TextIndo:   "",
						TextArabic: "",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSearch := new(MockSearch)
			mockSearch.On("FullTextSearch", tt.query, 0, 0, 1, 10).Return(tt.mockResults, tt.mockCount, tt.mockAggregate, tt.mockError)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest("GET", "/?q="+tt.query+"&p="+tt.page+"&juz="+tt.juz+"&surat="+tt.surat+"&n="+tt.rowsPerPage, nil)
			c.Request = req

			SearchHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.mockError == nil {
				var response Response
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, response.Data)
			}

			mockSearch.AssertExpectations(t)
		})
	}
}

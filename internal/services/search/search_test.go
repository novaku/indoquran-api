package search

import (
	"indoquran-api/internal/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
)

// buildSearchQuery constructs a SQL query based on the search term, juz, and surat filters.
func (s *DatabaseSearchRepository) buildSearchQuery(query string, juz int, surat int) *gorm.DB {
	sqlQuery := s.db.Table("quran_translation").Select("*")

	if query != "" {
		sqlQuery = sqlQuery.Where("MATCH(text_indo, text_arabic) AGAINST(? IN BOOLEAN MODE)", query)
	}
	if juz > 0 {
		sqlQuery = sqlQuery.Where("juz = ?", juz)
	}
	if surat > 0 {
		sqlQuery = sqlQuery.Where("surat = ?", surat)
	}

	return sqlQuery
}

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock
}

func setupTestRedis(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	return gormDB, mock
}

func TestSearch_FullTextSearch(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	tests := []struct {
		name        string
		searchText  string
		juz         int
		surat       int
		page        int
		limit       int
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectEmpty bool
	}{
		{
			name:       "Valid search with results",
			searchText: "bismillah",
			juz:        0,
			surat:      0,
			page:       1,
			limit:      10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"ayat_key", "surat", "ayat", "text", "juz"}).
					AddRow("0001001", 1, 1, "بِسْمِ ٱللَّهِ ٱلرَّحْمَٰنِ ٱلرَّحِيمِ", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat`")).
					WillReturnRows(rows)
			},
			expectError: false,
			expectEmpty: false,
		},
		{
			name:       "Search with juz filter",
			searchText: "bismillah",
			juz:        1,
			surat:      0,
			page:       1,
			limit:      10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"ayat_key", "surat", "ayat", "text", "juz"}).
					AddRow("0001001", 1, 1, "بِسْمِ ٱللَّهِ ٱلرَّحْمَٰنِ ٱلرَّحِيمِ", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat` WHERE juz = ?")).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectError: false,
			expectEmpty: false,
		},
		{
			name:       "Search with surat filter",
			searchText: "bismillah",
			juz:        0,
			surat:      1,
			page:       1,
			limit:      10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"ayat_key", "surat", "ayat", "text", "juz"}).
					AddRow("0001001", 1, 1, "بِسْمِ ٱللَّهِ ٱلرَّحْمَٰنِ ٱلرَّحِيمِ", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat` WHERE surat = ?")).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectError: false,
			expectEmpty: false,
		},
		{
			name:       "No results found",
			searchText: "nonexistent",
			juz:        0,
			surat:      0,
			page:       1,
			limit:      10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat`")).
					WillReturnRows(sqlmock.NewRows([]string{"ayat_key", "surat", "ayat", "text", "juz"}))
			},
			expectError: false,
			expectEmpty: true,
		},
		{
			name:        "Invalid page number",
			searchText:  "bismillah",
			juz:         0,
			surat:       0,
			page:        -1,
			limit:       10,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			service := NewSearchService(gormDB, nil)
			result, err := service.FullTextSearch(tt.searchText, tt.juz, tt.surat, tt.page, tt.limit)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectEmpty {
					assert.Empty(t, result.Results)
					assert.Zero(t, result.TotalCount)
					assert.Empty(t, result.Count)
				} else {
					assert.NotEmpty(t, result.Results)
					assert.Greater(t, result.TotalCount, int64(0))
					assert.NotNil(t, result.Count)
					assert.IsType(t, []*model.AyatDetail{}, result.Results)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBuildSQLLikeClause(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single word",
			input:    "test",
			expected: "quran_translation.translation LIKE '%test%'",
		},
		{
			name:     "Multiple words",
			input:    "hello world",
			expected: "quran_translation.translation LIKE '%hello%' AND quran_translation.translation LIKE '%world%'",
		},
		{
			name:     "More than 5 words (should limit to 5)",
			input:    "one two three four five six",
			expected: "quran_translation.translation LIKE '%one%' AND quran_translation.translation LIKE '%two%' AND quran_translation.translation LIKE '%three%' AND quran_translation.translation LIKE '%four%' AND quran_translation.translation LIKE '%five%'",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildSQLLikeClause(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// For testing purposes only
func (s *SearchService) GetQueryBuilder() *DatabaseQueryBuilder {
	return s.queryBuilder.(*DatabaseQueryBuilder)
}

func TestSearch_buildSearchQuery(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		juz         int
		surat       int
		expectedSQL string
	}{
		{
			name:        "Basic search query",
			query:       "test",
			juz:         0,
			surat:       0,
			expectedSQL: "MATCH(text_indo, text_arabic) AGAINST(? IN BOOLEAN MODE)",
		},
		{
			name:        "Search with juz filter",
			query:       "test",
			juz:         1,
			surat:       0,
			expectedSQL: "MATCH(text_indo, text_arabic) AGAINST(? IN BOOLEAN MODE) AND juz = ?",
		},
		{
			name:        "Search with surat filter",
			query:       "test",
			juz:         0,
			surat:       1,
			expectedSQL: "MATCH(text_indo, text_arabic) AGAINST(? IN BOOLEAN MODE) AND surat = ?",
		},
		{
			name:        "Search with both filters",
			query:       "test",
			juz:         1,
			surat:       1,
			expectedSQL: "MATCH(text_indo, text_arabic) AGAINST(? IN BOOLEAN MODE) AND juz = ? AND surat = ?",
		},
	}

	db, _ := setupTestDB(t)
	rdb := setupTestRedis(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSearchService(db, rdb).(*SearchService)
			query := service.GetQueryBuilder().BuildSearchQuery(tt.query, tt.juz, tt.surat)
			assert.Contains(t, query.Statement.SQL.String(), tt.expectedSQL)
		})
	}
}

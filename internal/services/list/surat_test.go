package list

import (
	"encoding/json"
	"indoquran-api/internal/model"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestSurat_GetSuratList(t *testing.T) {
	// Create SQL mock
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer sqlDB.Close()

	mockData := []model.IdMuntakhab{
		{
			Index: 1,
			Surat: 1,
			Ayat:  7,
			Text:  "Al-Fatihah",
		},
		{
			Index: 2,
			Surat: 2,
			Ayat:  286,
			Text:  "Al-Baqarah",
		},
	}

	tests := []struct {
		name        string
		suratID     string
		cacheSetup  func(*redis.Client)
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectEmpty bool
	}{
		{
			name:    "Get all surats from cache",
			suratID: "",
			cacheSetup: func(rdb *redis.Client) {
				data, _ := json.Marshal(mockData)
				rdb.Set("surat:all", data, 24*time.Hour)
			},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: false,
			expectEmpty: false,
		},
		{
			name:    "Get specific surat from cache",
			suratID: "1",
			cacheSetup: func(rdb *redis.Client) {
				data, _ := json.Marshal([]model.IdMuntakhab{mockData[0]})
				rdb.Set("surat:1", data, 24*time.Hour)
			},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: false,
			expectEmpty: false,
		},
		{
			name:       "Cache miss, DB success",
			suratID:    "",
			cacheSetup: func(rdb *redis.Client) {},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"index", "surat", "ayat", "text"})
				for _, d := range mockData {
					rows.AddRow(d.Index, d.Surat, d.Ayat, d.Text)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `id_muntakhab`")).
					WillReturnRows(rows)
			},
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "Invalid surat ID",
			suratID:     "invalid",
			cacheSetup:  func(rdb *redis.Client) {},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
			expectEmpty: true,
		},
		{
			name:       "Surat not found",
			suratID:    "999",
			cacheSetup: func(rdb *redis.Client) {},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `id_muntakhab`")).
					WithArgs(999).
					WillReturnRows(sqlmock.NewRows([]string{"index", "surat", "ayat", "text"}))
			},
			expectError: true,
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisMock := setupMockRedis()
			tt.cacheSetup(redisMock)
			tt.mockSetup(mock)

			service := &Surat{
				cache:    NewRedisCacheService(),
				database: NewGormDatabaseService(),
			}

			results, err := service.GetSuratList(tt.suratID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectEmpty {
					assert.Empty(t, results)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, results)
				if tt.suratID != "" {
					suratID, _ := strconv.Atoi(tt.suratID)
					for _, r := range results {
						assert.Equal(t, suratID, r.Surat)
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func setupMockRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestSurat_buildSuratKey(t *testing.T) {
	tests := []struct {
		name     string
		suratID  string
		expected string
	}{
		{
			name:     "Empty surat ID (all surats)",
			suratID:  "",
			expected: "surat:all",
		},
		{
			name:     "Specific surat",
			suratID:  "1",
			expected: "surat:1",
		},
		{
			name:     "Another surat",
			suratID:  "114",
			expected: "surat:114",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Surat{}
			result := service.buildSuratKey(tt.suratID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

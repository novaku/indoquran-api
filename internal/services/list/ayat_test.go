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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAyat_GetAyatList(t *testing.T) {
	// Create SQL mock
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer sqlDB.Close()

	// Create mock GORM DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	tests := []struct {
		name        string
		suratID     string
		page        int
		pageSize    int
		cacheSetup  func(*redis.Client)
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectEmpty bool
	}{
		{
			name:     "Get first page of ayat from cache",
			suratID:  "1",
			page:     1,
			pageSize: 5,
			cacheSetup: func(rdb *redis.Client) {
				ayats := []map[string]interface{}{
					{"id": 1, "juz": 1, "surat": 1, "ayat": 1, "text_indo": "Dengan nama Allah", "text_arabic": "بِسْمِ ٱللَّهِ"},
					{"id": 2, "juz": 1, "surat": 1, "ayat": 2, "text_indo": "Segala puji bagi Allah", "text_arabic": "ٱلْحَمْدُ لِلَّهِ"},
				}
				data, _ := json.Marshal(ayats)
				rdb.Set("ayat:1:1:5", data, 24*time.Hour)
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat`")).
					WithArgs(1, 5, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "juz", "surat", "ayat", "text_indo", "text_arabic"}).
						AddRow(1, 1, 1, 1, "Dengan nama Allah", "بِسْمِ ٱللَّهِ").
						AddRow(2, 1, 1, 2, "Segala puji bagi Allah", "ٱلْحَمْدُ لِلَّهِ"))
			},
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "Invalid surat ID",
			suratID:     "invalid",
			page:        1,
			pageSize:    10,
			cacheSetup:  func(rdb *redis.Client) {},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
			expectEmpty: true,
		},
		{
			name:        "Invalid page number",
			suratID:     "1",
			page:        -1,
			pageSize:    10,
			cacheSetup:  func(rdb *redis.Client) {},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
			expectEmpty: true,
		},
		{
			name:        "Invalid page size",
			suratID:     "1",
			page:        1,
			pageSize:    0,
			cacheSetup:  func(rdb *redis.Client) {},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
			expectEmpty: true,
		},
		{
			name:       "Cache miss, DB success",
			suratID:    "1",
			page:       1,
			pageSize:   5,
			cacheSetup: func(rdb *redis.Client) {},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat`")).
					WithArgs(1, 5, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "juz", "surat", "ayat", "text_indo", "text_arabic"}).
						AddRow(1, 1, 1, 1, "Dengan nama Allah", "بِسْمِ ٱللَّهِ"))
			},
			expectError: false,
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisMock := setupMockRedisClient()
			tt.cacheSetup(redisMock)
			tt.mockSetup(mock)

			service := &Ayat{
				repo:  NewDatabaseRepository(gormDB),
				cache: NewRedisCache(redisMock),
			}

			result, err := service.GetAyatList(tt.suratID, tt.page, tt.pageSize)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectEmpty {
					assert.Nil(t, result)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if len(result) > 0 {
					assert.IsType(t, &model.AyatDetail{}, result[0])
				}
			}

			// Verify all database expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func setupMockRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// buildAyatListKey generates a cache key for ayat list based on suratID, page, and pageSize.
func (a *Ayat) buildAyatListKey(suratID string, page, pageSize int) string {
	return "ayat:" + suratID + ":" + strconv.Itoa(page) + ":" + strconv.Itoa(pageSize)
}

func TestAyat_buildAyatListKey(t *testing.T) {
	tests := []struct {
		name     string
		suratID  string
		page     int
		pageSize int
		expected string
	}{
		{
			name:     "Standard parameters",
			suratID:  "1",
			page:     1,
			pageSize: 10,
			expected: "ayat:1:1:10",
		},
		{
			name:     "Different page and size",
			suratID:  "2",
			page:     3,
			pageSize: 5,
			expected: "ayat:2:3:5",
		},
		{
			name:     "Large numbers",
			suratID:  "114",
			page:     100,
			pageSize: 50,
			expected: "ayat:114:100:50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Ayat{}
			result := service.buildAyatListKey(tt.suratID, tt.page, tt.pageSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}

package detail

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

// Detail represents the service with database and Redis clients.
type Detail struct {
	DB  *gorm.DB
	RDS *redis.Client
}

// GetAyat retrieves an Ayat by its ID, first checking the cache, then the database.
func (d *Detail) GetAyat(ayatID string) (*testAyat, error) {
	var ayat testAyat

	// Check Redis cache
	cachedData, err := d.RDS.Get("ayat:" + ayatID).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &ayat); err == nil {
			return &ayat, nil
		}
	}

	// Query database
	result := d.DB.Raw("SELECT id, juz, surat, ayat, text_indo, text_arabic FROM `quran_ayat` WHERE id = ?", ayatID).Scan(&ayat)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	// Cache the result in Redis
	data, _ := json.Marshal(ayat)
	d.RDS.Set("ayat:"+ayatID, data, 24*time.Hour)

	return &ayat, nil
}

type testAyat struct {
	ID         int    `json:"id"`
	Juz        int    `json:"juz"`
	Surat      int    `json:"surat"`
	Ayat       int    `json:"ayat"`
	TextIndo   string `json:"text_indo"`
	TextArabic string `json:"text_arabic"`
}

func TestDetail_GetAyat(t *testing.T) {
	db, mock := setupTestDB(t)
	rdb := setupTestRedis(t)

	tests := []struct {
		name         string
		ayatID       string
		cacheSetup   func(*redis.Client)
		mockSetup    func()
		expectError  bool
		expectCached bool
		expectedAyat *testAyat
	}{
		{
			name:   "Get ayat from cache",
			ayatID: "1",
			cacheSetup: func(rdb *redis.Client) {
				ayat := &testAyat{
					ID: 1, Juz: 1, Surat: 1, Ayat: 1,
					TextIndo:   "Dengan nama Allah",
					TextArabic: "بِسْمِ ٱللَّهِ",
				}
				data, _ := json.Marshal(ayat)
				rdb.Set("ayat:1", data, 24*time.Hour)
			},
			mockSetup: func() {
				// No DB calls expected
			},
			expectError:  false,
			expectCached: true,
			expectedAyat: &testAyat{
				ID: 1, Juz: 1, Surat: 1, Ayat: 1,
				TextIndo:   "Dengan nama Allah",
				TextArabic: "بِسْمِ ٱللَّهِ",
			},
		},
		{
			name:   "Get ayat from database",
			ayatID: "1",
			cacheSetup: func(rdb *redis.Client) {
				// Empty cache
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "juz", "surat", "ayat", "text_indo", "text_arabic"}).
					AddRow(1, 1, 1, 1, "Dengan nama Allah", "بِسْمِ ٱللَّهِ")
				mock.ExpectQuery("^SELECT (.+) FROM `quran_ayat`").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectError:  false,
			expectCached: false,
			expectedAyat: &testAyat{
				ID: 1, Juz: 1, Surat: 1, Ayat: 1,
				TextIndo:   "Dengan nama Allah",
				TextArabic: "بِسْمِ ٱللَّهِ",
			},
		},
		{
			name:   "Ayat not found",
			ayatID: "999",
			cacheSetup: func(rdb *redis.Client) {
				// Empty cache
			},
			mockSetup: func() {
				mock.ExpectQuery("^SELECT (.+) FROM `quran_ayat`").
					WithArgs(999).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expectError:  true,
			expectCached: false,
			expectedAyat: nil,
		},
		{
			name:   "Invalid ayat ID",
			ayatID: "invalid",
			cacheSetup: func(rdb *redis.Client) {
				// Empty cache
			},
			mockSetup: func() {
				// No DB calls expected
			},
			expectError:  true,
			expectCached: false,
			expectedAyat: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test redis data
			tt.cacheSetup(rdb)

			// Setup mock DB expectations
			tt.mockSetup()

			service := &Detail{
				DB:  db,
				RDS: rdb,
			}

			result, err := service.GetAyat(tt.ayatID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				// For cached results, verify the cache was hit
				if tt.expectCached {
					cachedData, err := rdb.Get("ayat:" + tt.ayatID).Result()
					assert.NoError(t, err)
					assert.NotEmpty(t, cachedData)
				}

				// Verify the result matches expected data
				resultJSON, err := json.Marshal(result)
				assert.NoError(t, err)
				expectedJSON, err := json.Marshal(tt.expectedAyat)
				assert.NoError(t, err)
				assert.JSONEq(t, string(expectedJSON), string(resultJSON))
			}
		})
	}
}

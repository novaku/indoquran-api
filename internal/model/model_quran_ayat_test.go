package model

import (
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestQuranAyat(t *testing.T) {
	tests := []struct {
		name     string
		ayat     QuranAyat
		expected QuranAyat
	}{
		{
			name: "Valid QuranAyat entry with all fields",
			ayat: QuranAyat{
				AyatKey:    "0001001",
				AyatNumber: 1,
				Surat:      1,
				Ayat:       1,
				Text:       "بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ",
				Simple:     stringPtr("Bismillah"),
			},
			expected: QuranAyat{
				AyatKey:    "0001001",
				AyatNumber: 1,
				Surat:      1,
				Ayat:       1,
				Text:       "بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ",
				Simple:     stringPtr("Bismillah"),
			},
		},
		{
			name: "QuranAyat with null optional fields",
			ayat: QuranAyat{
				AyatKey:    "0002255",
				AyatNumber: 255,
				Surat:      2,
				Ayat:       255,
				Text:       "اللَّهُ لَا إِلَٰهَ إِلَّا هُوَ الْحَيُّ الْقَيُّومُ",
				Simple:     nil,
			},
			expected: QuranAyat{
				AyatKey:    "0002255",
				AyatNumber: 255,
				Surat:      2,
				Ayat:       255,
				Text:       "اللَّهُ لَا إِلَٰهَ إِلَّا هُوَ الْحَيُّ الْقَيُّومُ",
				Simple:     nil,
			},
		},
		{
			name: "Zero values with required fields only",
			ayat: QuranAyat{
				AyatKey:    "0000000",
				AyatNumber: 0,
				Surat:      0,
				Ayat:       0,
				Text:       "",
			},
			expected: QuranAyat{
				AyatKey:    "0000000",
				AyatNumber: 0,
				Surat:      0,
				Ayat:       0,
				Text:       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.AyatKey, tt.ayat.AyatKey)
			assert.Equal(t, tt.expected.AyatNumber, tt.ayat.AyatNumber)
			assert.Equal(t, tt.expected.Surat, tt.ayat.Surat)
			assert.Equal(t, tt.expected.Ayat, tt.ayat.Ayat)
			assert.Equal(t, tt.expected.Text, tt.ayat.Text)
			assert.Equal(t, tt.expected.Simple, tt.ayat.Simple)
		})
	}
}

func TestQuranAyat_GetAyat(t *testing.T) {
	// Create a new SQL mock
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	// Create a new GORM DB instance using the mock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	tests := []struct {
		name        string
		ayatNumber  string
		mockSetup   func(mock sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name:       "Successfully get ayat",
			ayatNumber: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat`")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"ayat_key", "surat", "ayat", "text"}).
						AddRow(1, 1, 1, "Sample ayat text"))
			},
			expectError: false,
		},
		{
			name:        "Invalid ayat number",
			ayatNumber:  "invalid",
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
		},
		{
			name:       "Ayat not found",
			ayatNumber: "999",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat`")).
					WithArgs(999).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup(mock)

			// Execute the function
			result, err := GetAyat(gormDB, 1, 1)

			// Verify expectations
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestQuranAyat_GetAyatInSurat(t *testing.T) {
	// Create a new SQL mock
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	// Create a new GORM DB instance using the mock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	tests := []struct {
		name        string
		suratID     string
		page        int
		pageSize    int
		mockSetup   func(mock sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name:     "Successfully get ayat in surat",
			suratID:  "1",
			page:     1,
			pageSize: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat`")).
					WithArgs(1, 10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"ayat_key", "surat", "ayat", "text"}).
						AddRow(1, 1, 1, "First ayat"))
			},
			expectError: false,
		},
		{
			name:        "Invalid surat ID",
			suratID:     "invalid",
			page:        1,
			pageSize:    10,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
		},
		{
			name:     "Surat not found",
			suratID:  "999",
			page:     1,
			pageSize: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat`")).
					WithArgs(999, 10, 0).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup(mock)

			// Convert string ID to int if valid
			suratIDInt, err := strconv.Atoi(tt.suratID)
			if err != nil && !tt.expectError {
				t.Errorf("Invalid surat ID for non-error test case: %v", err)
				return
			}

			// Execute the function
			result, err := GetAyatInSurat(gormDB, suratIDInt)

			// Verify expectations
			if tt.expectError {
				assert.Error(t, err)
				if err != nil && err.Error() == "record not found" {
					assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAyat(t *testing.T) {
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
		ayatNumber  string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name:       "Valid surat and ayat",
			suratID:    "1",
			ayatNumber: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "juz", "surat", "ayat", "text_indo", "text_arabic"}).
					AddRow(1, 1, 1, 1, "Dengan nama Allah", "بِسْمِ ٱللَّهِ")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat`")).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			expectError: false,
		},
		{
			name:        "Invalid surat ID",
			suratID:     "invalid",
			ayatNumber:  "1",
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
		},
		{
			name:        "Invalid ayat number",
			suratID:     "1",
			ayatNumber:  "invalid",
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
		},
		{
			name:       "Ayat not found",
			suratID:    "1",
			ayatNumber: "999",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat`")).
					WithArgs(1, 999).
					WillReturnRows(sqlmock.NewRows([]string{"id", "juz", "surat", "ayat", "text_indo", "text_arabic"}))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			// Convert string IDs to integers
			suratID, err := strconv.Atoi(tt.suratID)
			if err != nil && !tt.expectError {
				t.Fatalf("Test case error: invalid surat ID %q", tt.suratID)
			}

			ayatNumber, err := strconv.Atoi(tt.ayatNumber)
			if err != nil && !tt.expectError {
				t.Fatalf("Test case error: invalid ayat number %q", tt.ayatNumber)
			}

			result, err := GetAyat(gormDB, suratID, ayatNumber)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, suratID, result.Surat)
				assert.Equal(t, ayatNumber, result.Ayat)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAyatInSurat(t *testing.T) {
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
		mockSetup   func(sqlmock.Sqlmock)
		expectCount int
		expectError bool
	}{
		{
			name:    "Valid surat",
			suratID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "juz", "surat", "ayat", "text_indo", "text_arabic"}).
					AddRow(1, 1, 1, 1, "Dengan nama Allah", "بِسْمِ ٱللَّهِ").
					AddRow(2, 1, 1, 2, "Segala puji bagi Allah", "ٱلْحَمْدُ لِلَّهِ")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat` WHERE surat = ?")).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectCount: 2,
			expectError: false,
		},
		{
			name:        "Invalid surat ID",
			suratID:     "invalid",
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectCount: 0,
			expectError: true,
		},
		{
			name:    "No ayat found",
			suratID: "999",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `quran_ayat` WHERE surat = ?")).
					WithArgs(999).
					WillReturnRows(sqlmock.NewRows([]string{"id", "juz", "surat", "ayat", "text_indo", "text_arabic"}))
			},
			expectCount: 0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			// Convert string ID to integer
			suratID, err := strconv.Atoi(tt.suratID)
			if err != nil && !tt.expectError {
				t.Fatalf("Test case error: invalid surat ID %q", tt.suratID)
			}

			results, err := GetAyatInSurat(gormDB, suratID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, results)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectCount, len(results))
				if len(results) > 0 {
					for _, result := range results {
						assert.Equal(t, suratID, result.Surat)
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

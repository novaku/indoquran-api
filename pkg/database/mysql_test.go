package database

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitDatabase(t *testing.T) {
	tests := []struct {
		name        string
		dbConfig    map[string]string
		expectPanic bool
	}{
		{
			name: "Invalid database configuration",
			dbConfig: map[string]string{
				"DB_USER":     "invalid",
				"DB_PASSWORD": "invalid",
				"DB_HOST":     "invalid",
				"DB_PORT":     "3306",
				"DB_NAME":     "invalid",
			},
			expectPanic: true,
		},
		{
			name: "Valid database configuration",
			dbConfig: map[string]string{
				"DB_USER":     "root",
				"DB_PASSWORD": "root",
				"DB_HOST":     "localhost",
				"DB_PORT":     "3306",
				"DB_NAME":     "indoquran",
			},
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset db connection
			db = nil

			// Set viper config
			for k, v := range tt.dbConfig {
				viper.Set(k, v)
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					InitDatabase()
				})
			} else {
				InitDatabase()
				assert.NotNil(t, db)
			}
		})
	}
}

func TestGetDB(t *testing.T) {
	tests := []struct {
		name      string
		setupDB   bool
		expectNil bool
	}{
		{
			name:      "Returns nil when database not initialized",
			setupDB:   false,
			expectNil: true,
		},
		{
			name:      "Returns valid database connection when initialized",
			setupDB:   true,
			expectNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset db connection
			db = nil

			if tt.setupDB {
				viper.Set("DB_USER", "root")
				viper.Set("DB_PASSWORD", "root")
				viper.Set("DB_HOST", "localhost")
				viper.Set("DB_PORT", "3306")
				viper.Set("DB_NAME", "indoquran")
				InitDatabase()
			}

			result := GetDB()

			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

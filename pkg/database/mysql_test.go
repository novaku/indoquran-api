package database

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	tests := []struct {
		name        string
		dbConfig    map[string]string
		expectError bool
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
			expectError: true,
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
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set viper config
			for k, v := range tt.dbConfig {
				viper.Set(k, v)
			}

			mysqlDB := NewMySQLDatabase()
			dbManager := NewDatabaseManager(mysqlDB)
			err := dbManager.Initialize()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dbManager.GetDB())
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
			mysqlDB := NewMySQLDatabase()
			dbManager := NewDatabaseManager(mysqlDB)

			if tt.setupDB {
				viper.Set("DB_USER", "root")
				viper.Set("DB_PASSWORD", "root")
				viper.Set("DB_HOST", "localhost")
				viper.Set("DB_PORT", "3306")
				viper.Set("DB_NAME", "indoquran")
				err := dbManager.Initialize()
				assert.NoError(t, err)
			}

			result := dbManager.GetDB()

			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

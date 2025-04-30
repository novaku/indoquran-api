package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"indoquran-api/internal/constants"
	"indoquran-api/pkg/logger"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() {
	var (
		err   error
		sqlDB *sql.DB
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString(constants.DB_USER),
		viper.GetString(constants.DB_PASSWORD),
		viper.GetString(constants.DB_HOST),
		viper.GetString(constants.DB_PORT),
		viper.GetString(constants.DB_NAME),
	)

	// Retry logic to wait for the database to be ready
	for retries := 5; retries > 0; retries-- {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			// Ping the database to check if it's connected
			sqlDB, err = db.DB()
			if err == nil {
				err = sqlDB.Ping()
			}
		}

		if err == nil {
			logger.WriteLog(logger.LogLevelInfo, "Connected to the database successfully")
			return
		}

		logger.WriteLog(logger.LogLevelError, "Failed to connect to database: %#v Retrying in 5 seconds... (%d retries left)", err, retries)
		time.Sleep(5 * time.Second)
	}

	logger.WriteLog(logger.LogLevelFatal, "Failed to connect to database after retries: %v", err)
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	if os.Getenv("ENV") == "local" {
		return db.Debug()
	}

	return db
}

// SetDB sets the database instance - used for testing
func SetDB(instance *gorm.DB) {
	db = instance
}

// SetTestDB sets a test database instance - used for testing with sqlmock
func SetTestDB(sqlDB *sql.DB) {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Error creating test DB instance: %v", err))
	}

	db = gormDB
}

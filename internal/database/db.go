package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"indoquran-api/internal/config"
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
		viper.GetString(config.DB_USER),
		viper.GetString(config.DB_PASSWORD),
		viper.GetString(config.DB_HOST),
		viper.GetString(config.DB_PORT),
		viper.GetString(config.DB_NAME),
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

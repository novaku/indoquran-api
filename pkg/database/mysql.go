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

// Database interface defines the contract for database operations
type Database interface {
	Connect() error
	GetConnection() *gorm.DB
	Close() error
}

// MySQLDatabase implements the Database interface
type MySQLDatabase struct {
	DB *gorm.DB
}

// NewMySQLDatabase creates a new MySQL database instance
func NewMySQLDatabase() *MySQLDatabase {
	return &MySQLDatabase{}
}

// Connect establishes a connection to the MySQL database
func (m *MySQLDatabase) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString(constants.DB_USER),
		viper.GetString(constants.DB_PASSWORD),
		viper.GetString(constants.DB_HOST),
		viper.GetString(constants.DB_PORT),
		viper.GetString(constants.DB_NAME),
	)

	var err error
	for retries := 5; retries > 0; retries-- {
		m.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err := m.DB.DB()
			if err == nil {
				if err = sqlDB.Ping(); err != nil {
					continue
				}
			}
		}

		if err == nil {
			logger.WriteLog(logger.LogLevelInfo, "Connected to the database successfully")
			return nil
		}

		logger.WriteLog(logger.LogLevelError, "Failed to connect to database: %#v Retrying in 5 seconds... (%d retries left)", err, retries)
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("failed to connect to database after retries: %v", err)
}

// GetConnection returns the database connection
func (m *MySQLDatabase) GetConnection() *gorm.DB {
	if os.Getenv("ENV") == "local" {
		return m.DB.Debug()
	}
	return m.DB
}

// Close closes the database connection
func (m *MySQLDatabase) Close() error {
	if m.DB != nil {
		sqlDB, err := m.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// DatabaseManager manages the database connection
type DatabaseManager struct {
	db Database
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager(db Database) *DatabaseManager {
	return &DatabaseManager{db: db}
}

// Initialize initializes the database connection
func (dm *DatabaseManager) Initialize() error {
	return dm.db.Connect()
}

// GetDB returns the database connection
func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.db.GetConnection()
}

// SetTestDB sets a test database instance
func SetTestDB(sqlDB *sql.DB) {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Error creating test DB instance: %v", err))
	}

	// Create a new MySQLDatabase instance for testing
	mysqlDB := &MySQLDatabase{DB: gormDB}
	manager := NewDatabaseManager(mysqlDB)
	manager.Initialize()
}

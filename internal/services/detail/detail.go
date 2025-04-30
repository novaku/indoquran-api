package detail

import (
	"encoding/json"
	"fmt"
	"indoquran-api/pkg/cache"
	"indoquran-api/pkg/database"
	"indoquran-api/pkg/logger"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// DetailServiceInterface defines the contract for detail operations
type DetailServiceInterface interface {
	GetAyat(ayatID string) (interface{}, error)
}

// CacheServiceInterface defines the contract for caching operations
type CacheServiceInterface interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
}

// DatabaseServiceInterface defines the contract for database operations
type DatabaseServiceInterface interface {
	GetAyatByID(ayatID string) (interface{}, error)
}

// DetailService implements DetailServiceInterface
type DetailService struct {
	cacheService    CacheServiceInterface
	databaseService DatabaseServiceInterface
}

// RedisCacheService implements CacheServiceInterface
type RedisCacheService struct {
	client *redis.Client
}

// GormDatabaseService implements DatabaseServiceInterface
type GormDatabaseService struct {
	db *gorm.DB
}

// NewDetailService creates a new instance of DetailService
func NewDetailService(cacheService CacheServiceInterface, databaseService DatabaseServiceInterface) DetailServiceInterface {
	return &DetailService{
		cacheService:    cacheService,
		databaseService: databaseService,
	}
}

// NewRedisCacheService creates a new instance of RedisCacheService
func NewRedisCacheService() CacheServiceInterface {
	return &RedisCacheService{
		client: cache.GetRedis(),
	}
}

// NewGormDatabaseService creates a new instance of GormDatabaseService
func NewGormDatabaseService() DatabaseServiceInterface {
	return &GormDatabaseService{
		db: database.GetDB(),
	}
}

// Get implements CacheServiceInterface
func (r *RedisCacheService) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

// Set implements CacheServiceInterface
func (r *RedisCacheService) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(key, value, expiration).Err()
}

// GetAyatByID implements DatabaseServiceInterface
func (g *GormDatabaseService) GetAyatByID(ayatID string) (interface{}, error) {
	var ayat interface{}
	result := g.db.Table("quran_translation").
		Select("quran_translation.translation_id AS id, quran_ayat.juz, quran_ayat.surat, quran_ayat.ayat, quran_translation.translation AS text_indo, quran_ayat.text AS text_arabic").
		Joins("JOIN quran_ayat ON quran_ayat.ayat_number = quran_translation.translation_id").
		Where("quran_translation.translation_id = ?", ayatID).
		Find(&ayat)

	if result.Error != nil {
		logger.WriteLog(logger.LogLevelError, "Failed to fetch ayat: %v", result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("ayat not found")
	}

	return ayat, nil
}

// GetAyat implements DetailServiceInterface
func (d *DetailService) GetAyat(ayatID string) (interface{}, error) {
	// Create a Redis cache key based on ayatID
	cacheKey := fmt.Sprintf("ayat:%s", ayatID)

	// Try to get data from Redis cache
	cachedData, err := d.cacheService.Get(cacheKey)
	if err == nil {
		// Cache hit: unmarshal and return the data
		var ayat interface{}
		if err := json.Unmarshal([]byte(cachedData), &ayat); err != nil {
			logger.WriteLog(logger.LogLevelError, "Failed to unmarshal cached ayat: %v", err)
			return nil, fmt.Errorf("cache unmarshal error: %w", err)
		}
		return ayat, nil
	}

	// Cache miss or error: query the database
	ayat, err := d.databaseService.GetAyatByID(ayatID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Store in cache for future requests
	if data, err := json.Marshal(ayat); err == nil {
		if err := d.cacheService.Set(cacheKey, data, 24*time.Hour); err != nil {
			logger.WriteLog(logger.LogLevelError, "Failed to cache ayat: %v", err)
			// Don't return error here, just log it since we have the data
		}
	} else {
		logger.WriteLog(logger.LogLevelError, "Failed to marshal ayat for caching: %v", err)
		// Don't return error here, just log it since we have the data
	}

	return ayat, nil
}

// DefaultDetailService returns a new DetailService with default implementations
func DefaultDetailService() DetailServiceInterface {
	return NewDetailService(
		NewRedisCacheService(),
		NewGormDatabaseService(),
	)
}

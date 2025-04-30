package list

import (
	"encoding/json"
	"fmt"
	"indoquran-api/internal/model"
	"indoquran-api/pkg/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// Repository defines the interface for data access operations
type Repository interface {
	GetAyatList(suratID int, offset, limit int) ([]*model.AyatDetail, error)
}

// Cache defines the interface for caching operations
type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
}

// AyatService defines the interface for ayat-related operations
type AyatService interface {
	GetAyatList(suratID string, page, pageSize int) ([]*model.AyatDetail, error)
}

// DatabaseRepository implements the Repository interface
type DatabaseRepository struct {
	db *gorm.DB
}

func NewDatabaseRepository(db *gorm.DB) Repository {
	return &DatabaseRepository{db: db}
}

func (r *DatabaseRepository) GetAyatList(suratID int, offset, limit int) ([]*model.AyatDetail, error) {
	var ayatList []*model.AyatDetail

	query := r.db.Select("quran_translation.translation_id AS id, quran_ayat.juz, quran_ayat.surat, quran_ayat.ayat,quran_translation.translation AS text_indo,quran_ayat.text AS text_arabic").
		Table("quran_translation").
		Joins("JOIN quran_ayat ON quran_ayat.ayat_number = quran_translation.translation_id")

	if suratID > 0 {
		query = query.Where("quran_ayat.surat = ?", suratID)
	}

	err := query.Order("quran_ayat.surat, quran_ayat.ayat ASC").
		Offset(offset).
		Limit(limit).
		Scan(&ayatList).Error

	if err != nil {
		logger.WriteLog(logger.LogLevelError, "Error retrieving records: %#v", err)
		return nil, err
	}

	return ayatList, nil
}

// RedisCache implements the Cache interface
type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(key string) (string, error) {
	return c.client.Get(key).Result()
}

func (c *RedisCache) Set(key string, value string, expiration time.Duration) error {
	return c.client.Set(key, value, expiration).Err()
}

// Ayat implements the AyatService interface
type Ayat struct {
	repo  Repository
	cache Cache
}

func NewAyat(repo Repository, cache Cache) AyatService {
	return &Ayat{
		repo:  repo,
		cache: cache,
	}
}

// GetAyatList retrieves a list of ayat based on the provided suratID, page, and pageSize
func (a *Ayat) GetAyatList(suratID string, page, pageSize int) ([]*model.AyatDetail, error) {
	// Convert suratID to integer
	s, err := strconv.Atoi(suratID)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, "Error converting suratID to integer: %#v", err)
		return nil, err
	}

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Create a Redis cache key based on suratID, page, and pageSize
	cacheKey := fmt.Sprintf("ayatList:surat:%d:page:%d:pageSize:%d", s, page, pageSize)

	// Try to get data from cache
	cachedData, err := a.cache.Get(cacheKey)
	if err != nil {
		// Cache miss: query the database
		ayatList, err := a.repo.GetAyatList(s, offset, pageSize)
		if err != nil {
			return nil, err
		}

		// Serialize the data and store it in cache
		serializedData, err := json.Marshal(ayatList)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error serializing records: %#v", err)
			return nil, err
		}

		err = a.cache.Set(cacheKey, string(serializedData), 24*time.Hour)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error caching data: %#v", err)
			return nil, err
		}

		return ayatList, nil
	}

	// Cache hit: deserialize the data
	var ayatList []*model.AyatDetail
	err = json.Unmarshal([]byte(cachedData), &ayatList)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, "Error deserializing cached data: %#v", err)
		return nil, err
	}

	return ayatList, nil
}

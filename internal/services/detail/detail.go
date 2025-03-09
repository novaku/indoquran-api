package detail

import (
	"encoding/json"
	"fmt"
	"indoquran-api/internal/model"
	"indoquran-api/pkg/cache"
	"indoquran-api/pkg/database"
	"indoquran-api/pkg/logger"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type (
	IDetail interface {
		GetAyat(ayatID string) (*model.AyatDetail, error)
	}

	detail struct {
		db  *gorm.DB
		rds *redis.Client
	}
)

func NewDetail() IDetail {
	return &detail{
		db:  database.GetDB(),
		rds: cache.GetRedis(),
	}
}

// GetAyat retrieves ayat details by ayat ID
func (d *detail) GetAyat(ayatID string) (*model.AyatDetail, error) {
	var (
		ayat *model.AyatDetail
	)

	// Create a Redis cache key based on ayatID
	cacheKey := fmt.Sprintf("ayat:%s", ayatID)

	// Try to get data from Redis cache
	cachedData, err := d.rds.Get(cacheKey).Result()
	if err != nil {
		// Cache miss: query the database
		result := d.db.Select("quran_translation.translation_id AS id, quran_ayat.juz, quran_ayat.surat, quran_ayat.ayat,quran_translation.translation AS text_indo,quran_ayat.text AS text_arabic").
			Table("quran_translation").
			Joins("JOIN quran_ayat ON quran_ayat.ayat_number = quran_translation.translation_id").
			Where("quran_translation.translation_id = ?", ayatID).
			Find(&ayat)
		if result.Error != nil {
			logger.WriteLog(logger.LogLevelError, "Error retrieving records: %#v", result.Error)
			return nil, result.Error
		}

		// Serialize the data and store it in Redis with an expiration
		serializedData, err := json.Marshal(ayat)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error serializing record: %#v", err)
			return nil, err
		}
		err = d.rds.Set(cacheKey, serializedData, 24*time.Hour).Err() // Cache for 24 hours
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error caching data in Redis: %#v", err)
			return nil, err
		}
	} else {
		// Cache hit: deserialize the data
		err = json.Unmarshal([]byte(cachedData), &ayat)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error deserializing cached data: %#v", err)
			return nil, err
		}
	}

	return ayat, nil
}

package list

import (
	"encoding/json"
	"fmt"
	"indoquran-api/internal/model"
	"indoquran-api/pkg/cache"
	"indoquran-api/pkg/database"
	"indoquran-api/pkg/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type (
	Ayat struct {
		rds *redis.Client
		db  *gorm.DB
	}

	IAyat interface {
		GetAyatList(suratID string, page, pageSize int) ([]*model.AyatDetail, error)
	}
)

func NewAyat() IAyat {
	return &Ayat{
		rds: cache.GetRedis(),
		db:  database.GetDB(),
	}
}

// GetAyatList retrieves a list of ayat based on the provided suratID, page, and pageSize
func (a *Ayat) GetAyatList(suratID string, page, pageSize int) ([]*model.AyatDetail, error) {
	var (
		ayatList  []*model.AyatDetail
		sessQuery = a.db
	)

	sessQuery = sessQuery.Select("quran_translation.translation_id AS id, quran_ayat.juz, quran_ayat.surat, quran_ayat.ayat,quran_translation.translation AS text_indo,quran_ayat.text AS text_arabic").Table("quran_translation").Joins("JOIN quran_ayat ON quran_ayat.ayat_number = quran_translation.translation_id")

	// Convert suratID to integer
	s, err := strconv.Atoi(suratID)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, "Error converting suratID to integer: %#v", err)
	}

	if s > 0 {
		sessQuery = sessQuery.Where("quran_ayat.surat = ?", s)
	}

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Create a Redis cache key based on suratID, page, and pageSize
	cacheKey := fmt.Sprintf("ayatList:surat:%d:page:%d:pageSize:%d", s, page, pageSize)

	// Try to get data from Redis cache
	cachedData, err := a.rds.Get(cacheKey).Result()
	if err != nil {
		// Cache miss: query the database
		sessQuery = sessQuery.Order("quran_ayat.surat, quran_ayat.ayat ASC").
			Offset(offset).
			Limit(pageSize).
			Scan(&ayatList)
		if sessQuery.Error != nil {
			logger.WriteLog(logger.LogLevelError, "Error retrieving records: %#v", sessQuery.Error)
			return nil, sessQuery.Error
		}

		// Serialize the data and store it in Redis with an expiration
		serializedData, err := json.Marshal(ayatList)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error serializing records: %#v", err)
			return nil, err
		}
		err = a.rds.Set(cacheKey, serializedData, 24*time.Hour).Err() // Cache for 24 hours
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error caching data in Redis: %#v", err)
			return nil, err
		}
	} else {
		// Cache hit: deserialize the data
		err = json.Unmarshal([]byte(cachedData), &ayatList)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error deserializing cached data: %#v", err)
			return nil, err
		}
	}

	return ayatList, nil
}

package list

import (
	"encoding/json"
	"fmt"
	"time"

	"indoquran-api/internal/cache"
	"indoquran-api/internal/database"
	"indoquran-api/internal/model"
	"indoquran-api/pkg/logger"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type (
	Surat struct {
		rds *redis.Client
		db  *gorm.DB
	}

	ISurat interface {
		GetSuratList(suratID string) ([]model.IdMuntakhab, error)
	}
)

func NewSurat() ISurat {
	return &Surat{
		rds: cache.GetRedis(),
		db:  database.GetDB(),
	}
}

// GetSuratList retrieves a list of surat based on the provided suratID
func (s *Surat) GetSuratList(suratID string) ([]model.IdMuntakhab, error) {
	var (
		surats []model.IdMuntakhab
	)

	// Create a Redis cache key based on the suratID
	cacheKey := "surats:all"
	if suratID != "" {
		cacheKey = fmt.Sprintf("surats:%s", suratID)
	}

	// Try to get data from Redis cache
	cachedData, err := s.rds.Get(cacheKey).Result()
	if err != nil {
		logger.WriteLog(logger.LogLevelError, "Error fetching from Redis: %#v", err)
		querySession := s.db.Where("ayat = ?", 1).Order("surat ASC")
		if suratID != "" {
			querySession = querySession.Where("surat = ?", suratID)
		}
		result := querySession.Find(&surats)

		if result.Error != nil {
			logger.WriteLog(logger.LogLevelError, "Error retrieving records: %#v", result.Error)
			return nil, result.Error
		}

		// Serialize the data and store it in Redis with an expiration
		serializedData, err := json.Marshal(surats)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error serializing records: %#v", err)
			return nil, err
		}
		err = s.rds.Set(cacheKey, serializedData, 24*time.Hour).Err() // 24-hours expiration
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error caching data in Redis: %#v", err)
			return nil, err
		}
	} else {
		// Cache hit: deserialize the data
		err = json.Unmarshal([]byte(cachedData), &surats)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error deserializing cached data: %#v", err)
			return nil, err
		}
	}

	return surats, nil
}

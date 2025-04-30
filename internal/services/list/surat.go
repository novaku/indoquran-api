package list

import (
	"encoding/json"
	"indoquran-api/internal/model"
	"indoquran-api/pkg/logger"
	"time"
)

// CacheService defines the interface for caching operations
type CacheService interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
}

// DatabaseService defines the interface for database operations
type DatabaseService interface {
	GetSurats(suratID string) ([]model.IdMuntakhab, error)
}

// SuratService defines the interface for surat-related operations
type SuratService interface {
	GetSuratList(suratID string) ([]model.IdMuntakhab, error)
}

type (
	Surat struct {
		cache    CacheService
		database DatabaseService
	}

	ISurat interface {
		GetSuratList(suratID string) ([]model.IdMuntakhab, error)
	}
)

func NewSurat(cache CacheService, db DatabaseService) ISurat {
	return &Surat{
		cache:    cache,
		database: db,
	}
}

// GetSuratList retrieves a list of surat based on the provided suratID
func (s *Surat) GetSuratList(suratID string) ([]model.IdMuntakhab, error) {
	var surats []model.IdMuntakhab

	// Create a cache key based on the suratID
	cacheKey := s.buildSuratKey(suratID)

	// Try to get data from cache
	cachedData, err := s.cache.Get(cacheKey)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, "Error fetching from cache: %#v", err)

		// If cache miss, get from database
		surats, err = s.database.GetSurats(suratID)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error retrieving records: %#v", err)
			return nil, err
		}

		// Cache the results
		err = s.cache.Set(cacheKey, surats, 24*time.Hour)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error caching data: %#v", err)
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

func (g *gormDatabaseService) GetSurats(suratID string) ([]model.IdMuntakhab, error) {
	var surats []model.IdMuntakhab

	querySession := g.db.Where("ayat = ?", 1).Order("surat ASC")
	if suratID != "" {
		querySession = querySession.Where("surat = ?", suratID)
	}

	result := querySession.Find(&surats)
	if result.Error != nil {
		return nil, result.Error
	}

	return surats, nil
}

func (s *Surat) buildSuratKey(suratID string) string {
	if suratID == "" {
		return "surat:all"
	}
	return "surat:" + suratID
}

package search

import (
	"encoding/json"
	"fmt"
	"time"

	"indoquran-api/internal/cache"
	"indoquran-api/internal/model"
	"indoquran-api/pkg/logger"
)

// FullTextSearch performs a full-text search on the database.
func FullTextSearch(searchTerm string, juz, surat, pageNum, limit int) ([]*model.AyatDetail, int64, []*model.CountResult, error) {
	logger.WriteLog(logger.LogLevelInfo, "FullTextSearch called with searchTerm: %s, juz: %d, surat: %d, pageNum: %d, limit: %d", searchTerm, juz, surat, pageNum, limit)
	var (
		results     []*model.AyatDetail
		totalCount  int64
		offset      int
		err         error
		count       []*model.CountResult
		redisClient = cache.GetRedis()
	)

	// Get the search query from the URL query
	if searchTerm == "" {
		logger.WriteLog(logger.LogLevelError, "search query is required")
		return nil, 0, nil, fmt.Errorf("search query is required")
	}

	searchTerm = buildSQLLikeClause(searchTerm)

	// Ensure page and limit are positive
	if pageNum < 1 {
		pageNum = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Calculate offset
	offset = (pageNum - 1) * limit

	// Create a Redis cache key based on search parameters
	cacheKey := fmt.Sprintf("search:%s:juz:%d:surat:%d:page:%d:limit:%d", searchTerm, juz, surat, pageNum, limit)

	// Try to get data from Redis cache
	cachedData, err := redisClient.Get(cacheKey).Result()
	if err != nil {
		// Cache miss: query the database
		results, err = querySearch(searchTerm, limit, offset, juz, surat)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error fetching from database: %#v", err)
			return nil, 0, nil, err
		}

		// Get total count of results (without pagination)
		totalCount, err = queryTotalCount(searchTerm, juz, surat)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error fetching total count from database: %#v", err)
			return nil, 0, nil, err
		}

		// Get aggregate data
		count = queryAggregate(searchTerm, juz, surat)

		// Serialize the result and store it in Redis with an expiration
		cacheData := map[string]interface{}{
			"results":    results,
			"totalCount": totalCount,
			"count":      count,
		}
		serializedData, err := json.Marshal(cacheData)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error serializing data: %#v", err)
			return nil, 0, nil, err
		}
		err = redisClient.Set(cacheKey, serializedData, 24*time.Hour).Err() // Cache for 24 hours
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error caching data in Redis: %#v", err)
			return nil, 0, nil, err
		}
	} else {
		// Cache hit: deserialize the data
		var cachedResult map[string]interface{}
		err = json.Unmarshal([]byte(cachedData), &cachedResult)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error deserializing cached data: %#v", err)
			return nil, 0, nil, err
		}

		// Extract the cached results
		resultsBytes, _ := json.Marshal(cachedResult["results"])
		err = json.Unmarshal(resultsBytes, &results)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error deserializing results: %#v", err)
			return nil, 0, nil, err
		}

		totalCount = int64(cachedResult["totalCount"].(float64))

		countBytes, _ := json.Marshal(cachedResult["count"])
		err = json.Unmarshal(countBytes, &count)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, "Error deserializing count data: %#v", err)
			return nil, 0, nil, err
		}
	}

	return results, totalCount, count, nil
}

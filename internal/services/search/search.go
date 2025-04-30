package search

import (
	"encoding/json"
	"fmt"
	"indoquran-api/internal/model"
	"indoquran-api/pkg/logger"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// QueryBuilder defines the interface for building search queries
type QueryBuilder interface {
	BuildSearchQuery(searchTerm string, juz, surat int) *gorm.DB
	BuildCountQuery(searchTerm string, juz, surat int) *gorm.DB
	BuildAggregateQuery(searchTerm string, juz, surat int) *gorm.DB
}

// SearchRepository defines the interface for search operations
type SearchRepository interface {
	Search(query *gorm.DB, limit, offset int) ([]*model.AyatDetail, error)
	GetTotalCount(query *gorm.DB) (int64, error)
	GetAggregate(query *gorm.DB) []*model.CountResult
}

// CacheRepository defines the interface for caching operations
type CacheRepository interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
}

// SearchResult represents the complete search result
type SearchResult struct {
	Results    []*model.AyatDetail
	TotalCount int64
	Count      []*model.CountResult
}

// SearchServiceInterface defines the contract for search operations
type SearchServiceInterface interface {
	FullTextSearch(query string, juz, surat, page, rowsPerPage int) (*SearchResult, error)
}

// CacheServiceInterface defines the contract for caching operations
type CacheServiceInterface interface {
	GetFromCache(key string) (*SearchResult, error)
	SetInCache(key string, result *SearchResult, expiration time.Duration) error
}

// SearchService handles the business logic for search operations
type SearchService struct {
	queryBuilder QueryBuilder
	searchRepo   SearchRepository
	cacheService CacheServiceInterface
}

// NewSearchService creates a new instance of SearchService
func NewSearchService(db *gorm.DB, rds *redis.Client) SearchServiceInterface {
	return &SearchService{
		queryBuilder: NewDatabaseQueryBuilder(db),
		searchRepo:   NewDatabaseSearchRepository(db),
		cacheService: NewCacheService(rds),
	}
}

// CacheService implements CacheServiceInterface
type CacheService struct {
	cacheRepo CacheRepository
}

// NewCacheService creates a new instance of CacheService
func NewCacheService(rds *redis.Client) CacheServiceInterface {
	return &CacheService{
		cacheRepo: NewRedisCacheRepository(rds),
	}
}

// GetFromCache implements CacheServiceInterface
func (cs *CacheService) GetFromCache(key string) (*SearchResult, error) {
	cachedData, err := cs.cacheRepo.Get(key)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SetInCache implements CacheServiceInterface
func (cs *CacheService) SetInCache(key string, result *SearchResult, expiration time.Duration) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return cs.cacheRepo.Set(key, data, expiration)
}

// FullTextSearch implements SearchServiceInterface
func (s *SearchService) FullTextSearch(query string, juz, surat, page, rowsPerPage int) (*SearchResult, error) {
	// Build cache key
	cacheKey := fmt.Sprintf("search:%s:%d:%d:%d:%d", query, juz, surat, page, rowsPerPage)

	// Try to get from cache
	if result, err := s.cacheService.GetFromCache(cacheKey); err == nil {
		return result, nil
	}

	// Calculate pagination
	offset := (page - 1) * rowsPerPage

	// Build search query
	searchQuery := s.queryBuilder.BuildSearchQuery(query, juz, surat)
	countQuery := s.queryBuilder.BuildCountQuery(query, juz, surat)
	aggregateQuery := s.queryBuilder.BuildAggregateQuery(query, juz, surat)

	// Execute search
	results, err := s.searchRepo.Search(searchQuery, rowsPerPage, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	totalCount, err := s.searchRepo.GetTotalCount(countQuery)
	if err != nil {
		return nil, err
	}

	// Get aggregate results
	count := s.searchRepo.GetAggregate(aggregateQuery)

	// Create result
	result := &SearchResult{
		Results:    results,
		TotalCount: totalCount,
		Count:      count,
	}

	// Cache the result
	if err := s.cacheService.SetInCache(cacheKey, result, 24*time.Hour); err != nil {
		// Log cache error but don't fail the request
		logger.WriteLog(logger.LogLevelError, "Failed to cache search result: %v", err)
	}

	return result, nil
}

// DatabaseQueryBuilder implements QueryBuilder
type DatabaseQueryBuilder struct {
	db *gorm.DB
}

func NewDatabaseQueryBuilder(db *gorm.DB) *DatabaseQueryBuilder {
	return &DatabaseQueryBuilder{db: db}
}

func (b *DatabaseQueryBuilder) BuildSearchQuery(searchTerm string, juz, surat int) *gorm.DB {
	query := b.db.Table("quran_translation").Select("*")
	if searchTerm != "" {
		query = query.Where("MATCH(text_indo, text_arabic) AGAINST(? IN BOOLEAN MODE)", searchTerm)
	}
	if juz > 0 {
		query = query.Where("juz = ?", juz)
	}
	if surat > 0 {
		query = query.Where("surat = ?", surat)
	}
	return query
}

func (b *DatabaseQueryBuilder) BuildCountQuery(searchTerm string, juz, surat int) *gorm.DB {
	return b.BuildSearchQuery(searchTerm, juz, surat)
}

func (b *DatabaseQueryBuilder) BuildAggregateQuery(searchTerm string, juz, surat int) *gorm.DB {
	return b.BuildSearchQuery(searchTerm, juz, surat)
}

// DatabaseSearchRepository implements SearchRepository
type DatabaseSearchRepository struct {
	db *gorm.DB
}

func NewDatabaseSearchRepository(db *gorm.DB) *DatabaseSearchRepository {
	return &DatabaseSearchRepository{db: db}
}

func (r *DatabaseSearchRepository) Search(query *gorm.DB, limit, offset int) ([]*model.AyatDetail, error) {
	var results []*model.AyatDetail
	err := query.Limit(limit).Offset(offset).Find(&results).Error
	return results, err
}

func (r *DatabaseSearchRepository) GetTotalCount(query *gorm.DB) (int64, error) {
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (r *DatabaseSearchRepository) GetAggregate(query *gorm.DB) []*model.CountResult {
	var count []*model.CountResult
	query.Select("surat, COUNT(*) as count").Group("surat").Find(&count)
	return count
}

// RedisCacheRepository implements CacheRepository
type RedisCacheRepository struct {
	client *redis.Client
}

func NewRedisCacheRepository(client *redis.Client) *RedisCacheRepository {
	return &RedisCacheRepository{client: client}
}

func (r *RedisCacheRepository) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RedisCacheRepository) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(key, value, expiration).Err()
}

func normalizePagination(pageNum, limit int) (int, int) {
	if pageNum < 1 {
		pageNum = 1
	}
	if limit < 1 {
		limit = 10
	}
	return pageNum, limit
}

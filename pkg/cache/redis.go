package cache

import (
	"fmt"
	"time"

	"indoquran-api/internal/constants"
	"indoquran-api/pkg/logger"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// CacheConfig defines the interface for cache configuration
type CacheConfig interface {
	GetAddress() string
	GetPassword() string
	GetDB() int
}

// RedisConfig implements CacheConfig
type RedisConfig struct {
	host     string
	port     string
	password string
	db       int
}

// NewRedisConfig creates a new RedisConfig instance
func NewRedisConfig() CacheConfig {
	return &RedisConfig{
		host:     viper.GetString(constants.REDIS_HOST),
		port:     viper.GetString(constants.REDIS_PORT),
		password: viper.GetString(constants.REDIS_PASSWORD),
		db:       viper.GetInt(constants.REDIS_DB),
	}
}

func (c *RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

func (c *RedisConfig) GetPassword() string {
	return c.password
}

func (c *RedisConfig) GetDB() int {
	return c.db
}

// CacheClient defines the interface for cache operations
type CacheClient interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Ping() error
	Close() error
}

// RedisClient implements CacheClient
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new RedisClient instance
func NewRedisClient(config CacheConfig) CacheClient {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetAddress(),
		Password: config.GetPassword(),
		DB:       config.GetDB(),
	})

	return &RedisClient{client: client}
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(key, value, expiration).Err()
}

func (r *RedisClient) Ping() error {
	_, err := r.client.Ping().Result()
	return err
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

// CacheManager manages the cache client lifecycle
type CacheManager struct {
	client CacheClient
}

// NewCacheManager creates a new CacheManager instance
func NewCacheManager(client CacheClient) *CacheManager {
	return &CacheManager{client: client}
}

// Init initializes the cache connection
func (m *CacheManager) Init() error {
	if err := m.client.Ping(); err != nil {
		logger.WriteLog(logger.LogLevelFatal, "Failed to connect to Redis: %s", err)
		return err
	}
	logger.WriteLog(logger.LogLevelInfo, "Connected to Redis successfully")
	return nil
}

// GetClient returns the cache client
func (m *CacheManager) GetClient() CacheClient {
	return m.client
}

// Close closes the cache connection
func (m *CacheManager) Close() error {
	return m.client.Close()
}

// Global cache manager instance
var cacheManager *CacheManager

// InitRedis initializes the Redis client
func InitRedis() {
	config := NewRedisConfig()
	client := NewRedisClient(config)
	cacheManager = NewCacheManager(client)

	if err := cacheManager.Init(); err != nil {
		logger.WriteLog(logger.LogLevelFatal, "Failed to initialize Redis: %s", err)
	}
}

// GetRedis returns the Redis client
func GetRedis() CacheClient {
	return cacheManager.GetClient()
}

// SetRedis sets the Redis client instance - used for testing
func SetRedis(client CacheClient) {
	cacheManager = NewCacheManager(client)
}

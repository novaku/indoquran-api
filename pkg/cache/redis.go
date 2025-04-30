package cache

import (
	"fmt"

	"indoquran-api/internal/constants"
	"indoquran-api/pkg/logger"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var redisClient *redis.Client

// GetRedisConfig returns Redis configuration
func GetRedisConfig() *redis.Options {
	addr := fmt.Sprintf("%s:%s", viper.GetString(constants.REDIS_HOST), viper.GetString(constants.REDIS_PORT))
	return &redis.Options{
		Addr:     addr,
		Password: viper.GetString(constants.REDIS_PASSWORD),
		DB:       viper.GetInt(constants.REDIS_DB),
	}
}

// InitRedis initializes the Redis client
func InitRedis() {
	redisClient = redis.NewClient(GetRedisConfig())

	_, err := redisClient.Ping().Result()
	if err != nil {
		logger.WriteLog(logger.LogLevelFatal, "Failed to connect to Redis: %s", err)
	}

	logger.WriteLog(logger.LogLevelInfo, "Connected to Redis: %s", redisClient.Options().Addr)
}

// GetRedis returns the Redis client
func GetRedis() *redis.Client {
	return redisClient
}

// SetRedis sets the Redis client instance - used for testing
func SetRedis(client *redis.Client) {
	redisClient = client
}

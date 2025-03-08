package cache

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"indoquran.web.id/internal/config"
	"indoquran.web.id/pkg/logger"
)

var redisClient *redis.Client

// InitRedis initializes the Redis client
func InitRedis() {
	addr := fmt.Sprintf("%s:%s", viper.GetString(config.REDIS_HOST), viper.GetString(config.REDIS_PORT))

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: viper.GetString(config.REDIS_PASSWORD),
		DB:       viper.GetInt(config.REDIS_DB),
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		logger.WriteLog(logger.LogLevelFatal, "Failed to connect to Redis: %s, connection config:%#v", err, addr)
	}

	logger.WriteLog(logger.LogLevelInfo, "Connected to Redis: %s", addr)
}

// GetRedis returns the Redis client
func GetRedis() *redis.Client {
	return redisClient
}

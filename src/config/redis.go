package config

import (
	"github.com/go-redis/redis/v8"
)

// GetRedis : configuration for redis server
func GetRedis() (*redis.Options, error) {
	return redis.ParseURL(Config.Cache.URI)
}

package list

import (
	"indoquran-api/pkg/cache"
	"time"

	"github.com/go-redis/redis"
)

type redisCacheService struct {
	client *redis.Client
}

func NewRedisCacheService() CacheService {
	return &redisCacheService{
		client: cache.GetRedis(),
	}
}

func (r *redisCacheService) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *redisCacheService) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(key, value, expiration).Err()
}

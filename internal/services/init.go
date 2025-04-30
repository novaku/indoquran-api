package services

import (
	"indoquran-api/pkg/cache"
)

// InitServices initializes all required services
func InitServices() {
	// Initialize Redis
	cache.InitRedis()
}

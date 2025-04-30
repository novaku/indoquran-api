package services

import (
	"indoquran-api/pkg/cache"
	"indoquran-api/pkg/database"
)

// InitServices initializes all required services
func InitServices() {
	// Initialize Redis
	cache.InitRedis()

	// Initialize the database connection
	database.InitDatabase()
}

package main

import (
	"os"

	v1Router "indoquran-api/api/v1"
	"indoquran-api/internal/cache"
	"indoquran-api/internal/config"
	"indoquran-api/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize Redis
	cache.InitRedis()

	// Initialize the database connection
	database.InitDatabase()

	// Create a new Gin Engine
	r := gin.Default()

	// Set GIN_MODE to debug only for local environment
	if os.Getenv(config.ENV) == config.ENV_LOCAL {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize the router
	v1Router.NewServer(r).RunRouter()
}

package main

import (
	"os"

	"github.com/gin-gonic/gin"
	v1Router "indoquran.web.id/api/v1"
	"indoquran.web.id/internal/cache"
	"indoquran.web.id/internal/config"
	"indoquran.web.id/internal/database"
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

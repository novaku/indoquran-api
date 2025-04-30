package v1

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"indoquran-api/internal/controllers"
	"indoquran-api/internal/services/detail"
	"indoquran-api/pkg/cache"
	"indoquran-api/pkg/database"
	"indoquran-api/pkg/logger"
	"indoquran-api/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type (
	Server struct {
		g *gin.Engine
	}

	IServer interface {
		RunRouter()
	}
)

func NewServer(g *gin.Engine) IServer {
	return &Server{
		g: g,
	}
}

// RunRouter starts the router
func (s *Server) RunRouter() {
	// Initialize logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		logger.WriteLog(logger.LogLevelFatal, "Failed to initialize logger: %#v", err)
	}
	defer zapLogger.Sync()

	g := s.g

	// Initialize database and Redis
	db := database.GetDB()
	redisClient := cache.GetRedis()

	// Initialize services
	detailService := detail.NewDetailService(detail.NewRedisCacheService(), detail.NewGormDatabaseService())

	// Setup middleware
	g.Use(cors.Default()) // Default() Enable CORS for allows all origins
	g.Use(middleware.LoggingMiddleware(db, zapLogger))
	g.Use(middleware.TimeoutMiddleware(time.Minute))
	g.Use(middleware.ContentSecurityPolicy())

	// Setup rate limiting
	rateLimiter := middleware.NewRateLimiter(redisClient, map[string]middleware.RateLimit{
		"/api/v1/search": {
			Requests: 10,
			Period:   time.Second,
		},
		"/api/v1/surat": {
			Requests: 10,
			Period:   time.Second,
		},
		"/api/v1/surat/*": {
			Requests: 10,
			Period:   time.Second,
		},
		"/api/v1/ayat/*": {
			Requests: 10,
			Period:   time.Second,
		},
	})
	g.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Endpoint CSP report handler
	g.POST("/csp-report", middleware.CspReportHandler)

	v1 := g.Group("/api/v1")
	{
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Welcome to indoquran.web.id API v1.0",
			})
		})
		v1.GET("/search", controllers.SearchHandler)
		v1.GET("/surat", controllers.ListSurat)
		v1.GET("/surat/:id", controllers.ListAyatInSurat)
		v1.GET("/ayat/:id", func(c *gin.Context) {
			controllers.DetailAyat(c, detailService)
		})
	}

	port := ":" + viper.GetString("API_PORT")

	srv := &http.Server{
		Addr:    port,
		Handler: g,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WriteLog(logger.LogLevelFatal, "listen: %s\n", err)
		}
	}()

	gracefulShutdown(srv)
}

// gracefulShutdown gracefully shuts down the server
func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.WriteLog(logger.LogLevelInfo, "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WriteLog(logger.LogLevelFatal, "Server forced to shutdown: %s", err)
	}

	logger.WriteLog(logger.LogLevelInfo, "Server exiting")
}

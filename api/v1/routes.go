package v1

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"indoquran-api/internal/controllers"
	"indoquran-api/pkg/logger"
	"indoquran-api/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	// rate limiting
	f := viper.GetString("RATE_LIMIT_FILE")
	rateLimitConfig, err := middleware.LoadConfig(f)
	if err != nil {
		logger.WriteLog(logger.LogLevelFatal, "Failed to load rate limit config: %#v", err)
	}

	g := s.g

	g.Use(cors.Default()) // Default() Enable CORS for allows all origins
	g.Use(middleware.LoggingMiddleware())

	g.Use(middleware.NewRateLimiter(rateLimitConfig).RateLimitMiddleware())
	g.Use(middleware.TimeoutMiddleware(time.Minute))
	g.Use(middleware.ContentSecurityPolicy())

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
		v1.GET("/ayat/:id", controllers.DetailAyat)
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

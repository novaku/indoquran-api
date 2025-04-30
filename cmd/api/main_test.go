package main

import (
	"context"
	"indoquran-api/internal/config"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Interfaces for initialization functions
type (
	configLoader interface {
		Load()
	}

	redisInitializer interface {
		Init()
	}

	databaseInitializer interface {
		Init()
	}
)

// Mock implementations
type mockConfig struct {
	shouldSucceed bool
}

func (m *mockConfig) Load() {
	if !m.shouldSucceed {
		panic("mock config initialization failed")
	}
}

type mockRedis struct {
	shouldSucceed bool
}

func (m *mockRedis) Init() {
	if !m.shouldSucceed {
		panic("mock Redis initialization failed")
	}
}

type mockDatabase struct {
	shouldSucceed bool
}

func (m *mockDatabase) Init() {
	if !m.shouldSucceed {
		panic("mock database initialization failed")
	}
}

// Test initialization function that can be controlled
func testInit(cfg configLoader, rds redisInitializer, db databaseInitializer) {
	cfg.Load()
	rds.Init()
	db.Init()

	if os.Getenv(config.ENV) == config.ENV_LOCAL {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func TestMainInitialization(t *testing.T) {
	tests := []struct {
		name          string
		environment   string
		shouldSucceed bool
	}{
		{
			name:          "Local environment initialization",
			environment:   config.ENV_LOCAL,
			shouldSucceed: true,
		},
		{
			name:          "Production environment initialization",
			environment:   "production",
			shouldSucceed: true,
		},
		{
			name:          "Initialization failure",
			environment:   config.ENV_LOCAL,
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			cfg := &mockConfig{shouldSucceed: tt.shouldSucceed}
			rds := &mockRedis{shouldSucceed: tt.shouldSucceed}
			db := &mockDatabase{shouldSucceed: tt.shouldSucceed}

			// Set environment
			origEnv := os.Getenv(config.ENV)
			os.Setenv(config.ENV, tt.environment)
			defer func() {
				if origEnv != "" {
					os.Setenv(config.ENV, origEnv)
				} else {
					os.Unsetenv(config.ENV)
				}
			}()

			// Save original gin mode
			origMode := gin.Mode()
			defer gin.SetMode(origMode)

			// Test initialization
			if !tt.shouldSucceed {
				assert.Panics(t, func() {
					testInit(cfg, rds, db)
				})
			} else {
				assert.NotPanics(t, func() {
					testInit(cfg, rds, db)
				})

				// Verify gin mode
				if tt.environment == config.ENV_LOCAL {
					assert.Equal(t, gin.DebugMode, gin.Mode())
				} else {
					assert.Equal(t, gin.ReleaseMode, gin.Mode())
				}
			}
		})
	}
}

func TestEngineCreation(t *testing.T) {
	// Test if gin.Default() creates a valid engine
	engine := gin.Default()
	assert.NotNil(t, engine)
	assert.IsType(t, &gin.Engine{}, engine)

	// Verify middleware setup
	handlers := engine.Handlers
	assert.Greater(t, len(handlers), 0) // Should have at least logger and recovery middleware
}

func TestGracefulShutdown(t *testing.T) {
	// Create test server with short timeouts
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      gin.Default(),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Error during shutdown: %v", err)
	}
}

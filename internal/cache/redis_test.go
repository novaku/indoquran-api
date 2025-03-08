package cache

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitRedis(t *testing.T) {
	// Start miniredis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	tests := []struct {
		name          string
		redisHost     string
		redisPort     string
		redisPassword string
		redisDB       int
		shouldFail    bool
	}{
		{
			name:          "Successful connection",
			redisHost:     "localhost",
			redisPort:     mr.Port(),
			redisPassword: "",
			redisDB:       0,
			shouldFail:    false,
		},
		{
			name:          "Invalid port",
			redisHost:     "localhost",
			redisPort:     "99999",
			redisPassword: "",
			redisDB:       0,
			shouldFail:    true,
		},
		{
			name:          "With password",
			redisHost:     "localhost",
			redisPort:     mr.Port(),
			redisPassword: "testpass",
			redisDB:       1,
			shouldFail:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set viper config for test case
			viper.Set("REDIS_HOST", tt.redisHost)
			viper.Set("REDIS_PORT", tt.redisPort)
			viper.Set("REDIS_PASSWORD", tt.redisPassword)
			viper.Set("REDIS_DB", tt.redisDB)

			if tt.name == "With password" {
				mr.RequireAuth(tt.redisPassword)
			}

			if !tt.shouldFail {
				InitRedis()
				assert.NotNil(t, redisClient)

				// Verify connection
				pong, err := redisClient.Ping().Result()
				assert.NoError(t, err)
				assert.Equal(t, "PONG", pong)
			} else {
				assert.Panics(t, func() {
					InitRedis()
				})
			}
		})
	}
}
func TestGetRedis(t *testing.T) {
	// Start miniredis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	tests := []struct {
		name       string
		setupRedis bool
		expectNil  bool
	}{
		{
			name:       "Returns nil when redis not initialized",
			setupRedis: false,
			expectNil:  true,
		},
		{
			name:       "Returns valid redis client when initialized",
			setupRedis: true,
			expectNil:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset redis client before each test
			redisClient = nil

			if tt.setupRedis {
				viper.Set("REDIS_HOST", "localhost")
				viper.Set("REDIS_PORT", mr.Port())
				viper.Set("REDIS_PASSWORD", "")
				viper.Set("REDIS_DB", 0)
				InitRedis()
			}

			result := GetRedis()

			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				pong, err := result.Ping().Result()
				assert.NoError(t, err)
				assert.Equal(t, "PONG", pong)
			}
		})
	}
}

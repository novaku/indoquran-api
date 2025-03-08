package config

import (
	"os"

	"github.com/spf13/viper"
	"indoquran.web.id/pkg/logger"
)

// LoadConfig loads the configuration from the config file
func LoadConfig() {
	// Set the default environment to "local" if not set
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config/file")

	if err := viper.ReadInConfig(); err != nil {
		if _, notFound := err.(viper.ConfigFileNotFoundError); notFound {
			logger.WriteLog(logger.LogLevelFatal, "Config file not found: %v", err)
		} else {
			logger.WriteLog(logger.LogLevelFatal, "Failed to read config file: %v", err)
		}
	}

	logger.WriteLog(logger.LogLevelInfo, "Config loaded successfully file %s", viper.ConfigFileUsed())
}

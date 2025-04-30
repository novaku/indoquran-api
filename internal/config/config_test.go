package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name         string
		env          string
		setupFiles   bool
		expectedFile string
		expectPanic  bool
	}{
		{
			name:         "Default environment (heroku)",
			env:          "",
			setupFiles:   true,
			expectedFile: "heroku.yaml",
			expectPanic:  false,
		},
		{
			name:         "Custom environment",
			env:          "production",
			setupFiles:   true,
			expectedFile: "production.yaml",
			expectPanic:  false,
		},
		{
			name:         "Missing config directory",
			env:          "test",
			setupFiles:   false,
			expectedFile: "",
			expectPanic:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper config before each test
			viper.Reset()

			// Create temporary config directory and file if needed
			if tt.setupFiles {
				tmpDir := filepath.Join(os.TempDir(), "config_test")
				err := os.MkdirAll(filepath.Join(tmpDir, "internal/config/file"), 0755)
				assert.NoError(t, err)
				defer os.RemoveAll(tmpDir)

				configContent := []byte("test_key: test_value")
				err = os.WriteFile(filepath.Join(tmpDir, "internal/config/file", tt.expectedFile), configContent, 0644)
				assert.NoError(t, err)

				// Change working directory to temp directory
				originalWd, _ := os.Getwd()
				err = os.Chdir(tmpDir)
				assert.NoError(t, err)
				defer os.Chdir(originalWd)
			}

			// Set environment variable
			if tt.env != "" {
				os.Setenv("ENV", tt.env)
				defer os.Unsetenv("ENV")
			} else {
				os.Unsetenv("ENV")
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					LoadConfig()
				})
			} else {
				LoadConfig()
				assert.Contains(t, viper.ConfigFileUsed(), tt.expectedFile)
				assert.Equal(t, "test_value", viper.GetString("test_key"))
			}
		})
	}
}

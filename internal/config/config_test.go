package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		setup    func()
		cleanup  func()
		wantErr  bool
	}{
		{
			name:     "Load config successfully with local environment",
			envValue: "",
			setup: func() {
				os.Unsetenv("ENV")
				viper.Reset()
			},
			cleanup: func() {
				viper.Reset()
			},
			wantErr: false,
		},
		{
			name:     "Load config successfully with custom environment",
			envValue: "docker",
			setup: func() {
				os.Setenv("ENV", "development")
				viper.Reset()
			},
			cleanup: func() {
				os.Unsetenv("ENV")
				viper.Reset()
			},
			wantErr: false,
		},
		{
			name:     "Handle missing config file",
			envValue: "invalid",
			setup: func() {
				os.Setenv("ENV", "invalid")
				viper.Reset()
			},
			cleanup: func() {
				os.Unsetenv("ENV")
				viper.Reset()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()

			LoadConfig()

			if tt.envValue == "" {
				assert.Equal(t, "local.yaml", viper.ConfigFileUsed())
			} else if !tt.wantErr {
				assert.Equal(t, tt.envValue+".yaml", viper.ConfigFileUsed())
			}
		})
	}
}

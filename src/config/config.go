package config

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/indoquran-api/src/config/logger"
	"github.com/jbrodriguez/mlog"
	"gopkg.in/yaml.v2"
)

// Config : config variable
var Config *Struct

func init() {
	logger.InitLogger()

	env := os.Getenv("ENV")
	if env == "" {
		mlog.Warning("ENV variable is not set, set to default development")
		env = "development"
	}

	env = strings.ToLower(strings.TrimSpace(env))
	configFile := fmt.Sprintf("src/config/yaml/%s.yml", env)

	f, err := os.Open(configFile)
	if err != nil {
		mlog.Error(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&Config); err != nil {
		mlog.Error(err)
	}
}

package config

import (
	"fmt"
	"os"

	"bitbucket.org/indoquran-api/src/config/logger"
	"github.com/jbrodriguez/mlog"
	"gopkg.in/yaml.v2"
)

// Config : config variable
var Config *Struct

func init() {
	logger.InitLogger()

	env := "development"
	env = os.Getenv("ENV")
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

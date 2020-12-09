package logger

import (
	"github.com/jbrodriguez/mlog"
)

// InitLogger : initialize logger config
func InitLogger() {
	mlog.Start(mlog.LevelInfo, "logs/app.log")
}

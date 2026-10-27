package helpers

import (
	log "github.com/sirupsen/logrus"
)

func SetupLogLevels() map[string]log.Level {
	logLevels := make(map[string]log.Level, 7)

	logLevels["PANIC"] = log.PanicLevel
	logLevels["FATAL"] = log.FatalLevel
	logLevels["ERROR"] = log.ErrorLevel
	logLevels["WARN"] = log.WarnLevel
	logLevels["INFO"] = log.InfoLevel
	logLevels["DEBUG"] = log.DebugLevel
	logLevels["TRACE"] = log.TraceLevel

	return logLevels
}

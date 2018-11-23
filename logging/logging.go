package logging

import log "github.com/sirupsen/logrus"

// Logger defines basic interface for logging
type Logger interface {
	Errorf(string, ...interface{})
	Infof(string, ...interface{})
	Warningf(string, ...interface{})
}

// Info provides implementation of
// writing Info level message
func Info(title string) {
	log.WithFields(log.Fields{}).Info(title)
}

// Fatal provides implementation of
// writing fatal level message
func Fatal(title string) {
	log.WithFields(log.Fields{}).Fatal(title)
}

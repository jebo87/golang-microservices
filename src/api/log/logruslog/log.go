package logruslog

import (
	"fmt"
	"os"
	"strings"

	"github.com/jebo87/golang-microservices/src/api/config"
	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

func init() {

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.DebugLevel
	}

	Log = &logrus.Logger{
		Out:   os.Stdout,
		Level: level,
	}

	if config.IsProduction() {
		Log.Formatter = &logrus.JSONFormatter{}
	} else {
		Log.Formatter = &logrus.TextFormatter{}
	}

}

func Info(msg string, tags ...string) {
	if Log.Level < logrus.InfoLevel {
		return
	}
	Log.WithFields(parseFields(tags...)).Info(msg)
}
func Error(msg string, err error, tags ...string) {
	if Log.Level < logrus.InfoLevel {
		return
	}
	msg = fmt.Sprintf("%s - Error: %v", msg, err)
	Log.WithFields(parseFields(tags...)).Error(msg)
}
func Debug(msg string, tags ...string) {
	if Log.Level < logrus.InfoLevel {
		return
	}
	Log.WithFields(parseFields(tags...)).Debug(msg)
}

func parseFields(tags ...string) logrus.Fields {
	result := make(logrus.Fields, len(tags))
	for _, tag := range tags {
		els := strings.Split(tag, ":")
		result[strings.TrimSpace(els[0])] = strings.TrimSpace(els[1])
	}
	return result
}

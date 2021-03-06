package zaplog

import (
	"fmt"

	"github.com/jebo87/golang-microservices/src/api/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Zap *zap.Logger
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseColorLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	if config.IsProduction() {
		logConfig.Encoding = "json"
	} else {
		logConfig.Encoding = "console"
	}
	var err error
	Zap, err = logConfig.Build()
	if err != nil {
		panic(err)
	}
}

func Debug(msg string, tags ...zap.Field) {
	Zap.Debug(msg, tags...)
	Zap.Sync()
}
func Info(msg string, tags ...zap.Field) {
	Zap.Info(msg, tags...)
	Zap.Sync()
}
func Error(msg string, err error, tags ...zap.Field) {
	msg = fmt.Sprintf("%s - Error: %v", msg, err)
	Zap.Error(msg, tags...)
	Zap.Sync()
}

func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

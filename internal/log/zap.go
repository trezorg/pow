package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

func NewLocalZapLogger(levelS string) func() {
	level, err := zapcore.ParseLevel(levelS)
	if err != nil {
		panic(fmt.Sprintf("Unknown log level: %s", levelS))
	}
	logger, _ := zap.NewDevelopment(zap.WithCaller(true), zap.AddCallerSkip(1), zap.IncreaseLevel(level))
	log = logger.Sugar()
	return func() {
		_ = log.Sync()
	}
}

func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Infof(template string, args ...interface{}) {
	log.Infof(template, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Debugf(template string, args ...interface{}) {
	log.Debugf(template, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

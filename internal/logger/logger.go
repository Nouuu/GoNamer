package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger
var config zap.Config

func init() {
	config = zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	log = logger.Sugar()
}

func InjectLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if contextLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
		return contextLogger
	}
	InjectLogger(ctx, log)
	return log
}

func SetLoggerLevel(level zapcore.Level) {
	config.Level = zap.NewAtomicLevelAt(level)
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	log = logger.Sugar()
}

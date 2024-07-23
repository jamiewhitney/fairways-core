package logging

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func NewLogger(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableHTMLEscape: true,
	}
	logger.SetLevel(level)

	return logger
}

func NewLoggerFromEnv() *logrus.Logger {
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		return NewLogger(logrus.InfoLevel)
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil
	}

	return NewLogger(level)
}

func FromContext(ctx context.Context) *logrus.Logger {
	if logger, ok := ctx.Value("logger").(*logrus.Logger); ok {
		return logger
	}
	return NewLoggerFromEnv()
}

func WithLogger(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

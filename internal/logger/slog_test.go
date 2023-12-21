package logger_test

import (
	"log/slog"
	"testing"
	"zm/internal/logger"
)

func Test_SLogger(t *testing.T) {
	appLog := logger.NewAppSLogger("test")
	appLog.Info("test")

	customLogger1 := appLog.With(
		slog.String("keyA", "custom1"),
		slog.String("keyB", "custom2"),
	)
	customLogger1.Info("customLogger1 message")

	customLogger2 := appLog.With(
		slog.String("keyD", "custom3"),
		slog.Int64("number", 123),
	)
	customLogger2.Info("customLogger2 message")
}

package logger

import (
	"log/slog"
	"os"

	"github.com/empfaze/golang_grpc/sso/utils"
)

func SetupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case utils.LOCAL:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case utils.DEV:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case utils.PROD:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}

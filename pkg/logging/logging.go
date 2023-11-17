package logging

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	return logger
}

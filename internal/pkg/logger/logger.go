package logger

import (
	"log"
	"log/slog"
	"os"

	"go.uber.org/fx"
)

var Module = fx.Options(fx.Provide(New))

func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error while opening file : %v", err)
	}

	logger := slog.New(slog.NewTextHandler(file, opts))

	return logger
}

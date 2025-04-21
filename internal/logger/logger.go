package logger

import (
	"log/slog"
	"os"

	slogctx "github.com/veqryn/slog-context"
)

var Log *slog.Logger

func Init() {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	handler := slogctx.NewHandler(baseHandler, nil)
	Log = slog.New(handler)
	slog.SetDefault(Log)
}

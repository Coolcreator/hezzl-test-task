package slog

import (
	"os"

	"golang.org/x/exp/slog"
)

var slevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func NewLogger(level string) (logger *slog.Logger) {
	options := &slog.HandlerOptions{Level: toSlogLevel(level)}
	handler := slog.NewTextHandler(os.Stdout, options)
	logger = slog.New(handler)
	return
}

func Error(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func toSlogLevel(level string) (slevel slog.Level) {
	slevel, ok := slevels[level]
	if !ok {
		slevel = slog.LevelDebug
	}
	return
}

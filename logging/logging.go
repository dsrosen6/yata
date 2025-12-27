package logging

import (
	"fmt"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init(path string, debug bool) error {
	var err error
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("opening or creating log file: %w", err)
	}

	lvl := slog.LevelInfo
	if debug {
		lvl = slog.LevelDebug
	}

	Logger = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: lvl,
	}))

	return nil
}

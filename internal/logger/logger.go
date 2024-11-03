package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func GetLogger(level slog.Level, filePath string) (*slog.Logger, error) {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}

	log := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: level}))

	return log, nil
}

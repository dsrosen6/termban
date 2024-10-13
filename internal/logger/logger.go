package logger

import (
	"fmt"
	"log/slog"
	"os"
)

var log *slog.Logger

func init() {
	logFile, err := os.OpenFile("termban.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening log file: %v", err))
	}

	log = slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func GetLogger() *slog.Logger {
	return log
}

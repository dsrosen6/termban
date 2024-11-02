package logger

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"
)

const (
	logSubDir   string = "/Library/termban/Logs/"
	logFileName string = "termban.log"
)

func GetLogger(level slog.Level) (*slog.Logger, error) {
	usr, err := getCurrentUser()
	if err != nil {
		return nil, fmt.Errorf("getCurrentUser: %w", err)
	}

	d, err := makeLogDir(usr)
	if err != nil {
		return nil, fmt.Errorf("makeLogDir: %w", err)
	}

	f := d + logFileName
	logFile, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}

	log := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: level}))

	return log, nil
}

// Create termban log directory if it doesn't exist - also return path as string for later use
func makeLogDir(usr *user.User) (string, error) {
	h := usr.HomeDir
	d := h + logSubDir
	if err := os.MkdirAll(d, 0755); err != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", err)
	}

	return d, nil
}

// Get current user in order to construct log folder path
func getCurrentUser() (*user.User, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("user.Current: %w", err)
	}

	return usr, nil
}

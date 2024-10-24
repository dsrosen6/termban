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

var log *slog.Logger

func init() {
	usr, err := getCurrentUser()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ld, err := makeLogDir(usr)
	if err != nil {
		fmt.Println(err)
	}

	f := ld + logFileName
	logFile, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

// Create termban log directory if it doesn't exist - also return path as string for later use
func makeLogDir(usr *user.User) (string, error) {
	h := usr.HomeDir
	ld := h + logSubDir
	if err := os.MkdirAll(ld, 0755); err != nil {
		return "", fmt.Errorf("could not create termban log folder: %w", err)
	}

	return ld, nil
}

// Get current user in order to construct log folder path
func getCurrentUser() (*user.User, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("could not get current logged in user: %w", err)
	}

	return usr, nil
}

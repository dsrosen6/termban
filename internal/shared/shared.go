package shared

import (
	"fmt"
	"os"
	"os/user"
)

const (
	subDirPath = "/Library/termban" // directory path, which will follow the user's home folder
)

// GetTermbanDir constructs a path for the directory Termban keeps all of its files
func GetTermbanDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("user.Current: %w", err)
	}

	h := u.HomeDir

	d := h + subDirPath

	return d, nil
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

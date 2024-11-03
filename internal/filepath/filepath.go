package filepath

import (
	"fmt"
	"os"
	"os/user"
)

const (
	subDirPath = "/Library/termban" // directory path, which will follow the user's home folder
)

type FilePaths struct {
	MainDir string
	CfgFile string
	DBFile  string
	LogFile string
}

func GetFilePaths() (*FilePaths, error) {
	dirPath, err := getTermbanDir()
	if err != nil {
		return nil, fmt.Errorf("shared.GetTermbanDir: %w", err)
	}

	return &FilePaths{
		MainDir: dirPath,
		CfgFile: fmt.Sprintf("%s/%s", dirPath, "config.json"),
		DBFile:  dirPath,
		LogFile: fmt.Sprintf("%s/%s", dirPath, "termban.log"),
	}, nil
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// GetTermbanDir constructs a path for the directory Termban keeps all of its files
func getTermbanDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("user.Current: %w", err)
	}

	h := u.HomeDir

	d := h + subDirPath

	return d, nil
}

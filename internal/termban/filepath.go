package kanban

import (
	"fmt"
	"os"
	"os/user"
)

const (
	subDirPath = "/Library/ttask" // directory path, which will follow the user's home folder
)

type FilePaths struct {
	MainDir string
	CfgFile string
	DBFile  string
	LogFile string
}

func GetFilePaths() (*FilePaths, error) {
	dirPath, err := getTTaskDir()
	if err != nil {
		return nil, fmt.Errorf("shared.GetTTaskDir: %w", err)
	}

	return &FilePaths{
		MainDir: dirPath,
		CfgFile: fmt.Sprintf("%s/%s", dirPath, "config.json"),
		DBFile:  dirPath,
		LogFile: fmt.Sprintf("%s/%s", dirPath, "ttask.log"),
	}, nil
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func getTTaskDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("user.Current: %w", err)
	}

	h := u.HomeDir

	d := h + subDirPath

	return d, nil
}

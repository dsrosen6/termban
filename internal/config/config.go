package config

import (
	"encoding/json"
	"fmt"
	"github.com/dsrosen6/termban/internal/shared"
	"io"
	"os"
)

type Config struct {
	MColor     string `json:"main_color"`
	SColor     string `json:"secondary_color"`
	BorderType string `json:"border_type"`
	C1Name     string `json:"column_1_name"`
	C2Name     string `json:"column_2_name"`
	C3Name     string `json:"column_3_name"`
}

var DefaultConfig = Config{
	MColor:     "#FFFFFF",
	SColor:     "#00AFFF",
	BorderType: "rounded",
	C1Name:     "To Do",
	C2Name:     "Doing",
	C3Name:     "Done",
}

func LoadConfig() Config {
	c, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config - using default values.\nError: %v", err)
		return DefaultConfig
	}

	return c
}

func loadConfig() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return createConfigFile()
	}

	if !shared.FileExists(filePath) {
		return createConfigFile()
	}

	c := Config{}
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return createConfigFile()
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v", err)
		}
	}(jsonFile)

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return createConfigFile()
	}

	if err := json.Unmarshal(jsonBytes, &c); err != nil {
		return createConfigFile()
	}

	return c, nil
}

func createConfigFile() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return DefaultConfig, fmt.Errorf("getConfigFilePath: %w", err)
	}

	dir, err := shared.GetTermbanDir()
	if err != nil {
		return DefaultConfig, fmt.Errorf("shared.GetTermbanDir: %w", err)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return DefaultConfig, fmt.Errorf("os.MkdirAll: %w", err)
	}

	jsonFile, err := os.Create(filePath)
	if err != nil {
		return DefaultConfig, fmt.Errorf("os.Create: %w", err)
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v", err)
		}
	}(jsonFile)

	jsonBytes, err := json.MarshalIndent(DefaultConfig, "", "  ")
	if err != nil {
		return DefaultConfig, fmt.Errorf("json.Marshal: %w", err)
	}

	if _, err := jsonFile.Write(jsonBytes); err != nil {
		return DefaultConfig, fmt.Errorf("jsonFile.Write: %w", err)
	}

	return DefaultConfig, nil
}

func getConfigFilePath() (string, error) {
	dir, err := shared.GetTermbanDir()
	if err != nil {
		return "", fmt.Errorf("shared.GetTermbanDir: %w", err)
	}

	return fmt.Sprintf("%s/%s", dir, "config.json"), nil
}

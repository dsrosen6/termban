package config

import (
	"encoding/json"
	"fmt"
	"github.com/dsrosen6/termban/internal/filepath"
	"io"
	"log/slog"
	"os"
)

var (
	DefaultMainColor      = "#FFFFFF"
	DefaultSecondaryColor = "#00AFFF"
	DefaultBorderType     = "rounded"
	DefaultColumn1Name    = "To Do"
	DefaultColumn2Name    = "Doing"
	DefaultColumn3Name    = "Done"
)

type Config struct {
	log        *slog.Logger
	filePaths  *filepath.FilePaths
	DBLoc      string `json:"db_location"`
	MColor     string `json:"man_color"`
	SColor     string `json:"secondary_color"`
	BorderType string `json:"border_type"`
	C1Name     string `json:"column_1_name"`
	C2Name     string `json:"column_2_name"`
	C3Name     string `json:"column_3_name"`
}

// Load loads the config file, or creates it if it doesn't exist.
func Load(fp *filepath.FilePaths, log *slog.Logger) (*Config, error) {
	cfg := &Config{log: log, filePaths: fp}
	if err := cfg.load(); err != nil {
		log.Error("load config", "error", err)
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return cfg, nil
}

func (c *Config) load() error {
	c.log.Debug("loading config file")
	if !filepath.FileExists(c.filePaths.CfgFile) {
		c.log.Info("no config file found")
		return c.createDefaultCfg()
	}
	c.log.Debug("config file found", "path", c.filePaths.CfgFile)

	jsonFile, err := os.Open(c.filePaths.CfgFile)
	if err != nil {
		c.log.Error("could not open existing config file", "error", err)
		return c.createDefaultCfg()
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			c.log.Error("could not close existing config file", "error", err)
		}
	}(jsonFile)

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		c.log.Error("could not read existing config file", "error", err)
		return c.createDefaultCfg()
	}

	if err := json.Unmarshal(jsonBytes, &c); err != nil {
		c.log.Error("could not unmarshal existing config file", "error", err)
		return c.createDefaultCfg()
	}

	// Safeguard in case user has no db location set in their existing config
	if c.DBLoc == "" {
		c.log.Warn("no db location found in existing config, using default location")
		c.DBLoc = c.filePaths.DBFile
	}

	return nil
}

func (c *Config) createDefaultCfg() error {
	c.log.Info("creating new default config file")
	if err := os.MkdirAll(c.filePaths.MainDir, 0755); err != nil {
		c.log.Error("could not create main directory", "error", err)
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	jsonFile, err := os.Create(c.filePaths.CfgFile)
	if err != nil {
		c.log.Error("could not create new config file", "error", err)
		return fmt.Errorf("os.Create: %w", err)
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			c.log.Error("could not close new config file", "error", err)
		}
	}(jsonFile)

	c.setDefaults()
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		c.log.Error("could not marshal new config file", "error", err)
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if _, err := jsonFile.Write(jsonBytes); err != nil {
		c.log.Error("could not write to new config file", "error", err)
		return fmt.Errorf("jsonFile.Write: %w", err)
	}

	return nil
}

func (c *Config) setDefaults() {
	c.log.Debug("setting default config values")

	c.DBLoc = c.filePaths.DBFile
	c.MColor = DefaultMainColor
	c.SColor = DefaultSecondaryColor
	c.BorderType = DefaultBorderType
	c.C1Name = DefaultColumn1Name
	c.C2Name = DefaultColumn2Name
	c.C3Name = DefaultColumn3Name
}

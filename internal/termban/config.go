package termban

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
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
	filePaths  *FilePaths
	DBLoc      string `json:"db_location"`
	MColor     string `json:"man_color"`
	SColor     string `json:"secondary_color"`
	BorderType string `json:"border_type"`
	C1Name     string `json:"column_1_name"`
	C2Name     string `json:"column_2_name"`
	C3Name     string `json:"column_3_name"`
}

// Load loads the config file, or creates it if it doesn't exist.
func Load(fp *FilePaths, log *slog.Logger) (*Config, error) {
	cfg := &Config{log: log, filePaths: fp}
	if err := cfg.load(); err != nil {
		log.Error("load config", "error", err)
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return cfg, nil
}

func (c *Config) load() error {
	c.log.Debug("loading config file")
	if !FileExists(c.filePaths.CfgFile) {
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

func (c *Config) Border() lipgloss.Border {
	switch c.BorderType {
	case "normal":
		return lipgloss.NormalBorder()
	case "rounded":
		return lipgloss.RoundedBorder()
	case "thick":
		return lipgloss.ThickBorder()
	case "double":
		return lipgloss.DoubleBorder()
	default:
		return c.Border()
	}
}

// DBLocation returns the location of the database file as it's set in the config file.
func (c *Config) DBLocation() string {
	return c.DBLoc
}

// MainColor returns the main color as it's set in the config file as a lipgloss.Color type.
func (c *Config) MainColor() lipgloss.Color {
	if c.MColor == "" {
		c.log.Debug("no main color set, using default", "default", DefaultMainColor)
		return lipgloss.Color(DefaultMainColor)
	}

	return lipgloss.Color(c.MColor)
}

// SecondaryColor returns the secondary color as it's set in the config file as a lipgloss.Color type.
func (c *Config) SecondaryColor() lipgloss.Color {
	if c.SColor == "" {
		c.log.Debug("no secondary color set, using default", "default", DefaultSecondaryColor)
		return lipgloss.Color(DefaultSecondaryColor)
	}

	return lipgloss.Color(c.SColor)
}

// Column1Name returns the name of the first column as it's set in the config file.
func (c *Config) Column1Name() string {
	if c.C1Name == "" {
		c.log.Debug("no column 1 name set, using default", "default", DefaultColumn1Name)
		return DefaultColumn1Name
	}
	return c.C1Name
}

// Column2Name returns the name of the second column as it's set in the config file.
func (c *Config) Column2Name() string {
	if c.C2Name == "" {
		c.log.Debug("no column 2 name set, using default", "default", DefaultColumn2Name)
		return DefaultColumn2Name
	}

	return c.C2Name
}

// Column3Name returns the name of the third column as it's set in the config file.
func (c *Config) Column3Name() string {
	if c.C3Name == "" {
		c.log.Debug("no column 3 name set, using default", "default", DefaultColumn3Name)
		return DefaultColumn3Name
	}

	return c.C3Name
}

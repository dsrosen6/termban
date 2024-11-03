package config

import "github.com/charmbracelet/lipgloss"

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

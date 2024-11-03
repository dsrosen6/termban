package config

import "github.com/charmbracelet/lipgloss"

func (c Config) Border() lipgloss.Border {
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
		return DefaultConfig.Border()
	}
}

func (c Config) MainColor() lipgloss.Color {
	if c.MColor == "" {
		return DefaultConfig.MainColor()
	}

	return lipgloss.Color(c.MColor)
}

func (c Config) SecondaryColor() lipgloss.Color {
	if c.SColor == "" {
		return DefaultConfig.SecondaryColor()
	}

	return lipgloss.Color(c.SColor)
}

func (c Config) Column1Name() string {
	if c.C1Name == "" {
		return DefaultConfig.Column1Name()
	}
	return c.C1Name
}

func (c Config) Column2Name() string {
	if c.C2Name == "" {
		return DefaultConfig.Column2Name()
	}

	return c.C2Name
}

func (c Config) Column3Name() string {
	if c.C3Name == "" {
		return DefaultConfig.Column3Name()
	}

	return c.C3Name
}

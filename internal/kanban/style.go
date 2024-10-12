package kanban

import "github.com/charmbracelet/lipgloss"

var (
	regStyle   = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.HiddenBorder())
	focusStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder())
)

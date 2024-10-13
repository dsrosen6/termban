package termban

import "github.com/charmbracelet/lipgloss"

var (
	green = lipgloss.Color("085")
	grey  = lipgloss.Color("243")

	// used to get frame size
	dummyBorder = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
)

func (m *model) colWidth() int {
	return m.availWidth*1/3 - 2
}

func (m *model) colHeight() int {
	return m.availHeight * 9 / 10
}

func (m *model) InvisBorder() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Width(m.availWidth).Height(m.availHeight).MaxHeight(m.availHeight).Margin(0, 0)
}

func (m *model) FocusBorder() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(green).Width(m.colWidth()).Height(m.colHeight()).Padding(0, 1)
}

func (m *model) UnfocusedBorder() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(grey).Width(m.colWidth()).Height(m.colHeight()).Padding(0, 1)
}

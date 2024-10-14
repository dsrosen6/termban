package termban

import "github.com/charmbracelet/lipgloss"

var (
	green = lipgloss.Color("085")
	grey  = lipgloss.Color("243")

	// used to get frame size
	dummyBorder = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
)

func (m *model) colWidth() int {
	return m.availWidth*1/3 - 2
}

func (m *model) colHeight() int {
	return m.availHeight * 8 / 10
}

func (m *model) inputWidth() int {
	return lipgloss.Width(m.listsView()) - 2
}

func (m *model) FocusInputBorder() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(green).Width(m.inputWidth())
}

func (m *model) RegInputBorder() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(grey).Width(m.inputWidth())
}

func (m *model) FocusColumnView() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(green).Width(m.colWidth()).Height(m.colHeight()).Padding(1, 1)
}

func (m *model) RegColumnView() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(grey).Width(m.colWidth()).Height(m.colHeight()).Padding(1, 1)
}

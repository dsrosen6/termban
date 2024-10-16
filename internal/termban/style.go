package termban

import "github.com/charmbracelet/lipgloss"

var (
	green = lipgloss.Color("085")
	grey  = lipgloss.Color("243")

	// used to get frame size
	dummyBorder = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).MarginTop(2)
)

func (m *model) colWidth() int {
	return m.availWidth*1/3 - 2
}

func (m *model) colHeight() int {
	return m.availHeight * 2 / 3
}

func (m *model) inputWidth() int {
	h, _ := m.HiddenBorder().GetFrameSize()
	return lipgloss.Width(m.listsView()) - h
}

func (m *model) HiddenBorder() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Width(m.availWidth).Height(m.availHeight)
}

func (m *model) InputStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).BorderForeground(green).Width(m.inputWidth())
}

func (m *model) FocusColumnView() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(green).Width(m.colWidth()).Height(m.colHeight()).Padding(1, 1)
}

func (m *model) InactiveFocusColumnView() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(grey).Width(m.colWidth()).Height(m.colHeight()).Padding(1, 1)
}

func (m *model) RegColumnView() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).BorderForeground(green).Faint(true).Width(m.colWidth()).Height(m.colHeight()).Padding(1, 1)
}

func centerVertical(s string) string {
	return lipgloss.NewStyle().AlignVertical(lipgloss.Center).Render(s)
}

package termban

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	green = lipgloss.Color("085")
	blue  = lipgloss.Color("039")
	grey  = lipgloss.Color("243")

	// used to get frame size
	dummyBorder = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).MarginLeft(1).MarginRight(1)
)

func (m *model) colWidth() int {
	return m.availWidth*1/3 - 2
}

func (m *model) colHeight() int {
	return m.availHeight * 3 / 4
}

func (m *model) inputWidth() int {
	h, _ := m.HiddenBorder().GetFrameSize()
	return lipgloss.Width(m.listsView()) - h
}

func (m *model) inputHeight() int {
	ws := m.HiddenBorder().GetHeight()
	_, hf := m.HiddenBorder().GetFrameSize()
	_, cf := m.RegColumnView().GetFrameSize()

	ih := ws - hf - cf - m.colHeight()
	return ih
}

func (m *model) ModeColor() lipgloss.Color {
	color := green

	if m.mode == moveMode {
		color = blue
	}

	return color
}

func (m *model) HiddenBorder() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder()).
		Width(m.availWidth).
		Height(m.availHeight)
}

func (m *model) ListStyle() list.Styles {
	var s list.Styles
	s.TitleBar = lipgloss.NewStyle().
		Padding(0).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(m.ModeColor()).
		Align(lipgloss.Center).
		Width(m.colWidth() - 2)

	s.Title = lipgloss.NewStyle().
		Foreground(m.ModeColor()).
		Align(lipgloss.Center)
	return s
}

func (m *model) InputStyle() lipgloss.Style {
	return lipgloss.NewStyle().Width(m.inputWidth()).Height(m.inputHeight()).Padding(1, 1, 0, 1)
}

func (m *model) TemplateColumnView() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder()).
		Width(m.colWidth()).
		Height(m.colHeight()).
		Padding(0, 1)
}

func (m *model) FocusColumnView() lipgloss.Style {
	return m.TemplateColumnView().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.ModeColor())
}

func (m *model) InactiveFocusColumnView() lipgloss.Style {
	return m.TemplateColumnView().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(grey)
}

func (m *model) RegColumnView() lipgloss.Style {
	return m.TemplateColumnView()
}

func (m *model) FullyCenter(s string) string {
	return lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Width(m.availWidth).Height(m.availHeight).Render(s)
}

func CenterHorizontal(s string, width int) string {
	return lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Width(width).Render(s)
}

func CenterVertical(s string, width, height int) string {
	return lipgloss.NewStyle().AlignVertical(lipgloss.Center).Width(width).Height(height).Render(s)
}

func FormTheme() *huh.Theme {
	// TODO: customize this
	t := huh.ThemeBase()
	return t
}

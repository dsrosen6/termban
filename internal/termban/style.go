package termban

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	blue  = lipgloss.Color("039")
	white = lipgloss.Color("#FFFFFF")
	grey  = lipgloss.Color("243")

	// used to get frame size
	dummyBorder = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
)

func (m *model) inputWidth() int {
	h, _ := m.OuterFrame().GetFrameSize()
	return lipgloss.Width(m.listsView()) - h
}

func (m *model) inputHeight() int {
	ws := m.OuterFrame().GetHeight()
	_, hf := m.OuterFrame().GetFrameSize()
	_, cf := m.RegColumnView().GetFrameSize()

	ih := ws - hf - cf - m.colHeight
	return ih
}

func (m *model) OuterFrame() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(m.availWidth).
		Height(m.availHeight)
}

func (m *model) ListStyle() list.Styles {
	var s list.Styles
	s.TitleBar = lipgloss.NewStyle().
		Padding(0).
		Border(m.border, false, false, true, false).
		BorderForeground(m.mainColor).
		Align(lipgloss.Center).
		Width(m.colWidth)

	s.Title = lipgloss.NewStyle().
		Foreground(m.mainColor).
		Align(lipgloss.Center)

	s.NoItems = lipgloss.NewStyle().
		PaddingLeft(1)

	return s
}

func (m *model) SelectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(m.border, false, false, false, true).
		BorderForeground(m.GetModeColor()).
		Foreground(m.GetModeColor()).
		Padding(0, 0, 0, 1)
}

func (m *model) UnselectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(m.mainColor).
		Faint(true).
		Padding(0, 0, 0, 2)
}

func (m *model) InputStyle() lipgloss.Style {
	return lipgloss.NewStyle().Width(m.inputWidth()).Height(m.inputHeight()).Padding(1, 1, 0, 1)
}

func (m *model) TemplateColumnView() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder()).
		Width(m.colWidth).
		Height(m.colHeight).
		Padding(0, 0)
}

func (m *model) FocusColumnView() lipgloss.Style {
	return m.TemplateColumnView().
		Border(m.border).
		BorderForeground(m.mainColor)
}

func (m *model) InactiveFocusColumnView() lipgloss.Style {
	return m.TemplateColumnView().
		Border(m.border).
		BorderForeground(grey).
		Faint(true)
}

func (m *model) RegColumnView() lipgloss.Style {
	return m.TemplateColumnView().Faint(true)
}

func (m *model) GetModeColor() lipgloss.Color {
	if m.mode == moveMode {
		return m.secondaryColor
	}

	return m.mainColor
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

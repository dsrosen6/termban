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

func (m *Model) inputWidth() int {
	h, _ := m.outerFrame().GetFrameSize()
	return lipgloss.Width(m.listsView()) - h
}

func (m *Model) inputHeight() int {
	ws := m.outerFrame().GetHeight()
	_, hf := m.outerFrame().GetFrameSize()
	_, cf := m.regColumnView().GetFrameSize()

	ih := ws - hf - cf - m.colHeight
	return ih
}

func (m *Model) outerFrame() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(m.availWidth).
		Height(m.availHeight)
}

func (m *Model) customListStyle() list.Styles {
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

func (m *Model) selectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(m.border, false, false, false, true).
		BorderForeground(m.getModeColor()).
		Foreground(m.getModeColor()).
		Padding(0, 0, 0, 1)
}

func (m *Model) unselectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(m.mainColor).
		Faint(true).
		Padding(0, 0, 0, 2)
}

func (m *Model) inputStyle() lipgloss.Style {
	return lipgloss.NewStyle().Width(m.inputWidth()).Height(m.inputHeight()).Padding(1, 1, 0, 1)
}

func (m *Model) templateColumnView() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder()).
		Width(m.colWidth).
		Height(m.colHeight).
		Padding(0, 0)
}

func (m *Model) focusColumnView() lipgloss.Style {
	return m.templateColumnView().
		Border(m.border).
		BorderForeground(m.mainColor)
}

func (m *Model) inactiveFocusColumnView() lipgloss.Style {
	return m.templateColumnView().
		Border(m.border).
		BorderForeground(grey).
		Faint(true)
}

func (m *Model) regColumnView() lipgloss.Style {
	return m.templateColumnView().Faint(true)
}

func (m *Model) getModeColor() lipgloss.Color {
	if m.mode == moveMode {
		return m.secondaryColor
	}

	return m.mainColor
}

func formTheme() *huh.Theme {
	// TODO: customize this
	t := huh.ThemeBase()
	return t
}

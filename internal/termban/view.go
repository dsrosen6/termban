package termban

import "github.com/charmbracelet/lipgloss"

func (m *model) fullView() string {
	return m.HiddenBorder().Render(lipgloss.JoinVertical(lipgloss.Center, m.listsView(), m.inputView()))
}

func (m *model) inputView() string {
	if m.mode == inputMode {
		return m.InputStyle().Render(centerVertical(m.inputForm.View()))
	}

	return m.InputStyle().Render("")
}

func (m *model) listsView() string {
	focusStyle := m.getFocusColumnStyle()

	var views []string
	for i, list := range m.lists {
		if TaskStatus(i) == m.focused {
			views = append(views, focusStyle.Render(list.View()))
		} else {
			views = append(views, m.RegColumnView().Render(list.View()))
		}
	}

	lv := lipgloss.JoinHorizontal(lipgloss.Top, views...)
	return lv
}

func (m *model) getFocusColumnStyle() lipgloss.Style {
	var style lipgloss.Style

	switch m.mode {
	case listMode:
		style = m.FocusColumnView()
	case inputMode:
		style = m.InactiveFocusColumnView()
	}

	return style
}

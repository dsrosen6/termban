package termban

import "github.com/charmbracelet/lipgloss"

func (m *Model) fullView() string {
	return m.outerFrame().Render(lipgloss.JoinVertical(lipgloss.Center, m.listsView(), m.inputView()))
}

func (m *Model) inputView() string {
	if m.mode == inputMode {
		// log.Debug("putting input form in view")
		return m.inputStyle().Render(m.form.View())
	}

	// log.Debug("rendering view without input form")
	return m.inputStyle().Render("")
}

func (m *Model) listsView() string {
	focusStyle := m.getFocusColumnStyle()

	var views []string
	for i, list := range m.lists {
		if status(i) == m.focused {
			views = append(views, focusStyle.Render(list.View()))
		} else {
			views = append(views, m.regColumnView().Render(list.View()))
		}
	}

	lv := lipgloss.JoinHorizontal(lipgloss.Top, views...)
	return lv
}

func (m *Model) getFocusColumnStyle() lipgloss.Style {
	style := m.focusColumnView()

	if m.mode == inputMode {
		style = m.inactiveFocusColumnView()
	}

	return style
}

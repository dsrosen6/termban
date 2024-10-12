package kanban

import tea "github.com/charmbracelet/bubbletea"

func (m *model) initSetup() tea.Msg {

	return m.GetTasks
}

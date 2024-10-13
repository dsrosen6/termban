package termban

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskStatus int

const (
	ToDo TaskStatus = iota
	Doing
	Done
)

// These are all prepended with "Task" so as to not conflict with the other methods right below it.
// Sure, I didn't need to do this with ID, Description, or Status, but I have clinical OCD.
// I would never, ever stop thinking about it.
type Task struct {
	TaskID    int
	TaskTitle string
	TaskDesc  string
	TaskStatus
}

func (t Task) FilterValue() string { return t.TaskTitle }
func (t Task) Title() string       { return t.TaskTitle }
func (t Task) Description() string { return t.TaskDesc }

func (m *model) NextColumn() {
	if m.focused == Done {
		m.focused = ToDo
	} else {
		m.focused++
	}
}

func (m *model) PrevColumn() {
	if m.focused == ToDo {
		m.focused = Done
	} else {
		m.focused--
	}
}

func (m *model) initLists() tea.Msg {
	log.Debug("initializing lists")
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	titles := []string{"To Do", "Doing", "Done"}
	for i, title := range titles {
		m.lists[i].Title = title
	}

	log.Debug("lists successfully initialized")
	return tea.Msg("ListInit")
}

func (m *model) setListTasks() tea.Msg {
	items := map[TaskStatus][]list.Item{
		ToDo:  {},
		Doing: {},
		Done:  {},
	}

	for _, t := range m.tasks {
		items[t.TaskStatus] = append(items[t.TaskStatus], t)
	}

	for status, itemList := range items {
		m.lists[status].SetItems(itemList)
	}

	return tea.Msg("ListTasksSet")
}

func (m *model) getListStyles() string {
	var views []string
	for i, list := range m.lists {
		if TaskStatus(i) == m.focused {
			views = append(views, m.FocusBorder().Render(list.View()))
		} else {
			views = append(views, m.UnfocusedBorder().Render(list.View()))
		}
	}

	return m.InvisBorder().Render(lipgloss.JoinHorizontal(lipgloss.Left, views...))
}

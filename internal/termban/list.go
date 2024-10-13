package termban

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

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

type TaskStatus int

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

func (m *model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.lists[ToDo].Title = "To Do"
	m.lists[Doing].Title = "Doing"
	m.lists[Done].Title = "Done"
}

func (m *model) setListTasks() {
	todoItems := []list.Item{}
	doingItems := []list.Item{}
	doneItems := []list.Item{}

	for _, t := range m.tasks {
		switch t.TaskStatus {
		case ToDo:
			todoItems = append(todoItems, t)
		case Doing:
			doingItems = append(doingItems, t)
		case Done:
			doneItems = append(doneItems, t)
		}
	}

	m.lists[ToDo].SetItems(todoItems)
	m.lists[Doing].SetItems(doingItems)
	m.lists[Done].SetItems(doneItems)
}

func (m *model) getListStyles() string {
	switch m.focused {
	case ToDo:
		return m.InvisBorder().Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.FocusBorder().Render(m.lists[ToDo].View()),
			m.UnfocusedBorder().Render(m.lists[Doing].View()),
			m.UnfocusedBorder().Render(m.lists[Done].View())),
		)

	case Doing:
		return m.InvisBorder().Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.UnfocusedBorder().Render(m.lists[ToDo].View()),
			m.FocusBorder().Render(m.lists[Doing].View()),
			m.UnfocusedBorder().Render(m.lists[Done].View())),
		)
	case Done:
		return m.InvisBorder().Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.UnfocusedBorder().Render(m.lists[ToDo].View()),
			m.UnfocusedBorder().Render(m.lists[Doing].View()),
			m.FocusBorder().Render(m.lists[Done].View())),
		)
	}

	return ""
}

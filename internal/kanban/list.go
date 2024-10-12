package kanban

import (
	"github.com/charmbracelet/bubbles/list"
)

const (
	ToDo status = iota
	Doing
	Done
)

type Task struct {
	id    int
	title string
	desc  string
	status
}

type status int

func (t Task) FilterValue() string { return t.title }
func (t Task) Title() string       { return t.title }
func (t Task) Description() string { return t.desc }

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
		switch t.status {
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

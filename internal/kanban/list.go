package kanban

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	todo status = iota
	doing
	done
)

type task struct {
	id          int
	title       string
	description string
	status
}

type status int

func (t task) Title() string       { return t.title }
func (t task) Description() string { return t.description }
func (t task) Status() status      { return t.status }
func (t task) FilterValue() string { return t.title }

func (m *model) initList() tea.Msg {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 20, 20)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.lists[todo].Title = "To Do"
	m.lists[doing].Title = "Doing"
	m.lists[done].Title = "Done"

	todoItems := []list.Item{}
	doingItems := []list.Item{}
	doneItems := []list.Item{}

	for _, t := range m.tasks {
		switch t.status {
		case todo:
			todoItems = append(todoItems, t)
		case doing:
			doingItems = append(doingItems, t)
		case done:
			doneItems = append(doneItems, t)
		}
	}

	m.lists[todo].SetItems(todoItems)
	m.lists[doing].SetItems(doingItems)
	m.lists[done].SetItems(doneItems)

	return listsInitMsg(true)
}

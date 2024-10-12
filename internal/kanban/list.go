package kanban

import (
	"github.com/charmbracelet/bubbles/list"
)

const (
	ToDo Status = iota
	Doing
	Done
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status
}

type Status int

func (t Task) FilterValue() string { return t.Title }

func (m *model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
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
		switch t.Status {
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

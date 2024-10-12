package main

import (
	"database/sql"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

type task struct {
	id          int
	title       string
	description string
	status
}

type model struct {
	db     *sql.DB
	loaded bool
	tasks  []task
}

type (
	getTasksMsg bool
	errMsg      struct{ err error }
)

func (e errMsg) Error() string { return e.err.Error() }

func newModel() *model {
	var m model
	var err error
	m.db, err = openDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &m
}

func (m *model) Init() tea.Cmd {
	return m.getTasks
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		}
	case getTasksMsg:
		m.loaded = true
		return m, nil
	}

	return m, nil
}

func (m *model) View() string {
	if !m.loaded {
		return "Loading..."
	}

	var s string
	for _, t := range m.tasks {
		s += fmt.Sprintf("%s\n", t.title)
	}

	return s
}

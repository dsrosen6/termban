package kanban

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
)

type model struct {
	db      *sql.DB
	loaded  bool
	tasks   []task
	lists   []list.Model
	focused status
}

type (
	getTasksMsg  bool
	listsInitMsg bool
	errMsg       struct{ err error }
)

func (e errMsg) Error() string { return e.err.Error() }

func NewModel() *model {
	var m model
	var err error
	m.db, err = openDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m.focused = todo
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
		return m, m.initList
	case listsInitMsg:
		m.loaded = true
		return m, nil
	}

	return m, nil
}

func (m *model) View() string {
	if !m.loaded {
		return "Loading..."
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.lists[todo].View(),
		m.lists[doing].View(),
		m.lists[done].View(),
	)
}

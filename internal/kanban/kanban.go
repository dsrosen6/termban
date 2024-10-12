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
	db          *sql.DB
	fullyLoaded bool
	listLoaded  bool
	tasksLoaded bool
	tasks       []Task
	lists       []list.Model
	focused     Status
}

type (
	errMsg struct{ err error }
)

func (e errMsg) Error() string { return e.err.Error() }

func NewModel() *model {
	var m model
	var err error
	m.db, err = OpenDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m.focused = ToDo
	return &m
}

func (m *model) Init() tea.Cmd {
	return m.GetTasks
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.listLoaded {
			h, v := regStyle.GetFrameSize()
			m.initLists(msg.Width-h, msg.Height-v)
			m.listLoaded = true
		}
	}

	if m.listLoaded && m.tasksLoaded {
		m.setListTasks()
		return m, nil
	}

	return m, nil
}

func (m *model) View() string {
	if !m.listLoaded {
		return "Loading..."
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		regStyle.Render(m.lists[ToDo].View()),
		regStyle.Render(m.lists[Doing].View()),
		regStyle.Render(m.lists[Done].View()),
	)
}

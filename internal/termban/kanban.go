package termban

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

type model struct {
	db          *sql.DB
	fullyLoaded bool
	listLoaded  bool
	tasksLoaded bool
	size
	tasks   []Task
	lists   []list.Model
	focused TaskStatus
}

type size struct {
	availWidth, availHeight int
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
		case "left":
			m.PrevColumn()
			return m, nil
		case "right":
			m.NextColumn()
			return m, nil
		}

	case tea.WindowSizeMsg:
		// Set main border frame size
		h, v := dummyBorder.GetFrameSize()
		m.availWidth = msg.Width - h
		m.availHeight = msg.Height - v

		m.initLists(m.colWidth(), m.colHeight())
		m.setListTasks()
		m.listLoaded = true
	}

	if !m.fullyLoaded {
		if m.listLoaded && m.tasksLoaded {
			m.setListTasks()
			m.fullyLoaded = true
			return m, nil
		}

	} else {
		var cmd tea.Cmd
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) View() string {
	if !m.fullyLoaded {
		return "Loading..."
	}

	return m.getListStyles()
}

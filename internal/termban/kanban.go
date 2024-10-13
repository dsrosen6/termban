package termban

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/termban/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

const (
	minWidth  = 110
	minHeight = 34
)

var log *slog.Logger

type model struct {
	db           *sql.DB
	fullyLoaded  bool
	tasksLoaded  bool
	listInit     bool
	sizeObtained bool
	tooSmall     bool
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

func init() {
	log = logger.GetLogger()
}

func NewModel() *model {
	log.Debug("creating new model")
	var m model
	var err error

	log.Debug("opening db")
	m.db, err = OpenDB()
	if err != nil {
		log.Error("OpenDB", "error", err)
		fmt.Println(err)
		os.Exit(1)
	}

	log.Debug("setting focused to ToDo")
	m.focused = ToDo

	log.Info("model created")
	return &m
}

func (m *model) Init() tea.Cmd {
	log.Debug("initializing model")
	return tea.Batch(
		m.GetTasks,
		m.initLists,
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			log.Debug("user quit")
			return m, tea.Quit
		case "left":
			log.Debug("user moved left")
			m.PrevColumn()
			return m, nil
		case "right":
			log.Debug("user moved right")
			m.NextColumn()
			return m, nil
		}

	case tea.WindowSizeMsg:
		log.Debug("got window size message")
		h, v := dummyBorder.GetFrameSize()
		m.availWidth = msg.Width - h
		m.availHeight = msg.Height - v
		m.sizeObtained = true
		log.Debug("size obtained", "width", msg.Width, "height", msg.Height, "availWidth", m.availWidth, "availHeight", m.availHeight)

		if msg.Width < minWidth || msg.Height < minHeight {
			log.Debug("window too small")
			m.tooSmall = true
		} else {
			m.tooSmall = false
		}

	case tea.Msg:
		switch msg {
		case "TasksLoaded":
			m.tasksLoaded = true
		case "ListInit":
			m.listInit = true
		}
	}

	if !m.fullyLoaded {
		log.Debug("model not fully loaded")
		log.Debug("statuses", "tasksLoaded", m.tasksLoaded, "listInit", m.listInit, "sizeObtained", m.sizeObtained)

		if m.tasksLoaded && m.listInit && m.sizeObtained {
			log.Debug("tasks loaded, list init, and size obtained")
			for i := range m.lists {
				m.lists[i].SetSize(m.colWidth(), m.colHeight())
			}

			m.setListTasks()

			log.Info("model fully loaded")
			m.fullyLoaded = true
			return m, nil
		}

	} else {
		var cmd tea.Cmd
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
		return m, cmd
	}

	log.Debug("returning nil")
	return m, nil
}

func (m *model) View() string {
	if !m.fullyLoaded {
		return "Loading..."
	}

	if m.tooSmall {
		return "Window too small. Please resize."
	}

	return m.getListStyles()
}

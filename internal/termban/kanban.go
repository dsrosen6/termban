package termban

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/dsrosen6/termban/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

var log *slog.Logger

type model struct {
	db           *sql.DB
	fullyLoaded  bool
	tasksLoaded  bool
	listInit     bool
	sizeObtained bool
	mode
	size
	tasks     []Task
	lists     []list.Model
	focused   TaskStatus
	inputForm *huh.Form
}

type mode int

const (
	listMode mode = iota
	inputMode
)

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

func NewInputForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Placeholder("Title").
				Key("TaskTitle"),
		),
		huh.NewGroup(
			huh.NewInput().
				Placeholder("Description").
				Key("TaskDesc"),
		),
	).WithShowHelp(false)
}

func NewModel() *model {
	db, err := OpenDB()
	if err != nil {
		log.Error("OpenDB", "error", err)
		fmt.Println(err)
		os.Exit(1)
	}

	log.Info("model created")
	return &model{
		db:        db,
		mode:      listMode,
		focused:   ToDo,
		inputForm: NewInputForm(),
	}
}

func (m *model) Init() tea.Cmd {
	log.Debug("initializing model")
	return tea.Batch(
		m.GetTasks,
		m.initLists,
		m.inputForm.Init(),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case listMode:
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
			case "d":
				log.Debug("user deleted task")
				return m, m.deleteTask
			case "a":
				return m, m.setMode(inputMode)
			}

		case inputMode:
			switch msg.String() {
			case "esc":
				return m, m.setMode(listMode)
			}
		}

	case tea.WindowSizeMsg:
		log.Debug("got window size message")
		h, v := dummyBorder.GetFrameSize()
		m.availWidth = msg.Width - h
		m.availHeight = msg.Height - v
		m.sizeObtained = true
		log.Debug("size obtained", "width", msg.Width, "height", msg.Height, "availWidth", m.availWidth, "availHeight", m.availHeight)

		// if msg.Width < minWidth || msg.Height < minHeight {
		// 	log.Debug("window too small")
		// 	m.tooSmall = true
		// } else {
		// 	m.tooSmall = false
		// }

	case tea.Msg:
		switch msg {
		case "TasksLoaded":
			m.tasksLoaded = true
		case "ListInit":
			m.listInit = true
		// Sent by Create, Update, and Delete to initiate a refresh
		case "TasksRefreshNeeded":
			return m, m.GetTasks
		// Sent by GetTasks after tasks are loaded
		case "TasksRefreshed":
			// If tasks are loaded, update the lists
			return m, m.setListTasks
		case "ModeSet":
			// TODO: this is a buffer to make sure border colors change. is it necessary???
			return m, nil
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
		return m, nil
	}

	var cmd tea.Cmd

	switch m.mode {
	case listMode:
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	case inputMode:
		var form tea.Model
		form, cmd = m.inputForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.inputForm = f
		}
	}

	if m.inputForm.State == huh.StateCompleted {
		return m, tea.Batch(m.createTask, m.setMode(listMode), m.setListTasks)
	}

	return m, cmd
}

func (m *model) View() string {
	if !m.fullyLoaded {
		return "Loading..."
	}

	return m.fullView()
}

// setMode sets the mode!
func (m *model) setMode(mode mode) tea.Cmd {
	// mode
	return func() tea.Msg {
		m.mode = mode
		return tea.Msg("ModeSet")
	}
}

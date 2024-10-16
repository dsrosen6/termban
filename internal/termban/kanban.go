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
	cmdActive    bool
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
	log.Debug("setting fresh input form")
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
		m.DBGetTasks,
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
				return m, m.PrevColumn
			case "right":
				return m, m.NextColumn
			case "d":
				log.Debug("user deleted task")
				return m, m.deleteTask
			case "a":
				log.Debug("user pressed a to add task")
				return m, m.resetForm
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
		case "FullyLoaded":
			m.fullyLoaded = true
		case "TasksRefreshNeeded":
			log.Debug("task refresh needed")
			return m, m.DBGetTasks
		case "TasksRefreshed":
			// If tasks are loaded, update the lists
			log.Debug("tasks refreshed")
			m.cmdActive = false
			return m, m.setListTasks
		case "ListTasksSet":
			return m, m.setMode(listMode)
		case "FormInit":
			return m, tea.Batch(
				m.setMode(inputMode),
				m.inputForm.Init(),
			)
		}
	}

	if !m.fullyLoaded {
		if m.tasksLoaded && m.listInit && m.sizeObtained {
			for i := range m.lists {
				m.lists[i].SetSize(m.colWidth(), m.colHeight())
			}

			log.Info("model fully loaded")
			return m, tea.Batch(
				m.setFullyLoaded,
				m.setListTasks,
			)
		}
		log.Debug("init tasks not done",
			"tasksLoaded", m.tasksLoaded,
			"listInit", m.listInit,
			"sizeObtained", m.sizeObtained)

		return m, nil
	}

	var cmd tea.Cmd

	switch m.mode {
	case listMode:
		log.Debug("updating focused list", "listStatus", m.focused)
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)

	case inputMode:
		if m.inputForm != nil {
			// log.Debug("huh form state", "state", m.inputForm.State)
			switch m.inputForm.State {

			case huh.StateNormal:
				var form tea.Model
				// log.Debug("updating input form")
				form, cmd = m.inputForm.Update(msg)
				if f, ok := form.(*huh.Form); ok {
					m.inputForm = f
				}

			case huh.StateCompleted:
				if !m.cmdActive {
					log.Debug("setting cmdActive to true")
					m.cmdActive = true
					cmd = m.createTask
				}
			}

		} else {
			log.Debug("inputForm is nil")
		}
	}

	// log.Debug("sending cmd", "cmd", cmd)
	return m, cmd

}

func (m *model) View() string {
	if !m.fullyLoaded {
		return "Loading..."
	}

	return m.fullView()
}

func (m *model) resetForm() tea.Msg {
	m.inputForm = NewInputForm()
	log.Debug("form set")
	return tea.Msg("FormInit")
}

// setMode sets the mode!
func (m *model) setMode(mode mode) tea.Cmd {
	// mode
	return func() tea.Msg {
		m.mode = mode
		log.Debug("mode set", "mode", m.mode)
		return tea.Msg("ModeSet")
	}
}

func (m *model) setFullyLoaded() tea.Msg {
	return tea.Msg("FullyLoaded")
}

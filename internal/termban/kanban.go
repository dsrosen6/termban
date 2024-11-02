package termban

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/dsrosen6/termban/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

const (
	listMode mode = iota
	moveMode
	inputMode
)

var log *slog.Logger

type (
	errMsg struct{ err error }
)

type Model struct {
	dbHandler dbHandler
	cmdActive bool
	tasks     []task
	lists     []list.Model
	focused   status
	form      *huh.Form
	loadStatus
	size
	mode
	style
}

type loadStatus struct {
	fullyLoaded bool
	tasksLoaded bool
	listInit    bool
}

type size struct {
	sizeObtained     bool
	fullWindowWidth  int
	fullWindowHeight int
	xFrameSize       int
	yFrameSize       int
	availWidth       int
	availHeight      int
	colWidth         int
	colHeight        int
}

type mode int

type style struct {
	mainColor      lipgloss.Color
	secondaryColor lipgloss.Color
	border         lipgloss.Border
	listStyle      list.Styles
}

func (e errMsg) Error() string { return e.err.Error() }

func init() {
	log = logger.GetLogger()
}

func newInputForm() *huh.Form {
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
	).WithShowHelp(false).WithTheme(formTheme()).WithHeight(1)
}

func NewModel() *Model {
	dbHandler, err := newDBHandler()
	if err != nil {
		log.Error("OpenDB", "error", err)
		fmt.Println(err)
		os.Exit(1)
	}

	log.Info("model created")
	return &Model{
		dbHandler: *dbHandler,
		mode:      listMode,
		focused:   todo,
		form:      newInputForm(),
		style: style{
			mainColor:      white,
			secondaryColor: blue,
			border:         lipgloss.RoundedBorder(),
		},
	}
}

func (m *Model) Init() tea.Cmd {
	log.Debug("initializing model")
	return tea.Batch(
		m.refreshTasks,
		m.initLists,
		m.form.Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Shared keys
		switch msg.String() {
		case " ":
			switch m.mode {
			case listMode:
				return m, m.setMode(moveMode)
			case moveMode:
				return m, m.setMode(listMode)
			default:
				panic("unhandled default case")
			}
		}

		// Mode-dependant keys
		switch m.mode {
		case listMode:
			switch msg.String() {
			case "esc":
				log.Debug("user quit")
				return m, tea.Quit
			case "left":
				return m, m.changeFocusColumn(m.focused.prev())
			case "right":
				return m, m.changeFocusColumn(m.focused.next())
			case "d":
				log.Debug("user deleted task")
				return m, m.deleteTask
			case "a":
				log.Debug("user pressed a to add task")
				return m, m.resetForm
			}
		case moveMode:
			switch msg.String() {
			case "left":
				return m, m.moveTask(m.focused.prev())
			case "right":
				return m, m.moveTask(m.focused.next())
			case "esc":
				return m, m.setMode(listMode)
			}
		case inputMode:
			switch msg.String() {
			case "esc":
				return m, m.setMode(listMode)
			}
		}

	case tea.WindowSizeMsg:
		m.setDimensions(msg)
		m.sizeObtained = true
		m.listStyle = m.customListStyle()
		for i := range m.lists {
			m.lists[i].Styles = m.listStyle
		}

		log.Debug("size obtained", "width", msg.Width, "height", msg.Height, "availWidth", m.availWidth, "availHeight", m.availHeight)

		if m.fullyLoaded {
			for i := range m.lists {
				m.lists[i].SetSize(m.colWidth, m.colHeight)
			}
		}

	case taskMovedMsg:
		log.Debug("task moved", "status", msg.status)
		return m, tea.Batch(
			m.changeFocusColumn(msg.status),
			m.refreshTasks,
		)

	case tea.Msg:
		switch msg {
		case "FullyLoaded":
			m.fullyLoaded = true
		case "TasksRefreshNeeded":
			log.Debug("task refresh needed")
			return m, m.refreshTasks
		case "TasksRefreshed":
			// If tasks are loaded, update the lists
			log.Debug("tasks refreshed")
			m.cmdActive = false
			return m, m.setListTasks
		case "ListTasksSet":
			// If no task is selected, select the last task in the focused list
			// Used for when the last task in the list is deleted
			if m.selectedTask() == (task{}) && len(m.lists[m.focused].Items()) > 0 {
				m.lists[m.focused].Select(len(m.lists[m.focused].Items()) - 1)
			}

			log.Debug("selected task", "task", m.selectedTask())
			if m.mode == inputMode {
				return m, m.setMode(listMode)
			}
			return m, nil
		case "FormInit":
			return m, tea.Batch(
				m.setMode(inputMode),
				m.form.Init(),
			)
		}
	}

	if !m.fullyLoaded {
		if m.tasksLoaded && m.listInit && m.sizeObtained {
			for i := range m.lists {
				m.lists[i].SetSize(m.colWidth, m.colHeight)
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
	case listMode, moveMode:
		log.Debug("updating focused list", "listStatus", m.focused)
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)

	case inputMode:
		if m.form != nil {
			switch m.form.State {

			case huh.StateNormal:
				var form tea.Model
				form, cmd = m.form.Update(msg)
				if f, ok := form.(*huh.Form); ok {
					m.form = f
				}

			case huh.StateCompleted:
				if !m.cmdActive {
					log.Debug("setting cmdActive to true")
					m.cmdActive = true
					cmd = m.insertTask
				}
			default:
				panic("unhandled default case")
			}

		} else {
			log.Debug("inputForm is nil")
		}
	}

	return m, cmd

}

func (m *Model) View() string {
	if !m.fullyLoaded {
		return "Loading..."
	}

	return m.fullView()
}

func (m *Model) resetForm() tea.Msg {
	m.form = newInputForm()
	log.Debug("form set")
	return tea.Msg("FormInit")
}

// setMode sets the mode!
func (m *Model) setMode(mode mode) tea.Cmd {
	// mode
	return func() tea.Msg {
		m.mode = mode
		log.Debug("mode set", "mode", m.mode)
		for i := range m.lists {
			m.lists[i].Styles = m.customListStyle()
			m.setDelegate()
		}
		return tea.Msg("ModeSet")
	}
}

func (m *Model) setFullyLoaded() tea.Msg {
	return tea.Msg("FullyLoaded")
}

func (m *Model) setDimensions(msg tea.WindowSizeMsg) tea.Msg {
	m.xFrameSize, m.yFrameSize = dummyBorder.GetFrameSize()

	m.fullWindowWidth = msg.Width
	m.fullWindowHeight = msg.Height

	m.availWidth = msg.Width - m.xFrameSize
	m.availHeight = msg.Height - m.yFrameSize

	m.colWidth = m.availWidth/3 - m.xFrameSize
	m.colHeight = m.availHeight * 9 / 10

	return nil
}

package termban

import (
	"fmt"
	"github.com/dsrosen6/termban/internal/config"
	"log/slog"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
)

const (
	listMode mode = iota
	moveMode
	inputMode
)

type (
	errMsg struct{ err error }
)

type Model struct {
	log       *slog.Logger
	dbHandler dbHandler
	cmdActive bool
	tasks     []task
	lists     []list.Model
	columnNames
	focused status
	form    *huh.Form
	loadStatus
	size
	mode
	style
}

type columnNames struct {
	Column1Name, Column2Name, Column3Name string
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

func newInputForm() *huh.Form {
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

func NewModel(log *slog.Logger, cfg *config.Config) *Model {
	dbHandler, err := newDBHandler(log, cfg.DBLoc)
	if err != nil {
		log.Error("OpenDB", "error", err)
		fmt.Println(err)
		os.Exit(1)
	}

	log.Debug("model created")
	return &Model{
		log:       log,
		dbHandler: *dbHandler,
		mode:      listMode,
		columnNames: columnNames{
			Column1Name: cfg.Column1Name(),
			Column2Name: cfg.Column2Name(),
			Column3Name: cfg.Column3Name(),
		},
		focused: col1,
		form:    newInputForm(),
		style: style{
			mainColor:      cfg.MainColor(),
			secondaryColor: cfg.SecondaryColor(),
			border:         cfg.Border(),
		},
	}
}

func (m *Model) Init() tea.Cmd {
	m.log.Debug("initializing model")
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
				m.log.Debug("user quit")
				return m, tea.Quit
			case "left":
				return m, m.changeFocusColumn(m.focused.prev())
			case "right":
				return m, m.changeFocusColumn(m.focused.next())
			case "d":
				m.log.Debug("user deleted task")
				return m, m.deleteTask
			case "a":
				m.log.Debug("user pressed a to add task")
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

		m.log.Debug("size obtained", "width", msg.Width, "height", msg.Height, "availWidth", m.availWidth, "availHeight", m.availHeight)

		if m.fullyLoaded {
			for i := range m.lists {
				m.lists[i].SetSize(m.colWidth, m.colHeight)
			}
		}

	case taskMovedMsg:
		m.log.Debug("task moved", "status", msg.status)
		return m, tea.Batch(
			m.changeFocusColumn(msg.status),
			m.refreshTasks,
		)

	case tea.Msg:
		switch msg {
		case "FullyLoaded":
			m.fullyLoaded = true
		case "TasksRefreshNeeded":
			m.log.Debug("task refresh needed")
			return m, m.refreshTasks
		case "TasksRefreshed":
			// If tasks are loaded, update the lists
			m.log.Debug("tasks refreshed")
			m.cmdActive = false
			return m, m.setListTasks
		case "ListTasksSet":
			// If no task is selected, select the last task in the focused list
			// Used for when the last task in the list is deleted
			if m.selectedTask() == (task{}) && len(m.lists[m.focused].Items()) > 0 {
				m.lists[m.focused].Select(len(m.lists[m.focused].Items()) - 1)
			}

			m.log.Debug("selected task", "task", m.selectedTask())
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

			m.log.Debug("model fully loaded")
			return m, tea.Batch(
				m.setFullyLoaded,
				m.setListTasks,
			)
		}
		m.log.Debug("init tasks not done",
			"tasksLoaded", m.tasksLoaded,
			"listInit", m.listInit,
			"sizeObtained", m.sizeObtained)

		return m, nil
	}

	var cmd tea.Cmd

	switch m.mode {
	case listMode, moveMode:
		m.log.Debug("updating focused list", "listStatus", m.focused)
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)

	case inputMode:
		if m.form != nil {

			if m.form.State == huh.StateNormal {
				var form tea.Model
				form, cmd = m.form.Update(msg)
				if f, ok := form.(*huh.Form); ok {
					m.form = f
				}
			}

			if m.form.State == huh.StateCompleted {
				if !m.cmdActive {
					m.log.Debug("setting cmdActive to true")
					m.cmdActive = true
					cmd = m.insertTask
				}
			}

		} else {

			m.log.Debug("inputForm is nil")

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
	m.log.Debug("form set")
	return tea.Msg("FormInit")
}

// setMode sets the mode!
func (m *Model) setMode(mode mode) tea.Cmd {
	// mode
	return func() tea.Msg {
		m.mode = mode
		m.log.Debug("mode set", "mode", m.mode)
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

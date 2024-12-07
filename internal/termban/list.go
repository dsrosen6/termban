package termban

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const (
	col1 status = iota
	col2
	col3
)

type taskMovedMsg struct{ status status }

type task struct {
	id    int
	title string
	desc  string
	status
}

func (t task) FilterValue() string { return t.title }
func (t task) ID() int             { return t.id }
func (t task) Title() string       { return t.title }
func (t task) Description() string { return t.desc }
func (t task) Status() status      { return t.status }

func (s status) next() status {
	if s == col3 {
		return col1
	} else {
		return s + 1
	}
}

func (s status) prev() status {
	if s == col1 {
		return col3
	} else {
		return s - 1
	}
}

func (m *Model) changeFocusColumn(newStatus status) tea.Cmd {
	return func() tea.Msg {
		m.focused = newStatus
		m.setDelegate()
		return tea.Msg("ChangeFocusColumn")
	}
}

func (m *Model) focusedDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false

	d.Styles.SelectedTitle = m.selectedItemStyle()
	d.Styles.NormalTitle = m.unselectedItemStyle()

	return d
}

func (m *Model) normalDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false

	d.Styles.SelectedTitle = m.unselectedItemStyle()
	d.Styles.NormalTitle = m.unselectedItemStyle()

	return d
}

func (m *Model) initLists() tea.Msg {
	m.log.Debug("initializing lists")
	defaultList := list.New([]list.Item{}, m.normalDelegate(), m.colWidth, m.colHeight)
	defaultList.SetShowHelp(false)
	defaultList.Styles = m.customListStyle()

	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.setDelegate()

	titles := []string{m.Column1Name, m.Column2Name, m.Column3Name}

	for i, title := range titles {
		m.lists[i].Title = title
	}

	m.log.Debug("lists successfully initialized")

	m.listInit = true
	return tea.Msg("ListInit")
}

func (m *Model) setListTasks() tea.Msg {
	items := map[status][]list.Item{
		col1: {},
		col2: {},
		col3: {},
	}

	for _, t := range m.tasks {
		items[t.status] = append(items[t.status], t)
	}

	for status, itemList := range items {
		m.lists[status].SetItems(itemList)
		m.lists[status].SetShowStatusBar(false)
	}

	return tea.Msg("ListTasksSet")
}

func (m *Model) setDelegate() tea.Msg {
	for i := range m.lists {
		if i == int(m.focused) {
			m.lists[i].SetDelegate(m.focusedDelegate())
		} else {
			m.lists[i].SetDelegate(m.normalDelegate())
		}
	}
	return nil
}

func (m *Model) refreshTasks() tea.Msg {
	tasks, err := m.dbHandler.getTasks()
	if err != nil {
		return errMsg{err}
	}

	m.tasks = tasks
	if m.fullyLoaded {
		return tea.Msg("TasksRefreshed")
	}

	if !m.tasksLoaded {
		m.tasksLoaded = true
	}

	m.log.Debug("tasks successfully loaded")
	return tea.Msg("TasksLoaded")
}

func (m *Model) insertTask() tea.Msg {
	m.log.Debug("createTask called")
	task := task{
		title:  m.form.GetString("TaskTitle"),
		desc:   m.form.GetString("TaskDesc"),
		status: m.focused,
	}

	if err := m.dbHandler.insertTask(task); err != nil {
		return errMsg{err}
	}

	return tea.Msg("TasksRefreshNeeded")
}

func (m *Model) moveTask(newStatus status) tea.Cmd {
	return func() tea.Msg {
		st := m.selectedTask()
		st.status = newStatus
		if err := m.dbHandler.updateTask(st); err != nil {
			return errMsg{err}
		}

		return taskMovedMsg{newStatus}
	}
}

func (m *Model) deleteTask() tea.Msg {
	if err := m.dbHandler.deleteTask(m.selectedTask().id); err != nil {
		return errMsg{err}
	}

	return tea.Msg("TasksRefreshNeeded")
}

func (m *Model) selectedTask() task {
	st := m.lists[m.focused].SelectedItem()
	if st == nil {
		return task{}
	}

	return st.(task)
}

package termban

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const (
	todo status = iota
	doing
	done
)

type taskMovedMsg struct{ status status }

// These are all prepended with "task" so as to not conflict with the other methods right below it.
// Sure, I didn't need to do this with ID, Description, or Status, but I have clinical OCD.
// I would never, ever stop thinking about it.
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
	if s == done {
		return todo
	} else {
		return s + 1
	}
}

func (s status) prev() status {
	if s == todo {
		return done
	} else {
		return s - 1
	}
}

func (m *model) changeFocusColumn(newStatus status) tea.Cmd {
	return func() tea.Msg {
		m.focused = newStatus
		m.setDelegate()
		return tea.Msg("ChangeFocusColumn")
	}
}

func (m *model) focusedDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false

	d.Styles.SelectedTitle = m.selectedItemStyle()
	d.Styles.NormalTitle = m.unselectedItemStyle()

	return d
}

func (m *model) normalDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false

	d.Styles.SelectedTitle = m.unselectedItemStyle()
	d.Styles.NormalTitle = m.unselectedItemStyle()

	return d
}

func (m *model) initLists() tea.Msg {
	log.Debug("initializing lists")
	defaultList := list.New([]list.Item{}, m.normalDelegate(), m.colWidth, m.colHeight)
	defaultList.SetShowHelp(false)
	defaultList.Styles = m.customListStyle()

	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.setDelegate()

	titles := []string{
		"TO DO", "IN PROGRESS", "DONE"}

	for i, title := range titles {
		m.lists[i].Title = title
	}

	log.Debug("lists successfully initialized")

	m.listInit = true
	return tea.Msg("ListInit")
}

func (m *model) setListTasks() tea.Msg {
	items := map[status][]list.Item{
		todo:  {},
		doing: {},
		done:  {},
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

func (m *model) setDelegate() tea.Msg {
	for i := range m.lists {
		if i == int(m.focused) {
			m.lists[i].SetDelegate(m.focusedDelegate())
		} else {
			m.lists[i].SetDelegate(m.normalDelegate())
		}
	}
	return nil
}

func (m *model) refreshTasks() tea.Msg {
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

	log.Debug("tasks successfully loaded")
	return tea.Msg("TasksLoaded")
}

func (m *model) insertTask() tea.Msg {
	log.Debug("createTask called")
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

func (m *model) moveTask(newStatus status) tea.Cmd {
	return func() tea.Msg {
		st := m.selectedTask()
		st.status = newStatus
		if err := m.dbHandler.updateTask(st); err != nil {
			return errMsg{err}
		}

		return taskMovedMsg{newStatus}
	}
}

func (m *model) deleteTask() tea.Msg {
	if err := m.dbHandler.deleteTask(m.selectedTask().id); err != nil {
		return errMsg{err}
	}

	return tea.Msg("TasksRefreshNeeded")
}

func (m model) selectedTask() task {
	st := m.lists[m.focused].SelectedItem()
	if st == nil {
		return task{}
	}

	return st.(task)
}

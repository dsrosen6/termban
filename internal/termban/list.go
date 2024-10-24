package termban

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type TaskStatus int

const (
	ToDo TaskStatus = iota
	Doing
	Done
)

type TaskMovedMsg struct{ Status TaskStatus }

// These are all prepended with "Task" so as to not conflict with the other methods right below it.
// Sure, I didn't need to do this with ID, Description, or Status, but I have clinical OCD.
// I would never, ever stop thinking about it.
type Task struct {
	TaskID    int
	TaskTitle string
	TaskDesc  string
	TaskStatus
}

func (t Task) FilterValue() string { return t.TaskTitle }
func (t Task) ID() int             { return t.TaskID }
func (t Task) Title() string       { return t.TaskTitle }
func (t Task) Description() string { return t.TaskDesc }
func (t Task) Status() TaskStatus  { return t.TaskStatus }

func (s TaskStatus) Next() TaskStatus {
	if s == Done {
		return ToDo
	} else {
		return s + 1
	}
}

func (s TaskStatus) Prev() TaskStatus {
	if s == ToDo {
		return Done
	} else {
		return s - 1
	}
}

func (m *model) ChangeFocusColumn(newStatus TaskStatus) tea.Cmd {
	return func() tea.Msg {
		m.focused = newStatus
		m.setDelegate()
		return tea.Msg("ChangeFocusColumn")
	}
}

func (m *model) FocusedDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false

	d.Styles.SelectedTitle = m.SelectedItemStyle()
	d.Styles.NormalTitle = m.UnselectedItemStyle()

	return d
}

func (m *model) NormalDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false

	d.Styles.SelectedTitle = m.UnselectedItemStyle()
	d.Styles.NormalTitle = m.UnselectedItemStyle()

	return d
}

func (m *model) initLists() tea.Msg {
	log.Debug("initializing lists")
	defaultList := list.New([]list.Item{}, m.NormalDelegate(), m.colWidth, m.colHeight)
	defaultList.SetShowHelp(false)
	defaultList.Styles = m.ListStyle()

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
	items := map[TaskStatus][]list.Item{
		ToDo:  {},
		Doing: {},
		Done:  {},
	}

	for _, t := range m.tasks {
		items[t.TaskStatus] = append(items[t.TaskStatus], t)
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
			m.lists[i].SetDelegate(m.FocusedDelegate())
		} else {
			m.lists[i].SetDelegate(m.NormalDelegate())
		}
	}
	return nil
}

func (m *model) refreshTasks() tea.Msg {
	tasks, err := m.dbHandler.DBGetTasks()
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
	task := Task{
		TaskTitle:  m.form.GetString("TaskTitle"),
		TaskDesc:   m.form.GetString("TaskDesc"),
		TaskStatus: m.focused,
	}

	if err := m.dbHandler.DBInsertTask(task); err != nil {
		return errMsg{err}
	}

	return tea.Msg("TasksRefreshNeeded")
}

func (m *model) moveTask(newStatus TaskStatus) tea.Cmd {
	return func() tea.Msg {
		st := m.selectedTask()
		st.TaskStatus = newStatus
		if err := m.dbHandler.DBUpdateTask(st); err != nil {
			return errMsg{err}
		}

		return TaskMovedMsg{newStatus}
	}
}

func (m *model) deleteTask() tea.Msg {
	if err := m.dbHandler.DBDeleteTask(m.selectedTask().TaskID); err != nil {
		return errMsg{err}
	}

	return tea.Msg("TasksRefreshNeeded")
}

func (m model) selectedTask() Task {
	st := m.lists[m.focused].SelectedItem()
	if st == nil {
		return Task{}
	}

	return st.(Task)
}

package termban

import (
	"database/sql"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./tasks.db") // TODO: Permanent location
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT,
		status INTEGER
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not exec db: %w", err)
	}

	return db, nil
}

func (m *model) DBInsertTask(task Task) error {
	log.Debug("new task received",
		"title", task.TaskTitle,
		"desc", task.TaskDesc,
		"status", task.TaskStatus)
	stmt, err := m.db.Prepare("INSERT INTO tasks(title, description, status) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.TaskTitle, task.TaskDesc, task.TaskStatus)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	log.Info("task added to db",
		"title", task.TaskTitle,
		"desc", task.TaskDesc,
		"status", task.TaskStatus)

	return nil
}

func (m *model) DBGetTasks() tea.Msg {
	log.Debug("getting tasks from db")
	rows, err := m.db.Query("SELECT * FROM tasks")
	if err != nil {
		return errMsg{err}
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.TaskID, &task.TaskTitle, &task.TaskDesc, &task.TaskStatus)
		if err != nil {
			return errMsg{err}
		}
		tasks = append(tasks, task)
	}

	m.tasks = tasks
	if m.fullyLoaded {
		return tea.Msg("TasksRefreshed")
	}

	m.tasksLoaded = true
	log.Debug("tasks successfully loaded")
	return tea.Msg("TasksLoaded")
}

func (m *model) DBUpdateTask(task Task) error {
	stmt, err := m.db.Prepare("UPDATE tasks SET title=?, description=?, status=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.TaskTitle, task.TaskDesc, task.TaskStatus, task.TaskID)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

func (m *model) DBDeleteTask(id int) error {
	stmt, err := m.db.Prepare("DELETE FROM tasks WHERE id=?")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

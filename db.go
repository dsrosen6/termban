package main

import (
	"database/sql"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./tasks.db") // TODO: Permanent location
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT,
		status TEXT
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not exec db: %w", err)
	}

	return db, nil
}

func (m *model) createTask(task task) error {
	stmt, err := m.db.Prepare("INSERT INTO tasks(title, description, status) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.title, task.description, task.status)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

func (m *model) getTasks() tea.Msg {
	rows, err := m.db.Query("SELECT * FROM tasks")
	if err != nil {
		return errMsg{err}
	}
	defer rows.Close()

	tasks := []task{}
	for rows.Next() {
		var task task
		err = rows.Scan(&task.id, &task.title, &task.description, &task.status)
		if err != nil {
			return errMsg{err}
		}
		tasks = append(tasks, task)
	}

	m.tasks = tasks
	return getTasksMsg(true)
}

func (m *model) updateTask(task task) error {
	stmt, err := m.db.Prepare("UPDATE tasks SET title=?, description=?, status=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.title, task.description, task.status, task.id)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

func (m *model) deleteTask(id int) error {
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

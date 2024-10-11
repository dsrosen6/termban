package main

import (
	"database/sql"
	"fmt"
)

type db struct {
	*sql.DB
}

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

func (d *db) createTask(task task) error {
	stmt, err := d.Prepare("INSERT INTO tasks(title, description, status) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.title, task.description, task.status)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

func (d *db) getTasks() ([]task, error) {
	rows, err := d.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("could not query db: %w", err)
	}
	defer rows.Close()

	tasks := []task{}
	for rows.Next() {
		var task task
		err = rows.Scan(&task.id, &task.title, &task.description, &task.status)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (d *db) updateTask(task task) error {
	stmt, err := d.Prepare("UPDATE tasks SET title=?, description=?, status=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.title, task.description, task.status, task.id)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

func (d *db) deleteTask(id int) error {
	stmt, err := d.Prepare("DELETE FROM tasks WHERE id=?")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

package termban

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
)

type dbHandler struct {
	*sql.DB
}

func NewDBHandler() (*dbHandler, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, fmt.Errorf("OpenDB: %w", err)
	}

	return &dbHandler{db}, err
}

func OpenDB() (*sql.DB, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("could not get current user: %w", err)
	}
	h := usr.HomeDir

	dbFolder := fmt.Sprintf("%s/Library/termban/db", h)
	if err := os.MkdirAll(dbFolder, 0755); err != nil {
		return nil, fmt.Errorf("could not create db folder: %w", err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/Library/termban/db/tasks.db", h)) // TODO: Permanent location
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

func (db *dbHandler) DBInsertTask(task Task) error {
	log.Debug("new task received",
		"title", task.TaskTitle,
		"desc", task.TaskDesc,
		"status", task.TaskStatus)
	stmt, err := db.Prepare("INSERT INTO tasks(title, description, status) VALUES(?, ?, ?)")
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

func (db *dbHandler) DBGetTasks() ([]Task, error) {
	log.Debug("getting tasks from db")
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.TaskID, &task.TaskTitle, &task.TaskDesc, &task.TaskStatus)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (db *dbHandler) DBUpdateTask(task Task) error {
	stmt, err := db.Prepare("UPDATE tasks SET title=?, description=?, status=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.TaskTitle, task.TaskDesc, task.TaskStatus, task.TaskID)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

func (db *dbHandler) DBDeleteTask(id int) error {
	stmt, err := db.Prepare("DELETE FROM tasks WHERE id=?")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

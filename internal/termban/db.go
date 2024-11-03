package termban

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type dbHandler struct {
	log *slog.Logger
	*sql.DB
}

func newDBHandler(log *slog.Logger, dbPath string) (*dbHandler, error) {
	db, err := openDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("OpenDB: %w", err)
	}

	return &dbHandler{log, db}, err
}

func openDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/tasks.db", dbPath))
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
		err := db.Close()
		if err != nil {
			return nil, fmt.Errorf("could not close db: %w", err)
		}
		return nil, fmt.Errorf("could not exec db: %w", err)
	}

	return db, nil
}

func (db *dbHandler) insertTask(task task) error {
	db.log.Debug("new task received",
		"title", task.title,
		"desc", task.desc,
		"status", task.status)
	stmt, err := db.Prepare("INSERT INTO tasks(title, description, status) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.title, task.desc, task.status)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	db.log.Debug("task added to db",
		"title", task.title,
		"desc", task.desc,
		"status", task.status)

	return nil
}

func (db *dbHandler) getTasks() ([]task, error) {
	db.log.Debug("getting tasks from db")
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			db.log.Error("could not close rows", "error", err)
		}
	}(rows)

	var tasks []task
	for rows.Next() {
		var task task
		err = rows.Scan(&task.id, &task.title, &task.desc, &task.status)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (db *dbHandler) updateTask(task task) error {
	stmt, err := db.Prepare("UPDATE tasks SET title=?, description=?, status=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.title, task.desc, task.status, task.id)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

func (db *dbHandler) deleteTask(id int) error {
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

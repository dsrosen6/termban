package termban

import (
	"database/sql"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type DBHandler interface {
	insertTask(task task) error
	getTasks() ([]task, error)
	updateTask(task task) error
	deleteTask(id int) error
}

type SQLiteHandler struct {
	log *slog.Logger
	*sql.DB
}

type mongoHandler struct {
	log    *slog.Logger
	client *mongo.Client
}

func NewSQLiteHandler(log *slog.Logger, dbPath string) (*SQLiteHandler, error) {
	db, err := openSQLiteDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("OpenDB: %w", err)
	}

	return &SQLiteHandler{log, db}, err
}

func openSQLiteDB(dbPath string) (*sql.DB, error) {
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

func (sq *SQLiteHandler) insertTask(task task) error {
	sq.log.Debug("new task received",
		"title", task.title,
		"desc", task.desc,
		"status", task.status)
	stmt, err := sq.Prepare("INSERT INTO tasks(title, description, status) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.title, task.desc, task.status)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	sq.log.Debug("task added to db",
		"title", task.title,
		"desc", task.desc,
		"status", task.status)

	return nil
}

func (sq *SQLiteHandler) getTasks() ([]task, error) {
	sq.log.Debug("getting tasks from db")
	rows, err := sq.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			sq.log.Error("could not close rows", "error", err)
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

func (sq *SQLiteHandler) updateTask(task task) error {
	stmt, err := sq.Prepare("UPDATE tasks SET title=?, description=?, status=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(task.title, task.desc, task.status, task.id)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

func (sq *SQLiteHandler) deleteTask(id int) error {
	stmt, err := sq.Prepare("DELETE FROM tasks WHERE id=?")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("could not exec statement: %w", err)
	}

	return nil
}

package main

import (
	"fmt"
	"os"
	"term-kanban/internal/kanban"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(kanban.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

// uncomment below to make a test db
// func main() {
// 	db, err := openDB()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer db.Close()

// 	tasks := []task{
// 		{title: "Write project proposal", status: todo},
// 		{title: "Set up development environment", status: done},
// 		{title: "Design database schema", status: doing},
// 		{title: "Implement authentication", status: todo},
// 		{title: "Create user interface mockups", status: done},
// 		{title: "Write unit tests", status: doing},
// 		{title: "Deploy to staging server", status: todo},
// 	}

// 	for _, t := range tasks {
// 		if err := createTask(t, db); err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 	}

// 	// read tasks
// 	rows, err := db.Query("SELECT * FROM tasks")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var t task
// 		err := rows.Scan(&t.id, &t.title, &t.description, &t.status)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		fmt.Printf("Task: %s Status: %s\n", t.title, t.status)
// 	}

// }

// func createTask(task task, db *sql.DB) error {
// 	stmt, err := db.Prepare("INSERT INTO tasks(title, description, status) VALUES(?, ?, ?)")
// 	if err != nil {
// 		return fmt.Errorf("could not prepare statement: %w", err)
// 	}

// 	_, err = stmt.Exec(task.title, task.description, task.status)
// 	if err != nil {
// 		return fmt.Errorf("could not exec statement: %w", err)
// 	}

// 	return nil
// }

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
// 	db, err := kanban.OpenDB()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer db.Close()

// 	tasks := []kanban.Task{
// 		{Title: "Write project proposal", Status: kanban.ToDo, Description: "Things"},
// 		{Title: "Set up development environment", Status: kanban.Done, Description: "Stuff"},
// 		{Title: "Design database schema", Status: kanban.Doing, Description: "Stuff"},
// 		{Title: "Implement authentication", Status: kanban.ToDo, Description: "Things"},
// 		{Title: "Create user interface mockups", Status: kanban.Done, Description: "Stuff"},
// 		{Title: "Write unit tests", Status: kanban.Doing, Description: "Things"},
// 		{Title: "Deploy to staging server", Status: kanban.ToDo, Description: "Things"},
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
// 		var t kanban.Task
// 		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		fmt.Printf("Task: %s Status: %d\n", t.Title, t.Status)
// 	}

// }

// func createTask(task kanban.Task, db *sql.DB) error {
// 	stmt, err := db.Prepare("INSERT INTO tasks(title, description, status) VALUES(?, ?, ?)")
// 	if err != nil {
// 		return fmt.Errorf("could not prepare statement: %w", err)
// 	}

// 	_, err = stmt.Exec(task.Title, task.Description, task.Status)
// 	if err != nil {
// 		return fmt.Errorf("could not exec statement: %w", err)
// 	}

// 	return nil
// }

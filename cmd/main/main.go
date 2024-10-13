package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/termban/internal/termban"
)

func main() {
	log.Println("Starting Termban")
	p := tea.NewProgram(termban.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

// uncomment below to make a test db
// func main() {
// 	db, err := termban.OpenDB()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer db.Close()

// 	tasks := []termban.Task{
// 		{TaskTitle: "Write project proposal", TaskStatus: termban.ToDo, TaskDesc: "Things"},
// 		{TaskTitle: "Set up development environment", TaskStatus: termban.Done, TaskDesc: "Stuff"},
// 		{TaskTitle: "Design database schema", TaskStatus: termban.Doing, TaskDesc: "Stuff"},
// 		{TaskTitle: "Implement authentication", TaskStatus: termban.ToDo, TaskDesc: "Things"},
// 		{TaskTitle: "Create user interface mockups", TaskStatus: termban.Done, TaskDesc: "Stuff"},
// 		{TaskTitle: "Write unit tests", TaskStatus: termban.Doing, TaskDesc: "Things"},
// 		{TaskTitle: "Deploy to staging server", TaskStatus: termban.ToDo, TaskDesc: "Things"},
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
// 		var t termban.Task
// 		err := rows.Scan(&t.TaskID, &t.TaskTitle, &t.TaskDesc, &t.TaskStatus)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		fmt.Printf("Task: %s Status: %d\n", t.TaskTitle, t.TaskStatus)
// 	}

// }

// func createTask(task termban.Task, db *sql.DB) error {
// 	stmt, err := db.Prepare("INSERT INTO tasks(title, description, status) VALUES(?, ?, ?)")
// 	if err != nil {
// 		return fmt.Errorf("could not prepare statement: %w", err)
// 	}

// 	_, err = stmt.Exec(task.Title, task.Description, task.TaskStatus)
// 	if err != nil {
// 		return fmt.Errorf("could not exec statement: %w", err)
// 	}

// 	return nil
// }

package main

import (
	"fmt"
	"github.com/dsrosen6/ttask/internal/kanban"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fp, err := kanban.GetFilePaths()
	if err != nil {
		fmt.Printf("Error getting file paths: %v", err)
		os.Exit(1)
	}

	lev := slog.LevelInfo
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-d", "--debug":
			lev = slog.LevelDebug
		}
	}

	log, err := kanban.GetLogger(lev, fp.LogFile)
	if err != nil {
		fmt.Printf("Error getting logger: %v", err)
		os.Exit(1)
	}

	cfg, err := kanban.Load(fp, log)
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(kanban.NewModel(log, cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error occured: %v", err)
		os.Exit(1)
	}
}

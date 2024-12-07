package main

import (
	"fmt"
	"github.com/dsrosen6/termban/internal/termban"

	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fp, err := termban.GetFilePaths()
	if err != nil {
		fmt.Printf("Error getting file paths: %v", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(fp.MainDir, 0755); err != nil {
		fmt.Printf("Error creating main directory: %v", err)
		os.Exit(1)
	}
	
	lev := slog.LevelInfo
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-d", "--debug":
			lev = slog.LevelDebug
		}
	}

	log, err := termban.GetLogger(lev, fp.LogFile)
	if err != nil {
		fmt.Printf("Error getting logger: %v", err)
		os.Exit(1)
	}

	cfg, err := termban.Load(fp, log)
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(termban.NewModel(log, cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error occured: %v", err)
		os.Exit(1)
	}
}

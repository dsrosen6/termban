package main

import (
	"fmt"
	"github.com/dsrosen6/termban/internal/config"
	"github.com/dsrosen6/termban/internal/filepath"
	"github.com/dsrosen6/termban/internal/logger"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/termban/internal/termban"
)

func main() {
	fp, err := filepath.GetFilePaths()
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

	log, err := logger.GetLogger(lev, fp.LogFile)
	if err != nil {
		fmt.Printf("Error getting logger: %v", err)
		os.Exit(1)
	}

	cfg, err := config.Load(fp, log)
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

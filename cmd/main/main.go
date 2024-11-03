package main

import (
	"fmt"
	"github.com/dsrosen6/termban/internal/logger"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/termban/internal/termban"
)

func main() {

	logLev := slog.LevelInfo
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-d", "--debug":
			logLev = slog.LevelDebug
		}
	}

	log, err := logger.GetLogger(logLev)
	if err != nil {
		fmt.Printf("Error getting logger: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(termban.NewModel(log), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

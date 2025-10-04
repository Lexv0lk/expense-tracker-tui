package main

import (
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/menu"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	m, err := menu.InitialModel()

	if err != nil {
		log.Fatalf("Error initializing new menu: %v", err)
	} else {
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			log.Fatalf("Error running program: %v", err)
		}
	}
}

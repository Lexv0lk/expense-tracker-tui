package main

import (
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/menu"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(menu.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Ошибка:", err)
	}
}

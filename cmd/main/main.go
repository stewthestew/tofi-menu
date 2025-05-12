package main

import (
	"fmt"
	"tofi/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error happened when opening prog %v", err)
	}
}

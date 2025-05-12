package ui

// MainWindow for searching andd displaying executables

import (
	"fmt"
	"os"
	"tofi/internal/backend"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Header lipgloss.Style
}

type model struct {
	cursorPos    int
	applications []string
	theme        Theme
}

func Log(s string, args ...any) {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", fmt.Sprintf(s, args...))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer f.Close()
	}
}

func NewModel() model {

	v, err := backend.List(os.Getenv("PATH"))
	if err != nil {
		Log("Error: %s", err)
	}

	return model{
		cursorPos:    0,
		applications: v,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	var returnVal string
	for i, v := range m.applications {
		returnVal += fmt.Sprintf("%d. %s\n", i+1, v)
	}
	return fmt.Sprintf("\n%s\n", returnVal)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			if m.cursorPos >= 0 {
				m.cursorPos++
			}

		case "up":
			if m.cursorPos > 0 {
				m.cursorPos--
			}

		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

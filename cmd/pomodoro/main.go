package main

import (
	"fmt"
	"os"

	tuitea "pomodoro/cmd/pomodoro/tui-tea"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := loadConfig()
	m := newModel(cfg, tuitea.NewPanelRenderer())
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

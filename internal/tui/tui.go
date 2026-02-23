package tui

import tea "github.com/charmbracelet/bubbletea"
import "github.com/H-strangeone/todo/internal/cache"

// Run starts the TUI.
func Run(c *cache.Cache) error {
	p := tea.NewProgram(New(c), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
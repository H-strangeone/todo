package tui

import tea "github.com/charmbracelet/bubbletea"

// Update handles keyboard input and window resize events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.err = nil
			}

		case "down", "j":
			if len(m.todos) > 0 && m.cursor < len(m.todos)-1 {
				m.cursor++
				m.err = nil
			}

		case "enter", " ":
			if len(m.todos) > 0 && m.cursor < len(m.todos) {
				todo := m.todos[m.cursor]
				if todo.Completed {
					m.err = m.cache.Uncomplete(todo.ID)
				} else {
					m.err = m.cache.Complete(todo.ID)
				}

				// Refresh todos after mutation
				m.todos = m.cache.All()

				// Keep cursor in bounds
				if m.cursor >= len(m.todos) {
					m.cursor = len(m.todos) - 1
				}
				if m.cursor < 0 {
					m.cursor = 0
				}
			}

		case "r":
			// Reload from storage
			m.err = m.cache.Reload()
			m.todos = m.cache.All()
			m.cursor = 0
		}
	}

	return m, nil
}
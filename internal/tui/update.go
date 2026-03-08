package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/H-strangeone/todo/internal/model"
)

// Update handles keyboard input and window resize events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// Handle different modes
		switch m.mode {
		case ViewList:
			return m.updateList(msg)
		case ViewAddForm:
			return m.updateAddForm(msg)
		case ViewDeleteConfirm:
			return m.updateDeleteConfirm(msg)
		}
	}

	return m, nil
}

// updateList handles input in list view
func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
			m.todos = m.cache.All()
			m.keepCursorInBounds()
		}

	case "a":
		// Open add form
		m.mode = ViewAddForm
		m.resetForm()

	case "d":
		// Open delete confirmation
		if len(m.todos) > 0 {
			m.mode = ViewDeleteConfirm
		}

	case "r":
		// Reload from storage
		m.err = m.cache.Reload()
		m.todos = m.cache.All()
		m.cursor = 0
	}

	return m, nil
}

// updateAddForm handles input in add form
func (m Model) updateAddForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel and return to list
		m.mode = ViewList
		m.err = nil

	case "tab", "shift+tab":
		// Navigate form fields
		if msg.String() == "tab" {
			m.formFocusIndex++
			if m.formFocusIndex > 4 {
				m.formFocusIndex = 0
			}
		} else {
			m.formFocusIndex--
			if m.formFocusIndex < 0 {
				m.formFocusIndex = 4
			}
		}

	case "ctrl+s":
		// Submit form
		return m.submitAddForm()

	case "backspace":
		m.deleteCharFromFocusedField()

	default:
		// Add character to focused field
		if len(msg.String()) == 1 {
			m.addCharToFocusedField(msg.String())
		}
	}

	return m, nil
}

// updateDeleteConfirm handles delete confirmation
func (m Model) updateDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Confirm delete
		if len(m.todos) > 0 && m.cursor < len(m.todos) {
			id := m.todos[m.cursor].ID
			m.err = m.cache.Delete(id)
			m.todos = m.cache.All()
			m.keepCursorInBounds()
		}
		m.mode = ViewList

	case "n", "N", "esc":
		// Cancel delete
		m.mode = ViewList
		m.err = nil
	}

	return m, nil
}

// Helper functions

func (m *Model) keepCursorInBounds() {
	if m.cursor >= len(m.todos) {
		m.cursor = len(m.todos) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *Model) resetForm() {
	m.formTitle = ""
	m.formDescription = ""
	m.formDue = ""
	m.formReminders = ""
	m.formNotify = "system"
	m.formFocusIndex = 0
	m.err = nil
}

func (m *Model) addCharToFocusedField(char string) {
	switch m.formFocusIndex {
	case 0:
		m.formTitle += char
	case 1:
		m.formDescription += char
	case 2:
		m.formDue += char
	case 3:
		m.formReminders += char
	case 4:
		m.formNotify += char
	}
}

func (m *Model) deleteCharFromFocusedField() {
	deleteLastChar := func(s string) string {
		if len(s) > 0 {
			return s[:len(s)-1]
		}
		return s
	}

	switch m.formFocusIndex {
	case 0:
		m.formTitle = deleteLastChar(m.formTitle)
	case 1:
		m.formDescription = deleteLastChar(m.formDescription)
	case 2:
		m.formDue = deleteLastChar(m.formDue)
	case 3:
		m.formReminders = deleteLastChar(m.formReminders)
	case 4:
		m.formNotify = deleteLastChar(m.formNotify)
	}
}
func (m *Model) submitAddForm() (tea.Model, tea.Cmd) {
	// Validate title
	if strings.TrimSpace(m.formTitle) == "" {
		m.err = fmt.Errorf("title cannot be empty")
		return m, nil
	}

	// Parse due date (default to end of today if empty)
	var dueAt time.Time
	if strings.TrimSpace(m.formDue) == "" {
		// Default to end of today
		now := time.Now()
		dueAt = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, time.Local)
	} else {
		var err error
		dueAt, err = parseFormTime(m.formDue)
		if err != nil {
			m.err = fmt.Errorf("invalid due date: %w", err)
			return m, nil
		}
	}

	// Create todo
	todo := model.NewTodo(m.formTitle, dueAt)
	todo.Description = m.formDescription

	// Parse reminders if provided
	if strings.TrimSpace(m.formReminders) != "" {
		reminders, err := parseReminders(m.formReminders)
		if err != nil {
			m.err = fmt.Errorf("invalid reminders: %w", err)
			return m, nil
		}
		todo.Reminders = reminders
	}

	// Parse notify type
	switch strings.ToLower(strings.TrimSpace(m.formNotify)) {
	case "email":
		todo.Notify = model.NotifyEmail
	case "both":
		todo.Notify = model.NotifyBoth
	case "system", "":
		todo.Notify = model.NotifySystem
	default:
		m.err = fmt.Errorf("invalid notify type (use: system, email, or both)")
		return m, nil
	}

	// Add to cache
	if err := m.cache.Add(todo); err != nil {
		m.err = err
		return m, nil
	}

	// Success - return to list
	m.todos = m.cache.All()
	m.mode = ViewList
	m.err = nil

	return m, nil
}

// parseFormTime parses time from form input (more flexible than CLI)
func parseFormTime(s string) (time.Time, error) {
	// Remove any quotes that user might type
	s = strings.Trim(s, "'\"")
	s = strings.ToLower(strings.TrimSpace(s))

	if s == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}

	// Try standard format first
	if t, err := time.Parse("2006-01-02 15:04", s); err == nil {
		return t, nil
	}

	// Try date only (set time to end of day)
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 0, 0, time.Local), nil
	}

	// Handle common shortcuts
	now := time.Now()
	
	switch s {
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, time.Local), nil
	
	case "tomorrow":
		tomorrow := now.Add(24 * time.Hour)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 0, 0, time.Local), nil
	
	default:
		// Try parsing as relative duration (e.g., "2h", "1d")
		if d, err := time.ParseDuration(s); err == nil {
			return now.Add(d), nil
		}
	}

	return time.Time{}, fmt.Errorf("use 'tomorrow', 'today', 'YYYY-MM-DD HH:MM', 'YYYY-MM-DD', or duration like '2h'")
}

// parseReminders parses comma-separated reminder durations
func parseReminders(s string) (model.DurationSlice, error) {
	parts := strings.Split(s, ",")
	reminders := model.DurationSlice{}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		d, err := time.ParseDuration(part)
		if err != nil {
			return nil, fmt.Errorf("invalid duration '%s': %w", part, err)
		}

		reminders = append(reminders, d)
	}

	return reminders, nil
}
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/H-strangeone/todo/internal/model"
)

// View renders the TUI.
func (m Model) View() string {
	if m.width < 60 || m.height < 15 {
		return "Terminal too small. Please resize."
	}

	// Calculate heights
	topHeight := (m.height * 2) / 3
	bottomHeight := m.height - topHeight - 4 // Reserve space for title and help

	// Render components
	title := titleStyle.Render("📝 TODO CLI")
	table := m.renderTable(topHeight)
	details := m.renderDetails(bottomHeight)
	help := m.renderHelp()

	// Join vertically
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		table,
		details,
		help,
	)

	return content
}

// renderTable renders the top pane with task list.
func (m Model) renderTable(maxHeight int) string {
	var rows []string

	// Header
	header := headerStyle.Render(
		fmt.Sprintf("%-5s %-40s %-10s", "ID", "TASK", "STATUS"),
	)
	rows = append(rows, header)

	// Empty state
	if len(m.todos) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Padding(2, 0).
			Render("No todos yet. Add some tasks to get started!")
		return tableBorderStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left, header, emptyMsg),
		)
	}

	// Task rows
	for i, todo := range m.todos {
		row := m.renderRow(i, todo)
		rows = append(rows, row)

		// Limit rows to fit height
		if len(rows) >= maxHeight-2 {
			break
		}
	}

	table := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return tableBorderStyle.Width(m.width - 4).Render(table)
}

// renderRow renders a single task row.
func (m Model) renderRow(index int, todo *model.Todo) string {
	cursor := "  "
	if index == m.cursor {
		cursor = "▶ "
	}

	// Use plain status icon for selected rows to avoid ANSI bleed inside
	// selectedRowStyle's Render call.
	isSelected := index == m.cursor
	var status string
	if isSelected {
		if todo.Completed {
			status = "✓"
		} else if todo.IsOverdue() {
			status = "✗"
		} else {
			status = "○"
		}
	} else {
		status = statusIcon(todo)
	}

	// Format row
	line := fmt.Sprintf("%s%-5d %-40s %s",
		cursor,
		todo.ID,
		truncate(todo.Title, 40),
		status,
	)

	// Apply style
	var style lipgloss.Style
	if index == m.cursor {
		style = selectedRowStyle
	} else if todo.Completed {
		style = completedRowStyle
	} else {
		style = normalRowStyle
	}

	return style.Render(line)
}

// renderDetails renders the bottom pane with task details.
func (m Model) renderDetails(maxHeight int) string {
	if len(m.todos) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(mutedColor)
		return detailBorderStyle.Width(m.width - 4).Render(
			emptyStyle.Render("No task selected"),
		)
	}

	todo := m.todos[m.cursor]

	var details []string

	details = append(details, detailTitleStyle.Render(fmt.Sprintf("Task #%d Details", todo.ID)))

	details = append(details, fmt.Sprintf("%s %s",
		detailLabelStyle.Render("Title:"),
		detailValueStyle.Render(todo.Title),
	))

	details = append(details, fmt.Sprintf("%s %s",
		detailLabelStyle.Render("Created:"),
		detailValueStyle.Render(todo.CreatedAt.Format("2006-01-02 15:04")),
	))

	details = append(details, fmt.Sprintf("%s %s",
		detailLabelStyle.Render("Due:"),
		detailValueStyle.Render(todo.DueAt.Format("2006-01-02 15:04")),
	))

	// Status
	statusText := "ONGOING"
	statusStyle := ongoingStyle
	if todo.Completed {
		statusText = "COMPLETED"
		statusStyle = completedStyle
	} else if todo.IsOverdue() {
		statusText = "OVERDUE"
		statusStyle = overdueStyle
	}
	details = append(details, fmt.Sprintf("%s %s",
		detailLabelStyle.Render("Status:"),
		statusStyle.Render(statusText),
	))

	details = append(details, fmt.Sprintf("%s %s",
		detailLabelStyle.Render("Notify:"),
		detailValueStyle.Render(todo.Notify.String()),
	))

	// Reminders
	if len(todo.Reminders) > 0 {
		reminderStrs := make([]string, len(todo.Reminders))
		for i, r := range todo.Reminders {
			reminderStrs[i] = r.String()
		}
		details = append(details, fmt.Sprintf("%s %s",
			detailLabelStyle.Render("Reminders:"),
			detailValueStyle.Render(strings.Join(reminderStrs, ", ")),
		))
	}

	// Description
	if todo.Description != "" {
		details = append(details, "")
		details = append(details, detailLabelStyle.Render("Description:"))
		details = append(details, detailValueStyle.Render(wrap(todo.Description, m.width-8)))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, details...)
	return detailBorderStyle.Width(m.width - 4).Render(content)
}

// renderHelp renders the help text.
func (m Model) renderHelp() string {
	help := "↑/↓: navigate • enter: toggle • r: reload • q: quit"
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(errorColor)
		help += " • " + errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}
	return helpStyle.Render(help)
}

// statusIcon returns the appropriate status icon for a todo.
// Completed rows must return a plain symbol — completedRowStyle applies
// strikethrough via its own Render call, and nesting a pre-rendered lipgloss
// string inside another Render causes the inner ANSI codes to leak as raw text.
func statusIcon(t *model.Todo) string {
	if t.Completed {
		return "✓"
	}
	if t.IsOverdue() {
		return overdueStyle.Render("✗")
	}
	return ongoingStyle.Render("○")
}

// truncate truncates a string to maxLen with ellipsis.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// wrap wraps text to the specified width.
func wrap(s string, width int) string {
	if len(s) <= width {
		return s
	}
	var result []string
	for len(s) > width {
		result = append(result, s[:width])
		s = s[width:]
	}
	if len(s) > 0 {
		result = append(result, s)
	}
	return strings.Join(result, "\n")
}
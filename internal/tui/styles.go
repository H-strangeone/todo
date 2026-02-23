package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("205")
	secondaryColor = lipgloss.Color("170")
	mutedColor     = lipgloss.Color("241")
	successColor   = lipgloss.Color("42")
	warningColor   = lipgloss.Color("214")
	errorColor     = lipgloss.Color("196")

	// Title
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1)

	// Table header
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	// Table rows
	selectedRowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("237")).
				Foreground(secondaryColor).
				Bold(true)

	normalRowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	completedRowStyle = lipgloss.NewStyle().
    Foreground(successColor).   // add this line
    Strikethrough(true).
    Padding(0, 1)
	// Status icons
	completedStyle = lipgloss.NewStyle().Foreground(successColor)
	overdueStyle   = lipgloss.NewStyle().Foreground(errorColor)
	ongoingStyle   = lipgloss.NewStyle().Foreground(warningColor)

	// Detail pane
	detailTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor).
				Padding(0, 0, 1, 0)

	detailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(secondaryColor).
				Width(15)

	detailValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255"))

	// Borders
	tableBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(1, 2)

	detailBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(1, 2)

	// Help text
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(1, 0, 0, 2)
)
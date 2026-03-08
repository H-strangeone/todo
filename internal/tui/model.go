package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/H-strangeone/todo/internal/cache"
	"github.com/H-strangeone/todo/internal/model"
)

// ViewMode represents different screens in the TUI
type ViewMode int

const (
	ViewList ViewMode = iota
	ViewAddForm
	ViewEditForm
	ViewDeleteConfirm
)

// Model is the TUI state.
type Model struct {
	cache  *cache.Cache
	todos  []*model.Todo
	cursor int
	width  int
	height int
	err    error

	mode ViewMode

	// Form fields (for add/edit)
	formTitle       string
	formDescription string
	formDue         string
	formReminders   string
	formNotify      string
	formFocusIndex  int
}

// New creates a new TUI model.
func New(c *cache.Cache) Model {
	todos := c.All()
	return Model{
		cache:  c,
		todos:  todos,
		cursor: 0,
		width:  80,
		height: 24,
		mode:   ViewList,
	}
}

// Init is called when the TUI starts.
func (m Model) Init() tea.Cmd {
	return nil
}
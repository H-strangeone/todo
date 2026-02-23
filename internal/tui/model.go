package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/H-strangeone/todo/internal/cache"
	"github.com/H-strangeone/todo/internal/model"
)

// Model is the TUI state.
// NOT to be confused with domain model (model.Todo).
type Model struct {
	cache  *cache.Cache
	todos  []*model.Todo
	cursor int
	width  int
	height int
	err    error
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
	}
}

// Init is called when the TUI starts.
func (m Model) Init() tea.Cmd {
	return nil
}
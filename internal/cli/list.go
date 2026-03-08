package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/H-strangeone/todo/internal/model"
)

var (
	listAll       bool
	listCompleted bool
	listOverdue   bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List todos",
	Long: `List todos with various filters.

Examples:
  todo list              # List pending todos
  todo list --all        # List all todos
  todo list --completed  # List completed todos
  todo list --overdue    # List overdue todos`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all todos")
	listCmd.Flags().BoolVarP(&listCompleted, "completed", "c", false, "Show only completed todos")
	listCmd.Flags().BoolVarP(&listOverdue, "overdue", "o", false, "Show only overdue todos")
}

func runList(cmd *cobra.Command, args []string) error {
	var todos []*model.Todo

	switch {
	case listAll:
		todos = todoCache.All()
	case listCompleted:
		todos = todoCache.Completed()
	case listOverdue:
		todos = todoCache.Overdue()
	default:
		todos = todoCache.Pending()
	}

	if len(todos) == 0 {
		fmt.Println("No todos found.")
		return nil
	}

	// Print header
	fmt.Printf("%-5s %-40s %-12s %-20s\n", "ID", "TITLE", "STATUS", "DUE")
	fmt.Println(strings.Repeat("─", 80))

	// Print todos
	for _, todo := range todos {
		status := "○ PENDING"
		if todo.Completed {
			status = "✓ DONE"
		} else if todo.IsOverdue() {
			status = "✗ OVERDUE"
		}

		dueStr := todo.DueAt.Format("Jan 02, 3:04 PM")

		title := todo.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}

		fmt.Printf("%-5d %-40s %-12s %-20s\n", todo.ID, title, status, dueStr)
	}

	fmt.Printf("\nTotal: %d\n", len(todos))

	return nil
}
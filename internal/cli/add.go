package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/H-strangeone/todo/internal/model"
)

var (
	addDue         string
	addDescription string
	addReminders   []string
	addNotify      string
)

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new todo",
	Long: `Add a new todo with optional deadline, description, and reminders.

Examples:
  todo add "Buy groceries"
  todo add "Finish report" --due "tomorrow 5pm"
  todo add "Meeting" --due "2026-02-25 14:00" --remind "1h,30m"
  todo add "Study" --due "Friday 6pm" --notify email`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addDue, "due", "d", "", "Due date (e.g., 'tomorrow 5pm', '2026-02-25 14:00')")
	addCmd.Flags().StringVarP(&addDescription, "description", "D", "", "Task description")
	addCmd.Flags().StringSliceVarP(&addReminders, "remind", "r", []string{}, "Reminders before due (e.g., '2h,30m,1d')")
	addCmd.Flags().StringVarP(&addNotify, "notify", "n", "system", "Notification type: system, email, both")

	addCmd.MarkFlagRequired("due")
}

func runAdd(cmd *cobra.Command, args []string) error {
	title := strings.Join(args, " ")

	// Parse due date
	dueAt, err := parseTime(addDue)
	if err != nil {
		return fmt.Errorf("invalid due date: %w", err)
	}

	// Create todo
	todo := model.NewTodo(title, dueAt)
	todo.Description = addDescription

	// Parse reminders
	if len(addReminders) > 0 {
		reminders := model.DurationSlice{}
		for _, r := range addReminders {
			r = strings.TrimSpace(r)
			d, err := time.ParseDuration(r)
			if err != nil {
				return fmt.Errorf("invalid reminder duration '%s': %w", r, err)
			}
			reminders = append(reminders, d)
		}
		todo.Reminders = reminders
	}

	// Parse notify type
	switch strings.ToLower(addNotify) {
	case "system":
		todo.Notify = model.NotifySystem
	case "email":
		todo.Notify = model.NotifyEmail
	case "both":
		todo.Notify = model.NotifyBoth
	default:
		return fmt.Errorf("invalid notify type '%s' (must be: system, email, both)", addNotify)
	}

	// Add to cache
	if err := todoCache.Add(todo); err != nil {
		return fmt.Errorf("failed to add todo: %w", err)
	}

	fmt.Printf("  Added todo #%d: %s\n", todo.ID, todo.Title)
	fmt.Printf("   Due: %s\n", todo.DueAt.Format("Mon Jan 02, 2006 at 3:04 PM"))
	if len(todo.Reminders) > 0 {
		fmt.Printf("   Reminders: %v\n", formatReminders(todo.Reminders))
	}

	return nil
}
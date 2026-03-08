package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [id]",
	Short: "Mark a todo as done",
	Long: `Mark a todo as completed.

Examples:
  todo done 3
  todo done 5`,
	Args: cobra.ExactArgs(1),
	RunE: runDone,
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

func runDone(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid todo ID: %s", args[0])
	}

	// Check if todo exists
	todo := todoCache.Get(id)
	if todo == nil {
		return fmt.Errorf("todo #%d not found", id)
	}

	if todo.Completed {
		fmt.Printf("⚠️  Todo #%d is already completed\n", id)
		return nil
	}

	// Mark as done
	if err := todoCache.Complete(id); err != nil {
		return fmt.Errorf("failed to mark todo as done: %w", err)
	}

	fmt.Printf("✅ Marked todo #%d as done: %s\n", id, todo.Title)

	return nil
}
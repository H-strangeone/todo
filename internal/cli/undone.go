package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var undoneCmd = &cobra.Command{
	Use:   "undone [id]",
	Short: "Mark a todo as not done",
	Long: `Mark a completed todo as incomplete.

Examples:
  todo undone 3`,
	Args: cobra.ExactArgs(1),
	RunE: runUndone,
}

func init() {
	rootCmd.AddCommand(undoneCmd)
}

func runUndone(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid todo ID: %s", args[0])
	}
	todo := todoCache.Get(id)
	if todo == nil {
		return fmt.Errorf("todo #%d not found", id)
	}
	if !todo.Completed {
		fmt.Printf("Todo #%d is not completed\n", id)
		return nil
	}

	if err := todoCache.Uncomplete(id); err != nil {
		return fmt.Errorf("failed to mark todo as undone: %w", err)
	}

	fmt.Printf("Marked todo #%d as not done: %s\n", id, todo.Title)

	return nil
}
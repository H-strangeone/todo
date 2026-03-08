package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Aliases: []string{"del", "rm"},
	Short:   "Delete a todo",
	Long: `Delete a todo permanently.

Examples:
  todo delete 3
  todo rm 5`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid todo ID: %s", args[0])
	}

	todo := todoCache.Get(id)
	if todo == nil {
		return fmt.Errorf("todo #%d not found", id)
	}

	if err := todoCache.Delete(id); err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	fmt.Printf("🗑️  Deleted todo #%d: %s\n", id, todo.Title)

	return nil
}
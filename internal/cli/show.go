package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/H-strangeone/todo/internal/tui"
)

var showCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"ui", "tui"},
	Short:   "Open the interactive TUI",
	Long:    `Open the interactive terminal user interface for managing todos.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.Run(todoCache); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
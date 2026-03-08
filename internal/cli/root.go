package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/H-strangeone/todo/internal/cache"
	"github.com/H-strangeone/todo/internal/storage"
	"github.com/H-strangeone/todo/internal/tui"
)

var (
	todoCache *cache.Cache
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A CLI todo app with reminders",
	Long: `Todo CLI - Manage your tasks with deadlines, reminders, and notifications.

Run 'todo' or 'todo show' to open the interactive TUI.
Use subcommands for quick CLI operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		// When no subcommand is provided, open TUI
		if err := tui.Run(todoCache); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	storagePath, err := storage.DefaultPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to determine storage path: %v\n", err)
		os.Exit(1)
	}

	store, err := storage.New(storagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize storage: %v\n", err)
		os.Exit(1)
	}

	todoCache, err = cache.New(store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize cache: %v\n", err)
		os.Exit(1)
	}
}
package main

import (
	"fmt"
	"os"
	"github.com/H-strangeone/todo/internal/storage"
	"github.com/H-strangeone/todo/internal/cache"
	"github.com/H-strangeone/todo/internal/tui"
)


func main() {
	// Get storage path
	storagePath, err := storage.DefaultPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to determine storage path: %v\n", err)
		os.Exit(1)
	}
	// Initialize storage
	store, err := storage.New(storagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize storage: %v\n", err)
		os.Exit(1)
	}
	// Initialize cache
	c, err := cache.New(store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize cache: %v\n", err)
		os.Exit(1)
	}
	// Run TUI
	if err := tui.Run(c); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
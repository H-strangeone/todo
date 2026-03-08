package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/H-strangeone/todo/internal/cache"
	"github.com/H-strangeone/todo/internal/model"
	"github.com/H-strangeone/todo/internal/notifier"
	"github.com/H-strangeone/todo/internal/scheduler"
	"github.com/H-strangeone/todo/internal/storage"
)

func main() {
	// Setup storage and cache
	storagePath, _ := storage.DefaultPath()
	store, _ := storage.New(storagePath)
	c, _ := cache.New(store)
	// Add a test todo with reminder in 10 seconds
	testTodo := model.NewTodo("Test reminder", time.Now().Add(30*time.Second))
	testTodo.Reminders = model.DurationSlice{10 * time.Second, 20 * time.Second}
	testTodo.Notify = model.NotifySystem
	if err := c.Add(testTodo); err != nil {
		fmt.Fprintf(os.Stderr, "Error adding todo: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   Added test todo")
	fmt.Println("   Reminders in: 10s, 20s")
	fmt.Println("   Deadline in: 30s")
	fmt.Println("\nScheduler running... (Ctrl+C to stop)")

	// Setup notifier and scheduler
	n := notifier.NewConsole()
	sched := scheduler.New(c, n)

	// Start scheduler
	if err := sched.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting scheduler: %v\n", err)
		os.Exit(1)
	}

	// Wait for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\n Shutting down")
	sched.Stop()
}
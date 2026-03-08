package notifier

import (
	"fmt"
	"time"
	"github.com/H-strangeone/todo/internal/model"
)
//prints notifications to stdout.
type ConsoleNotifier struct{}
func NewConsole() *ConsoleNotifier {
	return &ConsoleNotifier{}
}
func (n *ConsoleNotifier) Notify(todo *model.Todo, event EventType) error {
	timestamp := time.Now().Format("15:04:05")	
	fmt.Printf("\n(notif) [%s] %s: %s\n", timestamp, event.String(), todo.Title)
	fmt.Printf("   Due: %s\n", todo.DueAt.Format("2006-01-02 15:04"))
	
	if event == EventDeadline {
		fmt.Println("  (blud tf)  Deadline reached!")
	}
	fmt.Println()
	return nil
}
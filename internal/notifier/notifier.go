package notifier

import("github.com/H-strangeone/todo/internal/model"

"encoding/json")
//this will be our interface for sending notifications
type EventType int
const (
	EventReminder EventType=iota// beforeduedate
	EventDeadline //due date reached
)
func(e EventType) String() string{
	switch e {
	case EventReminder:
		return "reminder"
	case EventDeadline:
		return "deadline"
	default:
		return "unknown"
	}
}

type Notifier interface{
	Notify(todo *model.Todo, event EventType) error
}
func (e EventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}
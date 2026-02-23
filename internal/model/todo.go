package model

import (
	"time"
	"encoding/json"
	"fmt"
)

type NotifyType int
type DurationSlice []time.Duration

func (ds DurationSlice) MarshalJSON() ([]byte, error) {
	strings := make([]string, len(ds))
	for i, d := range ds {
		strings[i] = d.String()
	}
	return json.Marshal(strings)
}

func (ds *DurationSlice) UnmarshalJSON(data []byte) error {
	var strings []string
	if err := json.Unmarshal(data, &strings); err != nil {
		return err
	}
	*ds = make(DurationSlice, 0, len(strings))
	for _, s := range strings {
		d, err := time.ParseDuration(s)
		if err != nil {
			return fmt.Errorf("invalid duration %q: %w", s, err)
		}
		*ds = append(*ds, d)
	}
	return nil
}
const(
	NotifySystem NotifyType = iota
	NotifyEmail
	NotifyBoth
)
type Todo struct{
	ID int	`json:"id"`
	Title string	`json:"title"`
	Description string 	`json:"description,omitempty"`
	Completed bool	`json:"completed"`
	CreatedAt time.Time	`json:"created_at"`	
	CompletedAt *time.Time	`json:"completed_at,omitempty"`
	DueAt time.Time `json:"due_at"`
	Reminders DurationSlice `json:"reminders,omitempty"`
	Notify NotifyType	`json:"notify"`
}
func(nt NotifyType) String() string{
	switch nt {
	case NotifySystem:
		return "System"
	case NotifyEmail:
		return "Email"
	case NotifyBoth:
		return "Both"
	default:
		return fmt.Sprintf("Unknown(%d)", nt)
	}
}//human readable string representation of NotifyType


func (nt NotifyType) MarshalJSON() ([]byte, error){
	return json.Marshal(nt.String())
}//custom marshaler for NotifyType to return string representation in JSON

func(nt *NotifyType) UnmarshalJSON(data []byte) error{
	var s string
	if err:=json.Unmarshal(data, &s); err != nil{
		return err
	}
	switch s {
	case "System":
		*nt = NotifySystem
	case "Email":
		*nt = NotifyEmail
	case "Both":
		*nt = NotifyBoth
	default:
		return fmt.Errorf("invalid NotifyType: %s", s)
	}
	return nil
}

func(t *Todo) IsOverdue() bool{
	return !t.Completed && time.Now().After(t.DueAt)
}
func(t *Todo) TimeRemaining() time.Duration{
	return time.Until(t.DueAt)
}
func(t*Todo) Complete(){
	t.Completed = true
	now := time.Now()
	t.CompletedAt = &now
}
func(t *Todo) Uncomplete(){
	t.Completed = false
	t.CompletedAt = nil
}
func(t *Todo) Validate() error{
	if t.Title == ""{
		return fmt.Errorf("title cannot be empty")
	}
	// if t.DueAt.Before(time.Now()){
	// 	return fmt.Errorf("due date cannot be in the past")
	// }
	if t.DueAt.IsZero(){
		return fmt.Errorf("due date is required")
	}
	if t.CreatedAt.IsZero(){
	return fmt.Errorf("created_at is required")
	}
	if !t.DueAt.After(t.CreatedAt){
		return fmt.Errorf("due date cannot be after created_at")
	}
	for i, r := range t.Reminders {
		if r < 0 {
			return fmt.Errorf("reminder %d must be non-negative", i)
		}
		if t.DueAt.Add(-r).Before(t.CreatedAt) {
			return fmt.Errorf("reminder %d fires before creation time", i)
		}
	}
	return nil
}
func NewTodo(title string , dueAt time.Time) *Todo{
	return &Todo{
		Title: title,
		DueAt: dueAt,
		CreatedAt: time.Now(),
		Reminders: DurationSlice{},
		Notify: NotifySystem,
		Completed: false,
	}
}
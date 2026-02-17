package model

import (
	"time"
	"encoding/json"
	"fmt"
)

type NotifyType int

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
	Reminders []time.Duration `json:"reminders,omitempty"`
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
	case "system":
		*nt = NotifySystem
	case "email":
		*nt = NotifyEmail
	case "both":
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
	if t.DueAt.Before(time.Now()){
		return fmt.Errorf("due date cannot be in the past")
	}
	if t.DueAt.IsZero(){
		return fmt.Errorf("due date is required")
	}
	if t.CreatedAt.IsZero(){
	return fmt.Errorf("created_at is required")
	}
	if t.DueAt.Before(t.CreatedAt){
		return fmt.Errorf("due date cannot be before created_at")
	}
	for i, r := range t.Reminders{
		if r <= 0{
			return fmt.Errorf("reminder %d must be a positive duration", i)
		}
		if t.DueAt.Add(-r).Before(t.CreatedAt){
			return fmt.Errorf("reminder %d is before creation time",i)
		}
	}
	return nil
}
func NewTodo(title string , dueAt time.Time) *Todo{
	return &Todo{
		Title: title,
		DueAt: dueAt,
		CreatedAt: time.Now(),
		Reminders: []time.Duration{},
		Notify: NotifySystem,
		Completed: false,
	}
}
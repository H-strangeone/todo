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
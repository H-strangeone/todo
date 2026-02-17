package model

import "time"

type NotifyType int

const(
	NotifySystem NotifyType = iota
	NOtifyEmail
	NotifyBoth
)
type Todo struct{
	ID int	`json:"id"`
	Title string	`json:"title"`
	Description string 	`json:"description,omitempty"`
	Completed bool	`json:"completed"`
	CreatedAt time.Time	`json:"created_at"`	
	CompletedAt *time.Time	`json:"completed_at,omitempty"`
	DueAt *time.Time `json:"due_at,omitempty`
	Reminders []time.Duration `json:"reminders,omitempty"`
	Notify NotifyType	`json:"notify"`
}

 

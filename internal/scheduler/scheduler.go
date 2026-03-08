package scheduler
import (
	"fmt"
	"sync"
	"time"
	"github.com/H-strangeone/todo/internal/cache"
	"github.com/H-strangeone/todo/internal/model"
	"github.com/H-strangeone/todo/internal/notifier"
)
//to manage time based notifs for the todos
//rems
//timers are created back on reboot
//overdue are ignored
//cache hi bhagwaan hai
//plug the notifs
type Scheduler struct{
	cache *cache.Cache
	notifier notifier.Notifier
	timers map[int][]*time.Timer
	mu sync.RWMutex
	stopCh chan struct{}
}
func New(c*cache.Cache, n notifier.Notifier) *Scheduler{
	return &Scheduler{
		cache: c,
		notifier: n,
		timers: make(map[int][]*time.Timer),
		stopCh: make(chan struct{}),
	}
}
//schduler must be called with a lock
func (s *Scheduler) scheduleTodo(todo *model.Todo, now time.Time) {
    var timers []*time.Timer
    for _, reminderDuration := range todo.Reminders {
        reminderTime := todo.DueAt.Add(-reminderDuration)
        if reminderTime.After(now) {
            duration := reminderTime.Sub(now)
            t := todo
            timer := time.AfterFunc(duration, func() {
                if err := s.notifier.Notify(t, notifier.EventReminder); err != nil {
                    fmt.Println("notify error:", err)
                }
            }) 
            timers = append(timers, timer)
        }
    }
    if todo.DueAt.After(now) {
        duration := todo.DueAt.Sub(now)
        t := todo
        timer := time.AfterFunc(duration, func() {
            if err := s.notifier.Notify(t, notifier.EventDeadline); err != nil {
                fmt.Println("notify error:", err)
            }
        })
        timers = append(timers, timer)
    }
    s.timers[todo.ID] = timers
}
//load all pending ones,ignores overdue, schedules future rem and dead, handle missed ones
func(s* Scheduler) Start() error{
	s.mu.Lock()
	defer s.mu.Unlock()
	todos:=s.cache.Pending()
	now:= time.Now()
	for _,todo:= range todos{
		if todo.IsOverdue() {
			continue
		}
		s.scheduleTodo(todo,now)
	}
	return nil
}


//stop shuts cheduler
func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.mu.Lock()
	defer s.mu.Unlock()
	// Cancel all pending timers
	for _, timers := range s.timers {
		for _, timer := range timers {
			timer.Stop()
		}
	}
	s.timers = make(map[int][]*time.Timer)
}
func (s *Scheduler) Reschedule() error {
	s.Stop()
	return s.Start()
}
// Unschedule removes all timers for a specific todo.
// Called when a todo is completed or deleted.
func (s *Scheduler) Unschedule(todoID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if timers, exists := s.timers[todoID]; exists {
		for _, timer := range timers {
			timer.Stop()
		}
		delete(s.timers, todoID)
	}
}
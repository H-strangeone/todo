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
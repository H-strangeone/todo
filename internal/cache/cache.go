package cache

import (
	"fmt"
	"sync"

	"github.com/H-strangeone/todo/internal/model"
	"github.com/H-strangeone/todo/internal/storage"
)

type Cache struct {
	todos   map[int]*model.Todo
	storage storage.Storage
	nextID  int
	mu      sync.RWMutex
}

func New(s storage.Storage) (*Cache, error) {
	todos, err := s.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load from storage: %w", err)
		
	}
	cache:= &Cache{
		todos: make(map[int]*model.Todo),
		storage: s,
		nextID: 1,// cause didnt need to persist this as we can count it from persisted cache data using the logic below
	}
	for i:= range todos{
		t:=&todos[i]
		cache.todos[t.ID]=t
		if t.ID>= cache.nextID{
			cache.nextID=t.ID+1
		}
	}
	return cache, nil
}
func (c *Cache) Add(todo *model.Todo) error {//lock, validate, mutate memory, persist, unlock
	c.mu.Lock()
	defer c.mu.Unlock()
	if err:= todo.Validate(); err != nil{
		return fmt.Errorf("invalid todo: %w", err)
	}
	todo.ID = c.nextID
	c.nextID++
	c.todos[todo.ID] = todo
	return c.persist()// self explainable snippet
}

func(c* Cache) persist() error{
	todos:=make([]model.Todo,0,len(c.todos))
	for _,todo :=range c.todos{//key, value
		todos=append(todos,*todo)
	}
	if err:= c.storage.Save(todos); err != nil{
		return fmt.Errorf("failed to save to storage: %w", err)
	}
	return nil
}

func(c* Cache) Update(todo *model.Todo) error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists:= c.todos[todo.ID]; !exists{
		return fmt.Errorf("todo with ID %d does not exist", todo.ID)
	}
	if err:= todo.Validate(); err != nil{
		return fmt.Errorf("invalid todo: %w", err)
	}
	c.todos[todo.ID] = todo
	return c.persist()	
}
func(c* Cache) Delete(id int) error{
	c.mu.Lock()
	defer c.mu.Unlock()	
	if _, exists:= c.todos[id]; !exists{
		return fmt.Errorf("todo with ID %d does not exist", id)
	}
	delete(c.todos, id)
	return c.persist()
}
func(c* Cache) Complete(id int) error{
	c.mu.Lock()
	defer c.mu.Unlock()
	todo, exists:= c.todos[id]
	if !exists{
		return fmt.Errorf("todo with ID %d does not exist", id)
	}
	todo.Complete()
	return c.persist()
}
func(c* Cache) Uncomplete(id int) error{
	c.mu.Lock()
	defer c.mu.Unlock()
	todo, exists:= c.todos[id]
	if !exists{
		return fmt.Errorf("todo with ID %d does not exist", id)
	}
	todo.Uncomplete()
	return c.persist()
}
// now some queries
// lets keep queries read only
func(c *Cache) Get(id int) *model.Todo{
	c.mu.RLock()
	defer c.mu.RUnlock()
	todo, exists:= c.todos[id]
	if !exists{
		return nil
	}
	return todo
}
func (c *Cache) All() []*model.Todo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]*model.Todo, 0, len(c.todos))
	for _, todo := range c.todos {
		result = append(result, todo)
	}
	return result
}


func (c *Cache) Pending() []*model.Todo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var result []*model.Todo
	for _, todo := range c.todos {
		if !todo.Completed {
			result = append(result, todo)
		}
	}
	return result
}
func (c *Cache) Completed() []*model.Todo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var result []*model.Todo
	for _, todo := range c.todos {
		if todo.Completed {
			result = append(result, todo)
		}
	}
	return result
}
func (c *Cache) Overdue() []*model.Todo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var result []*model.Todo
	for _, todo := range c.todos {
		if todo.IsOverdue() {
			result = append(result, todo)
		}
	}
	return result
}
// Reload discards in-memory state and reloads from storage.
// Useful after:
// - External modifications to storage file
// - Server sync operations
// - Corruption recovery
//
// WARNING: Any unsaved in-memory changes will be lost.
func (c *Cache) Reload() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	todos, err := c.storage.Load()
	if err != nil {
		return fmt.Errorf("failed to reload from storage: %w", err)
	}

	// Clear and rebuild
	c.todos = make(map[int]*model.Todo)
	c.nextID = 1

	for i := range todos {
		t := &todos[i]
		c.todos[t.ID] = t
		if t.ID >= c.nextID {
			c.nextID = t.ID + 1
		}
	}

	return nil
}
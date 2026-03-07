package main

import (
	"fmt"
	"time"

	"github.com/H-strangeone/todo/internal/cache"
	"github.com/H-strangeone/todo/internal/model"
	"github.com/H-strangeone/todo/internal/storage"
)

func main() {
	path, err := storage.DefaultPath()
	if err != nil {
		panic(err)
	}
	store, err := storage.New(path)
	if err != nil {
		panic(err)
	}
	c, err := cache.New(store)
	if err != nil {
		panic(err)
	}
	todo := model.NewTodo(
		"Test cache add",
		time.Now().Add(2*time.Hour),
	)
	err = c.Add(todo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Added todo:", todo.ID)
	fmt.Println("\nAll todos:")
	for _, t := range c.All() {
		fmt.Printf("ID=%d | %s | completed=%v\n",t.ID, t.Title, t.Completed)
	}
	err = c.Complete(todo.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nAfter complete:")
	for _, t := range c.All() {
		fmt.Printf("ID=%d | %s | completed=%v\n",t.ID, t.Title, t.Completed)
	}
	err = c.Delete(todo.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nAfter delete:")
	fmt.Println("Todos:", len(c.All()))
}
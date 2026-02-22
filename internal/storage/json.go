package storage 

import(
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"github.com/H-strangeone/todo/internal/model"
)
type Storage interface{
	Load() ([]model.Todo, error)
	Save([]model.Todo) error
}

type JsonStorage struct{
	path string
	mu sync.Mutex
}

func New(path string) (*JsonStorage, error){
	dir:= filepath.Dir(path)
	if err:= os.MkdirAll(dir,0755); err!=nil{
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	return &JsonStorage{path: path}, nil
}
func (s* JsonStorage) Load() ([]model.Todo, error){
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err:= os.ReadFile(s.path)
	if err!=nil{
		if os.IsNotExist(err){
			return []model.Todo{}, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if len(data)==0{
		return []model.Todo{}, nil
	}
	var todos []model.Todo
	if err:= json.Unmarshal(data, &todos); err!=nil{
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}
	for i:=range todos{
		if err:=todos[i].Validate(); err!=nil{
			return nil, fmt.Errorf("invalid todo item at index %d: %w", i, err)
		}
	}
	return todos, nil
}
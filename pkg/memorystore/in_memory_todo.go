package memorystore

import (
	"sync"
	"todolist/pkg/model"
)

//repository implementation for in_memory DB

type TodoMap struct {
	mu    sync.RWMutex
	store map[string]model.Todo
}

// Factory Method to create new Todo
func NewTodoMap() *TodoMap {

	return &TodoMap{
		store: make(map[string]model.Todo),
	}
}

func (t *TodoMap) Create(Todo model.Todo) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.store[Todo.TID]; ok {
		return model.ErrObjectAlreadyExists

	}
	t.store[Todo.TID] = Todo
	// fmt.Println(m.store, "......")
	return nil
}

func (t *TodoMap) Delete(id string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.store[id]; ok {
		delete(t.store, id)
		return nil
	}
	return model.ErrObjectNotFound
}

func (t *TodoMap) Update(tid string, Todo model.Todo) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.store[tid]; ok {
		t.store[tid] = Todo
		return nil
	} else {
		return model.ErrObjectNotFound
	}

}

func (t *TodoMap) GetById(id string) (model.Todo, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.store[id]; ok {
		return t.store[id], nil
	}
	return model.Todo{}, model.ErrObjectNotFound
}

func (t *TodoMap) GetAll() ([]model.Todo, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	todo := make([]model.Todo, 0)
	if len(t.store) == 0 {
		return todo, nil
	}
	for _, v := range t.store {
		todo = append(todo, v)
	}
	return todo, nil
}

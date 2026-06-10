package memorystore

import (
	"errors"
	"sync"
	"todolist/model"
)

//repository implementation for in_memory DB

type TodoMap struct {
	mu    sync.RWMutex
	store map[string]model.TODO
}

// Factory Method to create new TODO
func NewTodoMap() *TodoMap {

	return &TodoMap{
		store: make(map[string]model.TODO),
	}
}

func (t *TodoMap) Create(todo model.TODO) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.store[todo.TID]; ok {
		return errors.New("ID already found")

	}
	t.store[todo.TID] = todo
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
	return errors.New("ID not found in the map ")
}

func (t *TodoMap) Update(todo model.TODO) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.store[todo.TID]; ok {
		t.store[todo.TID] = todo
		return nil
	} else {
		return errors.New("ID not found in the map ")
	}

}

func (t *TodoMap) GetById(id string) (model.TODO, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.store[id]; ok {
		return t.store[id], nil
	}
	return model.TODO{}, errors.New("Store is empty")
}

func (t *TodoMap) GetAll() ([]model.TODO, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	todo := make([]model.TODO, 0)
	if len(t.store) == 0 {
		return nil, errors.New("Store is empty")
	}
	for _, v := range t.store {
		todo = append(todo, v)
	}
	return todo, nil
}

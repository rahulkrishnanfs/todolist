package memorystore

import (
	"errors"
	"todolist/model"
)

//repository implementation for in_memory DB

type TodoMap struct {
	store map[string]model.TODO
}

// Factory Method to create new TODO
func NewTodoMap() *TodoMap {

	return &TodoMap{
		store: make(map[string]model.TODO),
	}
}

// Factory method to create new category object
func NewCategoryMap() *CategoryMap {

	return &CategoryMap{
		store: make(map[string]model.Category),
	}
}

func (m TodoMap) Create(todo model.TODO) error {
	if _, ok := m.store[todo.TID]; ok {
		return errors.New("ID already found")

	}
	m.store[todo.TID] = todo
	return nil
}

func (m TodoMap) Delete(id string) error {
	if _, ok := m.store[id]; ok {
		delete(m.store, id)
		return nil
	}
	return errors.New("ID not found in the map ")
}

func (m TodoMap) GetById(id string) (model.TODO, error) {
	if _, ok := m.store[id]; ok {
		return m.store[id], nil
	}
	return model.TODO{}, errors.New("Store is empty")
}

func (m TodoMap) GetAll() ([]model.TODO, error) {
	todo := make([]model.TODO, 5)
	if len(m.store) == 0 {
		return nil, errors.New("Store is empty")
	}
	for _, v := range m.store {
		todo = append(todo, v)
	}
	return todo, nil
}

type CategoryMap struct {
	store map[string]model.Category
}

func (c CategoryMap) Create(Category model.Category) error {
	if _, ok := c.store[Category.CID]; ok {
		return errors.New("ID already found")

	}
	c.store[Category.CID] = Category
	return nil
}

package memorystore

import (
	"sync"

	"todolist/model"
)

type CategoryMap struct {
	mu    sync.RWMutex
	store map[string]model.Category
}

// Factory method to create new category object
func NewCategoryMap() *CategoryMap {

	return &CategoryMap{
		store: make(map[string]model.Category),
	}
}

func (c *CategoryMap) Create(Category model.Category) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.store[Category.CID]; ok {
		return model.ErrObjectAlreadyExists
	}
	c.store[Category.CID] = Category
	return nil
}

func (c *CategoryMap) Update(category model.Category) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.store[category.CID]; ok {
		c.store[category.CID] = category
		return nil
	}
	return model.ErrObjectNotFound

}
func (c *CategoryMap) Delete(cid string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.store[cid]; ok {
		delete(c.store, cid)
		return nil
	}
	return model.ErrObjectNotFound
}

func (c *CategoryMap) GetByID(cid string) (model.Category, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.store[cid]; ok {
		return c.store[cid], nil
	}
	return model.Category{}, model.ErrObjectNotFound
}

func (c *CategoryMap) GetAll() ([]model.Category, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	category := make([]model.Category, 0)
	if len(c.store) == 0 {
		return nil, model.ErrStoreEmpty
	}
	for _, v := range c.store {
		category = append(category, v)
	}
	return category, nil
}

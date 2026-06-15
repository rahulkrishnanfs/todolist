package model

import (
	"errors"
	"time"
)

var (
	ErrObjectAlreadyExists = errors.New("object already exists")
	ErrObjectNotFound      = errors.New("object not found")
	ErrStoreEmpty          = errors.New("store is empty")
)

// Database indepdendent (Domain model)
// Users can track the individual Items
type Todo struct {
	TID          string    `json:"tid"`
	Activity     string    `json:"activity"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creation_date"`
	IsDone       bool      `json:"is_done"`
	CategoryID   string    `json:"category_id"`
}

// Entity
// Groups TODO together
type Category struct {
	CID         string `json:"cid"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

// interface Abstraction
// Abstracts persistance for TODO Items
type TodoRepository interface {
	Create(Todo) error
	Update(string, Todo) error
	Delete(string) error
	GetById(string) (Todo, error)
	GetAll() ([]Todo, error)
}

// Intermediate between Datasource and domain model
// abstracts persistance of Category Items
type CategoryRepository interface {
	Create(Category) error
	Update(string, Category) error
	Delete(string) error
	GetById(string) (Category, error)
	GetAll() ([]Category, error)
}

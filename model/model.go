package model

import "time"

// Database indepdendent (Domain model)
// Users can track the individual Items
type TODO struct {
	TID          string
	Activity     string
	Description  string
	CreationDate time.Time
	IsDone       bool
	CategoryID   string
}

// Entity
// Groups TODO together
type Category struct {
	CID         string
	Description string
	Name        string
}

// interface Abstraction
// Abstracts persistance for TODO Items
type ToDoRepository interface {
	Create(TODO) error
	// Update(string, TODO) error
	Delete(string) error
	GetById(string) (TODO, error)
	GetAll() ([]TODO, error)
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

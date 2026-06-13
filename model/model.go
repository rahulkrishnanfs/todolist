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
type TODO struct {
	TID          string    `json:"tid"`
	Activity     string    `json:"activity"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creationdate"`
	IsDone       bool      `json:"isdone"`
	CategoryID   string    `json:"cayegoryid"`
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
type ToDoRepository interface {
	Create(TODO) error
	Update(TODO) error
	Delete(string) error
	GetById(string) (TODO, error)
	GetAll() ([]TODO, error)
}

// Intermediate between Datasource and domain model
// abstracts persistance of Category Items
type CategoryRepository interface {
	Create(Category) error
	Update(Category) error
	Delete(string) error
	GetByID(string) (Category, error)
	GetAll() ([]Category, error)
}

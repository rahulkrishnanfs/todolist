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
type (
	Todo struct {
		TID          string    `json:"tid"`
		Activity     string    `json:"activity"`
		Description  string    `json:"description"`
		CreationDate time.Time `json:"creation_date"`
		IsDone       bool      `json:"is_done"`
		CategoryID   string    `json:"category_id"`
		UserID       string    `json:"user_id"`
	}

	// Entity
	// Groups TODO together
	Category struct {
		CID         string `json:"cid"`
		Description string `json:"description"`
		Name        string `json:"name"`
		UID         string `json:"uid"`
	}

	User struct {
		UID          string `json:"uid"           jsonschema:"the user id of the user who signup"`
		Username     string `json:"username"      jsonschema:"the name of the person who signup"`
		Password     string `json:"password"      jsonschema:"the password of the person who signup"`
		EmailAddress string `json:"email_address" jsonschema:"the email address of the person who signup"`
	}
)

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

type UserRepository interface {
	Create(User) error
	Login(string, string) (bool, error)
}

package main

import "todolist/model"

type TODOController struct {
	//property for abstraction
	store model.ToDoRepository
}

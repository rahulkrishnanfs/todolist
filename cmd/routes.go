package main

import (
	"net/http"

	"todolist/controller"
)

func ToDoRoutes(todo controller.TODOController, mux *http.ServeMux) http.Handler {
	mux.HandleFunc("POST /api/todo/create", todo.Create)
	mux.HandleFunc("POST /api/todo/update", todo.Update)
	mux.HandleFunc("POST /api/todo/delete/{id}", todo.Delete)
	mux.HandleFunc("GET /api/todo/getbyid/{id}", todo.GetById)
	mux.HandleFunc("GET /api/todo/getall", todo.GetAll)
	return mux
}

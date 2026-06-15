package main

import (
	"context"
	"log/slog"
	"net/http"

	"todolist/controller"
)

func ToDoRoutes(todo controller.TodoController, mux *http.ServeMux, logger *slog.Logger) {
	logger.LogAttrs(context.Background(), slog.LevelDebug, "Adding http handler to the route for TODO")
	mux.HandleFunc("POST /api/v1/todos", todo.Create)
	mux.HandleFunc("GET /api/v1/todos", todo.GetAll)
	mux.HandleFunc("GET /api/v1/todos/{id}", todo.GetById)
	mux.HandleFunc("PUT /api/v1/todos/{id}", todo.Update)
	mux.HandleFunc("DELETE /api/v1/todos/{id}", todo.Delete)

}

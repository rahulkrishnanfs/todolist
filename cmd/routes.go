package main

import (
	"context"
	"log/slog"
	"net/http"

	"todolist/controller"
)

func ToDoRoutes(todo controller.TODOController, mux *http.ServeMux, logger *slog.Logger) {
	logger.LogAttrs(context.Background(), slog.LevelDebug, "Adding http handler to the route for TODO")
	mux.HandleFunc("POST /api/todo/create", todo.Create)
	mux.HandleFunc("PUT /api/todo/update", todo.Update)
	mux.HandleFunc("POST /api/todo/delete/{id}", todo.Delete)
	mux.HandleFunc("GET /api/todo/getbyid/{id}", todo.GetById)
	mux.HandleFunc("GET /api/todo/getall", todo.GetAll)
}

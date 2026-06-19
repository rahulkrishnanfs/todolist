package router

import (
	"context"
	"log/slog"
	"net/http"

	"todolist/pkg/auth"
	"todolist/pkg/controller"
)

func SetTodoRoutes(todo controller.TodoController, mux *http.ServeMux, logger *slog.Logger, auth *auth.Authenticator) {
	logger.LogAttrs(context.Background(), slog.LevelDebug, "Adding http handler to the route for TODO")
	mux.Handle("POST /api/v1/todos", auth.AuthorizeRequest(http.HandlerFunc(todo.Create)))
	mux.Handle("GET /api/v1/todos", auth.AuthorizeRequest(http.HandlerFunc(todo.GetAll)))
	mux.Handle("GET /api/v1/todos/{id}", auth.AuthorizeRequest(http.HandlerFunc(todo.GetById)))
	mux.Handle("PUT /api/v1/todos/{id}", auth.AuthorizeRequest(http.HandlerFunc(todo.Update)))
	mux.Handle("DELETE /api/v1/todos/{id}", auth.AuthorizeRequest(http.HandlerFunc(todo.Delete)))

}

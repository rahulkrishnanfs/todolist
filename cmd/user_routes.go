package main

import (
	"context"
	"log/slog"
	"net/http"

	"todolist/controller"
)

func UserRoutes(user controller.UserController, mux *http.ServeMux, logger *slog.Logger) {
	logger.LogAttrs(context.Background(), slog.LevelDebug, "Adding http handler to the route for TODO")
	mux.HandleFunc("POST /api/v1/users/signup", user.Create)
	mux.HandleFunc("POST /api/v1/users/login", user.Login)

}

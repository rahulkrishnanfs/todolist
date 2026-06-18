package main

import (
	"context"
	"log/slog"
	"net/http"

	"todolist/auth"
	"todolist/controller"
)

func CategoryRoutes(category controller.CategoryController, mux *http.ServeMux, logger *slog.Logger, auth *auth.Authenticator) {
	logger.LogAttrs(context.Background(), slog.LevelDebug, "Adding http handler to the route for Catgeory")
	mux.Handle("POST /api/v1/categories", auth.AuthorizeRequest(http.HandlerFunc(category.Create)))
	mux.Handle("GET /api/v1/categories", auth.AuthorizeRequest(http.HandlerFunc(category.GetAll)))
	mux.Handle("GET /api/v1/categories/{id}", auth.AuthorizeRequest(http.HandlerFunc(category.GetById)))
	mux.Handle("PUT /api/v1/categories/{id}", auth.AuthorizeRequest(http.HandlerFunc(category.Update)))
	mux.Handle("DELETE /api/v1/categories/{id}", auth.AuthorizeRequest(http.HandlerFunc(category.Delete)))

	// mux.HandleFunc("DELETE /api/v1/categories/{id}", category.Delete)
}

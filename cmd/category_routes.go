package main

import (
	"context"
	"log/slog"
	"net/http"

	"todolist/controller"
)

func CategoryRoutes(category controller.CategoryController, mux *http.ServeMux, logger *slog.Logger) {
	logger.LogAttrs(context.Background(), slog.LevelDebug, "Adding http handler to the route for Catgeory")
	mux.HandleFunc("POST /api/v1/categories", category.Create)
	mux.HandleFunc("GET /api/v1/categories", category.GetAll)
	mux.HandleFunc("GET /api/v1/categories/{id}", category.GetById)
	mux.HandleFunc("PUT /api/v1/categories/{id}", category.Update)
	mux.HandleFunc("DELETE /api/v1/categories/{id}", category.Delete)

}

package main

import (
	"context"
	"log/slog"
	"net/http"

	"todolist/controller"
)

func CategoryRoutes(category controller.CategoryController, mux *http.ServeMux, logger *slog.Logger) {
	logger.LogAttrs(context.Background(), slog.LevelDebug, "Adding http handler to the route for Catgeory")
	mux.HandleFunc("POST /api/category/create", category.Create)
	mux.HandleFunc("PUT /api/category/update", category.Update)
	mux.HandleFunc("POST /api/category/delete/{id}", category.Delete)
	mux.HandleFunc("GET /api/category/getbyid/{id}", category.GetById)
	mux.HandleFunc("GET /api/category/getall", category.GetAll)
}

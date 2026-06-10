package main

import (
	"net/http"

	"todolist/controller"
)

func CategoryRoutes(category controller.CategoryController, mux *http.ServeMux) http.Handler {
	mux.HandleFunc("POST /api/category/create", category.Create)
	mux.HandleFunc("POST /api/category/update", category.Update)
	mux.HandleFunc("POST /api/category/delete/{id}", category.Delete)
	mux.HandleFunc("GET /api/category/getbyid/{id}", category.GetById)
	mux.HandleFunc("GET /api/category/getall", category.GetAll)
	return mux
}

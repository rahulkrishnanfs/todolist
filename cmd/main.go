package main

import (
	"log"
	"net/http"
	"todolist/memorystore"
)

func main() {

	todo := &TODOController{
		store: memorystore.NewTodoMap(),
	}

	category := CategoryController{
		store: memorystore.NewCategoryMap(),
	}

	router := initializeRoutes(todo, &category)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
func initializeRoutes(t *TODOController, c *CategoryController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/todo/create", t.Create)
	mux.HandleFunc("POST /api/todo/update", t.Update)
	mux.HandleFunc("POST /api/todo/delete/{id}", t.Delete)
	mux.HandleFunc("GET /api/todo/getbyid/{id}", t.GetById)
	mux.HandleFunc("GET /api/todo/getall", t.GetAll)

	mux.HandleFunc("POST /api/category/create", c.Create)
	mux.HandleFunc("POST /api/category/update", c.Update)
	mux.HandleFunc("POST /api/category/delete/{id}", c.Delete)
	mux.HandleFunc("GET /api/category/getbyid/{id}", c.GetById)
	mux.HandleFunc("GET /api/category/getall", c.GetAll)

	return mux
}

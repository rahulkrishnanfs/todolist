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

	router := initializeRoutes(todo)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Listening...")
	// Run the http server
	server.ListenAndServe()
}
func initializeRoutes(t *TODOController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/todo", t.Create)

	return mux
}

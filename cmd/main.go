package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"todolist/controller"
	"todolist/memorystore"
)

func main() {
	mux := http.NewServeMux()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	categorystore := memorystore.NewCategoryMap()
	todostore := memorystore.NewTodoMap()
	categoryController := controller.NewCategoryController(categorystore, logger)
	todoController := controller.NewTODOController(todostore, logger)
	CategoryRoutes(*categoryController, mux)
	ToDoRoutes(*todoController, mux)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Listening....")
	server.ListenAndServe() // Run the http server
}

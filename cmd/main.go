package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"todolist/auth"
	"todolist/controller"
	"todolist/memorystore"
	"todolist/utils"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	config, err := utils.InitConfig("./config/properties.toml", utils.NewConfig(), logger)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "could not load the config",
			slog.String("error", err.Error()))
	}
	logger.LogAttrs(context.Background(), slog.LevelInfo, "configuration fetched..")
	SecretObj := utils.NewSecret(config.Service.KeystoreFilePath, config.Service.KeystorePasswword, logger)
	privateKey, publicKey := SecretObj.Extract()
	auth := auth.NewAuthenticator(privateKey, publicKey)
	mux := http.NewServeMux()

	categoryStore := memorystore.NewCategoryMap()
	todoStore := memorystore.NewTodoMap()
	userStore := memorystore.NewUserMap()
	categoryController := controller.NewCategoryController(categoryStore, logger)
	todoController := controller.NewTodoController(todoStore, logger)
	userController := controller.NewUserController(userStore, logger, auth)

	CategoryRoutes(*categoryController, mux, logger, auth)
	ToDoRoutes(*todoController, mux, logger, auth)
	UserRoutes(*userController, mux, logger)

	server := &http.Server{
		Addr:    config.Service.Port,
		Handler: mux,
	}
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Listening on port ...", slog.Int("port", 8080))
	if err := server.ListenAndServe(); err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "http server stopped",
			slog.String("error", err.Error()))
		os.Exit(1)
	}
}

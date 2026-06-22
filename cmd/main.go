package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"todolist/pkg/auth"
	"todolist/pkg/controller"
	"todolist/pkg/mcptools"
	"todolist/pkg/memorystore"
	"todolist/pkg/router"
	"todolist/pkg/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
	auth := auth.NewAuthenticator(privateKey, publicKey, logger)

	mux := http.NewServeMux()

	categoryStore := memorystore.NewCategoryMap()
	todoStore := memorystore.NewTodoMap()
	userStore := memorystore.NewUserMap()

	categoryController := controller.NewCategoryController(categoryStore, logger)
	todoController := controller.NewTodoController(todoStore, logger)
	userController := controller.NewUserController(userStore, logger, auth)

	router.SetCategoryRoutes(*categoryController, mux, logger, auth)
	router.SetTodoRoutes(*todoController, mux, logger, auth)
	router.SetUserRoutes(*userController, mux, logger)

	// create mcp server for todolist
	mcpserver := mcp.NewServer(&mcp.Implementation{Name: "todolist", Title: "todolist", Version: "1.0.0"}, nil)
	mcptools.SetTools(mcpserver, mux, *userController)

	server := &http.Server{
		Addr:    config.Service.Port,
		Handler: mux,
	}
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Listening on port ...", slog.String("port", config.Service.Port))
	// if err := server.ListenAndServe(); err != nil {
	// reference : https://github.com/denji/golang-tls
	if err := server.ListenAndServeTLS(config.Service.ServerCert, config.Service.ServerKey); err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "http server stopped",
			slog.String("error", err.Error()))
		os.Exit(1)
	}

}

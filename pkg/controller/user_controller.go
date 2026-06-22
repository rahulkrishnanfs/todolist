package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"todolist/pkg/auth"
	"todolist/pkg/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type UserController struct {
	user   model.UserRepository
	logger *slog.Logger
	auth   *auth.Authenticator
}

func NewUserController(user model.UserRepository, logger *slog.Logger, auth *auth.Authenticator) *UserController {
	return &UserController{
		user:   user,
		logger: logger,
		auth:   auth,
	}
}

type Output struct {
	Status string `json:"status" jsonschema:"user has been created"`
}

func (u *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		u.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to decode the Todo object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = u.user.Create(user)
	if err != nil {
		u.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to create the user object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated) //Todo
	u.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"user object has been created with the id",
		slog.String("id", user.UID))

}

func (u *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		u.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to decode the User object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	status, err := u.user.Login(user.Username, user.Password)

	if !status {
		u.logger.LogAttrs(context.Background(), slog.LevelError,
			"invalid credential for the user",
			slog.String("error", err.Error()), slog.String("username", user.Username))
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := u.auth.GenerateJWT(user.Username, "testrole")

	if err != nil {
		u.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to generate token",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (u *UserController) CreateUserTool(ctx context.Context, req *mcp.CallToolRequest, user model.User) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	err := u.user.Create(user)
	if err != nil {
		u.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to create the user object",
			slog.String("error", err.Error()))
		return nil, Output{}, err

	}
	u.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"user has been created")
	return nil, Output{Status: "User has been created"}, nil
}

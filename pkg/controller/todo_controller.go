package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"todolist/pkg/model"
)

type TodoController struct {
	//property for abstraction
	store  model.TodoRepository
	logger *slog.Logger
}

func NewTodoController(store model.TodoRepository, logger *slog.Logger) *TodoController {
	return &TodoController{
		store:  store,
		logger: logger,
	}
}

func (t *TodoController) Create(w http.ResponseWriter, r *http.Request) {
	var Todolist model.Todo
	err := json.NewDecoder(r.Body).Decode(&Todolist)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to decode the Todo object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.store.Create(Todolist)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to create the Todo object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated) //Todo
	t.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"Todo object has been created with the id",
		slog.String("id", Todolist.TID))

}
func (t *TodoController) Update(w http.ResponseWriter, r *http.Request) {
	var Todolist model.Todo
	id := r.PathValue("id")
	err := json.NewDecoder(r.Body).Decode(&Todolist)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to decode the Todo object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.store.Update(id, Todolist)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to update the Todo object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) //Update
	t.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"Todo object with the id has been updated",
		slog.String("id", Todolist.TID))
}

func (t *TodoController) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := t.store.Delete(id)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to delete the Todo",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	t.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"Todo object with the requested id has been deleted",
		slog.String("id", id))

}

func (t *TodoController) GetAll(w http.ResponseWriter, r *http.Request) {
	Todolists, err := t.store.GetAll()
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to extract the Todo objects",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(Todolists)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to encode the Todo objects",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(j); err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to write the response",
			slog.String("error", err.Error()))
	}
	t.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"all requested Todo objects has been returned to the client")

}

func (t *TodoController) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	Todo, err := t.store.GetById(id)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to get Todo object by id",
			slog.String("error", err.Error()),
			slog.String("category_id", id))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(Todo)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to encode the Todo object)",
			slog.String("error", err.Error()),
			slog.String("id", Todo.TID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(j); err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to write the response",
			slog.String("error", err.Error()))
	}
	t.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"requested Todo object with id has been sent to the client",
		slog.String("id", id))

}

package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"todolist/model"
)

type TODOController struct {
	//property for abstraction
	store  model.ToDoRepository
	logger *slog.Logger
}

func NewTODOController(store model.ToDoRepository, logger *slog.Logger) *TODOController {
	return &TODOController{
		store:  store,
		logger: logger,
	}
}

func (t *TODOController) Create(w http.ResponseWriter, r *http.Request) {
	var todolist model.TODO
	err := json.NewDecoder(r.Body).Decode(&todolist)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to decode the todo object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.store.Create(todolist)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to create the todo object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated) //TODO
	t.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"todo object has been created with the id",
		slog.String("id", todolist.TID))

}
func (t *TODOController) Update(w http.ResponseWriter, r *http.Request) {
	var todolist model.TODO
	err := json.NewDecoder(r.Body).Decode(&todolist)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to decode the todo object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.store.Update(todolist)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to update the todo object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) //Update
	t.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"todo object with the id has been updated",
		slog.String("id", todolist.TID))
}

func (t *TODOController) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := t.store.Delete(id)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to delete the todo",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	t.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"todo object with the requested id has been deleted",
		slog.String("id", id))

}

func (t *TODOController) GetAll(w http.ResponseWriter, r *http.Request) {
	todolists, err := t.store.GetAll()
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to extract the todo objects",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(todolists)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to encode the todo objects",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	t.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"all requested todo objects has been returned to the client")

}

func (t *TODOController) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	todo, err := t.store.GetById(id)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to get todo object by id",
			slog.String("error", err.Error()),
			slog.String("category_id", id))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(todo)
	if err != nil {
		t.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to encode the todo object)",
			slog.String("error", err.Error()),
			slog.String("id", todo.TID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	t.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"requested todo object with id has been sent to the client",
		slog.String("id", id))

}

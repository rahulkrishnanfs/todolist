package controller

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("request came here")
	var todolist model.TODO
	err := json.NewDecoder(r.Body).Decode(&todolist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = t.store.Create(todolist)
	if err != nil {
		fmt.Println("somethings is having issue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated) //TODO

}
func (t *TODOController) Update(w http.ResponseWriter, r *http.Request) {
	var todolist model.TODO
	err := json.NewDecoder(r.Body).Decode(&todolist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = t.store.Update(todolist)
	if err != nil {
		fmt.Println("somethings is having issue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent) //Update
}

func (t *TODOController) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := t.store.Delete(id)
	if err != nil {
		fmt.Println("somethings is having issue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (t *TODOController) GetAll(w http.ResponseWriter, r *http.Request) {
	todolists, err := t.store.GetAll()
	if err != nil {
		fmt.Println("somethings is having issue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(todolists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func (t *TODOController) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	todo, err := t.store.GetById(id)
	if err != nil {
		fmt.Println("somethings is having issue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

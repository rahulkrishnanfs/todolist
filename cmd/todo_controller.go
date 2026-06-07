package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todolist/model"
)

type TODOController struct {
	//property for abstraction
	store model.ToDoRepository
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

}
func (t *TODOController) Update(w http.ResponseWriter, r *http.Request) {

}

func (t *TODOController) Delete(w http.ResponseWriter, r *http.Request) {

}

func (t *TODOController) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (t *TODOController) GetById(w http.ResponseWriter, r *http.Request) {

}

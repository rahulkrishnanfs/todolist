package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todolist/model"
)

type CategoryController struct {
	store model.CategoryRepository
}

func (t *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
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
func (t *CategoryController) Update(w http.ResponseWriter, r *http.Request) {

}

func (t *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {

}

func (t *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (t *CategoryController) GetById(w http.ResponseWriter, r *http.Request) {

}

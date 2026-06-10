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

func (c *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request came here")
	var category model.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = c.store.Create(category)
	if err != nil {
		fmt.Println("somethings is having issue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (c *CategoryController) Update(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = c.store.Update(category)
	if err != nil {
		fmt.Println("somethings is having issue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (c *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := c.store.Delete(id)
	if err != nil {
		fmt.Println("somethings is having issue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := c.store.GetAll()
	if err != nil {
		fmt.Println("somethings is having issue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (c *CategoryController) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	category, err := c.store.GetByID(id)
	if err != nil {
		fmt.Println("somethings is having issue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

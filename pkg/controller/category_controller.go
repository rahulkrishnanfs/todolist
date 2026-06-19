package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"todolist/pkg/model"
)

// dealing with http programing
type CategoryController struct {
	store  model.CategoryRepository // we are communicating by contract: Liskov substitution principle ; Design by contract
	logger *slog.Logger
}

func NewCategoryController(store model.CategoryRepository, logger *slog.Logger) *CategoryController {

	return &CategoryController{
		store:  store,
		logger: logger,
	}
}

func (c *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		// if (errors.Is(err,model.ErrObjectAlreadyExists)){
		//  Package level errors can be compared to idividual level errors
		// }
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to decode the category object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.store.Create(category)
	if err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to create the category object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//message followed by key value data -> structured logging
	c.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"category object has been created with the id",
		slog.String("category_id", category.CID))
	w.WriteHeader(http.StatusCreated)
}
func (c *CategoryController) Update(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	id := r.PathValue("id")
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to decode the category object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.store.Update(id, category)
	if err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to update the category object",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"catagory object with the id has been updated",
		slog.String("id", category.CID))
	w.WriteHeader(http.StatusNoContent)
}

func (c *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := c.store.Delete(id)
	if err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to delete the category",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"the requested category object has been deleted with id",
		slog.String("id", id))
	w.WriteHeader(http.StatusNoContent)
}

func (c *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := c.store.GetAll()
	if err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to extract the category objects",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(categories)
	if err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to encode the category objects",
			slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"all requested category objects has been sent to the client")
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(j); err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to write the response",
			slog.String("error", err.Error()))
	}
}

func (c *CategoryController) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	category, err := c.store.GetById(id)
	if err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to get category object by id",
			slog.String("error", err.Error()),
			slog.String("id", id))

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(category)
	if err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to encode the category object",
			slog.String("error", err.Error()),
			slog.String("id", category.CID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"requested category object with id has been sent to the client",
		slog.String("id", id))
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(j); err != nil {
		c.logger.LogAttrs(context.Background(), slog.LevelError,
			"failed to write the response",
			slog.String("error", err.Error()))
	}
}

package main

import (
	"fmt"
	"time"

	"todolist/memorystore"
	"todolist/model"
)

func main() {

	todo := TODOController{
		store: memorystore.NewTodoMap(),
	}
	todo.store.Create(model.TODO{
		TID:          "1",
		Activity:     "Going to Shop",
		Description:  "",
		CreationDate: time.Now(),
		IsDone:       false,
		CategoryID:   "1"})

	fmt.Println(todo.store.GetAll())
}

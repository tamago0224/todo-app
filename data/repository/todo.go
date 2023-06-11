package repository

import (
	"github.com/tamago0224/rest-app-backend/data/model"
)

type TodoRepository interface {
	GetAllTodo(userId int64) ([]model.Todo, error)
	GetTodo(model.Todo) (model.Todo, error)
	AddTodo(model.Todo) (model.Todo, error)
	DeleteTodo(model.Todo) (model.Todo, error)
}

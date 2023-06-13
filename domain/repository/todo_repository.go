package repository

import (
	"github.com/tamago0224/rest-app-backend/domain/model"
)

type TodoRepository interface {
	GetAllTodo(userId int) ([]model.Todo, error)
	GetTodo(userID, todoID int) (model.Todo, error)
	AddTodo(model.Todo) (model.Todo, error)
	DeleteTodo(model.Todo) (model.Todo, error)
}

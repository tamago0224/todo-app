package repository

import (
	"github.com/tamago0224/rest-app-backend/models"
)

type TodoRepository interface {
	GetAllTodo() ([]models.Todo, error)
	GetTodo(models.Todo) (models.Todo, error)
	AddTodo(models.Todo) (models.Todo, error)
	DeleteTodo(models.Todo) (models.Todo, error)
}

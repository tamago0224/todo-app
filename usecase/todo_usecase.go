package usecase

import (
	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/domain/service"
)

type ITodoUsecase interface {
	GetTodoList(userID int) ([]model.Todo, error)
	GetTodo(userID, todoID int) (model.Todo, error)
	AddTodo(user model.Todo) (model.Todo, error)
	DeleteTodo(model.Todo) (model.Todo, error)
}

type todoUsecase struct {
	svc service.ITodoService
}

func NewTodoUsecase(svc service.ITodoService) ITodoUsecase {
	return &todoUsecase{svc: svc}
}

func (t *todoUsecase) GetTodoList(userID int) ([]model.Todo, error) {
	todos, err := t.svc.FindAllTodo(userID)
	if err != nil {
		return []model.Todo{}, err
	}

	return todos, nil
}

func (t *todoUsecase) GetTodo(userID, todoID int) (model.Todo, error) {
	todo, err := t.svc.FindByID(userID, todoID)
	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

func (t *todoUsecase) AddTodo(todo model.Todo) (model.Todo, error) {
	todo, err := t.svc.CreateTodo(int(todo.UserId), todo)
	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

func (t *todoUsecase) DeleteTodo(todo model.Todo) (model.Todo, error) {
	todo, err := t.svc.DeleteTodo(int(todo.UserId), todo)
	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

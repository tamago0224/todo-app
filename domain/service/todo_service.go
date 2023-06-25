package service

import (
	"database/sql"
	"errors"

	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/domain/repository"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
)

type ITodoService interface {
	FindAllTodo(userID int) ([]model.Todo, error)
	FindByID(userID, todoID int) (model.Todo, error)
	CreateTodo(userID int, todo model.Todo) (model.Todo, error)
	DeleteTodo(userID int, todo model.Todo) (model.Todo, error)
}

type todoService struct {
	userRepo repository.UserRepository
	todoRepo repository.TodoRepository
}

func NewTodoService(userRepo repository.UserRepository, todoRepo repository.TodoRepository) ITodoService {
	return &todoService{
		userRepo: userRepo,
		todoRepo: todoRepo,
	}
}

func (tc *todoService) FindAllTodo(userID int) ([]model.Todo, error) {
	_, err := tc.userRepo.SelectByID(userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return []model.Todo{}, ErrUserNotFound
		}
		return []model.Todo{}, err
	}

	todos, err := tc.todoRepo.GetAllTodo(userID)
	if err != nil {
		return []model.Todo{}, err
	}

	return todos, nil
}

func (tc *todoService) FindByID(userID, todoID int) (model.Todo, error) {
	_, err := tc.userRepo.SelectByID(userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return model.Todo{}, ErrUserNotFound
		}
		return model.Todo{}, err
	}

	todo, err := tc.todoRepo.GetTodo(userID, todoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Todo{}, ErrTodoNotFound
		}
		return model.Todo{}, err
	}

	return todo, nil
}

func (tc *todoService) CreateTodo(userID int, todo model.Todo) (model.Todo, error) {
	_, err := tc.userRepo.SelectByID(userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return model.Todo{}, ErrUserNotFound
		}
		return model.Todo{}, err
	}

	todo.UserId = int64(userID)
	newTodo, err := tc.todoRepo.AddTodo(todo)
	if err != nil {
		return model.Todo{}, err
	}

	return newTodo, nil
}

func (tc *todoService) DeleteTodo(userID int, todo model.Todo) (model.Todo, error) {
	_, err := tc.userRepo.SelectByID(userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return model.Todo{}, ErrUserNotFound
		}
		return model.Todo{}, err
	}

	todo.UserId = int64(userID)
	oldTodo, err := tc.todoRepo.DeleteTodo(todo)
	if err != nil {
		return model.Todo{}, err
	}

	return oldTodo, nil
}

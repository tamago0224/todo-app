package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tamago0224/rest-app-backend/models"
	"github.com/tamago0224/rest-app-backend/repository"
)

type TodoController struct {
	todoRepo repository.TodoRepository
}

func NewTodoController(todoRepo repository.TodoRepository) *TodoController {
	return &TodoController{todoRepo: todoRepo}
}

func (tc *TodoController) GetTodoList(c echo.Context) error {
	userId := LoginUserId(c)
	todos, err := tc.todoRepo.GetAllTodo(userId)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, todos)
}

func (tc *TodoController) GetTodo(c echo.Context) error {
	userId := LoginUserId(c)
	todoId := c.Param("id")
	id, err := strconv.Atoi(todoId)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	todo, err := tc.todoRepo.GetTodo(models.Todo{Id: int64(id), UserId: userId})
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, todo)
}

func (tc *TodoController) AddTodo(c echo.Context) error {
	userId := LoginUserId(c)
	var todo models.Todo

	err := c.Bind(&todo)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid todo body")
	}
	// ログインユーザのIDを追加する
	todo.UserId = userId
	addTodo, err := tc.todoRepo.AddTodo(todo)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusCreated, addTodo)
}

func (tc *TodoController) DeleteTodo(c echo.Context) error {
	userId := LoginUserId(c)
	todoId := c.Param("id")
	id, err := strconv.Atoi(todoId)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	todo, err := tc.todoRepo.DeleteTodo(models.Todo{Id: int64(id), UserId: userId})
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusNoContent, todo)
}

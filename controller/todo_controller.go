package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/usecase"
)

type TodoController struct {
	usecase usecase.ITodoUsecase
}

func NewTodoController(todoUsecase usecase.ITodoUsecase) *TodoController {
	return &TodoController{usecase: todoUsecase}
}

func (tc *TodoController) GetTodoList(c echo.Context) error {
	userID := LoginUserId(c)
	todos, err := tc.usecase.GetTodoList(userID)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, todos)
}

func (tc *TodoController) GetTodo(c echo.Context) error {
	userID := LoginUserId(c)
	id := c.Param("id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	todo, err := tc.usecase.GetTodo(userID, todoID)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, todo)
}

func (tc *TodoController) AddTodo(c echo.Context) error {
	userID := LoginUserId(c)
	var todo model.Todo

	err := c.Bind(&todo)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid todo body")
	}
	todo.UserId = int64(userID)
	addTodo, err := tc.usecase.AddTodo(todo)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusCreated, addTodo)
}

func (tc *TodoController) DeleteTodo(c echo.Context) error {
	userID := LoginUserId(c)
	id := c.Param("id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	todo, err := tc.usecase.DeleteTodo(model.Todo{Id: int64(todoID), UserId: int64(userID)})
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusNoContent, todo)
}

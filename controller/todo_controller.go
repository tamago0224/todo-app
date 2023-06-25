package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/domain/service"
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

		var apiError APIError
		if errors.Is(err, service.ErrUserNotFound) {
			apiError = APIError{Code: http.StatusForbidden, Message: http.StatusText(http.StatusForbidden)}
		} else {
			apiError = APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		}
		return c.JSON(apiError.Code, apiError)
	}

	return c.JSON(http.StatusOK, todos)
}

func (tc *TodoController) GetTodo(c echo.Context) error {
	userID := LoginUserId(c)
	id := c.Param("id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return APIError{Code: http.StatusBadRequest, Message: "invalid id format"}
	}

	todo, err := tc.usecase.GetTodo(userID, todoID)
	if err != nil {
		log.Print(err)

		var apiError APIError
		if errors.Is(err, service.ErrUserNotFound) {
			apiError = APIError{Code: http.StatusForbidden, Message: http.StatusText(http.StatusForbidden)}
		} else if errors.Is(err, service.ErrTodoNotFound) {
			apiError = APIError{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}
		} else {
			apiError = APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		}

		return c.JSON(apiError.Code, apiError)
	}

	return c.JSON(http.StatusOK, todo)
}

func (tc *TodoController) AddTodo(c echo.Context) error {
	userID := LoginUserId(c)
	var todo model.Todo

	err := c.Bind(&todo)
	if err != nil {
		return APIError{Code: http.StatusBadRequest, Message: "invalid todo body"}
	}
	todo.UserId = int64(userID)
	addTodo, err := tc.usecase.AddTodo(todo)
	if err != nil {
		log.Print(err)

		var apiError APIError
		if errors.Is(err, service.ErrUserNotFound) {
			apiError = APIError{Code: http.StatusForbidden, Message: http.StatusText(http.StatusForbidden)}
		} else {
			apiError = APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		}
		return c.JSON(apiError.Code, apiError)
	}

	return c.JSON(http.StatusCreated, addTodo)
}

func (tc *TodoController) DeleteTodo(c echo.Context) error {
	userID := LoginUserId(c)
	id := c.Param("id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		apiError := APIError{Code: http.StatusBadRequest, Message: "invalid id format"}
		return c.JSON(apiError.Code, apiError)
	}

	_, err = tc.usecase.DeleteTodo(model.Todo{Id: int64(todoID), UserId: int64(userID)})
	if err != nil {
		log.Print(err)

		var apiError APIError
		if errors.Is(err, service.ErrUserNotFound) {
			apiError = APIError{Code: http.StatusForbidden, Message: http.StatusText(http.StatusForbidden)}
		} else {
			apiError = APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		}
		return c.JSON(apiError.Code, apiError)
	}

	return c.NoContent(http.StatusNoContent)
}

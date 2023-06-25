package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/domain/service"
	"github.com/tamago0224/rest-app-backend/usecase"
)

type UserController struct {
	usecase usecase.IUserUsecase
}

func NewUserController(userUsecase usecase.IUserUsecase) *UserController {
	return &UserController{usecase: userUsecase}
}

func (uc *UserController) SearchUser(c echo.Context) error {
	name := c.QueryParam("name")

	user, err := uc.usecase.SearchUser(name)
	if err != nil {
		log.Print(err)

		var apiError APIError
		if errors.Is(err, service.ErrUserNotFound) {
			apiError = APIError{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}
		} else {
			apiError = APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		}

		return c.JSON(apiError.Code, apiError)
	}

	return c.JSON(http.StatusOK, user)
}

func (uc *UserController) CreateUser(c echo.Context) error {
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		log.Print(err)

		apiError := APIError{Code: http.StatusBadRequest, Message: "invalid create user request body"}
		return c.JSON(apiError.Code, apiError)
	}

	createdUser, err := uc.usecase.CreateUser(user)
	if err != nil {
		log.Print(err)

		var apiError APIError
		if errors.Is(err, service.ErrUserAlreadyExist) {
			apiError = APIError{Code: http.StatusConflict, Message: http.StatusText(http.StatusConflict)}
		} else {
			apiError = APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		}

		return c.JSON(apiError.Code, apiError)
	}

	return c.JSON(http.StatusCreated, createdUser)
}

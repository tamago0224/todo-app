package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tamago0224/rest-app-backend/data/model"
	"github.com/tamago0224/rest-app-backend/data/repository"
)

type UserController struct {
	userRepo repository.UserRepository
}

func NewUserController(userRepo repository.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

func (uc *UserController) SearchUser(c echo.Context) error {
	name := c.QueryParam("name")

	user, err := uc.userRepo.SearchUser(name)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user %s not found.", name))
		}
		return InternalServerError(err)
	}

	return c.JSON(http.StatusOK, user)
}

func (uc *UserController) CreateUser(c echo.Context) error {
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		return InternalServerError(err)
	}

	createdUser, err := uc.userRepo.CreateUser(user)
	if err != nil {
		return InternalServerError(err)
	}

	return c.JSON(http.StatusCreated, createdUser)
}

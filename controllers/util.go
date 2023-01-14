package controllers

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func InternalServerError(err error) error {
	log.Print(err)
	return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
}

func LoginUserId(c echo.Context) int64 {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	userId := claims.Id

	return userId
}

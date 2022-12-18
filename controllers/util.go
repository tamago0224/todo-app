package controllers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InternalServerError(err error) error {
	log.Print(err)
	return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
}

package controller

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func InternalServerError(err error) error {
	log.Print(err)
	return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
}

func LoginUserId(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	userId := claims.Id

	return int(userId)
}

func randString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(c echo.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   c.IsTLS(),
		HttpOnly: true,
	}

	c.SetCookie(cookie)
}

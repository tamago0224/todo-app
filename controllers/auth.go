package controllers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/tamago0224/rest-app-backend/models"
	"github.com/tamago0224/rest-app-backend/repository"
)

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type AuthController struct {
	userRepo repository.UserRepository
}

func NewAuthController(userRepo repository.UserRepository) AuthController {
	return AuthController{userRepo: userRepo}
}

func (ac *AuthController) Login(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	// nameをキーにDBからユーザ名を取得し、パスワードが一致することをチェックする
	u, err := ac.userRepo.SearchUser(user.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user anauthorized")
	}
	if user.Password != u.Password {
		return echo.NewHTTPError(http.StatusUnauthorized, "user anauthorized")
	}

	// ユーザ名、パスワードが一致したらJWTトークンを生成する
	claims := &jwtCustomClaims{Name: u.Name, StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return InternalServerError(err)
	}

	c.SetCookie(&http.Cookie{Name: "auth_token", Value: t, HttpOnly: true, Secure: true})
	return c.JSON(http.StatusOK, nil)
}

package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/tamago0224/rest-app-backend/data/model"
	"github.com/tamago0224/rest-app-backend/data/repository"
)

type JwtCustomClaims struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
	jwt.StandardClaims
}

type AuthController struct {
	userRepo repository.UserRepository
}

func NewAuthController(userRepo repository.UserRepository) AuthController {
	return AuthController{userRepo: userRepo}
}

func (ac *AuthController) Login(c echo.Context) error {
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	// nameをキーにDBからユーザ名を取得し、パスワードが一致することをチェックする
	u, err := ac.userRepo.SearchUser(user.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user unauthorized")
	}
	if user.Password != u.Password {
		return echo.NewHTTPError(http.StatusUnauthorized, "user unauthorized")
	}

	// ユーザ名、パスワードが一致したらJWTトークンを生成する
	claims := &JwtCustomClaims{Name: u.Name, Id: u.Id, StandardClaims: jwt.StandardClaims{
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

func (ac *AuthController) RegistUser(c echo.Context) error {
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	// ユーザを登録する
	u, err := ac.userRepo.CreateUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("user already exist. %s", user.Name))
	}

	// ユーザの作成に成功すればログイな使いにするのでCookieをセットする
	claims := &JwtCustomClaims{Name: u.Name, Id: u.Id, StandardClaims: jwt.StandardClaims{
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

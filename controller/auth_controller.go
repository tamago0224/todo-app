package controller

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/domain/service"
	"github.com/tamago0224/rest-app-backend/usecase"
)

type JwtCustomClaims struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
	jwt.StandardClaims
}

type AuthController struct {
	usecase usecase.IUserUsecase
}

func NewAuthController(userUsecase usecase.IUserUsecase) AuthController {
	return AuthController{usecase: userUsecase}
}

func (ac *AuthController) Login(c echo.Context) error {
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		log.Print(err)

		apiError := APIError{Code: http.StatusBadRequest, Message: "invalid login request body"}
		return c.JSON(apiError.Code, apiError)
	}

	// nameをキーにDBからユーザ名を取得し、パスワードが一致することをチェックする
	u, err := ac.usecase.SearchUser(user.Name)
	if err != nil {
		log.Print(err)

		var apiError APIError
		if errors.Is(err, service.ErrUserNotFound) {
			apiError = APIError{Code: http.StatusUnauthorized}
		} else {
			apiError = APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		}
		return c.JSON(apiError.Code, apiError)
	}
	if user.Password != u.Password {
		apiError := APIError{Code: http.StatusUnauthorized}
		return c.JSON(apiError.Code, apiError)
	}

	// ユーザ名、パスワードが一致したらJWTトークンを生成する
	claims := &JwtCustomClaims{Name: u.Name, Id: u.Id, StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Print(err)

		apiError := APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		return c.JSON(apiError.Code, apiError)
	}

	c.SetCookie(&http.Cookie{Name: "auth_token", Value: t, HttpOnly: true, Secure: true})
	return c.JSON(http.StatusOK, nil)
}

func (ac *AuthController) RegistUser(c echo.Context) error {
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		log.Print(err)

		apiError := APIError{Code: http.StatusBadRequest, Message: "invalid login request body"}
		return c.JSON(apiError.Code, apiError)

	}

	// ユーザを登録する
	u, err := ac.usecase.CreateUser(user)
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

	// ユーザの作成に成功すればログイン済みの扱いにするのでCookieをセットする
	claims := &JwtCustomClaims{Name: u.Name, Id: u.Id, StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Print(err)

		apiError := APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		return c.JSON(apiError.Code, apiError)

	}

	c.SetCookie(&http.Cookie{Name: "auth_token", Value: t, HttpOnly: true, Secure: true})
	return c.JSON(http.StatusOK, nil)
}

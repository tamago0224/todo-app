package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/domain/service"
	"github.com/tamago0224/rest-app-backend/usecase"
	"golang.org/x/oauth2"
)

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	signingKey   = os.Getenv("SIGNING_KEY")
)

type JwtCustomClaims struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
	jwt.StandardClaims
}

type AuthConfig struct {
	verifier     *oidc.IDTokenVerifier
	oauth2Config oauth2.Config
}

func NewOAuthConfig() (AuthConfig, error) {

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return AuthConfig{}, err
	}
	oidcConfig := &oidc.Config{ClientID: clientID}
	verifier := provider.Verifier(oidcConfig)

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://127.0.0.1:8080/callback",
		Scopes:       []string{oidc.ScopeOpenID, "email"},
	}

	return AuthConfig{verifier: verifier, oauth2Config: config}, nil
}

type AuthController struct {
	usecase    usecase.IUserUsecase
	authConfig AuthConfig
}

func NewAuthController(userUsecase usecase.IUserUsecase, config AuthConfig) AuthController {
	return AuthController{usecase: userUsecase, authConfig: config}
}

func (ac *AuthController) Login(c echo.Context) error {
	state, err := randString(16)
	if err != nil {
		return err
	}
	nonce, err := randString(16)
	if err != nil {
		log.Print(err)

		apiError := APIError{Code: http.StatusBadRequest, Message: "invalid login request body"}
		return c.JSON(apiError.Code, apiError)
	}

	// state, nonceをユーザのCookieに保存させる
	setCallbackCookie(c, "state", state)
	setCallbackCookie(c, "nonce", nonce)

	return c.Redirect(http.StatusFound, ac.authConfig.oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)))
}

func (ac *AuthController) Callback(c echo.Context) error {
	state, err := c.Cookie("state")
	if err != nil {
		log.Print(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to get cookie state"})
	}

	if c.QueryParam("state") != state.Value {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "state did not match"})
	}

	oauth2Token, err := ac.authConfig.oauth2Config.Exchange(c.Request().Context(), c.QueryParam("code"))
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to exchang token: "+err.Error())
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "no id_token field in oauth2 token.")
	}
	idToken, err := ac.authConfig.verifier.Verify(c.Request().Context(), rawIDToken)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to verify ID token: "+err.Error())
	}

	nonce, err := c.Cookie("nonce")
	if err != nil {
		log.Print(err)
		apiError := APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		return c.JSON(apiError.Code, apiError)
	}
	if idToken.Nonce != nonce.Value {
		log.Print(err)

		apiError := APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
		return c.JSON(apiError.Code, apiError)
	}

	// 検証完了後、トークンのemail属性のユーザが作成されているかどうかをチェックし、
	// 存在すれば認証成功、いなければユーザを作成し認証成功を返す
	var idTokenClaims struct {
		Email string `json:"email"`
	}
	if err = idToken.Claims(&idTokenClaims); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to unmarshal IDToken"})
	}
	user, err := ac.usecase.SearchOrCreate(idTokenClaims.Email)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to create user"})
	}

	// ユーザ情報を元にJWTを作成しユーザに返却する
	jwtClaims := jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := jwtToken.SignedString([]byte(signingKey))
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to generate jwt token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": tokenString})

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

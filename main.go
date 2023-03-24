package main

import (
	"database/sql"
	"log"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tamago0224/rest-app-backend/controllers"
	"github.com/tamago0224/rest-app-backend/repository"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	todoRepository := repository.NewTodoMariaDBRepository(db)
	todoController := controllers.NewTodoController(todoRepository)
	userRepository := repository.NewUserMariaDBRepository(db)
	authController := controllers.NewAuthController(userRepository)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
	}))

	e.POST("/login", authController.Login)
	e.POST("/register", authController.RegistUser)
	apiGroup := e.Group("/api/v1")
	apiGroup.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(controllers.JwtCustomClaims)
		},
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:auth_token",
	}))
	apiGroup.GET("/todos", todoController.GetTodoList)
	apiGroup.POST("/todos", todoController.AddTodo)
	apiGroup.GET("/todos/:id", todoController.GetTodo)
	apiGroup.DELETE("/todos/:id", todoController.DeleteTodo)

	e.Logger.Fatal(e.Start(":8080"))
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "todo:hogehoge@tcp(db:3306)/todo")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

package main

import (
	"github.com/labstack/echo/v4"
	"github.com/tamago0224/rest-app-backend/controllers"
)

func main() {
	e := echo.New()
	e.GET("/todos", controllers.GetTodo)
	e.POST("/todos", controllers.AddTodo)
	e.GET("/todos/:id", controllers.GetTodo)
	e.POST("/todos/:id", controllers.DeleteTodo)

	e.Logger.Fatal(e.Start(":8080"))
}

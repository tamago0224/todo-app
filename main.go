package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
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
	e := echo.New()
	e.GET("/todos", todoController.GetTodoList)
	e.POST("/todos", todoController.AddTodo)
	e.GET("/todos/:id", todoController.GetTodo)
	e.DELETE("/todos/:id", todoController.DeleteTodo)

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

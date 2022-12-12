package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tamago0224/rest-app-backend/models"
)

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "/var/rest-app/todo.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

func GetTodoList(c echo.Context) error {
	db, err := openDB()
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	var todos []models.Todo
	for rows.Next() {
		var id int64
		var title string
		var description string
		err = rows.Scan(&id, &title, &description)
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		todos = append(todos, models.Todo{Id: id, Title: title, Description: description})
	}

	return c.JSON(http.StatusOK, todos)
}

func AddTodo(c echo.Context) error {
	var todo models.Todo
	err := c.Bind(&todo)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid todo body")
	}

	db, err := openDB()
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	defer db.Close()

	// write todo to database.
	result, err := db.Exec("INSERT INTO todos (title, description) VALUES (?, ?)", todo.Title, todo.Description)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	if rows != 1 {
		log.Printf("expected to affect 1 row, affected %d", rows)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	todo.Id = id

	return c.JSON(http.StatusCreated, todo)
}

func GetTodo(c echo.Context) error {
	todoId := c.Param("id")

	db, err := openDB()
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	defer db.Close()

	var id int64
	var title string
	var description string
	err = db.QueryRow("SELECT * FROM todos WHERE id = ?", todoId).Scan(&id, &title, &description)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, models.Todo{Id: id, Title: title, Description: description})
}

func DeleteTodo(c echo.Context) error {
	todoId := c.Param("id")

	db, err := openDB()
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM todos WHERE id = ?", todoId)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	if rows != 1 {
		log.Printf("expected to affect 1 row, affected %d", rows)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusNoContent, nil)
}

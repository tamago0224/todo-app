package controllers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tamago0224/golang-rest-api/models"
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

func GetTodoList(rw http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var todos []models.Todo
	for rows.Next() {
		var id int64
		var title string
		var description string
		err = rows.Scan(&id, &title, &description)
		if err != nil {
			log.Print(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		todos = append(todos, models.Todo{Id: id, Title: title, Description: description})
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(todos)
}

func AddTodo(rw http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var todo models.Todo
	json.Unmarshal(body, &todo)

	db, err := openDB()
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// write todo to database.
	result, err := db.Exec("INSERT INTO todos (title, description) VALUES (?, ?)", todo.Title, todo.Description)
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rows != 1 {
		log.Printf("expected to affect 1 row, affected %d", rows)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo.Id = id

	json.NewEncoder(rw).Encode(todo)
}

func GetTodo(rw http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := openDB()
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var id int64
	var title string
	var description string
	err = db.QueryRow("SELECT * FROM todos WHERE id = ?", todoId).Scan(&id, &title, &description)
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(models.Todo{Id: id, Title: title, Description: description})
}

func DeleteTodo(rw http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := openDB()
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM todos WHERE id = ?", todoId)
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rows != 1 {
		log.Printf("expected to affect 1 row, affected %d", rows)
	}

	rw.WriteHeader(http.StatusNoContent)
}

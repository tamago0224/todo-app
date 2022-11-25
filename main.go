package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tamago0224/rest-app-backend/controllers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todos", controllers.GetTodoList).Methods(http.MethodGet)
	r.HandleFunc("/todos", controllers.AddTodo).Methods(http.MethodPost)
	r.HandleFunc("/todos/{id:[0-9]+}", controllers.GetTodo).Methods(http.MethodGet)
	r.HandleFunc("/todos/{id:[0-9]+}", controllers.DeleteTodo).Methods(http.MethodDelete)

	http.ListenAndServe(":8080", r)
}

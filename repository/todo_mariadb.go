package repository

import (
	"database/sql"

	"github.com/tamago0224/rest-app-backend/models"
)

type TodoMariaDB struct {
	db *sql.DB
}

func NewTodoMariaDBRepository(db *sql.DB) TodoRepository {
	return &TodoMariaDB{db: db}
}

func (t *TodoMariaDB) GetAllTodo() ([]models.Todo, error) {
	rows, err := t.db.Query("SELECT * FROM todos")
	if err != nil {
		return nil, err
	}

	var todos []models.Todo
	for rows.Next() {
		var id int64
		var title string
		var description string
		err = rows.Scan(&id, &title, &description)
		if err != nil {
			return nil, err
		}
		todos = append(todos, models.Todo{Id: id, Title: title, Description: description})
	}

	return todos, nil
}

func (t *TodoMariaDB) GetTodo(todo models.Todo) (models.Todo, error) {
	var id int64
	var title string
	var description string
	err := t.db.QueryRow("SELECT * FROM todos WHERE id = ?", todo.Id).Scan(&id, &title, &description)
	if err != nil {
		return models.Todo{}, err
	}

	return models.Todo{Id: id, Title: title, Description: description}, nil
}

func (t *TodoMariaDB) AddTodo(todo models.Todo) (models.Todo, error) {
	result, err := t.db.Exec("INSERT INTO todos (title, description) VALUES (?, ?)", todo.Title, todo.Description)
	if err != nil {
		return models.Todo{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return models.Todo{}, err
	}

	if rows != 1 {
		return models.Todo{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Todo{}, err
	}

	todo.Id = id

	return todo, nil
}

func (t *TodoMariaDB) DeleteTodo(todo models.Todo) (models.Todo, error) {
	result, err := t.db.Exec("DELETE FROM todos WHERE id = ?", todo.Id)
	if err != nil {
		return models.Todo{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return models.Todo{}, err
	}

	if rows != 1 {
		return models.Todo{}, err
	}

	return todo, nil
}

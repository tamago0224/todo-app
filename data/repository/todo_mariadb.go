package repository

import (
	"database/sql"

	"github.com/tamago0224/rest-app-backend/data/model"
)

type TodoMariaDB struct {
	db *sql.DB
}

func NewTodoMariaDBRepository(db *sql.DB) TodoRepository {
	return &TodoMariaDB{db: db}
}

func (t *TodoMariaDB) GetAllTodo(userId int) ([]model.Todo, error) {
	rows, err := t.db.Query("SELECT * FROM todos where user_id = ?", userId)
	if err != nil {
		return nil, err
	}

	todos := []model.Todo{}
	for rows.Next() {
		var id int64
		var userId int64
		var title string
		var description string
		err = rows.Scan(&id, &userId, &title, &description)
		if err != nil {
			return nil, err
		}
		todos = append(todos, model.Todo{Id: id, UserId: userId, Title: title, Description: description})
	}

	return todos, nil
}

func (t *TodoMariaDB) GetTodo(userID, todoID int) (model.Todo, error) {
	var id int64
	var userId int64
	var title string
	var description string
	err := t.db.QueryRow("SELECT * FROM todos WHERE id = ? AND user_id = ?", todoID, userID).Scan(&id, &userId, &title, &description)
	if err != nil {
		return model.Todo{}, err
	}

	return model.Todo{Id: id, Title: title, Description: description}, nil
}

func (t *TodoMariaDB) AddTodo(todo model.Todo) (model.Todo, error) {
	result, err := t.db.Exec("INSERT INTO todos (title, user_id, description) VALUES (?, ?, ?)", todo.Title, todo.UserId, todo.Description)
	if err != nil {
		return model.Todo{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return model.Todo{}, err
	}

	if rows != 1 {
		return model.Todo{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.Todo{}, err
	}

	todo.Id = id

	return todo, nil
}

func (t *TodoMariaDB) DeleteTodo(todo model.Todo) (model.Todo, error) {
	result, err := t.db.Exec("DELETE FROM todos WHERE id = ? AND user_id = ?", todo.Id, todo.UserId)
	if err != nil {
		return model.Todo{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return model.Todo{}, err
	}

	if rows != 1 {
		return model.Todo{}, err
	}

	return todo, nil
}

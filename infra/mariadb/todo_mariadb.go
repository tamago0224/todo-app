package mariadb

import (
	"database/sql"
	"time"

	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/domain/repository"
)

type TodoMariaDB struct {
	db *sql.DB
}

func NewTodoMariaDBRepository(db *sql.DB) repository.TodoRepository {
	return &TodoMariaDB{db: db}
}

func (t *TodoMariaDB) GetAllTodo(userId int) ([]model.Todo, error) {
	rows, err := t.db.Query("SELECT id, user_id, title, description, done, deadline, created FROM todos where user_id = ?", userId)
	if err != nil {
		return nil, err
	}

	todos := []model.Todo{}
	for rows.Next() {
		var id int64
		var userId int64
		var title string
		var description string
		var done bool
		var deadline *time.Time
		var created time.Time
		err = rows.Scan(&id, &userId, &title, &description, &done, &deadline, &created)
		if err != nil {
			return nil, err
		}
		todos = append(todos,
			model.Todo{
				Id:          id,
				UserId:      userId,
				Title:       title,
				Description: description,
				Done:        done,
				Deadline:    deadline,
				Created:     created,
			})
	}

	return todos, nil
}

func (t *TodoMariaDB) GetTodo(userID, todoID int) (model.Todo, error) {
	var id int64
	var userId int64
	var title string
	var description string
	var done bool
	var deadline *time.Time
	var created time.Time
	err := t.db.QueryRow("SELECT id, user_id, title, description, done, deadline, created FROM todos WHERE id = ? AND user_id = ?", todoID, userID).
		Scan(&id, &userId, &title, &description, &done, &deadline, &created)
	if err != nil {
		return model.Todo{}, err
	}

	return model.Todo{
		Id:          id,
		Title:       title,
		Description: description,
		Done:        done,
		Deadline:    deadline,
		Created:     created,
	}, nil
}

func (t *TodoMariaDB) AddTodo(todo model.Todo) (model.Todo, error) {
	result, err := t.db.Exec("INSERT INTO todos (user_id, title, description, done, deadline) VALUES (?, ?, ?, ?, ?)",
		todo.UserId, todo.Title, todo.Description, todo.Done, todo.Deadline)
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

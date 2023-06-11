package repository

import (
	"database/sql"
	"fmt"

	"github.com/tamago0224/rest-app-backend/model"
)

type UserMariaDB struct {
	db *sql.DB
}

func NewUserMariaDBRepository(db *sql.DB) UserRepository {
	return &UserMariaDB{db: db}
}

func (u *UserMariaDB) SearchUser(name string) (model.User, error) {
	var userId int64
	var userName string
	var userPassword string
	err := u.db.QueryRow("SELECT * FROM users WHERE name = ?", name).Scan(&userId, &userName, &userPassword)
	if err != nil {
		return model.User{}, err
	}

	return model.User{Id: userId, Name: userName, Password: userPassword}, nil
}

func (u *UserMariaDB) CreateUser(user model.User) (model.User, error) {
	result, err := u.db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password)
	if err != nil {
		return model.User{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return model.User{}, err
	}
	if rows != 1 {
		return model.User{}, fmt.Errorf("expected effected rows is 1, but %d", rows)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.User{}, err
	}

	user.Id = id
	return user, nil
}

package repository

import (
	"database/sql"
	"fmt"

	"github.com/tamago0224/rest-app-backend/models"
)

type UserMariaDB struct {
	db *sql.DB
}

func NewUserMariaDBRepository(db *sql.DB) UserRepository {
	return &UserMariaDB{db: db}
}

func (u *UserMariaDB) SearchUser(name string) (models.User, error) {
	var userId int64
	var userName string
	var userPassword string
	err := u.db.QueryRow("SELECT * FROM users WHERE name = ?", name).Scan(&userId, &userName, &userPassword)
	if err != nil {
		return models.User{}, err
	}

	return models.User{Id: userId, Name: userName, Password: userPassword}, nil
}

func (u *UserMariaDB) CreateUser(user models.User) (models.User, error) {
	result, err := u.db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password)
	if err != nil {
		return models.User{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return models.User{}, err
	}
	if rows != 1 {
		return models.User{}, fmt.Errorf("expected effected rows is 1, but %d", rows)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return models.User{}, err
	}

	user.Id = id
	return user, nil
}

package repository

import "github.com/tamago0224/rest-app-backend/models"

type UserRepository interface {
	SearchUser(name string) (models.User, error)
	CreateUser(models.User) (models.User, error)
}

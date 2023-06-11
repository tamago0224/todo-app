package repository

import "github.com/tamago0224/rest-app-backend/model"

type UserRepository interface {
	SearchUser(name string) (model.User, error)
	CreateUser(model.User) (model.User, error)
}

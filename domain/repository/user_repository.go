package repository

import "github.com/tamago0224/rest-app-backend/domain/model"

type UserRepository interface {
	SelectByName(name string) (model.User, error)
	SelectByID(userID int) (model.User, error)
	CreateUser(model.User) (model.User, error)
}

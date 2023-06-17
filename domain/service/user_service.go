package service

import (
	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/domain/repository"
)

type IUserService interface {
	FindByName(name string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) FindByName(name string) (model.User, error) {
	user, err := u.userRepo.SelectByName(name)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userService) CreateUser(user model.User) (model.User, error) {
	user, err := u.userRepo.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

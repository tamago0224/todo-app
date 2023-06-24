package usecase

import (
	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/domain/service"
)

type IUserUsecase interface {
	SearchUser(name string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	SearchOrCreate(name string) (model.User, error)
}

type userUsecase struct {
	svc service.IUserService
}

func NewUserUsecase(svc service.IUserService) IUserUsecase {
	return &userUsecase{svc: svc}
}

func (u *userUsecase) SearchUser(name string) (model.User, error) {
	user, err := u.svc.FindByName(name)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userUsecase) CreateUser(user model.User) (model.User, error) {
	user, err := u.svc.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userUsecase) SearchOrCreate(name string) (model.User, error) {
	user, err := u.svc.FindBynameOrCreate(name)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

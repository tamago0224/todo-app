package service

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/tamago0224/rest-app-backend/domain/model"
	"github.com/tamago0224/rest-app-backend/domain/repository"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")
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
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (u *userService) CreateUser(user model.User) (model.User, error) {
	user, err := u.userRepo.CreateUser(user)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError); mysqlError.Number == 1062 {
			return model.User{}, ErrUserAlreadyExist
		}
		return model.User{}, err
	}

	return user, nil
}

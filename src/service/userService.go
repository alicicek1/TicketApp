package service

import (
	"github.com/pkg/errors"
	"strings"
	"ticketApp/src/repository"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
)

type UserServiceType struct {
	UserRepository repository.UserRepository
}

type UserService interface {
	UserServiceInsert(user entity.User) (string, error)
	UserServiceGetById(id string) (*entity.User, error)
	UserServiceDeleteById(id string) (bool, error)
}

func NewUserService(r repository.UserRepository) UserServiceType {
	return UserServiceType{UserRepository: r}
}

func (u UserServiceType) UserServiceInsert(user entity.User) (string, error) {
	if !strings.Contains(user.Email, "@") {
		return "_", errors.New("User e-mail address must contains '@'.")
	}

	user.Password = util.GetMD5Hash(user.Password)
	//user.Password = md5.Sum([]byte(user.Password))
	result, err := u.UserRepository.UserRepoInsert(user)

	return result, err
}

func (u UserServiceType) UserServiceGetById(id string) (*entity.User, error) {
	result, err := u.UserRepository.UserRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u UserServiceType) UserServiceDeleteById(id string) (bool, error) {
	result, err := u.UserRepository.UserRepoDeleteById(id)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}

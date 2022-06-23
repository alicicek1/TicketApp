package service

import (
	"github.com/google/uuid"
	"ticketApp/src/repository"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
	"time"
)

type UserServiceType struct {
	UserRepository repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserServiceType {
	return UserServiceType{UserRepository: r}
}

type UserService interface {
	UserServiceInsert(user entity.User) (*entity.UserPostResponseModel, *util.Error)
	UserServiceGetById(id string) (*entity.User, *util.Error)
	UserServiceDeleteById(id string) (util.DeleteResponseType, *util.Error)
	UserServiceGetAll(filter util.Filter) (*entity.UserGetResponseModel, *util.Error)
}

func (u UserServiceType) UserServiceInsert(user entity.User) (*entity.UserPostResponseModel, *util.Error) {
	if user.Id == "" {
		isSuccess, err := util.CheckUserModel(user)
		if !isSuccess {
			return nil, err
		}
	}

	if user.Password != "" {
		user.Password = util.GetMD5Hash(user.Password)
	}
	user.Type = entity.DEFAULT

	user.CreatedAt = time.Now()
	if user.Id == "" {
		user.Id = uuid.New().String()
	}
	user.UpdatedAt = time.Now()

	result, err := u.UserRepository.UserRepoInsert(user)

	return result, err
}
func (u UserServiceType) UserServiceGetById(id string) (*entity.User, *util.Error) {
	result, err := u.UserRepository.UserRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (u UserServiceType) UserServiceDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	result, err := u.UserRepository.UserRepoDeleteById(id)
	if err != nil || result.IsSuccess == false {
		return util.DeleteResponseType{IsSuccess: false}, err
	}
	return util.DeleteResponseType{IsSuccess: true}, nil
}
func (u UserServiceType) UserServiceGetAll(filter util.Filter) (*entity.UserGetResponseModel, *util.Error) {
	result, err := u.UserRepository.UserRepositoryGetAll(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

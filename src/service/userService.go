package service

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
	"ticketApp/src/repository"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
)

type UserServiceType struct {
	UserRepository repository.UserRepository
}

type UserService interface {
	UserServiceInsert(user entity.User) (map[string]interface{}, *util.Error)
	UserServiceGetById(id string) (*entity.User, *util.Error)
	UserServiceDeleteById(id string) (bool, *util.Error)
}

func NewUserService(r repository.UserRepository) UserServiceType {
	return UserServiceType{UserRepository: r}
}

func (u UserServiceType) UserServiceInsert(user entity.User) (map[string]interface{}, *util.Error) {
	//isSuccess, err := checkUserModel(user)
	//if !isSuccess {
	//	return nil, err
	//}

	user.Password = getMD5Hash(user.Password)
	user.Type = entity.DEFAULT
	result, err := u.UserRepository.UserRepoInsert(&user)

	return result, err
}

func (u UserServiceType) UserServiceGetById(id string) (*entity.User, *util.Error) {
	result, err := u.UserRepository.UserRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u UserServiceType) UserServiceDeleteById(id string) (bool, *util.Error) {
	result, err := u.UserRepository.UserRepoDeleteById(id)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}

func checkUserModel(user entity.User) (bool, *util.Error) {
	if user.Username == "" {
		return false, &util.Error{
			ApplicationName: "service",
			Operation:       "POST",
			Description:     "Username cannot be null.",
			StatusCode:      http.StatusBadRequest,
			ErrorCode:       4004,
		}
	}

	if user.Password == "" {
		return false, &util.Error{
			ApplicationName: "service",
			Operation:       "POST",
			Description:     "Password cannot be null.",
			StatusCode:      http.StatusBadRequest,
			ErrorCode:       4005,
		}
	}

	if !strings.Contains(user.Email, "@") {
		return false, &util.Error{
			ApplicationName: "service",
			Operation:       "POST",
			Description:     "E-mail address must contains a '@'.",
			StatusCode:      http.StatusBadRequest,
			ErrorCode:       4006,
		}
	}
	if user.Type == 0 {
		return false, &util.Error{
			ApplicationName: "service",
			Operation:       "POST",
			Description:     "User type cannot be zero",
			StatusCode:      http.StatusBadRequest,
			ErrorCode:       4007,
		}
	}
	return true, nil

}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strconv"
	"strings"
	"ticketApp/src/type/entity"
	"time"
)

func CheckCategoryModel(category entity.Category) (bool, *Error) {
	return true, nil
}
func CheckTicketModel(ticket entity.Ticket) (bool, *Error) {
	return true, nil
}
func CheckUserModel(user entity.User) (bool, *Error) {
	if user.Username == "" {
		return false, PostValidation.ModifyApplicationName("user service").ModifyDescription("Username cannot be null.").ModifyErrorCode(4024)
	}

	if user.Password == "" {
		return false, PostValidation.ModifyApplicationName("user service").ModifyDescription("Password cannot be null.").ModifyErrorCode(4025)
	}

	if !strings.Contains(user.Email, "@") {
		return false, PostValidation.ModifyApplicationName("user service").ModifyDescription("E-mail address must contains a '@'.").ModifyErrorCode(4026)
	}
	if user.Type == 0 {
		return false, PostValidation.ModifyApplicationName("user service").ModifyDescription("User type cannot be zero.").ModifyErrorCode(4027)
	}
	return true, nil

}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

func IsValidUUID(uuidStr string) bool {
	if _, err := uuid.Parse(uuidStr); err != nil {
		return false
	}
	return true
}

func ValidatePaginationFilters(page, pageSize string, maxLimit int) (int64, int64) {
	max := int64(maxLimit)

	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil || pageInt < 0 {
		pageInt = 0
	}

	pageSizeInt, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil || pageSizeInt <= 0 {
		pageSizeInt = max
	}
	return pageInt, pageSizeInt
}

func ValidateSortingFilters(entity any, sortingArea, SortingDirection string) (string, int) {
	sort := ""
	var direction int

	if strings.EqualFold(SortingDirection, "asc") || strings.EqualFold(SortingDirection, "1") {
		direction = 1
	} else if strings.EqualFold(SortingDirection, "desc") || strings.EqualFold(SortingDirection, "dsc") || strings.EqualFold(SortingDirection, "-1") {
		direction = -1
	} else {
		direction = 0
	}

	v := reflect.ValueOf(entity)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if sortingArea == typeOfS.Field(i).Name {
			sort = strings.ToLower(typeOfS.Field(i).Name)
			break
		}
	}

	return sort, direction
}

func CreateEqualFilter(value, field string) interface{} {
	return bson.M{"$regex": primitive.Regex{
		Pattern: value,
		Options: "i",
	}}
}

var operators []string = []string{"<=", ">=", "==", "||"}

func CreateFilter(model any, filters string) map[string]interface{} {
	responseFilter := map[string]interface{}{}
	v := reflect.ValueOf(model)
	typeOfS := v.Type()

	filterArr := strings.Split(filters, ",")
	for _, filter := range filterArr {
		for i := 0; i <= len(operators); i++ {
			if strings.Contains(filter, operators[i]) {
				splittedCurrentFilter := strings.Split(filter, operators[i])
				field := splittedCurrentFilter[0]
				value := splittedCurrentFilter[len(splittedCurrentFilter)-1]
				for j := 0; j < v.NumField(); j++ {
					if field == typeOfS.Field(j).Name {
						if strings.Contains(filter, ">=") {
							responseFilter[field] = bson.M{"$gte": value}
							break
						} else if strings.Contains(filter, "<=") {
							if prevFilter, exist := responseFilter[field]; !exist {
								responseFilter[field] = bson.M{"$lte": value}
							} else {
								fmt.Println(prevFilter)
								responseFilter["$and"] = bson.A{
									bson.M{field: value},
									prevFilter,
								}
							}
							break
						} else if strings.Contains(filter, "==") {
							responseFilter[field] = bson.M{"$regex": primitive.Regex{
								Pattern: value,
								Options: "i",
							}}
							break
						} else if strings.Contains(filter, "||") {
							break
						}
					}
				}
			}
		}
	}
	return nil
}

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(model entity.LoginRequestModel) (*entity.LoginResponseModel, error) {
	response := entity.LoginResponseModel{}

	expirationDate := time.Now().Add(time.Minute * 5)
	response.ExpiresDate = expirationDate

	claims := &Claims{
		Username: model.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationDate.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	response.IsSuccessful = true
	response.Token = tokenStr
	return &response, nil
}

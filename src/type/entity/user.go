package entity

import (
	"time"
)

type UserType int16

const (
	ADMIN   UserType = 1
	DEFAULT          = 2
)

func (u UserType) String() string {
	switch u {
	case ADMIN:
		return "admin"
	case DEFAULT:
		return "default"
	}
	return "unknown"
}

type User struct {
	Id        string    `json:"_id" bson:"_id,omitempty"`
	Username  string    `json:"username,omitempty" bson:"username,omitempty"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
	Age       int32     `json:"age,omitempty" bson:"age,omitempty"`
	Type      byte      `json:"type,omitempty" bson:"type,omitempty"`
}

type UserPostRequestModel struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Type     byte   `json:"type,omitempty"`
	Age      int32  `json:"age,omitempty"`
}

type UserPostResponseModel struct {
	Id string `json:"_id" bson:"_id,omitempty"`
}

type UserGetResponseModel struct {
	RowCount int64  `json:"rowCount"`
	Users    []User `json:"users"`
}

type LoginRequestModel struct {
	Username string  `json:"username"`
	Password *string `json:"password"`
}

type LoginResponseModel struct {
	IsSuccessful bool      `json:"isSuccessful"`
	Token        string    `json:"token"`
	ExpiresDate  time.Time `json:"-"`
}

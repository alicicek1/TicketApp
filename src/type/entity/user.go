package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	ADMIN   = 1
	DEFAULT = 2
)

type User struct {
	Id        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string              `json:"username,omitempty"`
	Password  string              `json:"password,omitempty"`
	Email     string              `json:"email,omitempty"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
	Type      byte                `json:"type,omitempty"`
}

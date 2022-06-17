package entity

type User struct {
	BaseEntity
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Type     byte   `json:"type"`
}

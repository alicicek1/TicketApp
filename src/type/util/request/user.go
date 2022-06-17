package request

type UserPostRequestModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Type     byte   `json:"type"`
}

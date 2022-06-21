package util

type UserPostRequestModel struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Type     byte   `json:"type,omitempty"`
}

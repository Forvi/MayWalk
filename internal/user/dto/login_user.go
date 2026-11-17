package dto

type UserLogin struct {
	Login      string  `json:"login"`
	Password   *string `json:"password"`
	RememberMe bool    `json:"remember_me"`
}

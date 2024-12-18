package model

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterUser struct {
	Email    string `json:"email" validate:"required,email"`
	Name string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

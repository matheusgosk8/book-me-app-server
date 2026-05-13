package models

type ServerStatus struct {
	Code    int
	Message string
}

type RegisterResponse struct {
	User    *RegisterUserResponse `json:"user"`
	Token   string                `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	Code    int                   `json:"code"`
	Message string                `json:"message"`
}

// RegisterUserResponse holds the minimal user fields returned on register.
type RegisterUserResponse struct {
	Id    string `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

type LoginUserResponse struct {
	Id    string `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
}

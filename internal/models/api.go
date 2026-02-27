package models

type ServerStatus struct {
	Code    int
	Message string
}

type User struct {
	Name  string
	Email string
	Id    string
}

type RegisterResponse struct {
	User    *User
	Code    int
	Message string
}

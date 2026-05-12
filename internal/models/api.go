package models

type ServerStatus struct {
	Code    int
	Message string
}

// type User struct {
// 	Id            string `json:"id"`
// 	Cep           string `json:"cep"`
// 	Cidade        string `json:"cidade"`
// 	Cnpj          string `json:"cnpj"`
// 	ConfirmaSenha string `json:"confirmaSenha"`
// 	Cpf           string `json:"cpf"`
// 	Email         string `json:"email"`
// 	Estado        string `json:"estado"`
// 	Logradouro    string `json:"logradouro"`
// 	Nome          string `json:"nome"`
// 	Rua           string `json:"rua"`
// 	Senha         string `json:"senha"`
// 	Telefone      string `json:"telefone"`
// 	UserType      string `json:"userType"`
// }

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

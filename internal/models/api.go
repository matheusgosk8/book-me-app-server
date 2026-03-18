package models

type ServerStatus struct {
	Code    int
	Message string
}

type User struct {
	Id            string `json:"id"`
	Cep           string `json:"cep"`
	Cidade        string `json:"cidade"`
	Cnpj          string `json:"cnpj"`
	ConfirmaSenha string `json:"confirmaSenha"`
	Cpf           string `json:"cpf"`
	Email         string `json:"email"`
	Estado        string `json:"estado"`
	Logradouro    string `json:"logradouro"`
	Nome          string `json:"nome"`
	Rua           string `json:"rua"`
	Senha         string `json:"senha"`
	Telefone      string `json:"telefone"`
	UserType      string `json:"userType"`
}

type RegisterResponse struct {
    User    *User  `json:"user"`
    Token   string `json:"token"` 
    Code    int    `json:"code"`
    Message string `json:"message"`
}

package auth

import (
	"context"
	"fmt"

	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/models"
	repositories "github.com/matheusgosk8/book-me-server/internal/repository"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

type RegisterInput struct {
	User    models.User
	Address models.Address
}

type RegisterOutput struct {
	User         *models.RegisterUserResponse
	AccessToken  string
	RefreshToken string
}

func (s *UserService) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	// 1. Hash da senha
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.User.Senha), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("falha ao processar senha: %w", err)
	}
	input.User.Senha = string(hashed)

	// 2. Geração de Tokens
	accessToken, refreshToken, err := utils.GenerateTokens(input.User.Id, input.User.UserType)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar tokens: %w", err)
	}

	// 3. Persistência via Repository
	_, _, err = repositories.CreateUserWithAddress(ctx, db.Client, input.User, input.Address)
	if err != nil {
		return nil, err // Retornamos o erro original para o handler tratar com o PGErrorHandler
	}

	return &RegisterOutput{
		User: &models.RegisterUserResponse{
			Id:    input.User.Id,
			Nome:  input.User.Nome,
			Email: input.User.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

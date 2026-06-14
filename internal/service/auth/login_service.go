package auth

import (
	"context"
	"errors"
	"time"

	"github.com/matheusgosk8/book-me-server/ent/user"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/models"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct{}

func NewLoginService() *LoginService {
	return &LoginService{}
}

type LoginInput struct {
	Email string
	Senha string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	User         models.LoginUserResponse
}

func (s *LoginService) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// 1. Busca o usuário
	u, err := db.Client.User.
		Query().
		Where(user.EmailEQ(input.Email)).
		Only(ctx)

	if err != nil {
		return nil, errors.New("e-mail ou senha incorretos")
	}

	// 2. Valida a senha
	err = bcrypt.CompareHashAndPassword([]byte(u.Senha), []byte(input.Senha))
	if err != nil {
		return nil, errors.New("e-mail ou senha incorretos")
	}

	// 3. Gera Tokens
	accessToken, refreshToken, err := utils.GenerateTokens(u.ID.String(), u.UserType)
	if err != nil {
		return nil, err
	}

	// 4. Cria Sessão
	_, err = db.Client.Session.
		Create().
		SetRefreshToken(refreshToken).
		SetLastLoginAt(time.Now()).
		SetExpiresAt(time.Now().Add(7 * 24 * time.Hour)).
		SetUserID(u.ID).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: models.LoginUserResponse{
			Id:    u.ID.String(),
			Nome:  u.Nome,
			Email: u.Email,
			Role:  u.UserType,
		},
	}, nil
}

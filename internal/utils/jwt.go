package utils

import (
	"os"
	"time"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	
)

type CustomClaims struct {
	UserId string `json:"user_id"`
	Type   string `json: "type"` //access ou refresh
	jwt.RegisteredClaims
}

func GenerateTokens(userID string) (string, string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	//token temporário
	acessClaims := &CustomClaims{
		UserId: userID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	at, err := jwt.NewWithClaims(jwt.SigningMethodHS256, acessClaims).SignedString(secret)
	if err != nil {
		return "", "", err
	}

	//refresh token
	refreshClaims := &CustomClaims{
		UserId: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	rt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret)
	if err != nil {
		return "", "", err
	}
	return at, rt, nil
}

// ValidateToken verifica se o token é válido e retorna os Claims
func ValidateToken(tokenString string) (*CustomClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token inválido")
}

// ParseUUID é uma função auxiliar para converter string para uuid.UUID
func ParseUUID(id string) uuid.UUID {
	parsedID, _ := uuid.Parse(id)
	return parsedID
}

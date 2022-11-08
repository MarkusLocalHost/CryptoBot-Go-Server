package service

import (
	"cryptocurrency/internal/models"
	"cryptocurrency/internal/utils/tokens"
)

type tokenService struct {
}

type TSConfig struct {
}

func NewTokenService(c *TSConfig) models.TokenService {
	return &tokenService{}
}

func (t tokenService) ValidateIDToken(tokenString string) (string, error) {
	jsonString, err := tokens.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}
	return jsonString, nil
}

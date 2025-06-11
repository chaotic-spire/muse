package service

import (
	"backend/internal/domain/dto"
	"backend/internal/domain/utils/auth"
	"context"
	"github.com/spf13/viper"
	"time"
)

type tokenService struct {
}

func NewTokenService() *tokenService {
	return &tokenService{}
}

// GenerateToken is a method to generate a new token.
func (s *tokenService) GenerateToken(ctx context.Context, userID string, expires time.Time, tokenType string) (string, error) {
	jwtToken, err := auth.GenerateToken(userID, expires, tokenType)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// GenerateAuthTokens is a method to generate access and refresh tokens.
func (s *tokenService) GenerateAuthTokens(c context.Context, userID string) (*dto.AuthTokens, error) {
	authTokenExpiration := time.Now().UTC().Add(time.Minute * time.Duration(viper.GetInt("service.backend.jwt.access-token-expiration")))
	authToken, err := s.GenerateToken(
		c,
		userID,
		authTokenExpiration,
		auth.TokenTypeAccess,
	)
	if err != nil {
		return nil, err
	}

	refreshTokenExpiration := time.Now().UTC().Add(time.Minute * time.Duration(viper.GetInt("service.backend.jwt.refresh-token-expiration")))
	refreshToken, err := s.GenerateToken(
		c,
		userID,
		refreshTokenExpiration,
		auth.TokenTypeRefresh,
	)
	if err != nil {
		return nil, err
	}

	return &dto.AuthTokens{
		Access: dto.Token{
			Token:   authToken,
			Expires: authTokenExpiration,
		},
		Refresh: dto.Token{
			Token:   refreshToken,
			Expires: refreshTokenExpiration,
		},
	}, nil
}

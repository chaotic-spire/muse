package auth

import (
	"context"
	"time"
	"webTemplate/internal/domain/dto"
	"webTemplate/internal/domain/entity"
)

type TokenService interface {
	GenerateToken(ctx context.Context, userID string, expires time.Time, tokenType string) (string, error)
	GenerateAuthTokens(c context.Context, userID string) (*dto.AuthTokens, error)
}

type UserService interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Create(ctx context.Context, registerReq *dto.UserReturn) (*entity.User, error)
}

type authUseсase struct {
	tokenService TokenService
	userService  UserService
}

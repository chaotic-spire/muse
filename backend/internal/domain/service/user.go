package service

import (
	"backend/internal/domain/common/errorz"
	"backend/internal/domain/dto"
	"backend/internal/domain/entity"
	"context"
)

type userStorage interface {
	Create(ctx context.Context, user entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id string) error
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type userService struct {
	storage userStorage
}

func NewUserService(storage userStorage) *userService {
	return &userService{storage: storage}
}

func (s *userService) Create(ctx context.Context, registerReq dto.UserReturn) (*entity.User, error) {
	if user, err := s.storage.GetByUsername(ctx, registerReq.Username); err == nil && user != nil {
		return nil, errorz.UserAlreadyExists
	}

	user := entity.User{
		ID:        registerReq.ID,
		Username:  registerReq.Username,
		Firstname: registerReq.Firstname,
		Lastname:  registerReq.Lastname,
		PhotoUrl:  registerReq.PhotoUrl,
	}

	return s.storage.Create(ctx, user)
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return s.storage.GetByUsername(ctx, username)
}

func (s *userService) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	return s.storage.Update(ctx, user)
}

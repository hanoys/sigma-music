package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type UserService struct {
	repository ports.IUserRepository
}

func NewUserService(repo ports.IUserRepository) *UserService {
	return &UserService{repository: repo}
}

func (us *UserService) Register(ctx context.Context, user ports.UserServiceCreateRequest) (domain.User, error) {
	_, err := us.repository.GetByName(ctx, user.Name)
	if err == nil {
		return domain.User{}, ports.ErrUserWithSuchNameAlreadyExists
	}

	_, err = us.repository.GetByEmail(ctx, user.Email)
	if err == nil {
		return domain.User{}, ports.ErrUserWithSuchEmailAlreadyExists
	}

	_, err = us.repository.GetByPhone(ctx, user.Phone)
	if err == nil {
		return domain.User{}, ports.ErrUserWithSuchPhoneAlreadyExists
	}

	createUser := domain.User{
		ID:       uuid.New(),
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Country:  user.Country,
	}

	return us.repository.Create(ctx, createUser)
}

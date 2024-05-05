package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type UserService struct {
	repository ports.IUserRepository
	hash       ports.IHashPasswordProvider
}

func NewUserService(repo ports.IUserRepository, hash ports.IHashPasswordProvider) *UserService {
	return &UserService{repository: repo, hash: hash}
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

	saltedPassword := us.hash.EncodePassword(user.Password)

	createUser := domain.User{
		ID:       uuid.New(),
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: saltedPassword.HashPassword,
		Salt:     saltedPassword.Salt,
		Country:  user.Country,
	}

	return us.repository.Create(ctx, createUser)
}

func (us *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	return us.repository.GetAll(ctx)
}

func (us *UserService) GetById(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	return us.repository.GetByID(ctx, userID)
}

func (us *UserService) GetByName(ctx context.Context, name string) (domain.User, error) {
	return us.repository.GetByName(ctx, name)
}

func (us *UserService) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return us.repository.GetByEmail(ctx, email)
}

func (us *UserService) GetByPhone(ctx context.Context, phone string) (domain.User, error) {
	return us.repository.GetByPhone(ctx, phone)
}

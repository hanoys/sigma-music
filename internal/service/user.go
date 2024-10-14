package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type UserService struct {
	repository ports.IUserRepository
	hash       ports.IHashPasswordProvider
	logger     *zap.Logger
}

func NewUserService(repo ports.IUserRepository, hash ports.IHashPasswordProvider,
	logger *zap.Logger) *UserService {
	return &UserService{
		repository: repo,
		hash:       hash,
		logger:     logger,
	}
}

func (us *UserService) Register(ctx context.Context, user ports.UserServiceCreateRequest) (domain.User, error) {
	_, err := us.repository.GetByName(ctx, user.Name)
	if err == nil {
		us.logger.Error("Failed to register user", zap.Error(err), zap.String("User Name", user.Name))
		return domain.User{}, ports.ErrUserWithSuchNameAlreadyExists
	}

	_, err = us.repository.GetByEmail(ctx, user.Email)
	if err == nil {
		us.logger.Error("Failed to register user", zap.Error(err), zap.String("User Email", user.Email))
		return domain.User{}, ports.ErrUserWithSuchEmailAlreadyExists
	}

	_, err = us.repository.GetByPhone(ctx, user.Phone)
	if err == nil {
		us.logger.Error("Failed to register user", zap.Error(err), zap.String("User Phone", user.Phone))
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

	u, err := us.repository.Create(ctx, createUser)
	if err != nil {
		us.logger.Error("Failed to register user", zap.Error(err))
		return domain.User{}, err
	}

	us.logger.Info("User successfully registered", zap.String("User ID", u.ID.String()))

	return u, nil
}

func (us *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := us.repository.GetAll(ctx)
	if err != nil {
		us.logger.Error("Failed to get all users", zap.Error(err))
		return nil, err
	}

	return users, nil
}

func (us *UserService) GetById(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	u, err := us.repository.GetByID(ctx, userID)
	if err != nil {
		us.logger.Error("Failed to get user by ID", zap.Error(err),
			zap.String("User ID", userID.String()))

		return domain.User{}, err
	}

	us.logger.Info("User successfully received by ID", zap.String("User ID", userID.String()))

	return u, nil
}

func (us *UserService) GetByName(ctx context.Context, name string) (domain.User, error) {
	u, err := us.repository.GetByName(ctx, name)
	if err != nil {
		us.logger.Error("Failed to gt user by name", zap.Error(err),
			zap.String("User name", name))

		return domain.User{}, err
	}

	us.logger.Info("User successfully received by name", zap.String("User name", name))

	return u, nil
}

func (us *UserService) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := us.repository.GetByEmail(ctx, email)
	if err != nil {
		us.logger.Error("Failed to gt user by email", zap.Error(err),
			zap.String("User email", email))

		return domain.User{}, err
	}

	us.logger.Info("User successfully received by email", zap.String("User email", email))

	return u, nil
}

func (us *UserService) GetByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := us.repository.GetByPhone(ctx, phone)
	if err != nil {
		us.logger.Error("Failed to gt user by phone", zap.Error(err),
			zap.String("User phone", phone))

		return domain.User{}, err
	}

	us.logger.Info("User successfully received by phone", zap.String("User phone", phone))

	return u, nil
}

package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrUserDuplicate      = errors.New("user duplicate error")
	ErrUserIDNotFound     = errors.New("user with such id not found")
	ErrUserNameNotFound   = errors.New("user with such name doesn't exists")
	ErrUserEmailNotFound  = errors.New("user with such email doesn't exists")
	ErrUserPhoneNotFound  = errors.New("user with such email doesn't exists")
	ErrUserUnknownCountry = errors.New("such country doesn't exists")
	ErrInternalUserRepo   = errors.New("user repository internal error")
)

type IUserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
	GetByName(ctx context.Context, name string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByPhone(ctx context.Context, phone string) (domain.User, error)
}

var (
	ErrUserWithSuchNameAlreadyExists  = errors.New("user with such name already exists")
	ErrUserWithSuchEmailAlreadyExists = errors.New("user with such email already exists")
	ErrUserWithSuchPhoneAlreadyExists = errors.New("user with such phone already exists")
)

type UserServiceCreateRequest struct {
	Name     string
	Email    string
	Phone    string
	Password string
	Country  string
}

type IUserService interface {
	Register(ctx context.Context, user UserServiceCreateRequest) (domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	GetById(ctx context.Context, userID uuid.UUID) (domain.User, error)
	GetByName(ctx context.Context, name string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByPhone(ctx context.Context, phone string) (domain.User, error)
}

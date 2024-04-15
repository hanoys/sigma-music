package ports

import (
	"context"
	"errors"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrUserDuplicate     = errors.New("")
	ErrUserIDNotFound    = errors.New("user with such id not found")
	ErrUserNameNotFound  = errors.New("user with such name doesn't exists")
	ErrUserEmailNotFound = errors.New("user with such email doesn't exists")
	ErrUserPhoneNotFound = errors.New("user with such email doesn't exists")
	ErrInternalUserRepo  = errors.New("user repository internal error")
)

var (
	ErrUserWithSuchNameAlreadyExists  = errors.New("user with such name already exists")
	ErrUserWithSuchEmailAlreadyExists = errors.New("user with such email already exists")
	ErrUserWithSuchPhoneAlreadyExists = errors.New("user with such phone already exists")
)

type IUserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByName(ctx context.Context, name string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByPhone(ctx context.Context, phone string) (domain.User, error)
}

type UserServiceCreateRequest struct {
	Name     string
	Email    string
	Phone    string
	Password string
	Country  string
}

type IUserService interface {
	Register(ctx context.Context, user UserServiceCreateRequest) (domain.User, error)
}

package ports

import (
	"context"
	"errors"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrMusicianWithSuchNameAlreadyExists  = errors.New("user with such name already exists")
	ErrMusicianWithSuchEmailAlreadyExists = errors.New("user with such email already exists")
	ErrMusicianWithSuchPhoneAlreadyExists = errors.New("user with such phone already exists")
	ErrMusicianRegister                   = errors.New("can't register user: internal error")
)

type IMusicianRepository interface {
	Create(ctx context.Context, user domain.Musician) (domain.Musician, error)
	GetByName(ctx context.Context, name string) (domain.Musician, error)
	GetByEmail(ctx context.Context, email string) (domain.Musician, error)
}

type MusicianServiceCreateRequest struct {
	Name        string
	Email       string
	Password    string
	Country     string
	Description string
}

type IMusicianService interface {
	Register(ctx context.Context, user UserServiceCreateRequest) (domain.Musician, error)
}

package ports

import (
	"context"
	"errors"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrMusicianDuplicate     = errors.New("")
	ErrMusicianIDNotFound    = errors.New("musician with such id not found")
	ErrMusicianNameNotFound  = errors.New("musician with such name doesn't exists")
	ErrMusicianEmailNotFound = errors.New("musician with such email doesn't exists")
	ErrInternalMusicianRepo  = errors.New("musician repository internal error")
)

var (
	ErrMusicianWithSuchNameAlreadyExists  = errors.New("musician with such name already exists")
	ErrMusicianWithSuchEmailAlreadyExists = errors.New("musician with such email already exists")
)

type IMusicianRepository interface {
	Create(ctx context.Context, musician domain.Musician) (domain.Musician, error)
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
	Register(ctx context.Context, musician UserServiceCreateRequest) (domain.Musician, error)
}

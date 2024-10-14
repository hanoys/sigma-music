package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrMusicianDuplicate      = errors.New("musician duplicate error")
	ErrMusicianIDNotFound     = errors.New("musician with such id not found")
	ErrMusicianNameNotFound   = errors.New("musician with such name doesn't exists")
	ErrMusicianEmailNotFound  = errors.New("musician with such email doesn't exists")
	ErrMusicianUnknownCountry = errors.New("such country doesn't exists")
	ErrInternalMusicianRepo   = errors.New("musician repository internal error")
)

type IMusicianRepository interface {
	Create(ctx context.Context, musician domain.Musician) (domain.Musician, error)
	GetAll(ctx context.Context) ([]domain.Musician, error)
	GetByID(ctx context.Context, musicianID uuid.UUID) (domain.Musician, error)
	GetByName(ctx context.Context, name string) (domain.Musician, error)
	GetByEmail(ctx context.Context, email string) (domain.Musician, error)
	GetByAlbumID(ctx context.Context, albumID uuid.UUID) (domain.Musician, error)
	GetByTrackID(ctx context.Context, trackID uuid.UUID) (domain.Musician, error)
}

var (
	ErrMusicianWithSuchNameAlreadyExists  = errors.New("musician with such name already exists")
	ErrMusicianWithSuchEmailAlreadyExists = errors.New("musician with such email already exists")
)

type MusicianServiceCreateRequest struct {
	Name        string
	Email       string
	Password    string
	Country     string
	Description string
}

type IMusicianService interface {
	Register(ctx context.Context, musician MusicianServiceCreateRequest) (domain.Musician, error)
	GetAll(ctx context.Context) ([]domain.Musician, error)
	GetByID(ctx context.Context, musicianID uuid.UUID) (domain.Musician, error)
	GetByName(ctx context.Context, name string) (domain.Musician, error)
	GetByEmail(ctx context.Context, email string) (domain.Musician, error)
	GetByAlbumID(ctx context.Context, albumID uuid.UUID) (domain.Musician, error)
	GetByTrackID(ctx context.Context, trackID uuid.UUID) (domain.Musician, error)
}

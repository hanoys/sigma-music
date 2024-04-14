package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrGenreIDNotFound   = errors.New("genre with such id not found")
	ErrGenreNotFound     = errors.New("can't find any genre")
	ErrInternalGenreRepo = errors.New("internal track repository error")
)

type IGenreRepository interface {
	GetAll(ctx context.Context) ([]domain.Genre, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Genre, error)
	AddForTrack(ctx context.Context, trackID uuid.UUID, genresID []uuid.UUID) error
}

type IGenreService interface {
	GetAll(ctx context.Context) ([]domain.Genre, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Genre, error)
	AddForTrack(ctx context.Context, trackID uuid.UUID, genresID []uuid.UUID) error
}

package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrGenreGetAll      = errors.New("genre: can't get all genres")
	ErrGenreGetByID     = errors.New("genre: can't get by id")
	ErrGenreSetForTrack = errors.New("genre: can't get by id")
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

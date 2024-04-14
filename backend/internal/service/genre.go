package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type GenreService struct {
	repository ports.IGenreRepository
}

func NewGenreService(repo ports.IGenreRepository) *GenreService {
	return &GenreService{repository: repo}
}

func (gs *GenreService) GetAll(ctx context.Context) ([]domain.Genre, error) {
	return gs.repository.GetAll(ctx)
}

func (gs *GenreService) GetByID(ctx context.Context, id uuid.UUID) (domain.Genre, error) {
	return gs.repository.GetByID(ctx, id)
}

func (gs *GenreService) AddForTrack(ctx context.Context, trackID uuid.UUID, genreID []uuid.UUID) error {
	return gs.repository.AddForTrack(ctx, trackID, genreID)
}

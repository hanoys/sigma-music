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
	genres, err := gs.repository.GetAll(ctx)
	if err != nil {
		return nil, ports.ErrGenreGetAll
	}

	return genres, nil
}

func (gs *GenreService) GetByID(ctx context.Context, id uuid.UUID) (domain.Genre, error) {
	genre, err := gs.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Genre{}, ports.ErrGenreGetByID
	}

	return genre, nil
}

func (gs *GenreService) AddForTrack(ctx context.Context, trackID uuid.UUID, genreID []uuid.UUID) error {
	// TODO: repository returns error genre id
	err := gs.repository.AddForTrack(ctx, trackID, genreID)
	if err != nil {
		return ports.ErrGenreSetForTrack
	}

	return nil
}

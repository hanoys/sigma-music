package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
	"strings"
)

type GenreService struct {
	repository ports.IGenreRepository
	logger     *zap.Logger
}

func NewGenreService(repo ports.IGenreRepository, logger *zap.Logger) *GenreService {
	return &GenreService{
		repository: repo,
		logger:     logger,
	}
}

func (gs *GenreService) GetAll(ctx context.Context) ([]domain.Genre, error) {
	genres, err := gs.repository.GetAll(ctx)
	if err != nil {
		gs.logger.Error("Failed to get all genres", zap.Error(err))
		return nil, err
	}

	return genres, nil
}

func (gs *GenreService) GetByID(ctx context.Context, id uuid.UUID) (domain.Genre, error) {
	genre, err := gs.repository.GetByID(ctx, id)
	if err != nil {
		gs.logger.Error("Failed to get genre by id", zap.Error(err), zap.String("Genre ID", id.String()))
		return domain.Genre{}, err
	}

	gs.logger.Info("Genre successfully received by ID", zap.String("Genre ID", id.String()))

	return genre, nil
}

func (gs *GenreService) AddForTrack(ctx context.Context, trackID uuid.UUID, genreID []uuid.UUID) error {
	err := gs.repository.AddForTrack(ctx, trackID, genreID)
	if err != nil {
		genresID := make([]string, len(genreID))
		for i, genre := range genreID {
			genresID[i] = genre.String()
		}

		gs.logger.Error("Failed to add genres for track", zap.Error(err),
			zap.String("Track ID", trackID.String()),
			zap.String("Genres IDs", strings.Join(genresID, " ")))

		return err
	}

	gs.logger.Info("Genres successfully added for track", zap.String("Track ID", trackID.String()))

	return nil
}
func (gs *GenreService) GetByTrackID(ctx context.Context, trackID uuid.UUID) ([]domain.Genre, error) {
	genres, err := gs.repository.GetByTrackID(ctx, trackID)
	if err != nil {
		gs.logger.Error("Failed to get track genres", zap.Error(err), zap.String("Track ID", trackID.String()))
		return nil, err
	}

	gs.logger.Info("Genres successfully received by track ID", zap.String("Track ID", trackID.String()))

	return genres, nil
}

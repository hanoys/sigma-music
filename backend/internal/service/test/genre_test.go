package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

func TestGenreServiceGetAll(t *testing.T) {
	tests := []struct {
		name           string
		repositoryMock func(repository *mocks.GenreRepository)
		expected       error
	}{
		{
			name: "get all success",
			repositoryMock: func(repository *mocks.GenreRepository) {
				repository.
					On("GetAll", context.Background()).
					Return([]domain.Genre{}, nil)
			},
			expected: nil,
		},
		{
			name: "genre not found",
			repositoryMock: func(repository *mocks.GenreRepository) {
				repository.
					On("GetAll", context.Background()).
					Return(nil, ports.ErrGenreNotFound)
			},
			expected: ports.ErrGenreNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			genreRepository := mocks.NewGenreRepository(t)
			genreService := service.NewGenreService(genreRepository, logger)
			test.repositoryMock(genreRepository)

			_, err := genreService.GetAll(context.Background())
			if !errors.Is(err, test.expected) {
				t.Errorf("got %v, want %v", err, test.expected)
			}
		})
	}
}

func TestGenreServiceGetByID(t *testing.T) {
	tests := []struct {
		name           string
		id             uuid.UUID
		repositoryMock func(repository *mocks.GenreRepository)
		expected       error
	}{
		{
			name: "get by id success",
			repositoryMock: func(repository *mocks.GenreRepository) {
				repository.
					On("GetByID", context.Background(), mock.AnythingOfType("uuid.UUID")).
					Return(domain.Genre{}, nil)
			},
			expected: nil,
		},
		{
			name: "genre not found",
			repositoryMock: func(repository *mocks.GenreRepository) {
				repository.
					On("GetByID", context.Background(), mock.AnythingOfType("uuid.UUID")).
					Return(domain.Genre{}, ports.ErrInternalGenreRepo)
			},
			expected: ports.ErrInternalGenreRepo,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			genreRepository := mocks.NewGenreRepository(t)
			genreService := service.NewGenreService(genreRepository, logger)
			test.repositoryMock(genreRepository)

			_, err := genreService.GetByID(context.Background(), uuid.New())
			if !errors.Is(err, test.expected) {
				t.Errorf("got: %v, expected: %v", err, test.expected)
			}
		})
	}
}

func TestGenreServiceAddForTrack(t *testing.T) {
	tests := []struct {
		name           string
		trackID        uuid.UUID
		genresID       []uuid.UUID
		repositoryMock func(repository *mocks.GenreRepository)
		expected       error
	}{
		{
			name:     "add for track success",
			trackID:  uuid.New(),
			genresID: []uuid.UUID{uuid.New()},
			repositoryMock: func(repository *mocks.GenreRepository) {
				repository.
					On("AddForTrack", context.Background(), mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("[]uuid.UUID")).
					Return(nil)
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			genreRepository := mocks.NewGenreRepository(t)
			genreService := service.NewGenreService(genreRepository, logger)
			test.repositoryMock(genreRepository)

			err := genreService.AddForTrack(context.Background(), test.trackID, test.genresID)
			if !errors.Is(err, test.expected) {
				t.Errorf("got: %v, expected: %v", err, test.expected)
			}
		})
	}
}

package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/zap"
)

type GenreSuite struct {
	suite.Suite
	logger *zap.Logger
}

func (s *GenreSuite) BeforeEach(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()
}

type GenreGetAllSuite struct {
	GenreSuite
}

func (s *GenreGetAllSuite) CorrectRepositoryMock(repository *mocks.GenreRepository) {
	repository.
		On("GetAll", context.Background()).
		Return([]domain.Genre{}, nil)
}

func (s *GenreGetAllSuite) TestCorrect(t provider.T) {
	t.Parallel()
	repository := mocks.NewGenreRepository(t)
	s.CorrectRepositoryMock(repository)
	genreService := service.NewGenreService(repository, s.logger)

	_, err := genreService.GetAll(context.Background())

	t.Assert().Nil(err)
}

func (s *GenreGetAllSuite) InternalErrorRepositoryMock(repository *mocks.GenreRepository) {
	repository.
		On("GetAll", context.Background()).
		Return(nil, ports.ErrInternalGenreRepo)
}

func (s *GenreGetAllSuite) TestInternalError(t provider.T) {
	t.Parallel()
	repository := mocks.NewGenreRepository(t)
	s.InternalErrorRepositoryMock(repository)
	genreService := service.NewGenreService(repository, s.logger)

	_, err := genreService.GetAll(context.Background())

	t.Assert().ErrorIs(err, ports.ErrInternalGenreRepo)
}

func TestGenreGetAllSuite(t *testing.T) {
	suite.RunSuite(t, new(GenreGetAllSuite))
}

type GenreGetByIDSuite struct {
	GenreSuite
}

func (s *GenreGetByIDSuite) CorrectRepositoryMock(repository *mocks.GenreRepository, genreID uuid.UUID) {
	repository.
		On("GetByID", context.Background(), genreID).
		Return(domain.Genre{}, nil)
}

func (s *GenreGetByIDSuite) TestCorrect(t provider.T) {
	t.Parallel()
	genreID := uuid.New()
	repository := mocks.NewGenreRepository(t)
	s.CorrectRepositoryMock(repository, genreID)
	genreService := service.NewGenreService(repository, s.logger)

	_, err := genreService.GetByID(context.Background(), genreID)

	t.Assert().Nil(err)
}

func (s *GenreGetByIDSuite) NotFoundRepositoryMock(repository *mocks.GenreRepository, genreID uuid.UUID) {
	repository.
		On("GetByID", context.Background(), genreID).
		Return(domain.Genre{}, ports.ErrGenreNotFound)
}

func (s *GenreGetByIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	genreID := uuid.New()
	repository := mocks.NewGenreRepository(t)
	s.NotFoundRepositoryMock(repository, genreID)
	genreService := service.NewGenreService(repository, s.logger)

	_, err := genreService.GetByID(context.Background(), genreID)

	t.Assert().ErrorIs(err, ports.ErrGenreNotFound)
}

func TestGenreGetByIDSuite(t *testing.T) {
	suite.RunSuite(t, new(GenreGetByIDSuite))
}

type GenreAddForTrackSuite struct {
	GenreSuite
}

func (s *GenreAddForTrackSuite) CorrectRepositoryMock(repository *mocks.GenreRepository, trackID uuid.UUID, genreID []uuid.UUID) {
	repository.
		On("AddForTrack", context.Background(), trackID, genreID).
		Return(nil)
}

func (s *GenreAddForTrackSuite) TestCorrect(t provider.T) {
	t.Parallel()
	trackID := uuid.New()
	genreID := []uuid.UUID{uuid.New()}
	repository := mocks.NewGenreRepository(t)
	s.CorrectRepositoryMock(repository, trackID, genreID)
	genreService := service.NewGenreService(repository, s.logger)

	err := genreService.AddForTrack(context.Background(), trackID, genreID)

	t.Assert().Nil(err)
}

func (s *GenreAddForTrackSuite) InternalErrorRepositoryMock(repository *mocks.GenreRepository, trackID uuid.UUID, genreID []uuid.UUID) {
	repository.
		On("AddForTrack", context.Background(), trackID, genreID).
		Return(ports.ErrInternalGenreRepo)
}

func (s *GenreAddForTrackSuite) TestInternalError(t provider.T) {
	t.Parallel()
	trackID := uuid.New()
	genreID := []uuid.UUID{uuid.New()}
	repository := mocks.NewGenreRepository(t)
	s.InternalErrorRepositoryMock(repository, trackID, genreID)
	genreService := service.NewGenreService(repository, s.logger)

	err := genreService.AddForTrack(context.Background(), trackID, genreID)

	t.Assert().ErrorIs(err, ports.ErrInternalGenreRepo)
}

func TestGenreAddForTrackSuite(t *testing.T) {
	suite.RunSuite(t, new(GenreAddForTrackSuite))
}

type GenreGetByTrackIDSuite struct {
	GenreSuite
}

func (s *GenreGetByTrackIDSuite) CorrectRepositoryMock(repository *mocks.GenreRepository, trackID uuid.UUID) {
	repository.
		On("GetByTrackID", context.Background(), trackID).
		Return(make([]domain.Genre, 0), nil)
}

func (s *GenreGetByTrackIDSuite) TestCorrect(t provider.T) {
	t.Parallel()
	trackID := uuid.New()
	repository := mocks.NewGenreRepository(t)
	s.CorrectRepositoryMock(repository, trackID)
	genreService := service.NewGenreService(repository, s.logger)

	genres, err := genreService.GetByTrackID(context.Background(), trackID)

	t.Assert().NotNil(genres)
	t.Assert().Nil(err)
}

func (s *GenreGetByTrackIDSuite) InternalErrorRepositoryMock(repository *mocks.GenreRepository, trackID uuid.UUID) {
	repository.
		On("GetByTrackID", context.Background(), trackID).
		Return(nil, ports.ErrInternalGenreRepo)
}

func (s *GenreGetByTrackIDSuite) TestInternalError(t provider.T) {
	t.Parallel()
	trackID := uuid.New()
	repository := mocks.NewGenreRepository(t)
	s.InternalErrorRepositoryMock(repository, trackID)
	genreService := service.NewGenreService(repository, s.logger)

	genres, err := genreService.GetByTrackID(context.Background(), trackID)

	t.Assert().Nil(genres)
	t.Assert().ErrorIs(err, ports.ErrInternalGenreRepo)
}

func TestGenreGetByTrackIDSuite(t *testing.T) {
	suite.RunSuite(t, new(GenreGetByTrackIDSuite))
}

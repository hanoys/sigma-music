package test

import (
	"context"
	"net/url"
	"testing"

	"github.com/google/uuid"
	mocks2 "github.com/hanoys/sigma-music/internal/adapters/miniostorage/mocks"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/hanoys/sigma-music/internal/service/test/builder"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type TrackSuite struct {
	suite.Suite
	logger *zap.Logger
}

func (s *TrackSuite) BeforeEach(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()
}

type TrackCreateSuite struct {
	TrackSuite
}

func (s *TrackCreateSuite) CorrectRepositoryMock(trackRepository *mocks.TrackRepository, trackStorage *mocks2.TrackObjectStorage, genreRepository *mocks.GenreRepository, track domain.Track, genreID []uuid.UUID) {
	trackStorage.
		On("PutTrack", context.Background(), mock.Anything).
		Return(url.URL{}, nil)

	trackRepository.
		On("Create", context.Background(), mock.AnythingOfType("domain.Track")).
		Return(track, nil)

	genreRepository.
		On("AddForTrack", context.Background(), mock.Anything, genreID).
		Return(nil)

}

func (s *TrackCreateSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Track create test correct")
	track := builder.NewTrackBuilder().Default().Build()
	genreID := []uuid.UUID{uuid.New()}
	createReq := builder.NewCreateTrackRequestBuilder().Default().SetGenresID(genreID).Build()
	genreRepository := mocks.NewGenreRepository(t)
	trackRepository := mocks.NewTrackRepository(t)
	trackStorage := mocks2.NewTrackObjectStorage(t)
	s.CorrectRepositoryMock(trackRepository, trackStorage, genreRepository, track, genreID)
	genreService := service.NewGenreService(genreRepository, s.logger)
	trackService := service.NewTrackService(trackRepository, trackStorage, genreService, s.logger)

	serviceTrack, err := trackService.Create(context.Background(), createReq)

	t.Assert().Equal(track, serviceTrack)
	t.Assert().Nil(err)
}

func (s *TrackCreateSuite) TrackDuplicateRepositoryMock(trackRepository *mocks.TrackRepository, trackStorage *mocks2.TrackObjectStorage, genreRepository *mocks.GenreRepository, track domain.Track, genreID []uuid.UUID) {
	trackStorage.
		On("PutTrack", context.Background(), mock.Anything).
		Return(url.URL{}, nil)

	trackRepository.
		On("Create", context.Background(), mock.AnythingOfType("domain.Track")).
		Return(domain.Track{}, ports.ErrTrackDuplicate)
}

func (s *TrackCreateSuite) TestTrackDuplicate(t provider.T) {
	t.Parallel()
	t.Title("Track create test duplicate")
	track := builder.NewTrackBuilder().Default().Build()
	genreID := []uuid.UUID{uuid.New()}
	createReq := builder.NewCreateTrackRequestBuilder().Default().SetGenresID(genreID).Build()
	genreRepository := mocks.NewGenreRepository(t)
	trackRepository := mocks.NewTrackRepository(t)
	trackStorage := mocks2.NewTrackObjectStorage(t)
	s.TrackDuplicateRepositoryMock(trackRepository, trackStorage, genreRepository, track, genreID)
	genreService := service.NewGenreService(genreRepository, s.logger)
	trackService := service.NewTrackService(trackRepository, trackStorage, genreService, s.logger)

	_, err := trackService.Create(context.Background(), createReq)

	t.Assert().ErrorIs(err, ports.ErrTrackDuplicate)
}

func TestTrackCreateSuite(t *testing.T) {
	suite.RunSuite(t, new(TrackCreateSuite))
}

type TrackGetAllSuite struct {
	TrackSuite
}

func (s *TrackGetAllSuite) CorrectRepositoryMock(repository *mocks.TrackRepository) {
	repository.
		On("GetAll", context.Background()).
		Return(make([]domain.Track, 0), nil)
}

func (s *TrackGetAllSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Track get all test correct")
	genreRepository := mocks.NewGenreRepository(t)
	trackRepository := mocks.NewTrackRepository(t)
	trackStorage := mocks2.NewTrackObjectStorage(t)
	s.CorrectRepositoryMock(trackRepository)
	genreService := service.NewGenreService(genreRepository, s.logger)
	trackService := service.NewTrackService(trackRepository, trackStorage, genreService, s.logger)

	tracks, err := trackService.GetAll(context.Background())

	t.Assert().NotNil(tracks)
	t.Assert().Nil(err)
}

func (s *TrackGetAllSuite) InternalErrorRepositoryMock(repository *mocks.TrackRepository) {
	repository.
		On("GetAll", context.Background()).
		Return(nil, ports.ErrInternalTrackRepo)
}

func (s *TrackGetAllSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Track get all test internal error")
	genreRepository := mocks.NewGenreRepository(t)
	trackRepository := mocks.NewTrackRepository(t)
	trackStorage := mocks2.NewTrackObjectStorage(t)
	s.InternalErrorRepositoryMock(trackRepository)
	genreService := service.NewGenreService(genreRepository, s.logger)
	trackService := service.NewTrackService(trackRepository, trackStorage, genreService, s.logger)

	tracks, err := trackService.GetAll(context.Background())

	t.Assert().Nil(tracks)
	t.Assert().ErrorIs(err, ports.ErrInternalTrackRepo)
}

func TestTrackGetAllSuite(t *testing.T) {
	suite.RunSuite(t, new(TrackGetAllSuite))
}

type TrackGetByIDSuite struct {
	TrackSuite
}

func (s *TrackGetByIDSuite) CorrectRepositoryMock(repository *mocks.TrackRepository, track domain.Track) {
	repository.
		On("GetByID", context.Background(), track.ID).
		Return(track, nil)
}

func (s *TrackGetByIDSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Track get by id test correct")
	track := builder.NewTrackBuilder().Default().Build()
	genreRepository := mocks.NewGenreRepository(t)
	trackRepository := mocks.NewTrackRepository(t)
	trackStorage := mocks2.NewTrackObjectStorage(t)
	s.CorrectRepositoryMock(trackRepository, track)
	genreService := service.NewGenreService(genreRepository, s.logger)
	trackService := service.NewTrackService(trackRepository, trackStorage, genreService, s.logger)

	tracks, err := trackService.GetByID(context.Background(), track.ID)

	t.Assert().NotNil(tracks)
	t.Assert().Nil(err)
}

func (s *TrackGetByIDSuite) NotFoundRepositoryMock(repository *mocks.TrackRepository, track domain.Track) {
	repository.
		On("GetByID", context.Background(), track.ID).
		Return(domain.Track{}, ports.ErrTrackIDNotFound)
}

func (s *TrackGetByIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	t.Title("Track get by id test not found")
	track := builder.NewTrackBuilder().Default().Build()
	genreRepository := mocks.NewGenreRepository(t)
	trackRepository := mocks.NewTrackRepository(t)
	trackStorage := mocks2.NewTrackObjectStorage(t)
	s.NotFoundRepositoryMock(trackRepository, track)
	genreService := service.NewGenreService(genreRepository, s.logger)
	trackService := service.NewTrackService(trackRepository, trackStorage, genreService, s.logger)

	_, err := trackService.GetByID(context.Background(), track.ID)

	t.Assert().ErrorIs(err, ports.ErrTrackIDNotFound)
}

func TestTrackGetByIDSuite(t *testing.T) {
	suite.RunSuite(t, new(TrackGetByIDSuite))
}

type TrackDeleteSuite struct {
	TrackSuite
}

func (s *TrackDeleteSuite) CorrectRepositoryMock(repository *mocks.TrackRepository, trackStorage *mocks2.TrackObjectStorage, track domain.Track) {
	repository.
		On("Delete", context.Background(), track.ID).
		Return(track, nil)

	trackStorage.
		On("DeleteTrack", context.Background(), track.ID).
		Return(nil)
}

func (s *TrackDeleteSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Track delete test correct")
	track := builder.NewTrackBuilder().Default().Build()
	genreRepository := mocks.NewGenreRepository(t)
	trackRepository := mocks.NewTrackRepository(t)
	trackStorage := mocks2.NewTrackObjectStorage(t)
	s.CorrectRepositoryMock(trackRepository, trackStorage, track)
	genreService := service.NewGenreService(genreRepository, s.logger)
	trackService := service.NewTrackService(trackRepository, trackStorage, genreService, s.logger)

	serviceTrack, err := trackService.Delete(context.Background(), track.ID)

	t.Assert().Equal(track, serviceTrack)
	t.Assert().Nil(err)
}

func (s *TrackDeleteSuite) NotFoundRepositoryMock(repository *mocks.TrackRepository, trackStorage *mocks2.TrackObjectStorage, track domain.Track) {
	repository.
		On("Delete", context.Background(), track.ID).
		Return(domain.Track{}, ports.ErrTrackIDNotFound)
}

func (s *TrackDeleteSuite) TestNotFound(t provider.T) {
	t.Parallel()
	t.Title("Track delete test not found")
	track := builder.NewTrackBuilder().Default().Build()
	genreRepository := mocks.NewGenreRepository(t)
	trackRepository := mocks.NewTrackRepository(t)
	trackStorage := mocks2.NewTrackObjectStorage(t)
	s.NotFoundRepositoryMock(trackRepository, trackStorage, track)
	genreService := service.NewGenreService(genreRepository, s.logger)
	trackService := service.NewTrackService(trackRepository, trackStorage, genreService, s.logger)

	_, err := trackService.Delete(context.Background(), track.ID)

	t.Assert().ErrorIs(err, ports.ErrTrackIDNotFound)
}

func TestTrackDeleteSuite(t *testing.T) {
	suite.RunSuite(t, new(TrackDeleteSuite))
}

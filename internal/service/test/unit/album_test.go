package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
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

type AlbumSuite struct {
	suite.Suite
	logger *zap.Logger
}

func (s *AlbumSuite) BeforeEach(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()
}

type AlbumCreateSuite struct {
	AlbumSuite
}

func (s *AlbumCreateSuite) CorrectRepositoryMock(repository *mocks.AlbumRepository, musicianID uuid.UUID) {
	repository.
		On("Create", context.Background(), mock.AnythingOfType("domain.Album"), musicianID).
		Return(domain.Album{}, nil)
}

func (s *AlbumCreateSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Album create test correct")
	req := builder.NewCreateAlbumServiceRequestBuilder().Default().Build()
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	s.CorrectRepositoryMock(repository, req.MusicianID)

	_, err := albumService.Create(context.Background(), req)

	t.Assert().Nil(err)
}

func (s *AlbumCreateSuite) DuplicateRepositoryMock(repository *mocks.AlbumRepository, musicianID uuid.UUID) {
	repository.
		On("Create", context.Background(), mock.AnythingOfType("domain.Album"), musicianID).
		Return(domain.Album{}, ports.ErrAlbumDuplicate)
}

func (s *AlbumCreateSuite) TestDuplicate(t provider.T) {
	t.Parallel()
	t.Title("Album create test duplicate")
	req := builder.NewCreateAlbumServiceRequestBuilder().Default().Build()
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	s.DuplicateRepositoryMock(repository, req.MusicianID)

	_, err := albumService.Create(context.Background(), req)

	t.Assert().ErrorIs(err, ports.ErrAlbumDuplicate)
}

func TestAlbumCreateSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumCreateSuite))
}

type AlbumPublishSuite struct {
	AlbumSuite
}

func (s *AlbumPublishSuite) CorrectRepositoryMock(repository *mocks.AlbumRepository, albumID uuid.UUID) {
	repository.
		On("Publish", context.Background(), albumID).
		Return(nil)
}

func (s *AlbumPublishSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Album publish test correct")
	albumID := uuid.New()
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	s.CorrectRepositoryMock(repository, albumID)

	err := albumService.Publish(context.Background(), albumID)

	t.Assert().Nil(err)
}

func (s *AlbumPublishSuite) ErrorPublishRepositoryMock(repository *mocks.AlbumRepository, albumID uuid.UUID) {
	repository.
		On("Publish", context.Background(), albumID).
		Return(ports.ErrAlbumPublish)
}

func (s *AlbumPublishSuite) TestErrorPublish(t provider.T) {
	t.Parallel()
	t.Title("Album publish test error")
	albumID := uuid.New()
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	s.ErrorPublishRepositoryMock(repository, albumID)

	err := albumService.Publish(context.Background(), albumID)

	t.Assert().ErrorIs(err, ports.ErrAlbumPublish)
}

func TestAlbumPublishSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumPublishSuite))
}

type AlbumGetAllSuite struct {
	AlbumSuite
}

func (s *AlbumGetAllSuite) CorrectRepositoryMock(repository *mocks.AlbumRepository) {
	repository.
		On("GetAll", context.Background()).
		Return([]domain.Album{}, nil)
}

func (s *AlbumGetAllSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Album get all test correct")
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	s.CorrectRepositoryMock(repository)

	_, err := albumService.GetAll(context.Background())

	t.Assert().Nil(err)
}

func (s *AlbumGetAllSuite) InternalErrorRepositoryMock(repository *mocks.AlbumRepository) {
	repository.
		On("GetAll", context.Background()).
		Return(nil, ports.ErrInternalAlbumRepo)
}

func (s *AlbumGetAllSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Album get all test internal error")
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	s.InternalErrorRepositoryMock(repository)

	_, err := albumService.GetAll(context.Background())

	t.Assert().ErrorIs(err, ports.ErrInternalAlbumRepo)
}

func TestAlbumGetAllSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumGetAllSuite))
}

type AlbumGetByMusicianIDSuite struct {
	AlbumSuite
}

func (s *AlbumGetByMusicianIDSuite) CorrectRepositoryMock(repository *mocks.AlbumRepository, musicianID uuid.UUID) {
	repository.
		On("GetByMusicianID", context.Background(), musicianID).
		Return([]domain.Album{}, nil)
}

func (s *AlbumGetByMusicianIDSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Album get by musician id test correct")
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	musicianID := uuid.New()
	s.CorrectRepositoryMock(repository, musicianID)

	_, err := albumService.GetByMusicianID(context.Background(), musicianID)

	t.Assert().Nil(err)
}

func (s *AlbumGetByMusicianIDSuite) InternalErrorRepositoryMock(repository *mocks.AlbumRepository, musicianID uuid.UUID) {
	repository.
		On("GetByMusicianID", context.Background(), musicianID).
		Return(nil, ports.ErrInternalAlbumRepo)
}

func (s *AlbumGetByMusicianIDSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Album get by musician id test internal error")
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	musicianID := uuid.New()
	s.InternalErrorRepositoryMock(repository, musicianID)

	_, err := albumService.GetByMusicianID(context.Background(), musicianID)

	t.Assert().ErrorIs(err, ports.ErrInternalAlbumRepo)
}

func TestAlbumGetByMusicianIDSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumGetByMusicianIDSuite))
}

type AlbumGetOwnSuite struct {
	AlbumSuite
}

func (s *AlbumGetOwnSuite) CorrectRepositoryMock(repository *mocks.AlbumRepository, musicianID uuid.UUID) {
	repository.
		On("GetOwn", context.Background(), musicianID).
		Return([]domain.Album{}, nil)
}

func (s *AlbumGetOwnSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Album get own id test correct")
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	musicianID := uuid.New()
	s.CorrectRepositoryMock(repository, musicianID)

	_, err := albumService.GetOwn(context.Background(), musicianID)

	t.Assert().Nil(err)
}

func (s *AlbumGetOwnSuite) InternalErrorRepositoryMock(repository *mocks.AlbumRepository, musicianID uuid.UUID) {
	repository.
		On("GetOwn", context.Background(), musicianID).
		Return(nil, ports.ErrInternalAlbumRepo)
}

func (s *AlbumGetOwnSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Album get own id test internal error")
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	musicianID := uuid.New()
	s.InternalErrorRepositoryMock(repository, musicianID)

	_, err := albumService.GetOwn(context.Background(), musicianID)

	t.Assert().ErrorIs(err, ports.ErrInternalAlbumRepo)
}

func TestAlbumGetOwnSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumGetOwnSuite))
}

type AlbumGetByIDSuite struct {
	AlbumSuite
}

func (s *AlbumGetByIDSuite) CorrectRepositoryMock(repository *mocks.AlbumRepository, albumID uuid.UUID) {
	repository.
		On("GetByID", context.Background(), albumID).
		Return(domain.Album{}, nil)
}

func (s *AlbumGetByIDSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Album get by id test correct")
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	albumID := uuid.New()
	s.CorrectRepositoryMock(repository, albumID)

	_, err := albumService.GetByID(context.Background(), albumID)

	t.Assert().Nil(err)
}

func (s *AlbumGetByIDSuite) NotFoundRepositoryMock(repository *mocks.AlbumRepository, albumID uuid.UUID) {
	repository.
		On("GetByID", context.Background(), albumID).
		Return(domain.Album{}, ports.ErrAlbumIDNotFound)
}

func (s *AlbumGetByIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	t.Title("Album get by id test not found")
	repository := mocks.NewAlbumRepository(t)
	albumService := service.NewAlbumService(repository, s.logger)
	albumID := uuid.New()
	s.NotFoundRepositoryMock(repository, albumID)

	_, err := albumService.GetByID(context.Background(), albumID)

	t.Assert().ErrorIs(err, ports.ErrAlbumIDNotFound)
}

func TestAlbumGetByIDSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumGetByIDSuite))
}

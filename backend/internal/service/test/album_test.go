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
	logger       *zap.Logger
	repository   *mocks.AlbumRepository
	albumService *service.AlbumService
}

func (s *AlbumSuite) BeforeEach(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()

	s.repository = mocks.NewAlbumRepository(t)
	s.albumService = service.NewAlbumService(s.repository, s.logger)
}

type AlbumCreateSuite struct {
	AlbumSuite
}

func (s *AlbumCreateSuite) CorrectRepositoryMock(musicianID uuid.UUID) {
	s.repository.
		On("Create", context.Background(), mock.AnythingOfType("domain.Album"), musicianID).
		Return(domain.Album{}, nil)
}

func (s *AlbumCreateSuite) TestCorrect(t provider.T) {
	req := builder.NewCreateAlbumServiceRequestBuilder().Default().Build()
	s.CorrectRepositoryMock(req.MusicianID)

	_, err := s.albumService.Create(context.Background(), req)

	t.Assert().Nil(err)
}

func (s *AlbumCreateSuite) DuplicateRepositoryMock(musicianID uuid.UUID) {
	s.repository.
		On("Create", context.Background(), mock.AnythingOfType("domain.Album"), musicianID).
		Return(domain.Album{}, ports.ErrAlbumDuplicate)
}

func (s *AlbumCreateSuite) TestDuplicate(t provider.T) {
	req := builder.NewCreateAlbumServiceRequestBuilder().Default().Build()
	s.DuplicateRepositoryMock(req.MusicianID)

	_, err := s.albumService.Create(context.Background(), req)

	t.Assert().ErrorIs(err, ports.ErrAlbumDuplicate)
}

func TestAlbumCreateSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumCreateSuite))
}

type AlbumPublishSuite struct {
	AlbumSuite
}

func (s *AlbumPublishSuite) CorrectRepositoryMock(albumID uuid.UUID) {
	s.repository.
		On("Publish", context.Background(), albumID).
		Return(nil)
}

func (s *AlbumPublishSuite) TestCorrect(t provider.T) {
	albumID := uuid.New()
	s.CorrectRepositoryMock(albumID)

	err := s.albumService.Publish(context.Background(), albumID)

	t.Assert().Nil(err)
}

func (s *AlbumPublishSuite) ErrorPublishRepositoryMock(albumID uuid.UUID) {
	s.repository.
		On("Publish", context.Background(), albumID).
		Return(ports.ErrAlbumPublish)
}

func (s *AlbumPublishSuite) TestErrorPublish(t provider.T) {
	albumID := uuid.New()
	s.ErrorPublishRepositoryMock(albumID)

	err := s.albumService.Publish(context.Background(), albumID)

	t.Assert().ErrorIs(err, ports.ErrAlbumPublish)
}

func TestAlbumPublishSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumPublishSuite))
}

type AlbumGetAllSuite struct {
	AlbumSuite
}

func (s *AlbumGetAllSuite) CorrectRepositoryMock() {
	s.repository.
		On("GetAll", context.Background()).
		Return([]domain.Album{}, nil)
}

func (s *AlbumGetAllSuite) TestCorrect(t provider.T) {
	s.CorrectRepositoryMock()

	_, err := s.albumService.GetAll(context.Background())

	t.Assert().Nil(err)
}

func (s *AlbumGetAllSuite) InternalErrorRepositoryMock() {
	s.repository.
		On("GetAll", context.Background()).
		Return(nil, ports.ErrInternalAlbumRepo)
}

func (s *AlbumGetAllSuite) TestInternalError(t provider.T) {
	s.InternalErrorRepositoryMock()

	_, err := s.albumService.GetAll(context.Background())

	t.Assert().ErrorIs(err, ports.ErrInternalAlbumRepo)
}

func TestAlbumGetAllSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumGetAllSuite))
}

type AlbumGetByMusicianIDSuite struct {
	AlbumSuite
}

func (s *AlbumGetByMusicianIDSuite) CorrectRepositoryMock(musicianID uuid.UUID) {
	s.repository.
		On("GetByMusicianID", context.Background(), musicianID).
		Return([]domain.Album{}, nil)
}

func (s *AlbumGetByMusicianIDSuite) TestCorrect(t provider.T) {
	musicianID := uuid.New()
	s.CorrectRepositoryMock(musicianID)

	_, err := s.albumService.GetByMusicianID(context.Background(), musicianID)

	t.Assert().Nil(err)
}

func (s *AlbumGetByMusicianIDSuite) InternalErrorRepositoryMock(musicianID uuid.UUID) {
	s.repository.
		On("GetByMusicianID", context.Background(), musicianID).
		Return(nil, ports.ErrInternalAlbumRepo)
}

func (s *AlbumGetByMusicianIDSuite) TestInternalError(t provider.T) {
	musicianID := uuid.New()
	s.InternalErrorRepositoryMock(musicianID)

	_, err := s.albumService.GetByMusicianID(context.Background(), musicianID)

	t.Assert().ErrorIs(err, ports.ErrInternalAlbumRepo)
}

func TestAlbumGetByMusicianIDSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumGetByMusicianIDSuite))
}

type AlbumGetOwnSuite struct {
	AlbumSuite
}

func (s *AlbumGetOwnSuite) CorrectRepositoryMock(musicianID uuid.UUID) {
	s.repository.
		On("GetOwn", context.Background(), musicianID).
		Return([]domain.Album{}, nil)
}

func (s *AlbumGetOwnSuite) TestCorrect(t provider.T) {
	musicianID := uuid.New()
	s.CorrectRepositoryMock(musicianID)

	_, err := s.albumService.GetOwn(context.Background(), musicianID)

	t.Assert().Nil(err)
}

func (s *AlbumGetOwnSuite) InternalErrorRepositoryMock(musicianID uuid.UUID) {
	s.repository.
		On("GetOwn", context.Background(), musicianID).
		Return(nil, ports.ErrInternalAlbumRepo)
}

func (s *AlbumGetOwnSuite) TestInternalError(t provider.T) {
	musicianID := uuid.New()
	s.InternalErrorRepositoryMock(musicianID)

	_, err := s.albumService.GetOwn(context.Background(), musicianID)

	t.Assert().ErrorIs(err, ports.ErrInternalAlbumRepo)
}

func TestAlbumGetOwnSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumGetOwnSuite))
}

type AlbumGetByIDSuite struct {
	AlbumSuite
}

func (s *AlbumGetByIDSuite) CorrectRepositoryMock(albumID uuid.UUID) {
	s.repository.
		On("GetByID", context.Background(), albumID).
		Return(domain.Album{}, nil)
}

func (s *AlbumGetByIDSuite) TestCorrect(t provider.T) {
	albumID := uuid.New()
	s.CorrectRepositoryMock(albumID)

	_, err := s.albumService.GetByID(context.Background(), albumID)

	t.Assert().Nil(err)
}

func (s *AlbumGetByIDSuite) NotFoundRepositoryMock(albumID uuid.UUID) {
	s.repository.
		On("GetByID", context.Background(), albumID).
		Return(domain.Album{}, ports.ErrAlbumIDNotFound)
}

func (s *AlbumGetByIDSuite) TestNotFound(t provider.T) {
	albumID := uuid.New()
	s.NotFoundRepositoryMock(albumID)

	_, err := s.albumService.GetByID(context.Background(), albumID)

	t.Assert().ErrorIs(err, ports.ErrAlbumIDNotFound)
}

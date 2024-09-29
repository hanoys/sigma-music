package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
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

type MusicianSuite struct {
	suite.Suite
	logger       *zap.Logger
	hashProvider *hash.HashPasswordProvider
}

func (s *MusicianSuite) BeforeEach(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()

	s.hashProvider = hash.NewHashPasswordProvider()
}

type MusicianRegisterSuite struct {
	MusicianSuite
}

func (s *MusicianRegisterSuite) CorrectRepositoryMock(repository *mocks.MusicianRepository, req ports.MusicianServiceCreateRequest) {
	repository.
		On("Create", context.Background(), mock.AnythingOfType("domain.Musician")).
		Return(domain.Musician{}, nil).
		On("GetByName", context.Background(), req.Name).
		Return(domain.Musician{}, ports.ErrMusicianNameNotFound).
		On("GetByEmail", context.Background(), req.Email).
		Return(domain.Musician{}, ports.ErrMusicianEmailNotFound)
}

func (s *MusicianRegisterSuite) TestCorrect(t provider.T) {
	req := builder.NewMusicianServiceCreateRequestBuilder().Default().Build()
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(repository, req)

	_, err := musicianService.Register(context.Background(), req)

	t.Assert().Nil(err)
}

func (s *MusicianRegisterSuite) NameExistsRepositoryMock(repository *mocks.MusicianRepository, req ports.MusicianServiceCreateRequest) {
	repository.
		On("GetByName", context.Background(), req.Name).
		Return(domain.Musician{}, nil)
}

func (s *MusicianRegisterSuite) TestNameExists(t provider.T) {
	req := builder.NewMusicianServiceCreateRequestBuilder().Default().Build()
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.NameExistsRepositoryMock(repository, req)

	_, err := musicianService.Register(context.Background(), req)

	t.Assert().ErrorIs(err, ports.ErrMusicianWithSuchNameAlreadyExists)
}

func (s *MusicianRegisterSuite) EmailExistsRepositoryMock(repository *mocks.MusicianRepository, req ports.MusicianServiceCreateRequest) {
	repository.
		On("GetByName", context.Background(), req.Name).
		Return(domain.Musician{}, ports.ErrMusicianNameNotFound).
		On("GetByEmail", context.Background(), req.Email).
		Return(domain.Musician{}, nil)
}

func (s *MusicianRegisterSuite) TestEmailExists(t provider.T) {
	req := builder.NewMusicianServiceCreateRequestBuilder().Default().Build()
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.EmailExistsRepositoryMock(repository, req)

	_, err := musicianService.Register(context.Background(), req)

	t.Assert().ErrorIs(err, ports.ErrMusicianWithSuchEmailAlreadyExists)
}

func TestMusicianRegisterSuite(t *testing.T) {
	suite.RunSuite(t, new(MusicianRegisterSuite))
}

type MusicianGetAllSuite struct {
	MusicianSuite
}

func (s *MusicianGetAllSuite) CorrectRepositoryMock(repository *mocks.MusicianRepository) {
	musician1 := builder.NewMusicianBuilder().Default().Build()
	musician2 := builder.NewMusicianBuilder().Default().Build()

	repository.
		On("GetAll", context.Background()).
		Return([]domain.Musician{musician1, musician2}, nil)
}

func (s *MusicianGetAllSuite) TestCorrect(t provider.T) {
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(repository)

	musicians, err := musicianService.GetAll(context.Background())

	t.Assert().Nil(err)
	t.Assert().Len(musicians, 2)
}

func (s *MusicianGetAllSuite) RepositoryErrorRepositoryMock(repository *mocks.MusicianRepository) {
	repository.
		On("GetAll", context.Background()).
		Return(nil, ports.ErrInternalMusicianRepo)
}

func (s *MusicianGetAllSuite) TestRepositoryError(t provider.T) {
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.RepositoryErrorRepositoryMock(repository)

	_, err := musicianService.GetAll(context.Background())

	t.Assert().ErrorIs(err, ports.ErrInternalMusicianRepo)
}

func TestMusicianGetAllSuite(t *testing.T) {
	suite.RunSuite(t, new(MusicianGetAllSuite))
}

type MusicianGetByIDSuite struct {
	MusicianSuite
}

func (s *MusicianGetByIDSuite) CorrectRepositoryMock(repository *mocks.MusicianRepository, musicianID uuid.UUID) {
	musician := builder.NewMusicianBuilder().Default().SetID(musicianID).Build()

	repository.
		On("GetByID", context.Background(), musicianID).
		Return(musician, nil)
}

func (s *MusicianGetByIDSuite) TestCorrect(t provider.T) {
	musicianID := uuid.New()
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(repository, musicianID)

	result, err := musicianService.GetByID(context.Background(), musicianID)

	t.Assert().Nil(err)
	t.Assert().Equal(musicianID, result.ID)
}

func (s *MusicianGetByIDSuite) NotFoundRepositoryMock(repository *mocks.MusicianRepository, musicianID uuid.UUID) {
	repository.
		On("GetByID", context.Background(), musicianID).
		Return(domain.Musician{}, ports.ErrMusicianIDNotFound)
}

func (s *MusicianGetByIDSuite) TestIDNotFound(t provider.T) {
	musicianID := uuid.New()
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.NotFoundRepositoryMock(repository, musicianID)

	_, err := musicianService.GetByID(context.Background(), musicianID)

	t.Assert().ErrorIs(err, ports.ErrMusicianIDNotFound)
}

func TestMusicianGetByIDSuite(t *testing.T) {
	suite.RunSuite(t, new(MusicianGetByIDSuite))
}

type MusicianGetByNameSuite struct {
	MusicianSuite
}

func (s *MusicianGetByNameSuite) CorrectRepositoryMock(repository *mocks.MusicianRepository, name string) {
	musician := builder.NewMusicianBuilder().Default().SetName(name).Build()

	repository.
		On("GetByName", context.Background(), name).
		Return(musician, nil)
}

func (s *MusicianGetByNameSuite) TestCorrect(t provider.T) {
	name := "Test Musician"
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(repository, name)

	result, err := musicianService.GetByName(context.Background(), name)

	t.Assert().Nil(err)
	t.Assert().Equal(name, result.Name)
}

func (s *MusicianGetByNameSuite) NotFoundRepositoryMock(repository *mocks.MusicianRepository, name string) {
	repository.
		On("GetByName", context.Background(), name).
		Return(domain.Musician{}, ports.ErrMusicianNameNotFound)
}

func (s *MusicianGetByNameSuite) TestNameNotFound(t provider.T) {
	name := "Test Musician"
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.NotFoundRepositoryMock(repository, name)

	_, err := musicianService.GetByName(context.Background(), name)

	t.Assert().ErrorIs(err, ports.ErrMusicianNameNotFound)
}

func TestMusicianGetByNameSuite(t *testing.T) {
	suite.RunSuite(t, new(MusicianGetByNameSuite))
}

type MusicianGetByEmailSuite struct {
	MusicianSuite
}

func (s *MusicianGetByEmailSuite) CorrectRepositoryMock(repository *mocks.MusicianRepository, email string) {
	musician := builder.NewMusicianBuilder().Default().SetEmail(email).Build()

	repository.
		On("GetByEmail", context.Background(), email).
		Return(musician, nil)
}

func (s *MusicianGetByEmailSuite) TestCorrect(t provider.T) {
	email := "test.musician@mail.com"
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(repository, email)

	result, err := musicianService.GetByEmail(context.Background(), email)

	t.Assert().Nil(err)
	t.Assert().Equal(email, result.Email)
}

func (s *MusicianGetByEmailSuite) NotFoundRepositoryMock(repository *mocks.MusicianRepository, email string) {
	repository.
		On("GetByEmail", context.Background(), email).
		Return(domain.Musician{}, ports.ErrMusicianEmailNotFound)
}

func (s *MusicianGetByEmailSuite) TestNotFound(t provider.T) {
	email := "test.musician@mail.com"
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.NotFoundRepositoryMock(repository, email)

	_, err := musicianService.GetByEmail(context.Background(), email)

	t.Assert().ErrorIs(err, ports.ErrMusicianEmailNotFound)
}

func TestMusicianGetByEmailSuite(t *testing.T) {
	suite.RunSuite(t, new(MusicianGetByEmailSuite))
}

type MusicianGetByAlbumIDSuite struct {
	MusicianSuite
}

func (s *MusicianGetByAlbumIDSuite) CorrectRepositoryMock(repository *mocks.MusicianRepository, albumID uuid.UUID) {
	musician := builder.NewMusicianBuilder().Default().Build()

	repository.
		On("GetByAlbumID", context.Background(), albumID).
		Return(musician, nil)
}

func (s *MusicianGetByAlbumIDSuite) TestCorrect(t provider.T) {
	albumID := uuid.New()
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(repository, albumID)

	_, err := musicianService.GetByAlbumID(context.Background(), albumID)

	t.Assert().Nil(err)
}

func (s *MusicianGetByAlbumIDSuite) NotFoundRepositoryMock(repository *mocks.MusicianRepository, albumID uuid.UUID) {
	musician := builder.NewMusicianBuilder().Default().Build()

	repository.
		On("GetByAlbumID", context.Background(), albumID).
		Return(musician, ports.ErrMusicianIDNotFound)
}

func (s *MusicianGetByAlbumIDSuite) TestNotFound(t provider.T) {
	albumID := uuid.New()
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.NotFoundRepositoryMock(repository, albumID)

	_, err := musicianService.GetByAlbumID(context.Background(), albumID)

	t.Assert().ErrorIs(err, ports.ErrMusicianIDNotFound)
}

func TestMusicianGetByAlbumIDSuite(t *testing.T) {
	suite.RunSuite(t, new(MusicianGetByAlbumIDSuite))
}

type MusicianGetByTrackIDSuite struct {
	MusicianSuite
}

func (s *MusicianGetByTrackIDSuite) CorrectRepositoryMock(repository *mocks.MusicianRepository, trackID uuid.UUID) {
	musician := builder.NewMusicianBuilder().Default().Build()

	repository.
		On("GetByTrackID", context.Background(), trackID).
		Return(musician, nil)
}

func (s *MusicianGetByTrackIDSuite) TestCorrect(t provider.T) {
	trackID := uuid.New()
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(repository, trackID)

	_, err := musicianService.GetByTrackID(context.Background(), trackID)

	t.Assert().Nil(err)
}

func (s *MusicianGetByTrackIDSuite) NotFoundRepositoryMock(repository *mocks.MusicianRepository, trackID uuid.UUID) {
	repository.
		On("GetByTrackID", context.Background(), trackID).
		Return(domain.Musician{}, ports.ErrMusicianIDNotFound)
}

func (s *MusicianGetByTrackIDSuite) TestNotFound(t provider.T) {
	trackID := uuid.New()
	repository := mocks.NewMusicianRepository(t)
	musicianService := service.NewMusicianService(repository, s.hashProvider, s.logger)
	s.NotFoundRepositoryMock(repository, trackID)

	_, err := musicianService.GetByTrackID(context.Background(), trackID)

	t.Assert().ErrorIs(err, ports.ErrMusicianIDNotFound)
}

func TestMusicianGetByTrackIDSuite(t *testing.T) {
	suite.RunSuite(t, new(MusicianGetByTrackIDSuite))
}

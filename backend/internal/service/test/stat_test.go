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

type StatSuite struct {
	suite.Suite
	logger       *zap.Logger
	hashProvider *hash.HashPasswordProvider
}

func (s *StatSuite) BeforeEach(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()
	s.hashProvider = hash.NewHashPasswordProvider()
}

type StatAddSuite struct {
	StatSuite
}

func (s *StatAddSuite) CorrectRepositoryMock(statRepository *mocks.StatRepository,
	musicianRepository *mocks.MusicianRepository, genreRepository *mocks.GenreRepository, userID uuid.UUID, trackID uuid.UUID) {
	statRepository.
		On("Add", context.Background(), mock.Anything, userID, trackID).
		Return(nil)
}

func (s *StatAddSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Stat add test correct")
	userID := uuid.New()
	trackID := uuid.New()
	statRepository := mocks.NewStatRepository(t)
	musicianRepository := mocks.NewMusicianRepository(t)
	genreRepository := mocks.NewGenreRepository(t)
	genreService := service.NewGenreService(genreRepository, s.logger)
	musicianService := service.NewMusicianService(musicianRepository, s.hashProvider, s.logger)
	statService := service.NewStatService(statRepository, genreService, musicianService, s.logger)
	s.CorrectRepositoryMock(statRepository, musicianRepository, genreRepository, userID, trackID)

	err := statService.Add(context.Background(), userID, trackID)

	t.Assert().Nil(err)
}

func (s *StatAddSuite) InternalErrorRepositoryMock(statRepository *mocks.StatRepository,
	musicianRepository *mocks.MusicianRepository, genreRepository *mocks.GenreRepository, userID uuid.UUID, trackID uuid.UUID) {
	statRepository.
		On("Add", context.Background(), mock.Anything, userID, trackID).
		Return(ports.ErrInternalStatRepo)
}

func (s *StatAddSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Stat add test internal error")
	userID := uuid.New()
	trackID := uuid.New()
	statRepository := mocks.NewStatRepository(t)
	musicianRepository := mocks.NewMusicianRepository(t)
	genreRepository := mocks.NewGenreRepository(t)
	genreService := service.NewGenreService(genreRepository, s.logger)
	musicianService := service.NewMusicianService(musicianRepository, s.hashProvider, s.logger)
	statService := service.NewStatService(statRepository, genreService, musicianService, s.logger)
	s.InternalErrorRepositoryMock(statRepository, musicianRepository, genreRepository, userID, trackID)

	err := statService.Add(context.Background(), userID, trackID)

	t.Assert().ErrorIs(err, ports.ErrInternalStatRepo)
}

func TestStatAddSuite(t *testing.T) {
	suite.RunSuite(t, new(StatAddSuite))
}

type StatFormReportSuite struct {
	StatSuite
}

func (s *StatFormReportSuite) CorrectRepositoryMock(statRepository *mocks.StatRepository,
	musicianRepository *mocks.MusicianRepository, genreRepository *mocks.GenreRepository, userID uuid.UUID, listenedMusicians []domain.UserMusiciansStat, listenedGenres []domain.UserGenresStat, musicians []domain.Musician, genres []domain.Genre) {
	statRepository.
		On("GetMostListenedMusicians", context.Background(), userID, 3).
		Return(listenedMusicians, nil).
		On("GetListenedGenres", context.Background(), userID).
		Return(listenedGenres, nil)

	musicianRepository.
		On("GetByID", context.Background(), listenedMusicians[0].MusicianID).
		Return(musicians[0], nil).
		On("GetByID", context.Background(), listenedMusicians[1].MusicianID).
		Return(musicians[1], nil).
		On("GetByID", context.Background(), listenedMusicians[2].MusicianID).
		Return(musicians[2], nil)

	genreRepository.
		On("GetByID", context.Background(), listenedGenres[0].GenreID).
		Return(genres[0], nil).
		On("GetByID", context.Background(), listenedGenres[1].GenreID).
		Return(genres[1], nil)
}

func (s *StatFormReportSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Stat form report test correct")
	userID := uuid.New()
	musicians := make([]domain.Musician, 0)
	for i := 0; i < 3; i++ {
		musicians = append(musicians, builder.NewMusicianBuilder().Default().Build())
	}
	musiciansStat := make([]domain.UserMusiciansStat, 0)
	for i := range musicians {
		musiciansStat = append(musiciansStat, builder.NewUserMusiciansStatBuilder().
			Default().
			SetMusicianID(musicians[i].ID).
			SetListenCount(1).
			Build())
	}
	genres := make([]domain.Genre, 0)
	for i := 0; i < 2; i++ {
		genres = append(genres, builder.NewGenreBuilder().Default().Build())
	}
	genresStat := make([]domain.UserGenresStat, 0)
	for i := range genres {
		genresStat = append(genresStat, builder.NewUserGenresStatBuilder().
			SetGenreID(genres[i].ID).
			SetListenCount(1).
			Build())
	}
	statRepository := mocks.NewStatRepository(t)
	musicianRepository := mocks.NewMusicianRepository(t)
	genreRepository := mocks.NewGenreRepository(t)
	genreService := service.NewGenreService(genreRepository, s.logger)
	musicianService := service.NewMusicianService(musicianRepository, s.hashProvider, s.logger)
	statService := service.NewStatService(statRepository, genreService, musicianService, s.logger)
	s.CorrectRepositoryMock(statRepository, musicianRepository, genreRepository, userID, musiciansStat, genresStat, musicians, genres)

	r, err := statService.FormReport(context.Background(), userID)

	t.Assert().Equal(r.ListenCount, int64(3))
	t.Assert().Nil(err)
}

func TestStatFormReportSuite(t *testing.T) {
	suite.RunSuite(t, new(StatFormReportSuite))
}

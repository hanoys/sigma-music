package integrationtest

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	testpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap"
)

type StatSuite struct {
	suite.Suite
	logger    *zap.Logger
	hash      *hash.HashPasswordProvider
	container *testpg.PostgresContainer
	db        *sqlx.DB
}

func (s *StatSuite) BeforeAll(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()

	ctx := context.Background()
	var err error
	s.container, err = newPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s.hash = hash.NewHashPasswordProvider()
}

func (s *StatSuite) BeforeEach(t provider.T) {
	url, err := s.container.ConnectionString(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	s.db, err = newPostgresDB(url)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *StatSuite) AfterAll(t provider.T) {
	if err := s.container.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func (s *StatSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *StatSuite) TestAdd(t provider.T) {
    t.Title("stat add integration test")
    if (isPreviousTestsFailed()) {
        t.Skip()
    }
	repo := postgres.NewPostgresStatRepository(s.db)
	musrepo := postgres.NewPostgresMusicianRepository(s.db)
	musicianService := service.NewMusicianService(musrepo, s.hash, s.logger)
	genrerepo := postgres.NewPostgresGenreRepository(s.db)
	genreService := service.NewGenreService(genrerepo, s.logger)
	statService := service.NewStatService(repo, genreService, musicianService, s.logger)
	userID, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fc")
	trackID, _ := uuid.Parse("41623ac1-b98d-4478-a10f-870a80c697b6")

	err := statService.Add(context.Background(), userID, trackID)

	t.Assert().Nil(err)
}

func (s *StatSuite) TestFormReport(t provider.T) {
    t.Title("stat form report integration test")
    if (isPreviousTestsFailed()) {
        t.Skip()
    }
	repo := postgres.NewPostgresStatRepository(s.db)
	musrepo := postgres.NewPostgresMusicianRepository(s.db)
	musicianService := service.NewMusicianService(musrepo, s.hash, s.logger)
	genrerepo := postgres.NewPostgresGenreRepository(s.db)
	genreService := service.NewGenreService(genrerepo, s.logger)
	statService := service.NewStatService(repo, genreService, musicianService, s.logger)
	userID, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fc")

	_, err := statService.FormReport(context.Background(), userID)

	t.Assert().Nil(err)
}

func TestStatSuite(t *testing.T) {
	suite.RunSuite(t, new(StatSuite))
}

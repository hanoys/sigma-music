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

type GenreSuite struct {
	suite.Suite
	logger    *zap.Logger
	hash      *hash.HashPasswordProvider
	container *testpg.PostgresContainer
	db        *sqlx.DB
}

func (s *GenreSuite) BeforeAll(t provider.T) {
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

func (s *GenreSuite) BeforeEach(t provider.T) {
	url, err := s.container.ConnectionString(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	s.db, err = newPostgresDB(url)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *GenreSuite) AfterAll(t provider.T) {
	if err := s.container.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func (s *GenreSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *GenreSuite) TestGetAll(t provider.T) {
    t.Title("genre get all integration test")
	repo := postgres.NewPostgresGenreRepository(s.db)
	genreService := service.NewGenreService(repo, s.logger)

	genres, err := genreService.GetAll(context.Background())

	t.Assert().Nil(err)
	t.Assert().NotNil(genres)
}

func (s *GenreSuite) TestGetByID(t provider.T) {
    t.Title("genre get by id integration test")
	repo := postgres.NewPostgresGenreRepository(s.db)
	genreService := service.NewGenreService(repo, s.logger)
	genreID, _ := uuid.Parse("32f24dfc-3823-41e4-a073-c3553c981db1")

	genre, err := genreService.GetByID(context.Background(), genreID)

	t.Assert().Nil(err)
	t.Assert().Equal(genreID, genre.ID)
}

func TestGenreSuite(t *testing.T) {
    suite.RunSuite(t, new(GenreSuite))
}

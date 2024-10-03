package integrationtest

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/hanoys/sigma-music/internal/service/test/builder"
	"github.com/jmoiron/sqlx"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	testpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap"
)

type AlbumSuite struct {
	suite.Suite
	logger    *zap.Logger
	hash      *hash.HashPasswordProvider
	container *testpg.PostgresContainer
	db        *sqlx.DB
}

func (s *AlbumSuite) BeforeAll(t provider.T) {
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

func (s *AlbumSuite) BeforeEach(t provider.T) {
	url, err := s.container.ConnectionString(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	s.db, err = newPostgresDB(url)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *AlbumSuite) AfterAll(t provider.T) {
	if err := s.container.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func (s *AlbumSuite) AfterEach(t provider.T) {
	err := s.container.Restore(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	s.db.Close()
}

func (s *AlbumSuite) TestCreate(t provider.T) {
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)
	musicianID, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fa")
	req := builder.NewCreateAlbumServiceRequestBuilder().
		Default().
		SetMusicianID(musicianID).
		Build()

	album, err := albumService.Create(context.Background(), req)

	t.Assert().Nil(err)
	t.Assert().Equal(album.Name, req.Name)
}

func (s *AlbumSuite) TestGetAll(t provider.T) {
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)

	albums, err := albumService.GetAll(context.Background())

	t.Assert().Nil(err)
	t.Assert().NotNil(albums)
}

func (s *AlbumSuite) TestGetByID(t provider.T) {
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)
	id, _ := uuid.Parse("b24fa8eb-9df6-406c-9b45-763d7b5a5078")

	_, err := albumService.GetByID(context.Background(), id)

	t.Assert().Nil(err)
}

func (s *AlbumSuite) TestGetOwn(t provider.T) {
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)
	id, _ := uuid.Parse("b24fa8eb-9df6-406c-9b45-763d7b5a5078")

	_, err := albumService.GetOwn(context.Background(), id)

	t.Assert().Nil(err)
}

func (s *AlbumSuite) TestGetByMusicianID(t provider.T) {
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)
	id, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fa")

	albums, err := albumService.GetByMusicianID(context.Background(), id)

	t.Assert().Nil(err)
    t.Assert().NotNil(albums)
}

func (s *AlbumSuite) TestPublish(t provider.T) {
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)
	id, _ := uuid.Parse("b24fa8eb-9df6-406c-9b45-763d7b5a5078")

	err := albumService.Publish(context.Background(), id)

	t.Assert().Nil(err)
}

func TestAlbumSuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumSuite))
}

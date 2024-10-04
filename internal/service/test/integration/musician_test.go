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

type MusicianSuite struct {
	suite.Suite
	logger    *zap.Logger
	hash      *hash.HashPasswordProvider
	container *testpg.PostgresContainer
	db        *sqlx.DB
}

func (s *MusicianSuite) BeforeAll(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()

	s.hash = hash.NewHashPasswordProvider()
}

func (s *MusicianSuite) BeforeEach(t provider.T) {
	ctx := context.Background()
	var err error
	s.container, err = newPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	url, err := s.container.ConnectionString(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	s.db, err = newPostgresDB(url)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *MusicianSuite) AfterAll(t provider.T) {
	if err := s.container.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func (s *MusicianSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *MusicianSuite) TestRegister(t provider.T) {
	t.Title("musician register integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresMusicianRepository(s.db)
	musicianService := service.NewMusicianService(repo, s.hash, s.logger)

	req := builder.NewMusicianServiceCreateRequestBuilder().
		Default().
		SetName("Test").
		SetEmail("test").
		Build()
	createdMusician, err := musicianService.Register(context.Background(), req)

	t.Assert().Nil(err)
	t.Assert().Equal(req.Name, createdMusician.Name)
}

func (s *MusicianSuite) TestGetByIDSuccess(t provider.T) {
	t.Title("musician get by id integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresMusicianRepository(s.db)
	musicianService := service.NewMusicianService(repo, s.hash, s.logger)

	id, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fa")
	foundMusician, err := musicianService.GetByID(context.Background(), id)

	t.Assert().Nil(err)
	t.Assert().Equal(id, foundMusician.ID)
}

func (s *MusicianSuite) TestGetByNameSuccess(t provider.T) {
	t.Title("musician get by name integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresMusicianRepository(s.db)
	musicianService := service.NewMusicianService(repo, s.hash, s.logger)

	name := "Timur"
	foundMusician, err := musicianService.GetByName(context.Background(), name)

	t.Assert().Nil(err)
	t.Assert().Equal(name, foundMusician.Name)
}

func (s *MusicianSuite) TestGetByEmailSuccess(t provider.T) {
	t.Title("musician get by email integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresMusicianRepository(s.db)
	musicianService := service.NewMusicianService(repo, s.hash, s.logger)

	email := "timur@mail.ru"
	foundMusician, err := musicianService.GetByEmail(context.Background(), email)

	t.Assert().Nil(err)
	t.Assert().Equal(email, foundMusician.Email)
}

func TestMusicianSuite(t *testing.T) {
	suite.RunSuite(t, new(MusicianSuite))
}

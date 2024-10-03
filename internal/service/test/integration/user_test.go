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

type UserSuite struct {
	suite.Suite
	logger    *zap.Logger
	hash      *hash.HashPasswordProvider
	container *testpg.PostgresContainer
	db        *sqlx.DB
}

func (s *UserSuite) BeforeAll(t provider.T) {
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

func (s *UserSuite) BeforeEach(t provider.T) {
	url, err := s.container.ConnectionString(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	s.db, err = newPostgresDB(url)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *UserSuite) AfterAll(t provider.T) {
	if err := s.container.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func (s *UserSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *UserSuite) TestCreateSuccess(t provider.T) {
    t.Title("user create integration test")
	repo := postgres.NewPostgresUserRepository(s.db)
	userService := service.NewUserService(repo, s.hash, s.logger)
	req := builder.NewUserServiceCreateRequestBuilder().
		Default().
		SetName("Test").
		SetEmail("test").
		SetPhone("+7").Build()
	createdUser, err := userService.Register(context.Background(), req)

	if err != nil {
		t.Errorf("unexcpected error: %v", err)
	}

	t.Assert().Equal(req.Name, createdUser.Name)
}

func (s *UserSuite) TestGetByIDSuccess(t provider.T) {
    t.Title("user get by id integration test")
	repo := postgres.NewPostgresUserRepository(s.db)
	userService := service.NewUserService(repo, s.hash, s.logger)

	id, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fc")
	foundUser, err := userService.GetById(context.Background(), id)

	t.Assert().Nil(err)
	t.Assert().Equal(id, foundUser.ID)
}

func (s *UserSuite) TestGetByNameSuccess(t provider.T) {
    t.Title("user get by name integration test")
	repo := postgres.NewPostgresUserRepository(s.db)
	userService := service.NewUserService(repo, s.hash, s.logger)

	name := "Timur"
	foundUser, err := userService.GetByName(context.Background(), name)

	t.Assert().Nil(err)
	t.Assert().Equal(name, foundUser.Name)
}

func (s *UserSuite) TestGetByEmailSuccess(t provider.T) {
    t.Title("user get by email integration test")
	repo := postgres.NewPostgresUserRepository(s.db)
	userService := service.NewUserService(repo, s.hash, s.logger)

	email := "timur@mail.ru"
	foundUser, err := userService.GetByEmail(context.Background(), email)

	t.Assert().Nil(err)
	t.Assert().Equal(email, foundUser.Email)
}

func (s *UserSuite) TestGetByPhoneSuccess(t provider.T) {
    t.Title("user get by phone integration test")
	repo := postgres.NewPostgresUserRepository(s.db)
	userService := service.NewUserService(repo, s.hash, s.logger)

	phone := "+79999999999"
	foundUser, err := userService.GetByPhone(context.Background(), phone)

	t.Assert().Nil(err)
	t.Assert().Equal(phone, foundUser.Phone)
}

func TestUserSuite(t *testing.T) {
	suite.RunSuite(t, new(UserSuite))
}

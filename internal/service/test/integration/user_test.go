package integrationtest

import (
	"context"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/hanoys/sigma-music/internal/service/test/builder"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *AllSuite) TestCreateSuccess(t provider.T) {
	t.Parallel()
	t.Title("user create integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
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

func (s *AllSuite) TestUserGetByIDSuccess(t provider.T) {
	t.Parallel()
	t.Title("user get by id integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresUserRepository(s.db)
	userService := service.NewUserService(repo, s.hash, s.logger)

	id, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fc")
	foundUser, err := userService.GetById(context.Background(), id)

	t.Assert().Nil(err)
	t.Assert().Equal(id, foundUser.ID)
}

func (s *AllSuite) TestUserGetByNameSuccess(t provider.T) {
	t.Title("user get by name integration test")
	t.Parallel()
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresUserRepository(s.db)
	userService := service.NewUserService(repo, s.hash, s.logger)

	name := "Timur"
	foundUser, err := userService.GetByName(context.Background(), name)

	t.Assert().Nil(err)
	t.Assert().Equal(name, foundUser.Name)
}

func (s *AllSuite) TestUserGetByEmailSuccess(t provider.T) {
	t.Parallel()
	t.Title("user get by email integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresUserRepository(s.db)
	userService := service.NewUserService(repo, s.hash, s.logger)

	email := "timur@mail.ru"
	foundUser, err := userService.GetByEmail(context.Background(), email)

	t.Assert().Nil(err)
	t.Assert().Equal(email, foundUser.Email)
}

func (s *AllSuite) TestGetByPhoneSuccess(t provider.T) {
	t.Parallel()
	t.Title("user get by phone integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresUserRepository(s.db)
	userService := service.NewUserService(repo, s.hash, s.logger)

	phone := "+79999999999"
	foundUser, err := userService.GetByPhone(context.Background(), phone)

	t.Assert().Nil(err)
	t.Assert().Equal(phone, foundUser.Phone)
}

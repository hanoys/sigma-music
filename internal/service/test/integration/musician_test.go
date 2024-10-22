package integrationtest

import (
	"context"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/hanoys/sigma-music/internal/service/test/builder"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *AllSuite) TestRegister(t provider.T) {
	t.Parallel()
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

func (s *AllSuite) TestGenreGetByIDSuccess(t provider.T) {
	t.Parallel()
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

func (s *AllSuite) TestGenreGetByNameSuccess(t provider.T) {
	t.Parallel()
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

func (s *AllSuite) TestGenreGetByEmailSuccess(t provider.T) {
	t.Parallel()
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

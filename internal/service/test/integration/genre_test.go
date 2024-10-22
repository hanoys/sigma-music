package integrationtest

import (
	"context"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *AllSuite) TestGenreGetAll(t provider.T) {
	t.Parallel()
	t.Title("genre get all integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresGenreRepository(s.db)
	genreService := service.NewGenreService(repo, s.logger)

	genres, err := genreService.GetAll(context.Background())

	t.Assert().Nil(err)
	t.Assert().NotNil(genres)
}

func (s *AllSuite) TestGenreGetByID(t provider.T) {
	t.Parallel()
	t.Title("genre get by id integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresGenreRepository(s.db)
	genreService := service.NewGenreService(repo, s.logger)
	genreID, _ := uuid.Parse("32f24dfc-3823-41e4-a073-c3553c981db1")

	genre, err := genreService.GetByID(context.Background(), genreID)

	t.Assert().Nil(err)
	t.Assert().Equal(genreID, genre.ID)
}

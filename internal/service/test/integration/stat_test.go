package integrationtest

import (
	"context"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *AllSuite) TestAdd(t provider.T) {
    t.Parallel()
	t.Title("stat add integration test")
	if isPreviousTestsFailed() {
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

func (s *AllSuite) TestFormReport(t provider.T) {
    t.Parallel()
	t.Title("stat form report integration test")
	if isPreviousTestsFailed() {
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

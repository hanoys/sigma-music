package integrationtest

import (
	"context"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/hanoys/sigma-music/internal/service/test/builder"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *AllSuite) TestAlbumCreate(t provider.T) {
	t.Parallel()
	t.Title("Album create integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
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

func (s *AllSuite) TestAlbumGetAll(t provider.T) {
	t.Parallel()
	t.Title("Album get all integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)

	albums, err := albumService.GetAll(context.Background())

	t.Assert().Nil(err)
	t.Assert().NotNil(albums)
}

func (s *AllSuite) TestAlbumGetByID(t provider.T) {
	t.Parallel()
	t.Title("Album get by id integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)
	id, _ := uuid.Parse("b24fa8eb-9df6-406c-9b45-763d7b5a5078")

	_, err := albumService.GetByID(context.Background(), id)

	t.Assert().Nil(err)
}

func (s *AllSuite) TestAlbumGetOwn(t provider.T) {
	t.Parallel()
	t.Title("Album get own integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)
	id, _ := uuid.Parse("b24fa8eb-9df6-406c-9b45-763d7b5a5078")

	_, err := albumService.GetOwn(context.Background(), id)

	t.Assert().Nil(err)
}

func (s *AllSuite) TestAlbumGetByMusicianID(t provider.T) {
	t.Parallel()
	t.Title("Album get by musician id integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)
	id, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fa")

	albums, err := albumService.GetByMusicianID(context.Background(), id)

	t.Assert().Nil(err)
	t.Assert().NotNil(albums)
}

func (s *AllSuite) TestAlbumPublish(t provider.T) {
	t.Title("Album publish integration test")
	t.Parallel()
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresAlbumRepository(s.db)
	albumService := service.NewAlbumService(repo, s.logger)
	id, _ := uuid.Parse("b24fa8eb-9df6-406c-9b45-763d7b5a5078")

	err := albumService.Publish(context.Background(), id)

	t.Assert().Nil(err)
}

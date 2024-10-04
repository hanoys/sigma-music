package e2e

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/miniostorage"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/hanoys/sigma-music/internal/service/test/builder"
	"github.com/jmoiron/sqlx"
	minio2 "github.com/minio/minio-go/v7"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/testcontainers/testcontainers-go/modules/minio"
	testpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap"
)

type E2ESuite struct {
	suite.Suite
	logger         *zap.Logger
	hash           *hash.HashPasswordProvider
	container      *testpg.PostgresContainer
	minioContainer *minio.MinioContainer
	db             *sqlx.DB
	minioClient    *minio2.Client
}

func (s *E2ESuite) BeforeAll(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()

	s.hash = hash.NewHashPasswordProvider()
}

func (s *E2ESuite) BeforeEach(t provider.T) {
	ctx := context.Background()
	var err error
	s.container, err = newPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s.minioContainer, err = newMinioContainer(context.Background())
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

	url, err = s.minioContainer.ConnectionString(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	s.minioClient, err = newMinioClient(url)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *E2ESuite) AfterAll(t provider.T) {
	if err := s.container.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func (s *E2ESuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *E2ESuite) TestE2E(t provider.T) {
	t.Title("user e2e integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresUserRepository(s.db)
	userService := service.NewUserService(repo, s.hash, s.logger)
	createUserReq := builder.NewUserServiceCreateRequestBuilder().
		Default().
		SetName("Test").
		SetEmail("test").
		SetPhone("+7").Build()
	genreRepo := postgres.NewPostgresGenreRepository(s.db)
	storage := miniostorage.NewTrackStorage(s.minioClient, "music")
	trackRepo := postgres.NewPostgresTrackRepository(s.db)
	trackService := service.NewTrackService(trackRepo, storage, genreRepo, s.logger)
	commentRepo := postgres.NewPostgresCommentRepository(s.db)
	commentService := service.NewCommentService(commentRepo, s.logger)

	albumID, _ := uuid.Parse("b24fa8eb-9df6-406c-9b45-763d7b5a5078")

	user, err := userService.Register(context.Background(), createUserReq)

	t.Assert().Nil(err)
	t.Assert().Equal(user.Name, createUserReq.Name)

	tracks, err := trackService.GetByAlbumID(context.Background(), albumID)

	t.Assert().Nil(err)
	t.Assert().NotNil(tracks)

	track := tracks[0]

	err = trackService.AddToUserFavorites(context.Background(), track.ID, user.ID)

	t.Assert().Nil(err)

	postCommentReq := builder.NewPostCommentServiceRequestBuilder().
		Default().
		SetUserID(user.ID).
		SetTrackID(track.ID).
		Build()

	comment, err := commentService.Post(context.Background(), postCommentReq)

	t.Assert().Nil(err)
	t.Assert().Equal(comment.UserID, postCommentReq.UserID)
	t.Assert().Equal(comment.TrackID, postCommentReq.TrackID)
}

func TestSuite(t *testing.T) {
	suite.RunSuite(t, new(E2ESuite))
}

package integrationtest

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres/test/builder"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	testpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap"
)

type CommentSuite struct {
	suite.Suite
	logger    *zap.Logger
	hash      *hash.HashPasswordProvider
	container *testpg.PostgresContainer
	db        *sqlx.DB
}

func (s *CommentSuite) BeforeAll(t provider.T) {
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

func (s *CommentSuite) BeforeEach(t provider.T) {
	url, err := s.container.ConnectionString(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	s.db, err = newPostgresDB(url)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *CommentSuite) AfterAll(t provider.T) {
	if err := s.container.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func (s *CommentSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *CommentSuite) TestPost(t provider.T) {
    t.Title("comment post integration test")
	repo := postgres.NewPostgresCommentRepository(s.db)
	commentService := service.NewCommentService(repo, s.logger)
	userID, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fc")
	trackID, _ := uuid.Parse("41623ac1-b98d-4478-a10f-870a80c697b6")
	req := builder.NewPostCommentServiceRequestBuilder().
		Default().
		SetUserID(userID).
		SetTrackID(trackID).
		Build()

	comment, err := commentService.Post(context.Background(), req)

	t.Assert().Nil(err)
	t.Assert().Equal(req.Text, comment.Text)
}

func (s *CommentSuite) TestGetUserComments(t provider.T) {
    t.Title("comment get user comments integration test")
	repo := postgres.NewPostgresCommentRepository(s.db)
	commentService := service.NewCommentService(repo, s.logger)
	userID, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fc")

	comments, err := commentService.GetUserComments(context.Background(), userID)

	t.Assert().Nil(err)
	t.Assert().NotNil(comments)
}

func (s *CommentSuite) TestGetCommentsOnTrack(t provider.T) {
    t.Title("comment get comments on track integration test")
	repo := postgres.NewPostgresCommentRepository(s.db)
	commentService := service.NewCommentService(repo, s.logger)
	trackID, _ := uuid.Parse("41623ac1-b98d-4478-a10f-870a80c697b6")

	comments, err := commentService.GetCommentsOnTrack(context.Background(), trackID)

	t.Assert().Nil(err)
	t.Assert().NotNil(comments)
}

func TestCommentSuite(t *testing.T) {
	suite.RunSuite(t, new(CommentSuite))
}

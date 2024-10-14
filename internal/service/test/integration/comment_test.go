package integrationtest

import (
	"context"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres/test/builder"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)


func (s *AllSuite) TestPost(t provider.T) {
    t.Parallel()
	t.Title("comment post integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
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

func (s *AllSuite) TestGetUserComments(t provider.T) {
    t.Parallel()
	t.Title("comment get user comments integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresCommentRepository(s.db)
	commentService := service.NewCommentService(repo, s.logger)
	userID, _ := uuid.Parse("1add32df-d439-4fd1-9d4c-bef946b4a1fc")

	comments, err := commentService.GetUserComments(context.Background(), userID)

	t.Assert().Nil(err)
	t.Assert().NotNil(comments)
}

func (s *AllSuite) TestGetCommentsOnTrack(t provider.T) {
    t.Parallel()
	t.Title("comment get comments on track integration test")
	if isPreviousTestsFailed() {
		t.Skip()
	}
	repo := postgres.NewPostgresCommentRepository(s.db)
	commentService := service.NewCommentService(repo, s.logger)
	trackID, _ := uuid.Parse("41623ac1-b98d-4478-a10f-870a80c697b6")

	comments, err := commentService.GetCommentsOnTrack(context.Background(), trackID)

	t.Assert().Nil(err)
	t.Assert().NotNil(comments)
}

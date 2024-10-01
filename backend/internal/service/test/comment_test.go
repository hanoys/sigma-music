package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/hanoys/sigma-music/internal/service/test/builder"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)


type CommentSuite struct {
	suite.Suite
	logger *zap.Logger
}

func (s *CommentSuite) BeforeEach(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()
}

type CommentPostSuite struct {
	CommentSuite
}

func (s *CommentPostSuite) CorrectRepositoryMock(repository *mocks.CommentRepository) {
	repository.
		On("Create", context.Background(), mock.AnythingOfType("domain.Comment")).
		Return(domain.Comment{}, nil)
}

func (s *CommentPostSuite) TestCorrect(t provider.T) {
    t.Parallel()
	t.Title("Comment post test correct")
	req := builder.NewPostCommentServiceRequestBuilder().Default().Build()
	repository := mocks.NewCommentRepository(t)
	commentService := service.NewCommentService(repository, s.logger)
	s.CorrectRepositoryMock(repository)

	_, err := commentService.Post(context.Background(), req)

	t.Assert().Nil(err)
}

func (s *CommentPostSuite) DuplicateRepositoryMock(repository *mocks.CommentRepository) {
	repository.
		On("Create", context.Background(), mock.AnythingOfType("domain.Comment")).
		Return(domain.Comment{}, ports.ErrCommentDuplicate)
}

func (s *CommentPostSuite) TestDuplicate(t provider.T) {
    t.Parallel()
	t.Title("Comment post test duplicate")
	req := builder.NewPostCommentServiceRequestBuilder().Default().Build()
	repository := mocks.NewCommentRepository(t)
	commentService := service.NewCommentService(repository, s.logger)
	s.DuplicateRepositoryMock(repository)

	_, err := commentService.Post(context.Background(), req)

	t.Assert().ErrorIs(err, ports.ErrCommentDuplicate)
}

func TestCommentPostSuite(t *testing.T) {
	suite.RunSuite(t, new(CommentPostSuite))
}

type CommentGetCommentsOnTrack struct {
	CommentSuite
}

func (s *CommentGetCommentsOnTrack) CorrectRepositoryMock(repository *mocks.CommentRepository, trackID uuid.UUID) {
	repository.
		On("GetByTrackID", context.Background(), trackID).
		Return([]domain.Comment{}, nil)
}

func (s *CommentGetCommentsOnTrack) TestCorrect(t provider.T) {
    t.Parallel()
	t.Title("Comment get comments on track test correct")
	trackID := uuid.New()
	repository := mocks.NewCommentRepository(t)
	commentService := service.NewCommentService(repository, s.logger)
	s.CorrectRepositoryMock(repository, trackID)

	_, err := commentService.GetCommentsOnTrack(context.Background(), trackID)

	t.Assert().Nil(err)
}

func (s *CommentGetCommentsOnTrack) NotFoundRepositoryMock(repository *mocks.CommentRepository, trackID uuid.UUID) {
	repository.
		On("GetByTrackID", context.Background(), trackID).
		Return([]domain.Comment{}, ports.ErrCommentByTrackIDNotFound)
}

func (s *CommentGetCommentsOnTrack) TestIDNotFound(t provider.T) {
    t.Parallel()
	t.Title("Comment get comments on track test not found")
	trackID := uuid.New()
	repository := mocks.NewCommentRepository(t)
	commentService := service.NewCommentService(repository, s.logger)
	s.NotFoundRepositoryMock(repository, trackID)

	_, err := commentService.GetCommentsOnTrack(context.Background(), trackID)

	t.Assert().ErrorIs(err, ports.ErrCommentByTrackIDNotFound)
}

func TestCommentGetCommentsOnTrackSuite(t *testing.T) {
	suite.RunSuite(t, new(CommentGetCommentsOnTrack))
}

type CommentGetUserComments struct {
	CommentSuite
}

func (s *CommentGetUserComments) CorrectRepositoryMock(repository *mocks.CommentRepository, userID uuid.UUID) {
	repository.
		On("GetByUserID", context.Background(), userID).
		Return([]domain.Comment{}, nil)
}

func (s *CommentGetUserComments) TestCorrect(t provider.T) {
    t.Parallel()
	t.Title("Comment get user comments test correct")
	userID := uuid.New()
	repository := mocks.NewCommentRepository(t)
	commentService := service.NewCommentService(repository, s.logger)
	s.CorrectRepositoryMock(repository, userID)

	_, err := commentService.GetUserComments(context.Background(), userID)

	t.Assert().Nil(err)
}

func (s *CommentGetUserComments) NotFoundRepositoryMock(repository *mocks.CommentRepository, userID uuid.UUID) {
	repository.
		On("GetByUserID", context.Background(), userID).
		Return(nil, ports.ErrCommentByUserIDNotFound)
}

func (s *CommentGetUserComments) TestIDNotFound(t provider.T) {
    t.Parallel()
	t.Title("Comment get user comments test not found")
	userID := uuid.New()
	repository := mocks.NewCommentRepository(t)
	commentService := service.NewCommentService(repository, s.logger)
	s.NotFoundRepositoryMock(repository, userID)

	_, err := commentService.GetUserComments(context.Background(), userID)

	t.Assert().ErrorIs(err, ports.ErrCommentByUserIDNotFound)
}

func TestCommentGetUserCommentsSuite(t *testing.T) {
	suite.RunSuite(t, new(CommentGetUserComments))
}

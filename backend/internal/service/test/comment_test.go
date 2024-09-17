package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

var postCommentReq = ports.PostCommentServiceReq{
	UserID:  uuid.New(),
	TrackID: uuid.New(),
	Stars:   5,
	Text:    "Test Comment",
}

func TestCommentServicePost(t *testing.T) {
	tests := []struct {
		name           string
		repositoryMock func(repository *mocks.CommentRepository)
		req            ports.PostCommentServiceReq
		expected       error
	}{
		{
			name: "post comment success",
			req:  postCommentReq,
			repositoryMock: func(repository *mocks.CommentRepository) {
				repository.
					On("Create", context.Background(), mock.AnythingOfType("domain.Comment")).
					Return(domain.Comment{}, nil)
			},
			expected: nil,
		},
		{
			name: "post comment failure",
			req:  postCommentReq,
			repositoryMock: func(repository *mocks.CommentRepository) {
				repository.
					On("Create", context.Background(), mock.AnythingOfType("domain.Comment")).
					Return(domain.Comment{}, ports.ErrCommentDuplicate)
			},
			expected: ports.ErrCommentDuplicate,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			commentRepository := mocks.NewCommentRepository(t)
			commentService := service.NewCommentService(commentRepository, logger)
			test.repositoryMock(commentRepository)

			_, err := commentService.Post(context.Background(), test.req)
			if !errors.Is(err, test.expected) {
				t.Errorf("got %v, want %v", err, test.expected)
			}
		})
	}
}

func TestCommentServiceGetCommentsOnTrack(t *testing.T) {
	tests := []struct {
		name           string
		repositoryMock func(repository *mocks.CommentRepository)
		id             uuid.UUID
		expected       error
	}{
		{
			name: "post comment success",
			id:   uuid.New(),
			repositoryMock: func(repository *mocks.CommentRepository) {
				repository.
					On("GetByTrackID", context.Background(), mock.AnythingOfType("uuid.UUID")).
					Return([]domain.Comment{}, nil)
			},
			expected: nil,
		},
		{
			name: "post comment failure",
			id:   uuid.New(),
			repositoryMock: func(repository *mocks.CommentRepository) {
				repository.
					On("GetByTrackID", context.Background(), mock.AnythingOfType("uuid.UUID")).
					Return([]domain.Comment{}, ports.ErrCommentByTrackIDNotFound)
			},
			expected: ports.ErrCommentByTrackIDNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			commentRepository := mocks.NewCommentRepository(t)
			commentService := service.NewCommentService(commentRepository, logger)
			test.repositoryMock(commentRepository)

			_, err := commentService.GetCommentsOnTrack(context.Background(), test.id)
			if !errors.Is(err, test.expected) {
				t.Errorf("got %v, want %v", err, test.expected)
			}
		})
	}
}

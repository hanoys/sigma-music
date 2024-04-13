package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrPostComment = errors.New("can't post comment: internal error")
	ErrGetComments = errors.New("can't get comments: internal error")
)

type ICommentRepository interface {
	Create(ctx context.Context, comment domain.Comment) (domain.Comment, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Comment, error)
	GetByTrackID(ctx context.Context, trackID uuid.UUID) ([]domain.Comment, error)
}

type PostCommentServiceReq struct {
	UserID  uuid.UUID
	TrackID uuid.UUID
	Stars   int
	Text    string
}

type ICommentService interface {
	Post(ctx context.Context, comment PostCommentServiceReq) (domain.Comment, error)
	GetCommentsOnTrack(ctx context.Context, trackID uuid.UUID) ([]domain.Comment, error)
}

package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrCommentDuplicate         = errors.New("comment duplicate error")
	ErrCommentIDNotFound        = errors.New("comment with such id not found")
	ErrCommentByTrackIDNotFound = errors.New("comment with such track id not found")
	ErrInternalCommentRepo      = errors.New("comment repository internal error")
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
	GetUserComments(ctx context.Context, userID uuid.UUID) ([]domain.Comment, error)
}

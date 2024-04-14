package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type CommentService struct {
	repository ports.ICommentRepository
}

func NewCommentService(repo ports.ICommentRepository) *CommentService {
	return &CommentService{repository: repo}
}

func (cs *CommentService) Post(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	return cs.repository.Create(ctx, comment)
}

func (cs *CommentService) GetCommentsOnTrack(ctx context.Context, trackID uuid.UUID) ([]domain.Comment, error) {
	return cs.repository.GetByTrackID(ctx, trackID)
}

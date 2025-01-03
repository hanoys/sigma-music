package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type CommentService struct {
	repository ports.ICommentRepository
	logger     *zap.Logger
}

func NewCommentService(repo ports.ICommentRepository, logger *zap.Logger) *CommentService {
	return &CommentService{
		repository: repo,
		logger:     logger,
	}
}

func (cs *CommentService) Post(ctx context.Context, comment ports.PostCommentServiceReq) (domain.Comment, error) {
	postComment := domain.Comment{
		ID:      uuid.New(),
		UserID:  comment.UserID,
		TrackID: comment.TrackID,
		Stars:   comment.Stars,
		Text:    comment.Text,
	}

	comments, err := cs.GetUserComments(ctx, comment.UserID)
    if err != nil {
        return domain.Comment{}, err
    }

	for _, postedComment := range comments {
		if postedComment.UserID == comment.UserID && postedComment.TrackID == comment.TrackID {
			_, err = cs.repository.Delete(ctx, comment.UserID, postedComment.TrackID)
            if err != nil {
                return domain.Comment{}, err
            }
		}
	}

	comm, err := cs.repository.Create(ctx, postComment)
	if err != nil {
		cs.logger.Error("Failed to post comment", zap.Error(err), zap.String("Comment ID", postComment.ID.String()),
			zap.String("User ID", postComment.UserID.String()), zap.String("Track ID", postComment.TrackID.String()))

		return domain.Comment{}, err
	}

	cs.logger.Info("Comment successfully posted", zap.String("Comment ID", postComment.ID.String()),
		zap.String("User ID", postComment.UserID.String()), zap.String("Track ID", postComment.TrackID.String()))

	return comm, nil
}

func (cs *CommentService) Delete(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) (domain.Comment, error) {
	comment, err := cs.repository.Delete(ctx, userID, trackID)
	return comment, err
}

func (cs *CommentService) GetCommentsOnTrack(ctx context.Context, trackID uuid.UUID) ([]domain.Comment, error) {
	comments, err := cs.repository.GetByTrackID(ctx, trackID)
	if err != nil {
		cs.logger.Error("Failed to get comments on track", zap.Error(err))
		return nil, err
	}

	cs.logger.Info("Comments successfully received by track ID", zap.String("Track ID", trackID.String()))

	return comments, nil
}

func (cs *CommentService) GetUserComments(ctx context.Context, userID uuid.UUID) ([]domain.Comment, error) {
	comments, err := cs.repository.GetByUserID(ctx, userID)
	if err != nil {
		cs.logger.Error("Failed to get user comments", zap.Error(err), zap.String("User ID", userID.String()))
		return nil, err
	}

	cs.logger.Info("Comments successfully received by user ID", zap.String("User ID", userID.String()))

	return comments, nil
}

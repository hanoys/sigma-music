package dto

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type CommentDTO struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	TrackID uuid.UUID `json:"track_id"`
	Stars   int       `json:"stars"`
	Text    string    `json:"text"`
}

func CommentFromDomain(comment domain.Comment) CommentDTO {
	return CommentDTO{
		ID:      comment.ID,
		UserID:  comment.UserID,
		TrackID: comment.TrackID,
		Stars:   comment.Stars,
		Text:    comment.Text,
	}
}

type PostCommentDTO struct {
	Stars int    `json:"stars"`
	Text  string `json:"text"`
}

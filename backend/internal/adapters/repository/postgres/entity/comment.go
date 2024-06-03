package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type PgComment struct {
	ID      uuid.UUID `db:"id"`
	UserID  uuid.UUID `db:"user_id"`
	TrackID uuid.UUID `db:"track_id"`
	Stars   int       `db:"stars"`
	Text    string    `db:"comment_text"`
}

func (c *PgComment) ToDomain() domain.Comment {
	return domain.Comment{
		ID:      c.ID,
		UserID:  c.UserID,
		TrackID: c.TrackID,
		Stars:   c.Stars,
		Text:    c.Text,
	}
}

func NewPgComment(comment domain.Comment) PgComment {
	return PgComment{
		ID:      comment.ID,
		UserID:  comment.UserID,
		TrackID: comment.TrackID,
		Stars:   comment.Stars,
		Text:    comment.Text,
	}
}

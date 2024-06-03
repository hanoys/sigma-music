package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type MongoComment struct {
	ID      uuid.UUID `bson:"_id"`
	UserID  uuid.UUID `bson:"user_id"`
	TrackID uuid.UUID `bson:"track_id"`
	Stars   int       `bson:"stars"`
	Text    string    `bson:"comment_text"`
}

func (c *MongoComment) ToDomain() domain.Comment {
	return domain.Comment{
		ID:      c.ID,
		UserID:  c.UserID,
		TrackID: c.TrackID,
		Stars:   c.Stars,
		Text:    c.Text,
	}
}

func NewMongoComment(comment domain.Comment) MongoComment {
	return MongoComment{
		ID:      comment.ID,
		UserID:  comment.UserID,
		TrackID: comment.TrackID,
		Stars:   comment.Stars,
		Text:    comment.Text,
	}
}

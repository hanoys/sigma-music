package builder

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/ports"
)

type PostCommentServiceRequestBuilder struct {
	obj ports.PostCommentServiceReq
}

func NewPostCommentServiceRequestBuilder() *PostCommentServiceRequestBuilder {
	return new(PostCommentServiceRequestBuilder)
}

func (b *PostCommentServiceRequestBuilder) Build() ports.PostCommentServiceReq {
	return b.obj
}

func (b *PostCommentServiceRequestBuilder) Default() *PostCommentServiceRequestBuilder {
	b.obj = ports.PostCommentServiceReq{
		UserID:  uuid.New(),
		TrackID: uuid.New(),
		Stars:   5,
		Text:    "Test comment",
	}

	return b
}

func (b *PostCommentServiceRequestBuilder) SetUserID(userID uuid.UUID) *PostCommentServiceRequestBuilder {
	b.obj.UserID = userID
	return b
}

func (b *PostCommentServiceRequestBuilder) SetTrackID(trackID uuid.UUID) *PostCommentServiceRequestBuilder {
	b.obj.TrackID = trackID
	return b
}

func (b *PostCommentServiceRequestBuilder) SetStars(stars int) *PostCommentServiceRequestBuilder {
	b.obj.Stars = stars
	return b
}

func (b *PostCommentServiceRequestBuilder) SetText(text string) *PostCommentServiceRequestBuilder {
	b.obj.Text = text
	return b
}

package builder

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type CommentBuilder struct {
	obj domain.Comment
}

func NewCommentBuilder() *CommentBuilder {
	return new(CommentBuilder)
}

func (b *CommentBuilder) Build() domain.Comment {
	return b.obj
}

func (b *CommentBuilder) Default() *CommentBuilder {
	b.obj = domain.Comment{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		TrackID: uuid.New(),
		Stars:   0,
		Text:    "text",
	}
	return b
}

func (b *CommentBuilder) SetID(id uuid.UUID) *CommentBuilder {
	b.obj.ID = id
	return b
}

func (b *CommentBuilder) SetUserID(id uuid.UUID) *CommentBuilder {
	b.obj.UserID = id
	return b
}

func (b *CommentBuilder) SetTrackID(id uuid.UUID) *CommentBuilder {
	b.obj.TrackID = id
	return b
}

func (b *CommentBuilder) SetStars(stars int) *CommentBuilder {
	b.obj.Stars = stars
	return b
}

func (b *CommentBuilder) SetText(text string) *CommentBuilder {
	b.obj.Text = text
	return b
}

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

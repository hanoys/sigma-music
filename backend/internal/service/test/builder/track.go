package builder

import (
	"io"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type CreateTrackRequestBuilder struct {
	obj ports.CreateTrackReq
}

func NewCreateTrackRequestBuilder() *CreateTrackRequestBuilder {
	return new(CreateTrackRequestBuilder)
}

func (b *CreateTrackRequestBuilder) Build() ports.CreateTrackReq {
	return b.obj
}

func (b *CreateTrackRequestBuilder) Default() *CreateTrackRequestBuilder {
	b.obj = ports.CreateTrackReq{
		AlbumID:   uuid.New(),
		Name:      "track",
		TrackBLOB: nil,
		GenresID:  []uuid.UUID{},
	}
	return b
}

func (b *CreateTrackRequestBuilder) SetAlbumID(id uuid.UUID) *CreateTrackRequestBuilder {
	b.obj.AlbumID = id
	return b
}

func (b *CreateTrackRequestBuilder) SetName(name string) *CreateTrackRequestBuilder {
	b.obj.Name = name
	return b
}

func (b *CreateTrackRequestBuilder) SetTrackBLOB(blob io.Reader) *CreateTrackRequestBuilder {
	b.obj.TrackBLOB = blob
	return b
}

func (b *CreateTrackRequestBuilder) SetGenresID(ids []uuid.UUID) *CreateTrackRequestBuilder {
	b.obj.GenresID = ids
	return b
}

type PutTrackRequestBuilder struct {
	obj ports.PutTrackReq
}

func NewPutTrackReqeustBuilder() *PutTrackRequestBuilder {
	return new(PutTrackRequestBuilder)
}

func (b *PutTrackRequestBuilder) Build() ports.PutTrackReq {
	return b.obj
}

func (b *PutTrackRequestBuilder) Default() *PutTrackRequestBuilder {
	b.obj = ports.PutTrackReq{
		TrackID:   "trakid",
		TrackBLOB: nil,
	}
	return b
}

func (b *PutTrackRequestBuilder) SetTrackID(id string) *PutTrackRequestBuilder {
	b.obj.TrackID = id
	return b
}

func (b *PutTrackRequestBuilder) SetTrackBLOB(blob io.Reader) *PutTrackRequestBuilder {
	b.obj.TrackBLOB = blob
	return b
}

type TrackBuilder struct {
	obj domain.Track
}

func NewTrackBuilder() *TrackBuilder {
	return new(TrackBuilder)
}

func (b *TrackBuilder) Build() domain.Track {
	return b.obj
}

func (b *TrackBuilder) Default() *TrackBuilder {
	b.obj = domain.Track{
		ID:      uuid.New(),
		AlbumID: uuid.New(),
		Name:    "trackname",
		URL:     "url",
	}
	return b
}

func (b *TrackBuilder) SetID(id uuid.UUID) *TrackBuilder {
	b.obj.ID = id
	return b
}

func (b *TrackBuilder) SetAlbumID(id uuid.UUID) *TrackBuilder {
	b.obj.AlbumID = id
	return b
}

func (b *TrackBuilder) SetName(name string) *TrackBuilder {
	b.obj.Name = name
	return b
}

func (b *TrackBuilder) SetURL(url string) *TrackBuilder {
	b.obj.URL = url
	return b
}

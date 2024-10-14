package builder

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/ports"
)

type CreateAlbumServiceRequestBuilder struct {
	obj ports.CreateAlbumServiceReq
}

func NewCreateAlbumServiceRequestBuilder() *CreateAlbumServiceRequestBuilder {
	return new(CreateAlbumServiceRequestBuilder)
}

func (b *CreateAlbumServiceRequestBuilder) Build() ports.CreateAlbumServiceReq {
	return b.obj
}

func (b *CreateAlbumServiceRequestBuilder) Default() *CreateAlbumServiceRequestBuilder {
	b.obj = ports.CreateAlbumServiceReq{
		MusicianID:  uuid.New(),
		Name:        "Album name",
		Description: "Album description",
	}

	return b
}

func (b *CreateAlbumServiceRequestBuilder) SetMusicianID(musicianID uuid.UUID) *CreateAlbumServiceRequestBuilder {
	b.obj.MusicianID = musicianID
	return b
}

func (b *CreateAlbumServiceRequestBuilder) SetName(name string) *CreateAlbumServiceRequestBuilder {
	b.obj.Name = name
	return b
}

func (b *CreateAlbumServiceRequestBuilder) SetDescription(description string) *CreateAlbumServiceRequestBuilder {
	b.obj.Description = description
	return b
}

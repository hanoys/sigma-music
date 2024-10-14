package builder

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/hanoys/sigma-music/internal/domain"
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

type AlbumBuilder struct {
	obj domain.Album
}

func NewAlbumBuilder() *AlbumBuilder {
	return new(AlbumBuilder)
}

func (b *AlbumBuilder) Build() domain.Album {
	return b.obj
}

func (b *AlbumBuilder) Default() *AlbumBuilder {
	b.obj = domain.Album{
		ID:          uuid.New(),
		Name:        "name",
		Description: "description",
		Published:   false,
		ReleaseDate: null.Time{},
	}
	return b
}

func (b *AlbumBuilder) SetID(id uuid.UUID) *AlbumBuilder {
	b.obj.ID = id
	return b
}

func (b *AlbumBuilder) SetName(name string) *AlbumBuilder {
	b.obj.Name = name
	return b
}

func (b *AlbumBuilder) SetDescription(description string) *AlbumBuilder {
	b.obj.Description = description
	return b
}

func (b *AlbumBuilder) SetPublished(published bool) *AlbumBuilder {
	b.obj.Published = published
	return b
}

func (b *AlbumBuilder) SetReleaseDate(releaseDate null.Time) *AlbumBuilder {
	b.obj.ReleaseDate = releaseDate
	return b
}

package builder

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type MusicianServiceCreateRequestBuilder struct {
	obj ports.MusicianServiceCreateRequest
}

func NewMusicianServiceCreateRequestBuilder() *MusicianServiceCreateRequestBuilder {
	return new(MusicianServiceCreateRequestBuilder)
}

func (b *MusicianServiceCreateRequestBuilder) Build() ports.MusicianServiceCreateRequest {
	return b.obj
}

func (b *MusicianServiceCreateRequestBuilder) Default() *MusicianServiceCreateRequestBuilder {
	b.obj = ports.MusicianServiceCreateRequest{
		Name:        "test musician",
		Email:       "test.musician@mail.com",
		Password:    "testpassword",
		Country:     "USA",
		Description: "Test description",
	}
	return b
}

func (b *MusicianServiceCreateRequestBuilder) SetName(name string) {
	b.obj.Name = name
}

func (b *MusicianServiceCreateRequestBuilder) SetEmail(email string) {
	b.obj.Email = email
}

func (b *MusicianServiceCreateRequestBuilder) SetPassword(password string) {
	b.obj.Password = password
}

func (b *MusicianServiceCreateRequestBuilder) SetCountry(country string) {
	b.obj.Country = country
}

func (b *MusicianServiceCreateRequestBuilder) SetDescription(description string) {
	b.obj.Description = description
}

type MusicianBuilder struct {
	obj domain.Musician
}

func NewMusicianBuilder() *MusicianBuilder {
	return new(MusicianBuilder)
}

func (b *MusicianBuilder) Build() domain.Musician {
	return b.obj
}

func (b *MusicianBuilder) Default() *MusicianBuilder {
	b.obj = domain.Musician{
		ID:          uuid.New(),
		Name:        "musician",
		Email:       "musician@mail.com",
		Password:    "musician",
		Salt:        "salt",
		Country:     "pass",
		Description: "description",
	}

	return b
}

func (b *MusicianBuilder) SetID(id uuid.UUID) *MusicianBuilder {
	b.obj.ID = id
	return b
}

func (b *MusicianBuilder) SetName(name string) *MusicianBuilder {
	b.obj.Name = name
	return b
}

func (b *MusicianBuilder) SetEmail(email string) *MusicianBuilder {
	b.obj.Email = email
	return b
}

func (b *MusicianBuilder) SetPassword(password string) *MusicianBuilder {
	b.obj.Password = password
	return b
}

func (b *MusicianBuilder) SetSalt(salt string) *MusicianBuilder {
	b.obj.Salt = salt
	return b
}

func (b *MusicianBuilder) SetCountry(country string) *MusicianBuilder {
	b.obj.Country = country
	return b
}

func (b *MusicianBuilder) SetDescription(description string) *MusicianBuilder {
	b.obj.Description = description
	return b
}

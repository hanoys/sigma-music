package builder

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type GenreBuilder struct {
	obj domain.Genre
}

func NewGenreBuilder() *GenreBuilder {
	return new(GenreBuilder)
}

func (b *GenreBuilder) Build() domain.Genre {
	return b.obj
}

func (b *GenreBuilder) Default() *GenreBuilder {
	b.obj = domain.Genre{
		ID:   uuid.New(),
		Name: "name",
	}
	return b
}

func (b *GenreBuilder) SetID(id uuid.UUID) *GenreBuilder {
	b.obj.ID = id
	return b
}

func (b *GenreBuilder) SetName(name string) *GenreBuilder {
	b.obj.Name = name
	return b
}

package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type PgGenre struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (g *PgGenre) ToDomain() domain.Genre {
	return domain.Genre{
		ID:   g.ID,
		Name: g.Name,
	}
}

func NewPgGenre(genre domain.Genre) PgGenre {
	return PgGenre{
		ID:   genre.ID,
		Name: genre.Name,
	}
}

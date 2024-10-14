package entity

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/hanoys/sigma-music/internal/domain"
)

type PgAlbum struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Published   bool      `db:"published"`
	ReleaseDate null.Time `db:"release_date"`
}

func (a *PgAlbum) ToDomain() domain.Album {
	return domain.Album{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		Published:   a.Published,
		ReleaseDate: a.ReleaseDate,
	}
}

func NewPgAlbum(album domain.Album) PgAlbum {
	return PgAlbum{
		ID:          album.ID,
		Name:        album.Name,
		Description: album.Description,
		Published:   album.Published,
		ReleaseDate: album.ReleaseDate,
	}
}

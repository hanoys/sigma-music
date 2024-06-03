package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type PgTrack struct {
	ID      uuid.UUID `db:"id"`
	AlbumID uuid.UUID `db:"album_id"`
	Name    string    `db:"name"`
	URL     string    `db:"url"`
}

func (t *PgTrack) ToDomain() domain.Track {
	return domain.Track{
		ID:      t.ID,
		AlbumID: t.AlbumID,
		Name:    t.Name,
		URL:     t.URL,
	}
}

func NewPgTrack(track domain.Track) PgTrack {
	return PgTrack{
		ID:      track.ID,
		AlbumID: track.AlbumID,
		Name:    track.Name,
		URL:     track.URL,
	}
}

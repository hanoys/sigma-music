package dto

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type TrackDTO struct {
	ID      uuid.UUID `json:"id"`
	AlbumID uuid.UUID `json:"album_id"`
	Name    string    `json:"name"`
	URL     string    `json:"url"`
}

func TrackFromDomain(track domain.Track) TrackDTO {
	return TrackDTO{
		ID:      track.ID,
		AlbumID: track.AlbumID,
		Name:    track.Name,
		URL:     track.URL,
	}
}

type CreateTrackDTO struct {
	Name     string   `json:"name" binding:"required"`
	GenreIDs []string `json:"genres" binding:"omitempty"`
}

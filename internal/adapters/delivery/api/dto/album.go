package dto

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type AlbumDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Published   bool      `json:"published"`
	ReleaseDate string    `json:"release_date"`
	ImageURL    string    `json:"image_url"`
}

func AlbumFromDomain(album domain.Album) AlbumDTO {
	albumDTO := AlbumDTO{
		ID:          album.ID,
		Name:        album.Name,
		Description: album.Description,
		Published:   album.Published,
		ImageURL:    album.ImageURL.ValueOrZero(),
	}

	if album.ReleaseDate.IsZero() {
		albumDTO.ReleaseDate = "Not released"
	} else {
		albumDTO.ReleaseDate = album.ReleaseDate.Time.String()
	}

	return albumDTO
}

type CreateAlbumDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

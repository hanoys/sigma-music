package dto

import "github.com/hanoys/sigma-music/internal/domain"

type AlbumDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Published   bool   `json:"published"`
	ReleaseDate string `json:"release_date"`
}

func AlbumFromDomain(album domain.Album) AlbumDTO {
	albumDTO := AlbumDTO{
		ID:          album.ID.String(),
		Name:        album.Name,
		Description: album.Description,
		Published:   album.Published,
	}

	if album.ReleaseDate.IsZero() {
		albumDTO.ReleaseDate = "Not released"
	} else {
		albumDTO.ReleaseDate = album.ReleaseDate.Time.String()
	}

	return albumDTO
}

type CreateAlbumDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

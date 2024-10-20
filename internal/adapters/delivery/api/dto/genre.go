package dto

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type GenreDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func GenreFromDomain(genre domain.Genre) GenreDTO {
	return GenreDTO{
		ID:   genre.ID,
		Name: genre.Name,
	}
}

type AddForTrackDTO struct {
	GenreIDs []string `json:"genres" binding:"omitempty"`
}

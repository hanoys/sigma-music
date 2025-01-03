package dto

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type MusicianDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Country     string    `json:"country"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
}

func MusicianFromDomain(musician domain.Musician) MusicianDTO {
	return MusicianDTO{
		ID:          musician.ID,
		Name:        musician.Name,
		Email:       musician.Email,
		Country:     musician.Country,
		Description: musician.Description,
		ImageURL:    musician.ImageURL.ValueOrZero(),
	}
}

type RegisterMusicianDTO struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Country     string `json:"country" binding:"required"`
	Description string `json:"description" binding:"required"`
}

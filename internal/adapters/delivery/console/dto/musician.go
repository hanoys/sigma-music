package dto

import (
	"fmt"
	"github.com/hanoys/sigma-music/internal/domain"
)

type MusicianDTO struct {
	ID          string
	Name        string
	Email       string
	Country     string
	Description string
}

func NewMusicianDTO(user domain.Musician) MusicianDTO {
	return MusicianDTO{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		Country:     user.Country,
		Description: user.Description,
	}
}

func (u MusicianDTO) Print() {
	fmt.Println("ID:", u.ID)
	fmt.Println("Name:", u.Name)
	fmt.Println("Email:", u.Email)
	fmt.Println("Country:", u.Country)
	fmt.Println("Description:", u.Description)
}

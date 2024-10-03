package dto

import (
	"fmt"
	"github.com/hanoys/sigma-music/internal/domain"
)

type GenreDTO struct {
	ID   string
	Name string
}

func NewGenreDTO(g domain.Genre) GenreDTO {
	return GenreDTO{
		ID:   g.ID.String(),
		Name: g.Name,
	}
}

func (g GenreDTO) Print() {
	fmt.Println("ID:", g.ID)
	fmt.Println("Name:", g.Name)
}

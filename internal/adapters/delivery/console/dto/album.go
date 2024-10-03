package dto

import (
	"fmt"
	"github.com/hanoys/sigma-music/internal/domain"
)

type CreateAlbumDTO struct {
	Name        string
	Description string
}

func InputCreateAlbumDTO(albumDTO *CreateAlbumDTO) {
	fmt.Print("Name: ")
	fmt.Scan(&albumDTO.Name)

	fmt.Print("Description: ")
	fmt.Scan(&albumDTO.Description)
}

type AlbumDTO struct {
	ID          string
	Name        string
	Description string
	Published   bool
	ReleaseDate string
}

func NewAlbumDTO(album domain.Album) AlbumDTO {
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

func (a AlbumDTO) Print() {
	fmt.Println("ID:", a.ID)
	fmt.Println("Nane:", a.Name)
	fmt.Println("Description:", a.Description)
	fmt.Println("Published:", a.Published)
	fmt.Println("Release Date:", a.ReleaseDate)
}

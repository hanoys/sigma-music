package dto

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"io"
	"os"
	"path/filepath"
)

type CreateTrackDTO struct {
	AlbumID  uuid.UUID
	Track    io.Reader
	Name     string
	GenresID []uuid.UUID
}

func InputCreateTrackDTO(c *CreateTrackDTO) error {
	fmt.Print("Path to music: ")
	var path string
	fmt.Scan(&path)

	if filepath.Ext(path) != ".mp3" {
		return errors.New("invalid file extension")
	}

	f, err := os.Open(path)
	if err != nil {
		return errors.New("open file error")
	}

	c.Track = f

	var id string
	fmt.Print("Album ID:")
	fmt.Scan(&id)

	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("incorrect id")
	}

	c.AlbumID = uid

	fmt.Print("Name:")
	fmt.Scan(&c.Name)

	var genresID []uuid.UUID
	fmt.Println("Input genres id:")
	for {
		var genreID string
		fmt.Print("Genre ID:")
		_, err = fmt.Scan(&genreID)

		if err != nil {
			break
		}

		genreUUID, err := uuid.Parse(genreID)
		if err != nil {
			fmt.Println("incorrect genre id")
		} else {
			genresID = append(genresID, genreUUID)
		}
	}

	c.GenresID = genresID

	return nil
}

type TrackDTO struct {
	ID      string
	AlbumID string
	Name    string
	URL     string
}

func NewTrackDTO(track domain.Track) TrackDTO {
	return TrackDTO{
		ID:      track.ID.String(),
		AlbumID: track.AlbumID.String(),
		Name:    track.Name,
		URL:     track.URL,
	}
}

func (t TrackDTO) Print() {
	fmt.Println("ID:", t.ID)
	fmt.Println("Album ID:", t.AlbumID)
	fmt.Println("Name:", t.Name)
	fmt.Println("URL:", t.URL)
}

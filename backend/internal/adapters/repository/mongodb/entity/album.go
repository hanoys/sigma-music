package entity

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/hanoys/sigma-music/internal/domain"
)

type MongoAlbum struct {
	ID          uuid.UUID `bson:"_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Published   bool      `bson:"published"`
	ReleaseDate null.Time `bson:"release_date"`
}

func (a *MongoAlbum) ToDomain() domain.Album {
	return domain.Album{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		Published:   a.Published,
		ReleaseDate: a.ReleaseDate,
	}
}

func NewMongoAlbum(album domain.Album) MongoAlbum {
	return MongoAlbum{
		ID:          album.ID,
		Name:        album.Name,
		Description: album.Description,
		Published:   album.Published,
		ReleaseDate: album.ReleaseDate,
	}
}

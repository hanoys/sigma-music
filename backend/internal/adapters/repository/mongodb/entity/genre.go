package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type MongoGenre struct {
	ID   uuid.UUID `bson:"id"`
	Name string    `bson:"name"`
}

func (g *MongoGenre) ToDomain() domain.Genre {
	return domain.Genre{
		ID:   g.ID,
		Name: g.Name,
	}
}

func NewMongoGenre(genre domain.Genre) MongoGenre {
	return MongoGenre{
		ID:   genre.ID,
		Name: genre.Name,
	}
}

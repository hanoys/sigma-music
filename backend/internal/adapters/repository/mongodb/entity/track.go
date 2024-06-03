package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type MongoTrack struct {
	ID      uuid.UUID `bson:"_id"`
	AlbumID uuid.UUID `bson:"album_id"`
	Name    string    `bson:"name"`
	URL     string    `bson:"url"`
}

func (t *MongoTrack) ToDomain() domain.Track {
	return domain.Track{
		ID:      t.ID,
		AlbumID: t.AlbumID,
		Name:    t.Name,
		URL:     t.URL,
	}
}

func NewMongoTrack(track domain.Track) MongoTrack {
	return MongoTrack{
		ID:      track.ID,
		AlbumID: track.AlbumID,
		Name:    track.Name,
		URL:     track.URL,
	}
}

package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type MongoUserMusiciansStat struct {
	MusicianID  uuid.UUID `bson:"musician_id"`
	UserID      uuid.UUID `bson:"user_id"`
	ListenCount int64     `bson:"cnt"`
}

func (ums *MongoUserMusiciansStat) ToDomain() domain.UserMusiciansStat {
	return domain.UserMusiciansStat{
		MusicianID:  ums.MusicianID,
		UserID:      ums.UserID,
		ListenCount: ums.ListenCount,
	}
}

type MongoUserGenresStat struct {
	GenreID     uuid.UUID `bson:"genre_id"`
	UserID      uuid.UUID `bson:"user_id"`
	ListenCount int64     `bson:"cnt"`
}

func (ugs *MongoUserGenresStat) ToDomain() domain.UserGenresStat {
	return domain.UserGenresStat{
		GenreID:     ugs.GenreID,
		UserID:      ugs.UserID,
		ListenCount: ugs.ListenCount,
	}
}

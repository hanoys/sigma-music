package domain

import "github.com/google/uuid"

type UserMusiciansStat struct {
	MusicianID  uuid.UUID
	UserID      uuid.UUID
	ListenCount int64
}

type UserGenresStat struct {
	GenreID     uuid.UUID
	UserID      uuid.UUID
	ListenCount int64
}

type MusicianStat struct {
	MusicianID   uuid.UUID
	MusicianName string
	ListenCount  int64
}

type GenreStat struct {
	GenreID          uuid.UUID
	GenreName        string
	ListenPercentage int64
}

type ListenReport struct {
	UserID                  uuid.UUID
	MostListenedMusiciansID []MusicianStat
	ListenedGenres          []GenreStat
	ListenCount             int64
}

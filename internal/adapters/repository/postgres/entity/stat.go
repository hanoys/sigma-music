package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type PgUserMusiciansStat struct {
	MusicianID  uuid.UUID `db:"musician_id"`
	UserID      uuid.UUID `db:"user_id"`
	ListenCount int64     `db:"cnt"`
}

func (ums *PgUserMusiciansStat) ToDomain() domain.UserMusiciansStat {
	return domain.UserMusiciansStat{
		MusicianID:  ums.MusicianID,
		UserID:      ums.UserID,
		ListenCount: ums.ListenCount,
	}
}

type PgUserGenresStat struct {
	GenreID     uuid.UUID `db:"genre_id"`
	UserID      uuid.UUID `db:"user_id"`
	ListenCount int64     `db:"cnt"`
}

func (ugs *PgUserGenresStat) ToDomain() domain.UserGenresStat {
	return domain.UserGenresStat{
		GenreID:     ugs.GenreID,
		UserID:      ugs.UserID,
		ListenCount: ugs.ListenCount,
	}
}

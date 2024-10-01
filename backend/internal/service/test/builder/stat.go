package builder

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type UserMusiciansStatBuilder struct {
	obj domain.UserMusiciansStat
}

func NewUserMusiciansStatBuilder() *UserMusiciansStatBuilder {
	return new(UserMusiciansStatBuilder)
}

func (b *UserMusiciansStatBuilder) Build() domain.UserMusiciansStat {
	return b.obj
}

func (b *UserMusiciansStatBuilder) Default() *UserMusiciansStatBuilder {
	b.obj = domain.UserMusiciansStat{
		MusicianID:  uuid.New(),
		UserID:      uuid.New(),
		ListenCount: 0,
	}
	return b
}

func (b *UserMusiciansStatBuilder) SetMusicianID(id uuid.UUID) *UserMusiciansStatBuilder {
	b.obj.MusicianID = id
	return b
}

func (b *UserMusiciansStatBuilder) SetUserID(id uuid.UUID) *UserMusiciansStatBuilder {
	b.obj.UserID = id
	return b
}

func (b *UserMusiciansStatBuilder) SetListenCount(count int64) *UserMusiciansStatBuilder {
	b.obj.ListenCount = count
	return b
}

type UserGenresStatBuilder struct {
	obj domain.UserGenresStat
}

func NewUserGenresStatBuilder() *UserGenresStatBuilder {
	return new(UserGenresStatBuilder)
}

func (b *UserGenresStatBuilder) Build() domain.UserGenresStat {
	return b.obj
}

func (b *UserGenresStatBuilder) Default() *UserGenresStatBuilder {
	b.obj = domain.UserGenresStat{
		GenreID:     uuid.New(),
		UserID:      uuid.New(),
		ListenCount: 0,
	}
	return b
}

func (b *UserGenresStatBuilder) SetGenreID(id uuid.UUID) *UserGenresStatBuilder {
	b.obj.GenreID = id
	return b
}

func (b *UserGenresStatBuilder) SetUserID(id uuid.UUID) *UserGenresStatBuilder {
	b.obj.UserID = id
	return b
}

func (b *UserGenresStatBuilder) SetListenCount(count int64) *UserGenresStatBuilder {
	b.obj.ListenCount = count
	return b
}

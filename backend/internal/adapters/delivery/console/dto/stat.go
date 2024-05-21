package dto

import (
	"fmt"
	"github.com/hanoys/sigma-music/internal/domain"
)

type MusicianStatDTO struct {
	MusicianID   string
	MusicianName string
	ListenCount  int64
}

func NewMusicianStatDTO(m domain.MusicianStat) MusicianStatDTO {
	return MusicianStatDTO{
		MusicianID:   m.MusicianID.String(),
		MusicianName: m.MusicianName,
		ListenCount:  m.ListenCount,
	}
}

type GenreStatDTO struct {
	GenreID          string
	GenreName        string
	ListenPercentage int64
}

func NewGenreStatDTO(g domain.GenreStat) GenreStatDTO {
	return GenreStatDTO{
		GenreID:          g.GenreID.String(),
		GenreName:        g.GenreName,
		ListenPercentage: g.ListenPercentage,
	}
}

type ListenReportDTO struct {
	UserID                string
	MostListenedMusicians []MusicianStatDTO
	ListenedGenres        []GenreStatDTO
	ListenCount           int64
}

func NewListenReportDTO(l domain.ListenReport) ListenReportDTO {
	mostListenedMusiciansDTO := make([]MusicianStatDTO, len(l.MostListenedMusicians))
	listenedGenresDTO := make([]GenreStatDTO, len(l.ListenedGenres))

	for i, stat := range l.MostListenedMusicians {
		mostListenedMusiciansDTO[i] = NewMusicianStatDTO(stat)
	}

	for i, stat := range l.ListenedGenres {
		listenedGenresDTO[i] = NewGenreStatDTO(stat)
	}

	return ListenReportDTO{
		UserID:                l.UserID.String(),
		MostListenedMusicians: mostListenedMusiciansDTO,
		ListenedGenres:        listenedGenresDTO,
		ListenCount:           l.ListenCount,
	}
}

func (l ListenReportDTO) Print() {
	fmt.Println("User ID:", l.UserID)
	fmt.Println("Listen count:", l.ListenCount)
	fmt.Println()

	for _, stat := range l.MostListenedMusicians {
		fmt.Println("Musician:", stat.MusicianName, "Listen count:", stat.ListenCount)
	}

	for _, stat := range l.ListenedGenres {
		fmt.Println("Genre:", stat.GenreName, "Listen percentage:", stat.ListenPercentage)
	}
}

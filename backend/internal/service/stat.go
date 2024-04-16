package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type StatService struct {
	repository      ports.IStatRepository
	genreService    ports.IGenreService
	musicianService ports.IMusicianService
}

func NewStatService(repo ports.IStatRepository) *StatService {
	return &StatService{repository: repo}
}

func (ss *StatService) Add(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error {
	return ss.repository.Add(ctx, userID, trackID)
}

func (ss *StatService) FormReport(ctx context.Context, userID uuid.UUID) (domain.ListenReport, error) {
	var listenReport domain.ListenReport
	listenReport.UserID = userID

	listenedMusicians, err := ss.repository.GetMostListenedMusicians(ctx, userID, 3)
	if err != nil {
		return domain.ListenReport{}, err
	}

	listenedGenres, err := ss.repository.GetListenedGenres(ctx, userID)
	if err != nil {
		return domain.ListenReport{}, err
	}

	var listenCountSum int64
	for _, userMusicianStat := range listenedMusicians {
		listenCountSum += userMusicianStat.ListenCount
	}

	listenReport.ListenCount = listenCountSum

	listenReport.MostListenedMusiciansID = make([]domain.MusicianStat, len(listenedMusicians))
	for i, userMusicianStat := range listenedMusicians {
		musician, err := ss.musicianService.GetByID(ctx, userMusicianStat.MusicianID)
		if err != nil {
			return domain.ListenReport{}, err
		}

		listenReport.MostListenedMusiciansID[i] = domain.MusicianStat{
			MusicianID:   userMusicianStat.MusicianID,
			MusicianName: musician.Name,
			ListenCount:  userMusicianStat.ListenCount,
		}
	}

	var genresListenCountSum int64
	for _, userGenreStat := range listenedGenres {
		genresListenCountSum += userGenreStat.ListenCount
	}

	listenReport.ListenedGenres = make([]domain.GenreStat, len(listenedGenres))
	for i, userGenreStat := range listenedGenres {
		genre, err := ss.genreService.GetByID(ctx, userGenreStat.GenreID)
		if err != nil {
			return domain.ListenReport{}, err
		}

		listenReport.ListenedGenres[i] = domain.GenreStat{
			GenreID:          userGenreStat.GenreID,
			GenreName:        genre.Name,
			ListenPercentage: genresListenCountSum / userGenreStat.ListenCount,
		}
	}

	return listenReport, nil
}

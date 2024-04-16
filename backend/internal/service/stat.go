package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"math"
)

type StatService struct {
	repository      ports.IStatRepository
	genreService    ports.IGenreService
	musicianService ports.IMusicianService
}

func NewStatService(repo ports.IStatRepository, genreService ports.IGenreService,
	musService ports.IMusicianService) *StatService {

	return &StatService{
		repository:      repo,
		genreService:    genreService,
		musicianService: musService,
	}
}

func (ss *StatService) Add(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error {
	return ss.repository.Add(ctx, userID, trackID)
}

func (ss *StatService) fillMostListenedMusicians(ctx context.Context, listenReport *domain.ListenReport,
	listenedMusicians []domain.UserMusiciansStat) error {

	listenReport.MostListenedMusiciansID = make([]domain.MusicianStat, len(listenedMusicians))
	for i, userMusicianStat := range listenedMusicians {
		musician, err := ss.musicianService.GetByID(ctx, userMusicianStat.MusicianID)
		if err != nil {
			return err
		}

		listenReport.MostListenedMusiciansID[i] = domain.MusicianStat{
			MusicianID:   userMusicianStat.MusicianID,
			MusicianName: musician.Name,
			ListenCount:  userMusicianStat.ListenCount,
		}
	}

	return nil
}

func (ss *StatService) fillListenedGenres(ctx context.Context, listenReport *domain.ListenReport,
	listenedGenres []domain.UserGenresStat) error {

	var genresListenCountSum int64
	for _, userGenreStat := range listenedGenres {
		genresListenCountSum += userGenreStat.ListenCount
	}

	listenReport.ListenedGenres = make([]domain.GenreStat, len(listenedGenres))
	for i, userGenreStat := range listenedGenres {
		genre, err := ss.genreService.GetByID(ctx, userGenreStat.GenreID)
		if err != nil {
			return err
		}

		listenReport.ListenedGenres[i] = domain.GenreStat{
			GenreID:   userGenreStat.GenreID,
			GenreName: genre.Name,
			ListenPercentage: int64(math.Round(float64(userGenreStat.ListenCount) /
				float64(genresListenCountSum) * 100)),
		}
	}

	return nil
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

	err = ss.fillMostListenedMusicians(ctx, &listenReport, listenedMusicians)
	if err != nil {
		return domain.ListenReport{}, err
	}

	err = ss.fillListenedGenres(ctx, &listenReport, listenedGenres)
	if err != nil {
		return domain.ListenReport{}, err
	}

	var listenCountSum int64
	for _, userMusicianStat := range listenedMusicians {
		listenCountSum += userMusicianStat.ListenCount
	}

	listenReport.ListenCount = listenCountSum

	return listenReport, nil
}

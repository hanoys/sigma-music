package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
	"math"
)

type StatService struct {
	repository      ports.IStatRepository
	genreService    ports.IGenreService
	musicianService ports.IMusicianService
	logger          *zap.Logger
}

func NewStatService(repo ports.IStatRepository, genreService ports.IGenreService,
	musService ports.IMusicianService, logger *zap.Logger) *StatService {

	return &StatService{
		repository:      repo,
		genreService:    genreService,
		musicianService: musService,
		logger:          logger,
	}
}

func (ss *StatService) Add(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error {
	err := ss.repository.Add(ctx, uuid.New(), userID, trackID)
	if err != nil {
		ss.logger.Error("Failed to add track to statistics", zap.Error(err),
			zap.String("User ID", userID.String()), zap.String("Track ID", trackID.String()))

		return err
	}

	ss.logger.Info("Track successfully added to statistics",
		zap.String("User ID", userID.String()), zap.String("Track ID", trackID.String()))

	return nil
}

func (ss *StatService) fillMostListenedMusicians(ctx context.Context, listenReport *domain.ListenReport,
	listenedMusicians []domain.UserMusiciansStat) error {

	listenReport.MostListenedMusicians = make([]domain.MusicianStat, len(listenedMusicians))
	for i, userMusicianStat := range listenedMusicians {
		musician, err := ss.musicianService.GetByID(ctx, userMusicianStat.MusicianID)
		if err != nil {
			return err
		}

		listenReport.MostListenedMusicians[i] = domain.MusicianStat{
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
		ss.logger.Error("Error to create report", zap.Error(err))
		return domain.ListenReport{}, err
	}

	listenedGenres, err := ss.repository.GetListenedGenres(ctx, userID)
	if err != nil {
		ss.logger.Error("Error to create report", zap.Error(err))
		return domain.ListenReport{}, err
	}

	err = ss.fillMostListenedMusicians(ctx, &listenReport, listenedMusicians)
	if err != nil {
		ss.logger.Error("Error to create report", zap.Error(err))
		return domain.ListenReport{}, err
	}

	err = ss.fillListenedGenres(ctx, &listenReport, listenedGenres)
	if err != nil {
		ss.logger.Error("Error to create report", zap.Error(err))
		return domain.ListenReport{}, err
	}

	var listenCountSum int64
	for _, userMusicianStat := range listenedMusicians {
		listenCountSum += userMusicianStat.ListenCount
	}

	listenReport.ListenCount = listenCountSum

	ss.logger.Info("Report successfully created", zap.String("User ID", userID.String()))

	return listenReport, nil
}

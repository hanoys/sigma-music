package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type TrackService struct {
	repository   ports.ITrackRepository
	trackStorage ports.ITrackObjectStorage
	genreService ports.IGenreService
	logger       *zap.Logger
}

func NewTrackService(repo ports.ITrackRepository, storage ports.ITrackObjectStorage,
	genreService ports.IGenreService, logger *zap.Logger,
) *TrackService {
	return &TrackService{
		repository:   repo,
		trackStorage: storage,
		genreService: genreService,
		logger:       logger,
	}
}

func (ts *TrackService) Create(ctx context.Context, trackInfo ports.CreateTrackReq) (domain.Track, error) {
	trackID := uuid.New()

	url, err := ts.trackStorage.PutTrack(ctx, ports.PutTrackReq{
		TrackID:   trackID.String(),
		TrackBLOB: trackInfo.TrackBLOB,
	})
	if err != nil {
		ts.logger.Error("Failed to create track", zap.Error(err),
			zap.String("Album ID", trackInfo.AlbumID.String()), zap.String("Track name", trackInfo.Name))

		return domain.Track{}, err
	}

	track, err := ts.repository.Create(ctx, domain.Track{
		ID:      trackID,
		AlbumID: trackInfo.AlbumID,
		Name:    trackInfo.Name,
		URL:     url.String(),
	})
	if err != nil {
		ts.logger.Error("Failed to create track", zap.Error(err),
			zap.String("Track ID", track.ID.String()), zap.String("Album ID", trackInfo.AlbumID.String()),
			zap.String("Track name", trackInfo.Name), zap.String("Track URL", track.URL))

		return domain.Track{}, err
	}

	err = ts.genreService.AddForTrack(ctx, trackID, trackInfo.GenresID)
	if err != nil {
		ts.logger.Error("Failed to create track", zap.Error(err),
			zap.String("Track ID", track.ID.String()), zap.String("Album ID", trackInfo.AlbumID.String()),
			zap.String("Track name", trackInfo.Name), zap.String("Track URL", track.URL))

		return domain.Track{}, err
	}

	ts.logger.Info("Track successfully created",
		zap.String("Track ID", track.ID.String()), zap.String("Album ID", trackInfo.AlbumID.String()),
		zap.String("Track name", trackInfo.Name), zap.String("Track URL", track.URL))

	log.Println("TRACK DEBUG: created track: ", track.ID.String())

	return track, nil
}

func (ts *TrackService) GetAll(ctx context.Context) ([]domain.Track, error) {
	tracks, err := ts.repository.GetAll(ctx)
	if err != nil {
		ts.logger.Error("Failed to get all tracks", zap.Error(err))
		return nil, err
	}

	return tracks, nil
}

func (ts *TrackService) GetByID(ctx context.Context, trackID uuid.UUID) (domain.Track, error) {
	track, err := ts.repository.GetByID(ctx, trackID)
	if err != nil {
		ts.logger.Error("Failed to get track by ID", zap.Error(err),
			zap.String("Track ID", trackID.String()))

		return domain.Track{}, err
	}

	ts.logger.Info("Track successfully received by ID", zap.String("Track ID", trackID.String()))

	return track, nil
}

func (ts *TrackService) Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error) {
	trackInfo, err := ts.repository.Delete(ctx, trackID)
	if err != nil {
		ts.logger.Error("Failed to delete track", zap.Error(err), zap.String("Track ID", trackID.String()))
		return domain.Track{}, err
	}

	err = ts.trackStorage.DeleteTrack(ctx, trackID)
	if err != nil {
		ts.logger.Error("Failed to delete track", zap.Error(err), zap.String("Track ID", trackID.String()))
		return domain.Track{}, err
	}

	ts.logger.Info("Track successfully delted", zap.String("Track ID", trackID.String()))

	return trackInfo, nil
}

func (ts *TrackService) DeleteFavorite(ctx context.Context, trackID uuid.UUID, userID uuid.UUID) (domain.Track, error) {
	track, err := ts.repository.DeleteFavorite(ctx, trackID, userID)
	if err != nil {
		return domain.Track{}, nil
	}

	return track, nil
}

func (ts *TrackService) GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]domain.Track, error) {
	tracks, err := ts.repository.GetUserFavorites(ctx, userID)
	if err != nil {
		ts.logger.Error("Failed to get user favorites tracks", zap.Error(err),
			zap.String("User ID", userID.String()))

		return nil, err
	}

	ts.logger.Info("User favorites tracks successfully received", zap.String("User ID", userID.String()))

	return tracks, nil
}

func (ts *TrackService) AddToUserFavorites(ctx context.Context, trackID uuid.UUID, userID uuid.UUID) error {
	err := ts.repository.AddToUserFavorites(ctx, trackID, userID)
	if err != nil {
		ts.logger.Error("Failed to add track to user favorites", zap.Error(err),
			zap.String("Track ID", trackID.String()), zap.String("User ID", userID.String()))

		return err
	}

	ts.logger.Info("Track successfully added to user favorites", zap.String("User ID", userID.String()),
		zap.String("Track ID", trackID.String()))

	return nil
}

func (ts *TrackService) GetByAlbumID(ctx context.Context, albumID uuid.UUID) ([]domain.Track, error) {
	tracks, err := ts.repository.GetByAlbumID(ctx, albumID)
	if err != nil {
		ts.logger.Error("Failed to get tracks by album ID", zap.Error(err),
			zap.String("Album ID", albumID.String()))

		return nil, err
	}

	ts.logger.Info("Tracks successfully received by album ID", zap.String("Album ID", albumID.String()))

	return tracks, nil
}

func (ts *TrackService) GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error) {
	tracks, err := ts.repository.GetByMusicianID(ctx, musicianID)
	if err != nil {
		ts.logger.Error("Failed to get tracks by musician ID", zap.Error(err),
			zap.String("Musician ID", musicianID.String()))

		return nil, err
	}

	ts.logger.Info("Tracks successfully received by musician ID", zap.String("Musician ID", musicianID.String()))

	return tracks, nil
}

func (ts *TrackService) GetOwn(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error) {
	tracks, err := ts.repository.GetOwn(ctx, musicianID)
	if err != nil {
		ts.logger.Error("Failed to get own musician tracks", zap.Error(err),
			zap.String("Musician ID", musicianID.String()))

		return nil, err
	}

	ts.logger.Info("Own tracks successfully received by musician ID", zap.String("Musician ID", musicianID.String()))

	return tracks, nil
}

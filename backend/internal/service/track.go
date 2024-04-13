package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type TrackService struct {
	repository   ports.ITrackRepository
	trackStorage ports.ITrackObjectStorage
	GenreService ports.IGenreService
}

func NewTrackService(repo ports.ITrackRepository, storage ports.ITrackObjectStorage, genreService ports.IGenreService) *TrackService {
	return &TrackService{
		repository:   repo,
		trackStorage: storage,
		GenreService: genreService,
	}
}

func (ts *TrackService) Create(ctx context.Context, trackInfo ports.CreateTrackReq) (domain.Track, error) {
	trackID := uuid.New()

	track, err := ts.repository.Create(ctx, domain.Track{
		ID:      trackID,
		AlbumID: trackInfo.AlbumID,
		Name:    trackInfo.Name,
		URL:     trackID.String(),
	})

	if err != nil {
		return domain.Track{}, ports.ErrTrackCreate
	}

	err = ts.GenreService.AddForTrack(ctx, trackID, trackInfo.GenresID)
	if err != nil {
		return domain.Track{}, err
	}

	err = ts.trackStorage.PutTrack(ctx, ports.PutTrackReq{
		TrackID:   trackID.String(),
		TrackSize: trackInfo.TrackSize,
		TrackBLOB: trackInfo.TrackBLOB,
	})

	if err != nil {
		return domain.Track{}, ports.ErrTrackPut
	}

	return track, nil
}

func (ts *TrackService) Delete(ctx context.Context, trackID uuid.UUID) error {
	trackInfo, err := ts.repository.Delete(ctx, trackID)
	if err != nil {
		return ports.ErrTrackDelete
	}

	err = ts.trackStorage.DeleteTrack(ctx, trackID)
	if err != nil {
		_, _ = ts.repository.Create(ctx, trackInfo)
		return ports.ErrTrackDelete
	}

	return nil
}

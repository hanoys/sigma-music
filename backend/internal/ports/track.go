package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"io"
)

var (
	ErrTrackDuplicate    = errors.New("")
	ErrTrackIDNotFound   = errors.New("track with such id not found")
	ErrTrackDelete       = errors.New("can't delete track with such id")
	ErrInternalTrackRepo = errors.New("internal track repository error")
)

type ITrackRepository interface {
	Create(ctx context.Context, track domain.Track) (domain.Track, error)
	Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error)
}

type PutTrackReq struct {
	TrackID   string
	TrackSize int64
	TrackBLOB io.Reader
}

type ITrackObjectStorage interface {
	PutTrack(ctx context.Context, req PutTrackReq) error
	DeleteTrack(ctx context.Context, trackID uuid.UUID) error
}

type CreateTrackReq struct {
	AlbumID   uuid.UUID
	Name      string
	TrackBLOB io.Reader
	TrackSize int64
	GenresID  []uuid.UUID
}

type ITrackService interface {
	Create(ctx context.Context, trackInfo CreateTrackReq) (domain.Track, error)
	Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error)
}

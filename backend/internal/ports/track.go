package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"io"
	"net/url"
)

var (
	ErrTrackDuplicate    = errors.New("track duplicate error")
	ErrTrackIDNotFound   = errors.New("track with such id not found")
	ErrTrackDelete       = errors.New("can't delete track with such id")
	ErrInternalTrackRepo = errors.New("internal track repository error")
)

type ITrackRepository interface {
	Create(ctx context.Context, track domain.Track) (domain.Track, error)
	GetAll(ctx context.Context) ([]domain.Track, error)
	GetByID(ctx context.Context, trackID uuid.UUID) (domain.Track, error)
	Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error)
	GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]domain.Track, error)
	AddToUserFavorites(ctx context.Context, trackID uuid.UUID, userID uuid.UUID) error
	GetByAlbumID(ctx context.Context, albumID uuid.UUID) ([]domain.Track, error)
	GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error)
}

type PutTrackReq struct {
	TrackID   string
	TrackBLOB io.Reader
}

type ITrackObjectStorage interface {
	PutTrack(ctx context.Context, req PutTrackReq) (url.URL, error)
	DeleteTrack(ctx context.Context, trackID uuid.UUID) error
}

type CreateTrackReq struct {
	AlbumID   uuid.UUID
	Name      string
	TrackBLOB io.Reader
	GenresID  []uuid.UUID
}

type ITrackService interface {
	Create(ctx context.Context, trackInfo CreateTrackReq) (domain.Track, error)
	GetAll(ctx context.Context) ([]domain.Track, error)
	GetByID(ctx context.Context, trackID uuid.UUID) (domain.Track, error)
	Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error)
	GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]domain.Track, error)
	AddToUserFavorites(ctx context.Context, trackID uuid.UUID, userID uuid.UUID) error
	GetByAlbumID(ctx context.Context, albumID uuid.UUID) ([]domain.Track, error)
	GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error)
}

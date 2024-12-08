package ports

import (
	"context"
	"errors"
	"io"
	"net/url"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrTrackDuplicate    = errors.New("track duplicate error")
	ErrTrackIDNotFound   = errors.New("track with such id not found")
	ErrTrackDelete       = errors.New("can't delete track with such id")
	ErrInternalTrackRepo = errors.New("internal track repository error")
	ErrTrackUpdate       = errors.New("failed to update track")
)

type ITrackRepository interface {
	Create(ctx context.Context, track domain.Track) (domain.Track, error)
	Update(ctx context.Context, track domain.Track) (domain.Track, error)
	GetAll(ctx context.Context) ([]domain.Track, error)
	GetByID(ctx context.Context, trackID uuid.UUID) (domain.Track, error)
	Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error)
	GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]domain.Track, error)
	AddToUserFavorites(ctx context.Context, trackID uuid.UUID, userID uuid.UUID) error
	GetByAlbumID(ctx context.Context, albumID uuid.UUID) ([]domain.Track, error)
	GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error)
	GetOwn(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error)
}

type PutTrackReq struct {
	TrackID   string
	TrackBLOB io.Reader
}

type ITrackObjectStorage interface {
	PutTrack(ctx context.Context, req PutTrackReq) (url.URL, error)
	UploadImage(ctx context.Context, image io.Reader, id string) (url.URL, error)
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
	UploadImage(ctx context.Context, image io.Reader, id uuid.UUID, musician_id uuid.UUID) (domain.Track, error)
	GetAll(ctx context.Context) ([]domain.Track, error)
	GetByID(ctx context.Context, trackID uuid.UUID) (domain.Track, error)
	Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error)
	GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]domain.Track, error)
	AddToUserFavorites(ctx context.Context, trackID uuid.UUID, userID uuid.UUID) error
	GetByAlbumID(ctx context.Context, albumID uuid.UUID) ([]domain.Track, error)
	GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error)
	GetOwn(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error)
}

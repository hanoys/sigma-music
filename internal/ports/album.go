package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrAlbumDuplicate    = errors.New("album duplicate error")
	ErrAlbumIDNotFound   = errors.New("album with such id not found")
	ErrAlbumPublish      = errors.New("can't publish album with such id")
	ErrInternalAlbumRepo = errors.New("album repository internal error")
)

type IAlbumRepository interface {
	Create(ctx context.Context, album domain.Album, musicianID uuid.UUID) (domain.Album, error)
	GetAll(ctx context.Context) ([]domain.Album, error)
	GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error)
	GetOwn(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Album, error)
	Publish(ctx context.Context, id uuid.UUID) error
}

type CreateAlbumServiceReq struct {
	MusicianID  uuid.UUID
	Name        string
	Description string
}

type IAlbumService interface {
	Create(ctx context.Context, albumInfo CreateAlbumServiceReq) (domain.Album, error)
	GetAll(ctx context.Context) ([]domain.Album, error)
	GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error)
	GetOwn(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Album, error)
	Publish(ctx context.Context, albumID uuid.UUID) error
}

package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrCreateAlbum = errors.New("can't create the album: internal error")
	ErrPublish     = errors.New("can't publish the album: internal error")
)

type CreateAlbumRepositoryReq struct {
	ID          uuid.UUID
	Name        string
	Description string
	Published   bool
}

type IAlbumRepository interface {
	Create(ctx context.Context, album CreateAlbumRepositoryReq) (domain.Album, error)
	GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Album, error)
	Publish(ctx context.Context, id uuid.UUID) error
}

type CreateAlbumServiceReq struct {
	Name        string
	Description string
}

type IAlbumService interface {
	Create(ctx context.Context, albumInfo CreateAlbumServiceReq) (domain.Album, error)
	Publish(ctx context.Context, albumID uuid.UUID) error
}

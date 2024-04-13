package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type AlbumService struct {
	repository ports.IAlbumRepository
}

func NewAlbumService(repo ports.IAlbumRepository) *AlbumService {
	return &AlbumService{repository: repo}
}

func (as *AlbumService) Create(ctx context.Context, albumInfo ports.CreateAlbumServiceReq) (domain.Album, error) {
	createAlbum := domain.Album{
		ID:          uuid.New(),
		Name:        albumInfo.Name,
		Description: albumInfo.Description,
		Published:   false,
		ReleaseDate: null.Time{},
	}

	album, err := as.repository.Create(ctx, createAlbum)
	if err != nil {
		return domain.Album{}, ports.ErrCreateAlbum
	}

	return album, nil
}

func (as *AlbumService) Publish(ctx context.Context, albumID uuid.UUID) error {
	err := as.repository.Publish(ctx, albumID)
	if err != nil {
		return ports.ErrPublish
	}

	return nil
}

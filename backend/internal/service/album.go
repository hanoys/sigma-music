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

	return as.repository.Create(ctx, createAlbum)
}

func (as *AlbumService) Publish(ctx context.Context, albumID uuid.UUID) error {
	return as.repository.Publish(ctx, albumID)
}

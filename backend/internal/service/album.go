package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type AlbumService struct {
	repository ports.IAlbumRepository
	logger     *zap.Logger
}

func NewAlbumService(repo ports.IAlbumRepository, logger *zap.Logger) *AlbumService {
	return &AlbumService{repository: repo,
		logger: logger}
}

func (as *AlbumService) Create(ctx context.Context, albumInfo ports.CreateAlbumServiceReq) (domain.Album, error) {
	createAlbum := domain.Album{
		ID:          uuid.New(),
		Name:        albumInfo.Name,
		Description: albumInfo.Description,
		Published:   false,
		ReleaseDate: null.Time{},
	}

	album, err := as.repository.Create(ctx, createAlbum, albumInfo.MusicianID)
	if err != nil {
		as.logger.Error("Failed to create new album", zap.Error(err))
		return domain.Album{}, err
	}

	as.logger.Info("Album successfully created", zap.String("Album ID", createAlbum.ID.String()),
		zap.String("Album Name", createAlbum.Name))

	return album, nil
}

func (as *AlbumService) GetAll(ctx context.Context) ([]domain.Album, error) {
	albums, err := as.repository.GetAll(ctx)
	if err != nil {
		as.logger.Error("Failed to get all albums", zap.Error(err))
		return nil, err
	}

	return albums, nil
}

func (as *AlbumService) GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error) {
	albums, err := as.repository.GetByMusicianID(ctx, musicianID)
	if err != nil {
		as.logger.Error("Failed to get musician albums", zap.String("Musician ID", musicianID.String()))
		return nil, err
	}

	as.logger.Info("Successfully received musician albums", zap.String("Musician ID", musicianID.String()))

	return albums, nil
}

func (as *AlbumService) GetOwn(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error) {
	albums, err := as.repository.GetOwn(ctx, musicianID)
	if err != nil {
		as.logger.Error("Failed to get musician own albums", zap.Error(err))
	}

	as.logger.Info("Successfully received own musician albums", zap.String("Musician ID", musicianID.String()))

	return albums, nil
}

func (as *AlbumService) GetByID(ctx context.Context, id uuid.UUID) (domain.Album, error) {
	album, err := as.repository.GetByID(ctx, id)
	if err != nil {
		as.logger.Error("Failed to get album by ID", zap.Error(err), zap.String("Album ID", id.String()))
	}

	as.logger.Info("Successfully received album by ID", zap.String("Album ID", id.String()))

	return album, nil
}

func (as *AlbumService) Publish(ctx context.Context, albumID uuid.UUID) error {
	err := as.repository.Publish(ctx, albumID)
	if err != nil {
		as.logger.Error("Failed to publish album", zap.String("Album ID", albumID.String()))
	}

	as.logger.Info("Successfully published album", zap.String("Album ID", albumID.String()))

	return nil
}

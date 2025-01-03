package service

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type MusicianService struct {
	repository   ports.IMusicianRepository
	imageStorage ports.IMusicianImageStorage
	hash         ports.IHashPasswordProvider
	logger       *zap.Logger
}

func NewMusicianService(repo ports.IMusicianRepository, imageStorage ports.IMusicianImageStorage, hash ports.IHashPasswordProvider,
	logger *zap.Logger,
) *MusicianService {
	return &MusicianService{
		repository:   repo,
		imageStorage: imageStorage,
		hash:         hash,
		logger:       logger,
	}
}

func (ms *MusicianService) Register(ctx context.Context, musician ports.MusicianServiceCreateRequest) (domain.Musician, error) {
	_, err := ms.repository.GetByName(ctx, musician.Name)
	if err == nil {
		ms.logger.Error("Failed to register musician", zap.Error(err), zap.String("Musician Name", musician.Name))
		return domain.Musician{}, ports.ErrMusicianWithSuchNameAlreadyExists
	}

	_, err = ms.repository.GetByEmail(ctx, musician.Email)
	if err == nil {
		ms.logger.Error("Failed to register musician", zap.Error(err), zap.String("Musician Email", musician.Email))
		return domain.Musician{}, ports.ErrMusicianWithSuchEmailAlreadyExists
	}

	saltedPassword := ms.hash.EncodePassword(musician.Password)

	createMusician := domain.Musician{
		ID:          uuid.New(),
		Name:        musician.Name,
		Email:       musician.Email,
		Password:    saltedPassword.HashPassword,
		Salt:        saltedPassword.Salt,
		Country:     musician.Country,
		Description: musician.Description,
	}

	mus, err := ms.repository.Create(ctx, createMusician)
	if err != nil {
		ms.logger.Error("Failed to register musician", zap.Error(err))
		return domain.Musician{}, err
	}

	ms.logger.Info("Musician successfully registered", zap.String("Musician ID", createMusician.ID.String()))

	return mus, nil
}

func (ms *MusicianService) UploadImage(ctx context.Context, image io.Reader, id uuid.UUID) (domain.Musician, error) {
	url, err := ms.imageStorage.UploadImage(ctx, image, id.String())
	if err != nil {
		return domain.Musician{}, err
	}

	musician, err := ms.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Musician{}, err
	}

	musician.ImageURL = null.StringFrom(url.String())
	updatedMusician, err := ms.repository.Update(ctx, musician)
	if err != nil {
		return domain.Musician{}, err
	}

	return updatedMusician, nil
}

func (ms *MusicianService) GetAll(ctx context.Context) ([]domain.Musician, error) {
	musicians, err := ms.repository.GetAll(ctx)
	if err != nil {
		ms.logger.Error("Failed to get all musicians", zap.Error(err))
		return nil, err
	}

	return musicians, nil
}

func (ms *MusicianService) GetByID(ctx context.Context, musicianID uuid.UUID) (domain.Musician, error) {
	mus, err := ms.repository.GetByID(ctx, musicianID)
	if err != nil {
		ms.logger.Error("Failed to get musician by ID", zap.Error(err),
			zap.String("Musician ID", musicianID.String()))
		return domain.Musician{}, err
	}

	ms.logger.Info("Musician successfully received by id", zap.String("Musician ID", musicianID.String()))

	return mus, nil
}

func (ms *MusicianService) GetByName(ctx context.Context, name string) (domain.Musician, error) {
	mus, err := ms.repository.GetByName(ctx, name)
	if err != nil {
		ms.logger.Error("Failed to get musician by name", zap.Error(err), zap.String("Musician Name", name))
		return domain.Musician{}, err
	}

	ms.logger.Info("Musician successfully received by name", zap.String("Musician Name", name))

	return mus, nil
}

func (ms *MusicianService) GetByEmail(ctx context.Context, email string) (domain.Musician, error) {
	mus, err := ms.repository.GetByEmail(ctx, email)
	if err != nil {
		ms.logger.Error("Failed to get musician by email", zap.Error(err), zap.String("Musician Email", email))
		return domain.Musician{}, err
	}

	ms.logger.Info("Musician successfully received by email", zap.String("Musician Email", email))

	return mus, nil
}

func (ms *MusicianService) GetByAlbumID(ctx context.Context, albumID uuid.UUID) (domain.Musician, error) {
	mus, err := ms.repository.GetByAlbumID(ctx, albumID)
	if err != nil {
		ms.logger.Error("Failed to get musician by album ID", zap.Error(err),
			zap.String("Album ID", albumID.String()))

		return domain.Musician{}, err
	}

	ms.logger.Info("Musician successfully received by album ID", zap.String("Album ID", albumID.String()))

	return mus, nil
}

func (ms *MusicianService) GetByTrackID(ctx context.Context, trackID uuid.UUID) (domain.Musician, error) {
	mus, err := ms.repository.GetByTrackID(ctx, trackID)
	if err != nil {
		ms.logger.Error("Failed to get musician by track ID", zap.Error(err),
			zap.String("Album ID", trackID.String()))

		return domain.Musician{}, err
	}

	ms.logger.Info("Musician successfully received by track ID", zap.String("Track ID", trackID.String()))

	return mus, nil
}

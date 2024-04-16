package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type MusicianService struct {
	repository ports.IMusicianRepository
}

func NewMusicianService(repo ports.IMusicianRepository) *MusicianService {
	return &MusicianService{repository: repo}
}

func (ms *MusicianService) Register(ctx context.Context, musician ports.MusicianServiceCreateRequest) (domain.Musician, error) {
	_, err := ms.repository.GetByName(ctx, musician.Name)
	if err == nil {
		return domain.Musician{}, ports.ErrMusicianWithSuchNameAlreadyExists
	}

	_, err = ms.repository.GetByEmail(ctx, musician.Email)
	if err == nil {
		return domain.Musician{}, ports.ErrMusicianWithSuchEmailAlreadyExists
	}

	createMusician := domain.Musician{
		ID:          uuid.New(),
		Name:        musician.Name,
		Email:       musician.Email,
		Password:    musician.Password,
		Country:     musician.Country,
		Description: musician.Description,
	}

	return ms.repository.Create(ctx, createMusician)
}

func (ms *MusicianService) GetByID(ctx context.Context, musicianID uuid.UUID) (domain.Musician, error) {
	return ms.repository.GetByID(ctx, musicianID)
}

func (ms *MusicianService) GetByName(ctx context.Context, name string) (domain.Musician, error) {
	return ms.repository.GetByName(ctx, name)
}

func (ms *MusicianService) GetByEmail(ctx context.Context, email string) (domain.Musician, error) {
	return ms.repository.GetByEmail(ctx, email)
}

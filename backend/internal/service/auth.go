package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type AuthorizationService struct {
	userRepository     ports.IUserRepository
	musicianRepository ports.IMusicianRepository
	tokenProvider      ports.ITokenProvider
}

func NewAuthorizationService(userRepo ports.IUserRepository, musicianRepo ports.IMusicianRepository, tokenProvider ports.ITokenProvider) *AuthorizationService {
	return &AuthorizationService{
		userRepository:     userRepo,
		musicianRepository: musicianRepo,
		tokenProvider:      tokenProvider,
	}
}

func (a *AuthorizationService) authUser(ctx context.Context, name string, password string) (domain.User, error) {
	user, err := a.userRepository.GetByName(ctx, name)

	if errors.Is(err, ports.ErrUserNameNotFound) {
		return domain.User{}, ports.ErrIncorrectName
	} else if err != nil {
		return domain.User{}, ports.ErrInternalAuthRepo
	}

	if user.Password != password {
		return domain.User{}, ports.ErrIncorrectPassword
	}

	return user, nil
}

func (a *AuthorizationService) authMusician(ctx context.Context, name string, password string) (domain.Musician, error) {
	musician, err := a.musicianRepository.GetByName(ctx, name)

	if errors.Is(err, ports.ErrMusicianNameNotFound) {
		return domain.Musician{}, ports.ErrIncorrectName
	} else if err != nil {
		return domain.Musician{}, ports.ErrInternalAuthRepo
	}

	if musician.Password != password {
		return domain.Musician{}, ports.ErrIncorrectPassword
	}

	return musician, nil
}

func (a *AuthorizationService) LogIn(ctx context.Context, cred ports.LogInCredentials) (domain.TokenPair, error) {
	var id uuid.UUID

	switch cred.Role {
	case domain.UserRole:
		user, err := a.authUser(ctx, cred.Name, cred.Password)
		if err != nil {
			return domain.TokenPair{}, err
		}

		id = user.ID
	case domain.MusicianRole:
		musician, err := a.authMusician(ctx, cred.Name, cred.Password)
		if err != nil {
			return domain.TokenPair{}, err
		}

		id = musician.ID
	default:
		return domain.TokenPair{}, ports.ErrUnexpectedRole
	}

	payload := domain.Payload{
		UserID: id,
		Role:   cred.Role,
	}

	tokens, err := a.tokenProvider.NewSession(ctx, payload)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return tokens, nil
}

func (a *AuthorizationService) LogOut(ctx context.Context, tokenString string) error {
	err := a.tokenProvider.CloseSession(ctx, tokenString)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthorizationService) RefreshToken(ctx context.Context, refreshTokenString string) (domain.TokenPair, error) {
	payload, err := a.tokenProvider.RefreshSession(ctx, refreshTokenString)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return payload, nil
}

func (a *AuthorizationService) VerifyToken(ctx context.Context, tokenString string) (domain.Payload, error) {
	payload, err := a.tokenProvider.VerifyToken(ctx, tokenString)
	if err != nil {
		return domain.Payload{}, err
	}

	return payload, err
}

package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/auth"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type AuthorizationService struct {
	userRepository     ports.IUserRepository
	musicianRepository ports.IMusicianRepository
	tokenProvider      *auth.Provider
}

func NewAuthorizationService(userRepo ports.IUserRepository, musicianRepo ports.IMusicianRepository, tokenProvider *auth.Provider) *AuthorizationService {
	return &AuthorizationService{
		userRepository:     userRepo,
		musicianRepository: musicianRepo,
		tokenProvider:      tokenProvider,
	}
}

func (a *AuthorizationService) authUser(ctx context.Context, name string, password string) (domain.User, error) {
	user, err := a.userRepository.GetByName(ctx, name)
	if err != nil {
		return domain.User{}, ports.ErrIncorrectName
	}

	if user.Password != password {
		return domain.User{}, ports.ErrIncorrectPassword
	}

	return user, nil
}

func (a *AuthorizationService) authMusician(ctx context.Context, name string, password string) (domain.Musician, error) {
	musician, err := a.musicianRepository.GetByName(ctx, name)
	if err != nil {
		return domain.Musician{}, ports.ErrIncorrectName
	}

	if musician.Password != password {
		return domain.Musician{}, ports.ErrIncorrectPassword
	}

	return musician, nil
}

func (a *AuthorizationService) LogIn(ctx context.Context, cred ports.LogInCredentials) (*auth.TokenPair, error) {
	var id uuid.UUID

	switch cred.Role {
	case domain.UserRole:
		user, err := a.authUser(ctx, cred.Name, cred.Password)
		if err != nil {
			return nil, err
		}

		id = user.ID
	case domain.MusicianRole:
		musician, err := a.authMusician(ctx, cred.Name, cred.Password)
		if err != nil {
			return nil, err
		}

		id = musician.ID
	default:
		return nil, ports.ErrUnexpectedRole
	}

	tokenPayload, err := a.tokenProvider.NewPayload(id, cred.Role)
	if err != nil {
		return nil, err
	}

	session, err := a.tokenProvider.NewSession(ctx, tokenPayload)
	if err != nil {
		return nil, err
	}

	return session.Tokens, nil
}

func (a *AuthorizationService) LogOut(ctx context.Context, tokenString string) error {
	err := a.tokenProvider.CloseSession(ctx, tokenString)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthorizationService) RefreshToken(ctx context.Context, refreshTokenString string) (*auth.TokenPair, error) {
	session, err := a.tokenProvider.RefreshSession(ctx, refreshTokenString)
	if err != nil {
		return nil, err
	}

	return session.Tokens, nil
}

func (a *AuthorizationService) VerifyToken(ctx context.Context, tokenString string) (*auth.Payload, error) {
	payload, err := a.tokenProvider.VerifyToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return payload, err
}

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
	hash               ports.IHashPasswordProvider
}

func NewAuthorizationService(userRepo ports.IUserRepository, musicianRepo ports.IMusicianRepository,
	tokenProvider ports.ITokenProvider, hash ports.IHashPasswordProvider) *AuthorizationService {
	return &AuthorizationService{
		userRepository:     userRepo,
		musicianRepository: musicianRepo,
		tokenProvider:      tokenProvider,
		hash:               hash,
	}
}

func (a *AuthorizationService) authUser(ctx context.Context, name string, password string) (domain.User, error) {
	user, err := a.userRepository.GetByName(ctx, name)

	if errors.Is(err, ports.ErrUserNameNotFound) {
		return domain.User{}, ports.ErrIncorrectName
	} else if err != nil {
		return domain.User{}, ports.ErrInternalAuthRepo
	}

	saltedPassword := domain.SaltedPassword{
		HashPassword: user.Password,
		Salt:         user.Salt,
	}

	if !a.hash.ComparePasswordWithHash(password, saltedPassword) {
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

	saltedPassword := domain.SaltedPassword{
		HashPassword: musician.Password,
		Salt:         musician.Salt,
	}

	if !a.hash.ComparePasswordWithHash(password, saltedPassword) {
		return domain.Musician{}, ports.ErrIncorrectPassword
	}

	return musician, nil
}

func (a *AuthorizationService) LogIn(ctx context.Context, cred ports.LogInCredentials) (domain.TokenPair, error) {
	var id uuid.UUID
	var role int

	user, err := a.authUser(ctx, cred.Name, cred.Password)
	if err == nil {
		id = user.ID
		role = domain.UserRole
	} else {
		musician, err := a.authMusician(ctx, cred.Name, cred.Password)
		if err == nil {
			id = musician.ID
			role = domain.MusicianRole
		} else {
			return domain.TokenPair{}, err
		}
	}

	payload := domain.Payload{
		UserID: id,
		Role:   role,
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

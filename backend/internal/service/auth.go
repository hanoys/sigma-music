package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"go.uber.org/zap"
)

type AuthorizationService struct {
	userRepository     ports.IUserRepository
	musicianRepository ports.IMusicianRepository
	tokenProvider      ports.ITokenProvider
	hash               ports.IHashPasswordProvider
	logger             *zap.Logger
}

func NewAuthorizationService(userRepo ports.IUserRepository, musicianRepo ports.IMusicianRepository,
	tokenProvider ports.ITokenProvider, hash ports.IHashPasswordProvider, logger *zap.Logger) *AuthorizationService {
	return &AuthorizationService{
		userRepository:     userRepo,
		musicianRepository: musicianRepo,
		tokenProvider:      tokenProvider,
		hash:               hash,
		logger:             logger,
	}
}

func (a *AuthorizationService) authUser(ctx context.Context, name string, password string) (domain.User, error) {
	user, err := a.userRepository.GetByName(ctx, name)

	if errors.Is(err, ports.ErrUserNameNotFound) {
		return domain.User{}, ports.ErrIncorrectName
	} else if err != nil {
		return domain.User{}, util.WrapError(ports.ErrInternalAuthRepo, err)
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
		return domain.Musician{}, util.WrapError(ports.ErrInternalAuthRepo, err)
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
			a.logger.Error("Failed to authorize user", zap.Error(err), zap.String("User Name", cred.Name))
			return domain.TokenPair{}, err
		}
	}

	payload := domain.Payload{
		UserID: id,
		Role:   role,
	}

	var stringRole string
	if payload.Role == domain.UserRole {
		stringRole = "User"
	} else {
		stringRole = "Musician"
	}

	tokens, err := a.tokenProvider.NewSession(ctx, payload)
	if err != nil {
		a.logger.Error("Failed to create new session for user", zap.Error(err),
			zap.String("User ID", payload.UserID.String()), zap.String("User Role", stringRole))
		return domain.TokenPair{}, err
	}

	a.logger.Info("User successfully authorized",
		zap.String("User ID", payload.UserID.String()), zap.String("User Role", stringRole))

	return tokens, nil
}

func (a *AuthorizationService) LogOut(ctx context.Context, tokenString string) error {
	err := a.tokenProvider.CloseSession(ctx, tokenString)
	if err != nil {
		payload, errVerify := a.tokenProvider.VerifyToken(ctx, tokenString)
		if errVerify != nil {
			a.logger.Error("Failed to close session for user", zap.Error(errVerify))
		} else {
			var stringRole string
			if payload.Role == domain.UserRole {
				stringRole = "User"
			} else {
				stringRole = "Musician"
			}

			a.logger.Error("Failed to close user session", zap.Error(err),
				zap.String("User ID", payload.UserID.String()), zap.String("User Role", stringRole))
		}

		return err
	}

	return nil
}

func (a *AuthorizationService) RefreshToken(ctx context.Context, refreshTokenString string) (domain.TokenPair, error) {
	tokenPair, err := a.tokenProvider.RefreshSession(ctx, refreshTokenString)
	if err != nil {
		a.logger.Error("Failed to refresh token for user", zap.Error(err))
		return domain.TokenPair{}, err
	}

	return tokenPair, nil
}

func (a *AuthorizationService) VerifyToken(ctx context.Context, tokenString string) (domain.Payload, error) {
	payload, err := a.tokenProvider.VerifyToken(ctx, tokenString)
	var stringRole string
	if payload.Role == domain.UserRole {
		stringRole = "User"
	} else {
		stringRole = "Musician"
	}

	if err != nil {
		a.logger.Error("Failed to verify token for user", zap.Error(err),
			zap.String("User ID", payload.UserID.String()), zap.String("User Role", stringRole))
		return domain.Payload{}, err
	}

	a.logger.Info("Token successfully verified for user",
		zap.String("User ID", payload.UserID.String()), zap.String("User Role", stringRole))

	return payload, err
}

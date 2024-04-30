package ports

import (
	"context"
	"errors"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrTokenProviderInvalidToken = errors.New("invalid token")
	ErrTokenProviderExpiredToken = errors.New("token expired")
	ErrTokenProviderParsingToken = errors.New("can't parse token")
	ErrTokenProviderSignToken    = errors.New("can't sign token")
	ErrInternalTokenProvider     = errors.New("internal provider error ")
)

type ITokenProvider interface {
	NewSession(ctx context.Context, payload domain.Payload) (domain.TokenPair, error)
	CloseSession(ctx context.Context, refreshTokenString string) error
	RefreshSession(ctx context.Context, refreshTokenString string) (domain.TokenPair, error)
	VerifyToken(ctx context.Context, accessTokenString string) (domain.Payload, error)
}

type LogInCredentials struct {
	Name     string
	Password string
	Role     int
}

var (
	ErrIncorrectName     = errors.New("authentication error: incorrect name")
	ErrIncorrectPassword = errors.New("authentication error: incorrect password")
	ErrUnexpectedRole    = errors.New("authentication error: role doesn't exists")
	ErrInternalAuthRepo  = errors.New("authentication error: internal repository error")
)

type IAuthorizationService interface {
	LogIn(ctx context.Context, cred LogInCredentials) (domain.TokenPair, error)
	LogOut(ctx context.Context, accessTokenString string) error
	RefreshToken(ctx context.Context, refreshTokenString string) (domain.TokenPair, error)
	VerifyToken(ctx context.Context, accessTokenString string) (domain.Payload, error)
}

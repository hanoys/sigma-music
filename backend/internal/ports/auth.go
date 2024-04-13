package ports

import (
	"context"
	"errors"
	"github.com/hanoys/sigma-music/internal/auth"
)

var (
	ErrIncorrectName     = errors.New("authentication error: incorrect name")
	ErrIncorrectPassword = errors.New("authentication error: incorrect password")
	ErrUnexpectedRole    = errors.New("authentication error: role doesn't exists")
)

type LogInCredentials struct {
	Name     string
	Password string
	Role     int
}

type IAuthorizationService interface {
	LogIn(ctx context.Context, cred LogInCredentials) (*auth.TokenPair, error)
	LogOut(ctx context.Context, tokenString string) error
	RefreshToken(ctx context.Context, refreshTokenString string) (*auth.TokenPair, error)
	VerifyToken(ctx context.Context, tokenString string) (*auth.Payload, error)
}

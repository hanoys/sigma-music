package auth

import (
	"context"
	"github.com/hanoys/sigma-music/internal/adapters/auth/ports"
	"github.com/hanoys/sigma-music/internal/domain"
	serviceports "github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/utill"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenSession struct {
	Tokens         *TokenPair
	ExpirationTime time.Time
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Payload struct {
	UserID uuid.UUID
	Role   int
}

type JWTClaims struct {
	domain.Payload
	jwt.RegisteredClaims
}

type ProviderConfig struct {
	AccessTokenExpTime  int64
	RefreshTokenExpTime int64
	SecretKey           string
}

func NewProviderConfig(accessTime int64, refreshTime int64, secret string) *ProviderConfig {
	return &ProviderConfig{AccessTokenExpTime: accessTime,
		RefreshTokenExpTime: refreshTime,
		SecretKey:           secret}
}

type Provider struct {
	tokenStorage ports.ITokenStorage
	cfg          *ProviderConfig
}

func NewProvider(tokenStorage ports.ITokenStorage, cfg *ProviderConfig) *Provider {
	return &Provider{tokenStorage: tokenStorage,
		cfg: cfg}
}

func (p *Provider) newTokenWithExpiration(ctx context.Context, payload domain.Payload, exp time.Time) (string, error) {
	claims := &JWTClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(p.cfg.SecretKey))
	if err != nil {
		return "", serviceports.ErrTokenProviderSignToken
	}

	return tokenString, nil
}

func (p *Provider) NewPayload(userID uuid.UUID, role int) (*Payload, error) {
	return &Payload{
		UserID: userID,
		Role:   role,
	}, nil
}

func (p *Provider) NewSession(ctx context.Context, payload domain.Payload) (domain.TokenPair, error) {
	accessExpTime := time.Now().Add(time.Minute * time.Duration(p.cfg.AccessTokenExpTime))
	refreshExpTime := time.Now().Add(time.Minute * time.Duration(p.cfg.RefreshTokenExpTime))

	accessTokenString, err := p.newTokenWithExpiration(ctx, payload, accessExpTime)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refreshTokenString, err := p.newTokenWithExpiration(ctx, payload, refreshExpTime)
	if err != nil {
		return domain.TokenPair{}, err
	}

	tokenPair := domain.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	err = p.tokenStorage.Set(ctx, tokenPair.RefreshToken,
		payload, refreshExpTime.Sub(time.Now()))
	if err != nil {
		return domain.TokenPair{}, utill.WrapError(serviceports.ErrInternalTokenProvider, err)
	}

	return tokenPair, nil
}

func (p *Provider) RefreshSession(ctx context.Context, refreshTokenString string) (domain.TokenPair, error) {
	refreshClaims, err := p.parseToken(refreshTokenString)
	if err != nil {
		return domain.TokenPair{}, err
	}

	err = p.tokenStorage.Del(ctx, refreshTokenString)
	if err != nil {
		return domain.TokenPair{}, utill.WrapError(serviceports.ErrInternalTokenProvider, err)
	}

	payload := domain.Payload{
		UserID: refreshClaims.UserID,
		Role:   refreshClaims.Role,
	}

	return p.NewSession(ctx, payload)
}

func (p *Provider) CloseSession(ctx context.Context, refreshTokenString string) error {
	_, err := p.parseToken(refreshTokenString)
	if err != nil {
		return err
	}

	err = p.tokenStorage.Del(ctx, refreshTokenString)
	if err != nil {
		return utill.WrapError(serviceports.ErrInternalTokenProvider, err)
	}

	return nil
}

func (p *Provider) VerifyToken(ctx context.Context, tokenString string) (domain.Payload, error) {
	claims, err := p.parseToken(tokenString)
	if err != nil {
		return domain.Payload{}, err
	}

	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return domain.Payload{}, serviceports.ErrTokenProviderExpiredToken
	}

	return claims.Payload, nil
}

func (p *Provider) parseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(p.cfg.SecretKey), nil
		})

	if err != nil {
		return nil, serviceports.ErrTokenProviderParsingToken
	}

	if !token.Valid {
		return nil, serviceports.ErrTokenProviderInvalidToken
	}

	claims := token.Claims.(*JWTClaims)

	return claims, nil
}

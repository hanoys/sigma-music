package dto

import (
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type LoginDTO struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (l *LoginDTO) ToServiceRequest() ports.LogInCredentials {
	return ports.LogInCredentials{
		Name:     l.Name,
		Password: l.Password,
	}
}

type LoginResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func LoginResponseFromTokenPair(pair domain.TokenPair) LoginResponseDTO {
	return LoginResponseDTO{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}
}

type LogoutDTO struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type RefreshDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

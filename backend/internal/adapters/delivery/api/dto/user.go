package dto

import (
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type RegisterUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Country  string `json:"country"`
}

func (r *RegisterUserDTO) ToServiceRequest() ports.UserServiceCreateRequest {
	return ports.UserServiceCreateRequest{
		Name:     r.Name,
		Email:    r.Email,
		Phone:    r.Phone,
		Password: r.Password,
		Country:  r.Country,
	}
}

type LoginUserDTO struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (l *LoginUserDTO) ToServiceRequest() ports.LogInCredentials {
	return ports.LogInCredentials{
		Name:     l.Name,
		Password: l.Password,
	}
}

type LoginUserResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func LoginUserResponseFromTokenPair(pair domain.TokenPair) LoginUserResponseDTO {
	return LoginUserResponseDTO{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}
}

type LogoutUserDTO struct {
	AccessToken string `json:"access_token"`
}

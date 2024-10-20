package dto

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type UserDTO struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Phone   string    `json:"phone"`
	Country string    `json:"country"`
}

func UserFromDomain(user domain.User) UserDTO {
	return UserDTO{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Country: user.Country,
	}
}

type RegisterUserDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Country  string `json:"country" binding:"required"`
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


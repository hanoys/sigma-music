package dto

import (
	"fmt"
	"github.com/hanoys/sigma-music/internal/domain"
)

type UserDTO struct {
	ID      string
	Name    string
	Email   string
	Phone   string
	Country string
}

func NewUserDTO(user domain.User) UserDTO {
	return UserDTO{
		ID:      user.ID.String(),
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Country: user.Country,
	}
}

func (u UserDTO) Print() {
	fmt.Println("ID:", u.ID)
	fmt.Println("Name:", u.Name)
	fmt.Println("Email:", u.Email)
	fmt.Println("Phone:", u.Phone)
	fmt.Println("Country:", u.Country)
}

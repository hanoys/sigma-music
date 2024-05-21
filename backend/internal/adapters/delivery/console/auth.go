package console

import (
	"context"
	"fmt"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/console/dto"
	"github.com/hanoys/sigma-music/internal/ports"
)

func (h *Handler) LogIn(c *Console) {
	var logInDTO dto.LogInDTO
	dto.InputLogInDTO(&logInDTO)

	tokenPair, err := h.authService.LogIn(context.Background(), ports.LogInCredentials{
		Name:     logInDTO.Name,
		Password: logInDTO.Password,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	payload, _ := h.authService.VerifyToken(context.Background(), tokenPair.AccessToken)

	c.UserID = payload.UserID
	c.UserRole = payload.Role

	fmt.Println("Your ID:", c.UserID)
}

func (h *Handler) LogOut(c *Console) {
	c.UserRole = -1
}

func (h *Handler) SignUpUser(c *Console) {
	var signUpDTO dto.SignUpUserDTO
	err := dto.InputSignUpUserDTO(&signUpDTO)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = h.userService.Register(context.Background(), ports.UserServiceCreateRequest{
		Name:     signUpDTO.Name,
		Email:    signUpDTO.Email,
		Phone:    signUpDTO.Phone,
		Password: signUpDTO.Password,
		Country:  signUpDTO.County,
	})

	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) SignUpMusician(c *Console) {
	var signUpDTO dto.SignUpMusicianDTO
	err := dto.InputSignUpMusicianDTO(&signUpDTO)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = h.musicianService.Register(context.Background(), ports.MusicianServiceCreateRequest{
		Name:        signUpDTO.Name,
		Email:       signUpDTO.Email,
		Password:    signUpDTO.Password,
		Country:     signUpDTO.County,
		Description: signUpDTO.Description,
	})

	if err != nil {
		fmt.Println(err)
	}
}

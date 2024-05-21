package console

import (
	"context"
	"fmt"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/console/dto"
)

func (h *Handler) GetAllUsers(c *Console) {
	err := h.verifyAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	users, err := h.userService.GetAll(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, user := range users {
		dto.NewUserDTO(user).Print()
	}
}

func (h *Handler) GetByIdUsers(c *Console) {
	err := h.verifyAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	id, err := readID()
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := h.userService.GetById(context.Background(), id)
	if err != nil {
		fmt.Println(err)
		return
	}

	dto.NewUserDTO(user).Print()
}

func (h *Handler) GetByNameUser(c *Console) {
	err := h.verifyAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	var name string
	fmt.Scanf("%s", &name)

	user, err := h.userService.GetByName(context.Background(), name)
	if err != nil {
		fmt.Println(err)
		return
	}

	dto.NewUserDTO(user).Print()
}

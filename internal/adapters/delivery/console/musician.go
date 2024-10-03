package console

import (
	"context"
	"fmt"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/console/dto"
)

func (h *Handler) GetAllMusicians(c *Console) {
	err := h.verifyAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	users, err := h.musicianService.GetAll(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, user := range users {
		dto.NewMusicianDTO(user).Print()
        fmt.Println("-----------------------")
	}
}

func (h *Handler) GetByIdMusician(c *Console) {
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

	user, err := h.musicianService.GetByID(context.Background(), id)
	if err != nil {
		fmt.Println(err)
		return
	}

	dto.NewMusicianDTO(user).Print()
}

func (h *Handler) GetByNameMusician(c *Console) {
	err := h.verifyAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	var name string
	fmt.Scanf("%s", &name)

	user, err := h.musicianService.GetByName(context.Background(), name)
	if err != nil {
		fmt.Println(err)
		return
	}

	dto.NewMusicianDTO(user).Print()
}

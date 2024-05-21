package console

import (
	"context"
	"fmt"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/console/dto"
)

func (h *Handler) GetAllGenre(c *Console) {
	err := h.verifyAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	genres, err := h.genreService.GetAll(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, genre := range genres {
		dto.NewGenreDTO(genre).Print()
	}
}

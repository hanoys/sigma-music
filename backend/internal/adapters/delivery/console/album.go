package console

import (
	"context"
	"fmt"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/console/dto"
	"github.com/hanoys/sigma-music/internal/ports"
)

func (h *Handler) CreateAlbum(c *Console) {
	err := h.verifyMusicianAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	var createAlbumDTO dto.CreateAlbumDTO
	dto.InputCreateAlbumDTO(&createAlbumDTO)

	_, err = h.albumService.Create(context.Background(), ports.CreateAlbumServiceReq{
		Name:        createAlbumDTO.Name,
		Description: createAlbumDTO.Description,
		MusicianID:  c.UserID,
	})

	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) GetAllAlbums(c *Console) {
	err := h.verifyAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	albums, err := h.albumService.GetAll(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, album := range albums {
		dto.NewAlbumDTO(album).Print()
	}
}

func (h *Handler) GetByIDAlbum(c *Console) {
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

	album, err := h.albumService.GetByID(context.Background(), id)

	if err != nil {
		fmt.Println(err)
		return
	}

	dto.NewAlbumDTO(album).Print()
}

func (h *Handler) GetByMusicianIDAlbum(c *Console) {
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

	albums, err := h.albumService.GetByMusicianID(context.Background(), id)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, album := range albums {
		dto.NewAlbumDTO(album).Print()
	}
}

func (h *Handler) GetOwn(c *Console) {
	err := h.verifyMusicianAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	albums, err := h.albumService.GetByMusicianID(context.Background(), c.UserID)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, album := range albums {
		dto.NewAlbumDTO(album).Print()
	}
}

func (h *Handler) PublishAlbum(c *Console) {
	err := h.verifyMusicianAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	id, err := readID()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = h.albumService.Publish(context.Background(), id)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Album published")
	}
}

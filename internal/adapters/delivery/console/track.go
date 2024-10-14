package console

import (
	"context"
	"fmt"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/console/dto"
	"github.com/hanoys/sigma-music/internal/ports"
)

func (h *Handler) CreateTrack(c *Console) {
	err := h.verifyMusicianAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	var createDTO dto.CreateTrackDTO
	err = dto.InputCreateTrackDTO(&createDTO)
	if err != nil {
		fmt.Println(err)
		return
	}

	mus, err := h.musicianService.GetByAlbumID(context.Background(), createDTO.AlbumID)
	if err != nil {
		fmt.Println("owner not found: ", err)
		return
	}

	if mus.ID != c.UserID {
		fmt.Println("Error: Forbidden")
		return
	}

	_, err = h.trackService.Create(context.Background(), ports.CreateTrackReq{
		AlbumID:   createDTO.AlbumID,
		Name:      createDTO.Name,
		TrackBLOB: createDTO.Track,
		GenresID:  createDTO.GenresID,
	})

	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) GetAllTrack(c *Console) {
	err := h.verifyAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	tracks, err := h.trackService.GetAll(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, track := range tracks {
		dto.NewTrackDTO(track).Print()
		fmt.Println("-----------------------")
	}
}

func (h *Handler) GetByIDTrack(c *Console) {
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

	track, err := h.trackService.GetByID(context.Background(), id)
	if err != nil {
		fmt.Println(err)
		return
	}

	dto.NewTrackDTO(track).Print()
}

func (h *Handler) DeleteTrack(c *Console) {
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

	mus, err := h.musicianService.GetByTrackID(context.Background(), id)
	if err != nil {
		fmt.Println(err)
		return
	}

	if mus.ID != c.UserID {
		fmt.Println("Error: Forbidden")
		return
	}

	_, err = h.trackService.Delete(context.Background(), id)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *Handler) GetUserFavoritesTrack(c *Console) {
	err := h.verifyUserAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	tracks, err := h.trackService.GetUserFavorites(context.Background(), c.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, track := range tracks {
		dto.NewTrackDTO(track).Print()
		fmt.Println("-----------------------")
	}
}

func (h *Handler) AddToUserFavoritesTrack(c *Console) {
	err := h.verifyUserAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	id, err := readID()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = h.trackService.AddToUserFavorites(context.Background(), id, c.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *Handler) GetByAlbumIDTrack(c *Console) {
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

	tracks, err := h.trackService.GetByAlbumID(context.Background(), id)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, track := range tracks {
		dto.NewTrackDTO(track).Print()
		fmt.Println("-----------------------")
	}
}

func (h *Handler) GetByMusicianIDTrack(c *Console) {
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

	tracks, err := h.trackService.GetByMusicianID(context.Background(), id)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, track := range tracks {
		dto.NewTrackDTO(track).Print()
		fmt.Println("-----------------------")
	}
}

func (h *Handler) GetOwnTrack(c *Console) {
	err := h.verifyMusicianAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	tracks, err := h.trackService.GetOwn(context.Background(), c.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, track := range tracks {
		dto.NewTrackDTO(track).Print()
		fmt.Println("-----------------------")
	}
}

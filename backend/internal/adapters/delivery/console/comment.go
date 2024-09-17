package console

import (
	"context"
	"fmt"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/console/dto"
	"github.com/hanoys/sigma-music/internal/ports"
)

func (h *Handler) PostComment(c *Console) {
	err := h.verifyUserAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	var postDTO dto.PostCommentDTO
	err = dto.InputPostCommentDTO(&postDTO)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = h.commentService.Post(context.Background(), ports.PostCommentServiceReq{
		UserID:  c.UserID,
		TrackID: postDTO.TrackID,
		Stars:   postDTO.Stars,
		Text:    postDTO.Text,
	})

	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) GetCommentsOnTrack(c *Console) {
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

	comments, err := h.commentService.GetCommentsOnTrack(context.Background(), id)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, comment := range comments {
		dto.NewCommentDTO(comment).Print()
        fmt.Println("-----------------------")
	}
}

func (h *Handler) GetUserComments(c *Console) {
	err := h.verifyUserAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	comments, err := h.commentService.GetUserComments(context.Background(), c.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, comment := range comments {
		dto.NewCommentDTO(comment).Print()
        fmt.Println("-----------------------")
	}
}

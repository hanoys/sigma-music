package dto

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type PostCommentDTO struct {
	TrackID uuid.UUID
	Stars   int
	Text    string
}

func InputPostCommentDTO(p *PostCommentDTO) error {
	var id string
	fmt.Print("Track ID:")
	fmt.Scan(&id)

	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("incorrect id")
	}

	p.TrackID = uid

	fmt.Print("Stars (1-5): ")
	_, err = fmt.Scan(&p.Stars)
	if err != nil || p.Stars < 1 || p.Stars > 5 {
		return errors.New("expected integer from 1 to 5")
	}

	fmt.Print("Text: ")
	fmt.Scan(&p.Text)

	return nil
}

type CommentDTO struct {
	ID      string
	UserID  string
	TrackID string
	Stars   int
	Text    string
}

func NewCommentDTO(c domain.Comment) CommentDTO {
	return CommentDTO{
		ID:      c.ID.String(),
		UserID:  c.UserID.String(),
		TrackID: c.TrackID.String(),
		Stars:   c.Stars,
		Text:    c.Text,
	}
}

func (c CommentDTO) Print() {
	fmt.Println("ID:", c.ID)
	fmt.Println("User ID:", c.UserID)
	fmt.Println("Track ID:", c.TrackID)
	fmt.Println("Stars:", c.Stars)
	fmt.Println("Text:", c.Text)
}

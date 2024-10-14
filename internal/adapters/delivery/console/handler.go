package console

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type Handler struct {
	albumService    ports.IAlbumService
	authService     ports.IAuthorizationService
	commentService  ports.ICommentService
	genreService    ports.IGenreService
	musicianService ports.IMusicianService
	statService     ports.IStatService
	trackService    ports.ITrackService
	userService     ports.IUserService
}

type HandlerParams struct {
	AlbumService    ports.IAlbumService
	AuthService     ports.IAuthorizationService
	CommentService  ports.ICommentService
	GenreService    ports.IGenreService
	MusicianService ports.IMusicianService
	StatService     ports.IStatService
	TrackService    ports.ITrackService
	UserService     ports.IUserService
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		albumService:    params.AlbumService,
		authService:     params.AuthService,
		commentService:  params.CommentService,
		genreService:    params.GenreService,
		musicianService: params.MusicianService,
		statService:     params.StatService,
		trackService:    params.TrackService,
		userService:     params.UserService,
	}
}

func (h *Handler) verifyAuth(c *Console) error {
	if c.UserRole == -1 {
		return errors.New("unauthorized")
	}

	return nil
}

func (h *Handler) verifyUserAuth(c *Console) error {
	if c.UserRole != domain.UserRole {
		return errors.New("unauthorized")
	}

	return nil
}

func (h *Handler) verifyMusicianAuth(c *Console) error {
	if c.UserRole != domain.MusicianRole {
		return errors.New("unauthorized")
	}

	return nil
}

func readID() (uuid.UUID, error) {
	var id string
	fmt.Print("ID: ")
	fmt.Scan(&id)

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, errors.New("incorrect id")
	}

	return uid, nil
}

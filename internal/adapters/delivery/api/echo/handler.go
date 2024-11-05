package echo

import (
	"errors"
	"net/http"

	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Services struct {
	AuthService     ports.IAuthorizationService
	AlbumService    ports.IAlbumService
	MusicianService ports.IMusicianService
	UserService     ports.IUserService
	TrackService    ports.ITrackService
	CommentService  ports.ICommentService
	GenreService    ports.IGenreService
}

type Handler struct {
	router      *echo.Echo
	logger      *zap.Logger
	services    *Services
	userHandler *UserHandler
}

func NewHandler(logger *zap.Logger) *Handler {
	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	return &Handler{
		router:   e,
		logger:   logger,
		services: nil,
	}
}

func (h *Handler) GetRouter() *echo.Echo {
	return h.router
}

func (h *Handler) SetServices(services *Services) {
	h.services = services
}

func (h *Handler) ConfigureHandlers() error {
	if h.services == nil {
		return errors.New("services are not set")
	}

	v1EchoRouter := h.router.Group("/api/echo/v1")
    h.userHandler = NewUserHandler(v1EchoRouter, h.logger, h.services)

	return nil
}

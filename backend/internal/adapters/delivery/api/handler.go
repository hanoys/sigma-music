package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hanoys/sigma-music/internal/ports"
	"net/http"
)

type Handler struct {
	router       *gin.Engine
	userHandler  *UserHandler
	albumHandler *AlbumHandler
}

func NewHandler() *Handler {
	return &Handler{router: gin.New()}
}

func (h *Handler) BuildUserHandler(userService ports.IUserService, authorizationService ports.IAuthorizationService) {
	h.userHandler = NewUserHandler(h.router.Group("/"), userService, authorizationService)
}

func (h *Handler) BuildAlbumHandler(albumService ports.IAlbumService) {
	h.albumHandler = NewAlbumHandler(h.router.Group("/"), albumService)
}

func (h *Handler) GetRouter() http.Handler {
	return h.router
}

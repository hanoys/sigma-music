package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/api/dto"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"net/http"
)

type AlbumHandler struct {
	router       *gin.RouterGroup
	albumService ports.IAlbumService
}

func NewAlbumHandler(router *gin.RouterGroup, service ports.IAlbumService) *AlbumHandler {
	albumHandler := &AlbumHandler{
		router:       router,
		albumService: service,
	}
	router.Group("/album")
	{
		router.POST("/", albumHandler.createAlbum)
		router.GET("/:id", albumHandler.getAlbum)
		router.GET("/", albumHandler.getAllAlbums)
		router.GET("/musician/:id", albumHandler.getAlbumsByMusician)
		router.GET("/musician", albumHandler.getOwnAlbums)
		router.POST("/publish/:id", albumHandler.publishAlbum)
	}

	return albumHandler
}

func (h *AlbumHandler) createAlbum(c *gin.Context) {
	var createAlbumDTO dto.CreateAlbumDTO

	if err := c.ShouldBindJSON(&createAlbumDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad json format"})
		return
	}

	val, ok := c.Get("UserInfo")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "server error"})
		return
	}

	payload, ok := val.(domain.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "server error"})
		return
	}

	if payload.Role != domain.MusicianRole {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "you are not able to create album"})
		return
	}

	album, err := h.albumService.Create(c.Request.Context(), ports.CreateAlbumServiceReq{
		MusicianID:  payload.UserID,
		Name:        createAlbumDTO.Name,
		Description: createAlbumDTO.Description,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create album: %v\n", err).Error()})
	}

	c.JSON(http.StatusOK, dto.AlbumFromDomain(album))
}

func (h *AlbumHandler) getAlbum(c *gin.Context) {
	albumID := c.Param("id")

	id, err := uuid.FromBytes([]byte(albumID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Errorf("can't get album: %v\n", err).Error()})
		return
	}

	album, err := h.albumService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create album: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AlbumFromDomain(album))
}

func (h *AlbumHandler) getAllAlbums(c *gin.Context) {

}

func (h *AlbumHandler) getAlbumsByMusician(c *gin.Context) {
	musicianID := c.Param("id")

	id, err := uuid.FromBytes([]byte(musicianID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Errorf("can't get album: %v\n", err).Error()})
		return
	}

	albums, err := h.albumService.GetByMusicianID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create album: %v\n", err).Error()})
		return
	}

	albumsDTO := make([]dto.AlbumDTO, len(albums))
	for i, album := range albums {
		albumsDTO[i] = dto.AlbumFromDomain(album)
	}

	c.JSON(http.StatusOK, albumsDTO)
}

func (h *AlbumHandler) getOwnAlbums(c *gin.Context) {
	val, ok := c.Get("UserInfo")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "server error"})
		return
	}

	payload, ok := val.(domain.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "server error"})
		return
	}

	if payload.Role != domain.MusicianRole {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "you are not able to get own albums"})
		return
	}

	albums, err := h.albumService.GetOwn(c.Request.Context(), payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create album: %v\n", err).Error()})
		return
	}

	albumsDTO := make([]dto.AlbumDTO, len(albums))
	for i, album := range albums {
		albumsDTO[i] = dto.AlbumFromDomain(album)
	}

	c.JSON(http.StatusOK, albumsDTO)
}

func (h *AlbumHandler) publishAlbum(c *gin.Context) {
	musicianID := c.Param("id")

	id, err := uuid.FromBytes([]byte(musicianID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Errorf("can't publish album: %v\n", err).Error()})
		return
	}

	if err = h.albumService.Publish(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Errorf("can't publish album: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "album published"})
}

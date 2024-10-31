package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/api/dto"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type AlbumHandler struct {
	router      *gin.RouterGroup
	logger      *zap.Logger
	s           *Services
	authHandler *AuthHandler
}

func NewAlbumHandler(router *gin.RouterGroup,
	logger *zap.Logger,
	services *Services,
	authHandler *AuthHandler,
) *AlbumHandler {
	albumHandler := &AlbumHandler{
		router:      router,
		s:           services,
		authHandler: authHandler,
	}

	router.PATCH("/albums/:album_id",
		authHandler.verifyToken,
		authHandler.verifyAlbumOwner,
		albumHandler.publish)

	router.GET("/albums/", albumHandler.getAll)
	router.GET("/albums/:album_id", albumHandler.getByID)

	router.GET("/musicians/:musician_id/albums",
		albumHandler.getByMusicianID)
	router.POST("/musicians/:musician_id/albums",
		authHandler.verifyToken,
		authHandler.verifyMusicianRole,
		authHandler.verifyMusicianID,
		albumHandler.create)
	router.GET("/musicians/me/albums",
		authHandler.verifyToken,
		authHandler.verifyMusicianRole,
		albumHandler.getOwn)

	return albumHandler
}

// @Summary CreateAlbum
// @Tags album
// @Security ApiKeyAuth
// @Description create album
// @Accept  json
// @Produce json
// @Param   musician_id   path    string  true  "musician id"
// @Param input body dto.CreateAlbumDTO true "create album info"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 403 {object} RestErrorForbidden
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 201 {object} dto.AlbumDTO
// @Router /musicians/{musician_id}/albums [post]
func (h *AlbumHandler) create(context *gin.Context) {
	var createAlbumDTO dto.CreateAlbumDTO
	err := context.ShouldBindJSON(&createAlbumDTO)
	if err != nil {
		errorResponse(context, err)
		return
	}

	id, err := getIdFromPath(context, "musician_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	album, err := h.s.AlbumService.Create(context.Request.Context(),
		ports.CreateAlbumServiceReq{
			MusicianID:  id,
			Name:        createAlbumDTO.Name,
			Description: createAlbumDTO.Description,
		})
	if err != nil {
		errorResponse(context, err)
		return
	}

	albumDTO := dto.AlbumFromDomain(album)
	createdResponse(context, albumDTO)
}

// @Summary GetAllAlbums
// @Tags album
// @Description get all albums
// @Accept  json
// @Produce json
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.AlbumDTO
// @Router /albums [get]
func (h *AlbumHandler) getAll(context *gin.Context) {
	albums, err := h.s.AlbumService.GetAll(context.Request.Context())
	if err != nil {
		errorResponse(context, err)
		return
	}

	albumDTOs := make([]dto.AlbumDTO, len(albums))
	for i := range albums {
		albumDTOs[i] = dto.AlbumFromDomain(albums[i])
	}

	successResponse(context, albumDTOs)
}

// @Summary GetOwn
// @Security ApiKeyAuth
// @Tags album
// @Description get own albums
// @Accept  json
// @Produce json
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.AlbumDTO
// @Router /musicians/me/albums [get]
func (h *AlbumHandler) getOwn(context *gin.Context) {
	musicianID, err := getIdFromRequestContext(context)
	if err != nil {
		errorResponse(context, err)
		return
	}

	albums, err := h.s.AlbumService.GetOwn(context.Request.Context(), musicianID)
	if err != nil {
		errorResponse(context, err)
		return
	}

	albumDTOs := make([]dto.AlbumDTO, len(albums))
	for i := range albums {
		albumDTOs[i] = dto.AlbumFromDomain(albums[i])
	}

	successResponse(context, albumDTOs)
}

// @Summary GetAlbumByID
// @Tags album
// @Description get album by id
// @Accept  json
// @Produce json
// @Param   id   path    string  true  "album id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.AlbumDTO
// @Router /albums/{id} [get]
func (h *AlbumHandler) getByID(context *gin.Context) {
	id, err := getIdFromPath(context, "album_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	album, err := h.s.AlbumService.GetByID(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	albumDTO := dto.AlbumFromDomain(album)
	successResponse(context, albumDTO)
}

// @Summary PublishAlbum
// @Tags album
// @Security ApiKeyAuth
// @Description publish album
// @Accept  json
// @Produce json
// @Param id   path    string  true  "album id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 403 {object} RestErrorForbidden
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200
// @Router /albums/{id} [patch]
func (h *AlbumHandler) publish(context *gin.Context) {
	id, err := getIdFromPath(context, "album_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	err = h.s.AlbumService.Publish(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	successResponse(context, struct{}{})
}

// @Summary GetAlbumByMusicianID
// @Tags album
// @Description get album by musician id
// @Accept  json
// @Produce json
// @Param id   path    string  true  "musician id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200
// @Router /musicians/{id}/albums [get]
func (h *AlbumHandler) getByMusicianID(context *gin.Context) {
	id, err := getIdFromPath(context, "musician_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	albums, err := h.s.AlbumService.GetByMusicianID(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	albumDTOs := make([]dto.AlbumDTO, len(albums))
	for i := range albums {
		albumDTOs[i] = dto.AlbumFromDomain(albums[i])
	}

	successResponse(context, albumDTOs)
}

package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/api/dto"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type MusicianHandler struct {
	router      *gin.RouterGroup
	logger      *zap.Logger
	authHandler *AuthHandler
	s           *Services
}

func NewMusicianHandler(router *gin.RouterGroup,
	logger *zap.Logger,
	services *Services,
	authHandler *AuthHandler,
) *MusicianHandler {
	musicianHandler := &MusicianHandler{
		router:      router,
		logger:      logger,
		authHandler: authHandler,
		s:           services,
	}

	musicianGroup := router.Group("/musicians")
	{
		musicianGroup.POST("/register", musicianHandler.register)
		musicianGroup.GET("/",
			musicianHandler.getAll)
		musicianGroup.GET("/:musician_id",
			musicianHandler.getByID)
	}

	router.GET("/albums/:album_id/musicians", musicianHandler.getByAlbumID)
	router.GET("/tracks/:track_id/musicians", musicianHandler.getByTrackID)

	router.PUT("/musicians/:musician_id/image",
		authHandler.verifyToken,
		authHandler.verifyMusicianRole,
		musicianHandler.uploadImage)

	return musicianHandler
}

// @Summary UploadImage
// @Tags musician
// @Security ApiKeyAuth
// @Description upload musician avatar
// @Accept  mpfd
// @Produce json
// @Param   musician_id   path    string  true  "musician id"
// @Param image formData file true "upload file"
// @Failure 500 {object} RestErrorInternalError
// @Success 201 {object} dto.MusicianDTO
// @Router /musicians/{musician_id}/image [put]
func (h *MusicianHandler) uploadImage(context *gin.Context) {
	musician_id, err := getIdFromRequestContext(context)
	if err != nil {
		errorResponse(context, err)
		return
	}

	fileheader, err := context.FormFile("image")
	if err != nil {
		errorResponse(context, err)
		return
	}

	file, err := fileheader.Open()
	if err != nil {
		errorResponse(context, err)
		return
	}

	musician, err := h.s.MusicianService.UploadImage(context.Request.Context(), file, musician_id)
	if err != nil {
        log.Println("MUS ERR: ", err)
		errorResponse(context, err)
		return
	}

	musicianDTO := dto.MusicianFromDomain(musician)
	successResponse(context, musicianDTO)
}

// @Summary MusicianRegister
// @Tags musician
// @Description registration
// @Accept  json
// @Produce json
// @Param input body dto.RegisterMusicianDTO true "musician information"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 409 {object} RestErrorConflict
// @Failure 500 {object} RestErrorInternalError
// @Success 201 {string} string "message"
// @Router /musicians/register [post]
func (h *MusicianHandler) register(context *gin.Context) {
	var registerDTO dto.RegisterMusicianDTO
	err := context.ShouldBindJSON(&registerDTO)
	if err != nil {
		errorResponse(context, err)
		return
	}

	musician, err := h.s.MusicianService.Register(
		context.Request.Context(),
		ports.MusicianServiceCreateRequest{
			Name:        registerDTO.Name,
			Email:       registerDTO.Email,
			Password:    registerDTO.Password,
			Country:     registerDTO.Country,
			Description: registerDTO.Description,
		},
	)
	if err != nil {
		errorResponse(context, err)
		return
	}

	musicianDTO := dto.MusicianFromDomain(musician)
	createdResponse(context, musicianDTO)
}

// @Summary GetAllMusicians
// @Tags musician
// @Description get all musicians
// @Accept  json
// @Produce json
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.MusicianDTO
// @Router /musicians [get]
func (h *MusicianHandler) getAll(context *gin.Context) {
	musicians, err := h.s.MusicianService.GetAll(context.Request.Context())
	if err != nil {
		errorResponse(context, err)
		return
	}

	musicianDTOs := make([]dto.MusicianDTO, len(musicians))
	for i := range musicians {
		musicianDTOs[i] = dto.MusicianFromDomain(musicians[i])
	}

	successResponse(context, musicianDTOs)
}

// @Summary GetMusicianByID
// @Tags musician
// @Description get musician by id
// @Accept  json
// @Produce json
// @Param   id   path    string  true  "musician id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.MusicianDTO
// @Router /musicians/{musician_id} [get]
func (h *MusicianHandler) getByID(context *gin.Context) {
	id, err := getIdFromPath(context, "musician_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	musician, err := h.s.MusicianService.GetByID(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	musicianDTO := dto.MusicianFromDomain(musician)
	successResponse(context, musicianDTO)
}

// @Summary GetByAlbumID
// @Tags musician
// @Description get musician by album id
// @Accept  json
// @Produce json
// @Param   album_id   path    string  true  "album id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.MusicianDTO
// @Router /albums/{album_id}/musicians [get]
func (h *MusicianHandler) getByAlbumID(context *gin.Context) {
	id, err := getIdFromPath(context, "album_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	musician, err := h.s.MusicianService.GetByAlbumID(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	musicianDTO := dto.MusicianFromDomain(musician)
	successResponse(context, musicianDTO)
}

// @Summary GetByTrackID
// @Tags musician
// @Description get musician by track id
// @Accept  json
// @Produce json
// @Param   track_id   path    string  true  "track id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.MusicianDTO
// @Router /tracks/{track_id}/musicians [get]
func (h *MusicianHandler) getByTrackID(context *gin.Context) {
	id, err := getIdFromPath(context, "track_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	musician, err := h.s.MusicianService.GetByTrackID(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	musicianDTO := dto.MusicianFromDomain(musician)
	successResponse(context, musicianDTO)
}

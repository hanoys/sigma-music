package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/api/dto"
	"go.uber.org/zap"
)

type GenreHandler struct {
	router      *gin.RouterGroup
	logger      *zap.Logger
	authHandler *AuthHandler
	s           *Services
}

func NewGenreHandler(router *gin.RouterGroup,
	logger *zap.Logger,
	services *Services,
	authHandler *AuthHandler) *GenreHandler {
	genreHandler := &GenreHandler{
		router:      router,
		logger:      logger,
		authHandler: authHandler,
		s:           services,
	}

	router.GET("/genres/", genreHandler.getAll)
	router.GET("/genres/:id", genreHandler.getByID)

	router.PATCH("/tracks/:track_id/genres",
		authHandler.verifyToken,
		genreHandler.addForTrack)
	router.GET("/tracks/:track_id/genres",
		genreHandler.getByTrackID)

	return genreHandler
}

// @Summary GetAllGenres
// @Tags genre
// @Description get all genres
// @Accept  json
// @Produce json
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.GenreDTO
// @Router /genres [get]
func (h *GenreHandler) getAll(context *gin.Context) {
	genres, err := h.s.GenreService.GetAll(context.Request.Context())
	if err != nil {
		errorResponse(context, err)
		return
	}

	genreDTOs := make([]dto.GenreDTO, len(genres))
	for i := range genres {
		genreDTOs[i] = dto.GenreFromDomain(genres[i])
	}

	successResponse(context, genreDTOs)
}

// @Summary GetGenreByID
// @Tags genre
// @Description get genre by id
// @Accept  json
// @Produce json
// @Param   id   path    string  true  "genre id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.GenreDTO
// @Router /genres/{id} [get]
func (h *GenreHandler) getByID(context *gin.Context) {
	id, err := getIdFromPath(context, "id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	genre, err := h.s.GenreService.GetByID(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	genreDTO := dto.GenreFromDomain(genre)
	successResponse(context, genreDTO)
}

// @Summary GetGenresByTrackID
// @Tags genre
// @Description get genres by track id
// @Accept  json
// @Produce json
// @Param   track_id   path    string  true  "track id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.GenreDTO
// @Router /tracks/{track_id}/genres [get]
func (h *GenreHandler) getByTrackID(context *gin.Context) {
	id, err := getIdFromPath(context, "track_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	genres, err := h.s.GenreService.GetByTrackID(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	genreDTOs := make([]dto.GenreDTO, len(genres))
	for i := range genres {
		genreDTOs[i] = dto.GenreFromDomain(genres[i])
	}

	successResponse(context, genreDTOs)
}

// @Summary SetGenresForTrack
// @Tags genre
// @Security ApiKeyAuth
// @Description set genres for track
// @Accept  json
// @Produce json
// @Param input body dto.AddForTrackDTO true "add genres payload"
// @Param track_id   path    string  true  "track id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 403 {object} RestErrorForbidden
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200
// @Router /tracks/{track_id}/genres [patch]
func (h *GenreHandler) addForTrack(context *gin.Context) {
	trackID, err := getIdFromPath(context, "track_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	var addGenresDTO dto.AddForTrackDTO
	err = context.ShouldBindJSON(&addGenresDTO)
	if err != nil {
		errorResponse(context, err)
		return
	}

	genresCount := 0
	if addGenresDTO.GenreIDs != nil {
		genresCount = len(addGenresDTO.GenreIDs)
	}

	genreIDs := make([]uuid.UUID, genresCount)
	for i, genre := range addGenresDTO.GenreIDs {
		id, err := uuid.Parse(genre)
		if err != nil {
			errorResponse(context, ParseGenreIDError)
			return
		}

		genreIDs[i] = id
	}

	err = h.s.GenreService.AddForTrack(context.Request.Context(), trackID, genreIDs)
	if err != nil {
		errorResponse(context, err)
		return
	}

	successResponse(context, struct{}{})
}

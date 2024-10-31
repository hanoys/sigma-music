package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/api/dto"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type TrackHandler struct {
	router      *gin.RouterGroup
	logger      *zap.Logger
	authHandler *AuthHandler
	s           *Services
}

func NewTrackHandler(router *gin.RouterGroup,
	logger *zap.Logger,
	services *Services,
	authHandler *AuthHandler,
) *TrackHandler {
	trackHandler := &TrackHandler{
		router:      router,
		logger:      logger,
		authHandler: authHandler,
		s:           services,
	}

	router.GET("/tracks/", trackHandler.getAll)
	router.GET("/tracks/:track_id", trackHandler.getByID)
	router.DELETE("/musicians/:musician_id/tracks/:track_id",
		authHandler.verifyToken,
		authHandler.verifyTrackOwner,
		trackHandler.delete)

	router.POST("/musicians/:musician_id/albums/:album_id/tracks",
		authHandler.verifyToken,
		authHandler.verifyMusicianAlbumOwner,
		trackHandler.create,
	)

	router.GET("/albums/:album_id/tracks",
		trackHandler.getByAlbumID)
	router.GET("/users/me/favorites",
		authHandler.verifyToken,
		authHandler.verifyUserRole,
		trackHandler.getFavorites)
	router.POST("/users/me/favorites/:track_id",
		authHandler.verifyToken,
		authHandler.verifyUserRole,
		trackHandler.addToFavorites)
	router.GET("/musicians/:musician_id/tracks",
		trackHandler.getByMusicianID)
	router.GET("/musicians/me/tracks",
		authHandler.verifyToken,
		authHandler.verifyMusicianRole,
		trackHandler.getOwn)

	return trackHandler
}

// @Summary CreateTrack
// @Tags track
// @Security ApiKeyAuth
// @Description create track
// @Accept  json
// @Produce json
// @Param   musician_id   path    string  true  "musician id"
// @Param   album_id   path    string  true  "album id"
// @Param input body dto.CreateTrackDTO true "create track info"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 403 {object} RestErrorForbidden
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 201 {object} dto.TrackDTO
// @Router /musicians/{musician_id}/albums/{album_id}/tracks [post]
func (h *TrackHandler) create(context *gin.Context) {
	albumID, err := getIdFromPath(context, "album_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	var createTrackDTO dto.CreateTrackDTO
	context.ShouldBindJSON(&createTrackDTO)

	genresCount := 0
	if createTrackDTO.GenreIDs != nil {
		genresCount = len(createTrackDTO.GenreIDs)
	}

	genreIDs := make([]uuid.UUID, genresCount)
	for i, genre := range createTrackDTO.GenreIDs {
		id, err := uuid.Parse(genre)
		if err != nil {
			errorResponse(context, ParseGenreIDError)
			return
		}

		genreIDs[i] = id
	}

	track, err := h.s.TrackService.Create(
		context.Request.Context(),
		ports.CreateTrackReq{
			AlbumID:   albumID,
			Name:      createTrackDTO.Name,
			TrackBLOB: strings.NewReader("fdf"),
			GenresID:  genreIDs,
		},
	)
	if err != nil {
		errorResponse(context, err)
		return
	}

	trackDTO := dto.TrackFromDomain(track)
	successResponse(context, trackDTO)
}

// @Summary GetAllTracks
// @Tags track
// @Description get all tracks
// @Accept  json
// @Produce json
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.TrackDTO
// @Router /tracks [get]
func (h *TrackHandler) getAll(context *gin.Context) {
	tracks, err := h.s.TrackService.GetAll(context.Request.Context())
	if err != nil {
		errorResponse(context, err)
		return
	}

	trackDTOs := make([]dto.TrackDTO, len(tracks))
	for i := range tracks {
		trackDTOs[i] = dto.TrackFromDomain(tracks[i])
	}

	successResponse(context, trackDTOs)
}

// @Summary getOwn
// @Tags track
// @Security ApiKeyAuth
// @Description get own tracks
// @Accept  json
// @Produce json
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.TrackDTO
// @Router /musicians/me/tracks [get]
func (h *TrackHandler) getOwn(context *gin.Context) {
	musicianID, err := getIdFromRequestContext(context)
	if err != nil {
		errorResponse(context, err)
		return
	}

	tracks, err := h.s.TrackService.GetOwn(context.Request.Context(), musicianID)
	if err != nil {
		errorResponse(context, err)
		return
	}

	trackDTOs := make([]dto.TrackDTO, len(tracks))
	for i := range tracks {
		trackDTOs[i] = dto.TrackFromDomain(tracks[i])
	}

	successResponse(context, trackDTOs)
}

// @Summary GetTrackByID
// @Tags track
// @Description get track by id
// @Accept  json
// @Produce json
// @Param   id   path    string  true  "track id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.TrackDTO
// @Router /tracks/{id} [get]
func (h *TrackHandler) getByID(context *gin.Context) {
	id, err := getIdFromPath(context, "track_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	track, err := h.s.TrackService.GetByID(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	trackDTO := dto.TrackFromDomain(track)
	successResponse(context, trackDTO)
}

// @Summary DeleteTrack
// @Tags track
// @Description get track by id
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param   track_id   path    string  true  "track id"
// @Param   musician_id   path    string  true  "musician id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 403 {object} RestErrorForbidden
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.TrackDTO
// @Router /musicians/{musician_id}/tracks/{track_id} [delete]
func (h *TrackHandler) delete(context *gin.Context) {
	id, err := getIdFromPath(context, "track_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	track, err := h.s.TrackService.Delete(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	trackDTO := dto.TrackFromDomain(track)
	successResponse(context, trackDTO)
}

// @Summary GetTracksByAlbumID
// @Tags track
// @Description get tracks by album id
// @Accept  json
// @Produce json
// @Param   id   path    string  true  "album id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.TrackDTO
// @Router /albums/{id}/tracks [get]
func (h *TrackHandler) getByAlbumID(context *gin.Context) {
	id, err := getIdFromPath(context, "album_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	tracks, err := h.s.TrackService.GetByAlbumID(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	trackDTOs := make([]dto.TrackDTO, len(tracks))
	for i := range tracks {
		trackDTOs[i] = dto.TrackFromDomain(tracks[i])
	}

	successResponse(context, trackDTOs)
}

// @Summary GetFavorites
// @Tags track
// @Security ApiKeyAuth
// @Description get user favorites tracks
// @Accept  json
// @Produce json
// @Failure 400 {object} RestErrorBadRequest
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 403 {object} RestErrorForbidden
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.TrackDTO
// @Router /users/me/favorites [get]
func (h *TrackHandler) getFavorites(context *gin.Context) {
	id, err := getIdFromRequestContext(context)
	if err != nil {
		errorResponse(context, err)
		return
	}

	tracks, err := h.s.TrackService.GetUserFavorites(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	trackDTOs := make([]dto.TrackDTO, len(tracks))
	for i := range tracks {
		trackDTOs[i] = dto.TrackFromDomain(tracks[i])
	}

	successResponse(context, trackDTOs)
}

// @Summary AddToFavorites
// @Tags track
// @Security ApiKeyAuth
// @Description add track to user favorites
// @Accept  json
// @Produce json
// @Param   id   path    string  true  "track id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 403 {object} RestErrorForbidden
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {string} string ""
// @Router /users/me/favorites/{id} [post]
func (h *TrackHandler) addToFavorites(context *gin.Context) {
	userID, err := getIdFromRequestContext(context)
	if err != nil {
		errorResponse(context, err)
		return
	}

	trackID, err := getIdFromPath(context, "track_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	err = h.s.TrackService.AddToUserFavorites(context.Request.Context(), trackID, userID)
	if err != nil {
		errorResponse(context, err)
		return
	}

	successResponse(context, struct{}{})
}

// @Summary GetTracksByMusicianID
// @Tags track
// @Description get tracks by musician id
// @Accept  json
// @Produce json
// @Param   id   path    string  true  "musician id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.TrackDTO
// @Router /musicians/{id}/tracks [get]
func (h *TrackHandler) getByMusicianID(context *gin.Context) {
	id, err := getIdFromPath(context, "musician_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	tracks, err := h.s.TrackService.GetByMusicianID(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	trackDTOs := make([]dto.TrackDTO, len(tracks))
	for i := range tracks {
		trackDTOs[i] = dto.TrackFromDomain(tracks[i])
	}

	successResponse(context, trackDTOs)
}

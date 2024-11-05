package gin

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/api/dto"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type AuthHandler struct {
	router *gin.RouterGroup
	logger *zap.Logger
	s      *Services
}

func NewAuthHandler(router *gin.RouterGroup, logger *zap.Logger, services *Services) *AuthHandler {
	authHandler := &AuthHandler{
		router: router,
		logger: logger,
		s:      services,
	}

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authHandler.login)
		authGroup.POST("/logout",
			authHandler.verifyToken,
			authHandler.logout)
		authGroup.POST("/refresh", authHandler.refresh)
	}

	return authHandler
}

// @Summary TokenRefresh
// @Tags auth
// @Description refresh
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param input body dto.RefreshDTO true "refresh payload"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.LoginResponseDTO
// @Router /auth/refresh [post]
func (h *AuthHandler) refresh(context *gin.Context) {
	var refreshDTO dto.RefreshDTO
	err := context.ShouldBindJSON(&refreshDTO)
	if err != nil {
		errorResponse(context, err)
		return
	}

	tokenPair, err := h.s.AuthService.RefreshToken(context.Request.Context(), refreshDTO.RefreshToken)
	if err != nil {
		errorResponse(context, err)
		return
	}

	response := dto.LoginResponseFromTokenPair(tokenPair)
	successResponse(context, response)
}

// @Summary LogIn
// @Tags auth
// @Description login
// @Accept  json
// @Produce json
// @Param input body dto.LoginDTO true "credentials"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.LoginResponseDTO
// @Router /auth/login [post]
func (h *AuthHandler) login(context *gin.Context) {
	var loginDTO dto.LoginDTO
	err := context.ShouldBindJSON(&loginDTO)
	if err != nil {
		errorResponse(context, err)
		return
	}

	tokenPair, err := h.s.AuthService.LogIn(
		context.Request.Context(),
		ports.LogInCredentials{
			Name:     loginDTO.Name,
			Password: loginDTO.Password,
		},
	)
	if err != nil {
		errorResponse(context, err)
		return
	}

	response := dto.LoginResponseFromTokenPair(tokenPair)
	successResponse(context, response)
}

// @Summary Logout
// @Tags auth
// @Security ApiKeyAuth
// @Description logout
// @Accept  json
// @Produce json
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 500 {object} RestErrorInternalError
// @Success 200
// @Router /auth/logout [post]
func (h *AuthHandler) logout(context *gin.Context) {
	tokenString, err := h.extractAuthToken(context)
	if err != nil {
		//h.logger.Error("Failed to extart auth token", zap.Error(err))
		errorResponse(context, err)
		return
	}

	err = h.s.AuthService.LogOut(context.Request.Context(), tokenString)
	if err != nil {
		errorResponse(context, err)
		return
	}

	successResponse(context, struct{}{})
}

func (h *AuthHandler) verifyMusicianRole(context *gin.Context) {
	role, err := getRoleFromRequestContext(context)
	if err != nil {
		errorResponse(context, err)
		return
	}

	if role != domain.MusicianRole {
		errorResponse(context, ForbiddenError)
		return
	}
}

func (h *AuthHandler) verifyUserRole(context *gin.Context) {
	role, err := getRoleFromRequestContext(context)
	if err != nil {
		errorResponse(context, err)
		return
	}

	if role != domain.UserRole {
		errorResponse(context, ForbiddenError)
		return
	}
}

func (h *AuthHandler) verifyMusicianID(context *gin.Context) {
	id, err := getIdFromRequestContext(context)
	if err != nil {
		errorResponse(context, err)
		return
	}

	expectedID, err := getIdFromPath(context, "musician_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	_, err = h.s.MusicianService.GetByID(context.Request.Context(), expectedID)
	if err != nil {
		errorResponse(context, err)
		return
	}

	if id != expectedID {
		errorResponse(context, ForbiddenError)
		return
	}
}

func (h *AuthHandler) verifyMusicianAlbumOwner(context *gin.Context) {
	albumID, err := getIdFromPath(context, "album_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	musicianID, err := getIdFromPath(context, "musician_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	musician, err := h.s.MusicianService.GetByAlbumID(context.Request.Context(), albumID)
	if err != nil {
		errorResponse(context, err)
		return
	}

	if musicianID != musician.ID {
		errorResponse(context, ForbiddenError)
		return
	}
}

func (h *AuthHandler) verifyAlbumOwner(context *gin.Context) {
	albumID, err := getIdFromPath(context, "album_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	userID, err := getIdFromRequestContext(context)
	if err != nil {
		errorResponse(context, UnauthorizedError)
	}

	musician, err := h.s.MusicianService.GetByAlbumID(context.Request.Context(), albumID)
	if err != nil {
		errorResponse(context, err)
		return
	}

	if userID != musician.ID {
		errorResponse(context, ForbiddenError)
		return
	}
}

func (h *AuthHandler) verifyTrackOwner(context *gin.Context) {
	trackID, err := getIdFromPath(context, "track_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	userID, err := getIdFromRequestContext(context)
	if err != nil {
		errorResponse(context, UnauthorizedError)
	}

	tracks, err := h.s.TrackService.GetOwn(context.Request.Context(), userID)
	if err != nil {
		errorResponse(context, err)
		return
	}

	var track domain.Track
	found := false
	for _, trackEn := range tracks {
		if trackEn.ID.String() == trackID.String() {
			found = true
			track = trackEn
			break
		}
	}

	if !found {
		errorResponse(context, ports.ErrTrackIDNotFound)
		return
	}

	musician, err := h.s.MusicianService.GetByAlbumID(context.Request.Context(), track.AlbumID)
	if err != nil {
		errorResponse(context, err)
		return
	}

	if userID != musician.ID {
		errorResponse(context, ForbiddenError)
		return
	}
}

func (h *AuthHandler) verifyToken(context *gin.Context) {
	tokenString, err := h.extractAuthToken(context)
	if err != nil {
		//h.logger.Error("Failed to extart auth token", zap.Error(err))
		errorResponse(context, err)
		return
	}

	payload, err := h.s.AuthService.VerifyToken(context.Request.Context(), tokenString)
	if err != nil {
		//h.logger.Error("Failed to verify token", zap.Error(err))
		errorResponse(context, err)
		return
	}

	context.Set("UserID", payload.UserID.String())
	context.Set("UserRole", payload.Role)
}

func (h *AuthHandler) extractAuthToken(context *gin.Context) (string, error) {
	authHeader := context.GetHeader("Authorization")
	if authHeader == "" {
		//h.logger.Error("No authorization header in request", zap.Error(UnauthorizedError))
		return "", UnauthorizedError
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		//h.logger.Error("No {Bearer} string before the token", zap.Error(UnauthorizedError))
		return "", UnauthorizedError
	}

	if len(headerParts[1]) == 0 {
		//h.logger.Error("No token after {Bearer} string", zap.Error(UnauthorizedError))
		return "", UnauthorizedError
	}

	return headerParts[1], nil
}

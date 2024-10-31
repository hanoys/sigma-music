package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/api/dto"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type CommentHandler struct {
	router      *gin.RouterGroup
	logger      *zap.Logger
	authHandler *AuthHandler
	s           *Services
}

func NewCommentHandler(router *gin.RouterGroup,
	logger *zap.Logger,
	services *Services,
	authHandler *AuthHandler) *CommentHandler {
	commentHandler := &CommentHandler{
		router:      router,
		logger:      logger,
		authHandler: authHandler,
		s:           services,
	}

	router.POST("/tracks/:track_id/comments",
		authHandler.verifyToken,
		authHandler.verifyUserRole,
		commentHandler.post)
	router.GET("/tracks/:track_id/comments",
		commentHandler.getOnTrack)
	router.GET("/users/me/comments",
		authHandler.verifyToken,
		commentHandler.getUsers)

	return commentHandler
}

// @Summary PostComment
// @Tags comment
// @Security ApiKeyAuth
// @Description post comment
// @Accept  json
// @Produce json
// @Param id   path    string  true  "track id"
// @Param input body dto.PostCommentDTO true "post comment info"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 401 {object} RestErrorUnauthorized
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.CommentDTO
// @Router /tracks/{id}/comments [post]
func (h *CommentHandler) post(context *gin.Context) {
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

	var postDTO dto.PostCommentDTO
	err = context.ShouldBindJSON(&postDTO)
	if err != nil {
		errorResponse(context, err)
		return
	}

	comment, err := h.s.CommentService.Post(
		context.Request.Context(),
		ports.PostCommentServiceReq{
			UserID:  userID,
			TrackID: trackID,
			Stars:   postDTO.Stars,
			Text:    postDTO.Text,
		},
	)
	if err != nil {
		errorResponse(context, err)
		return
	}

	commentDTO := dto.CommentFromDomain(comment)
	successResponse(context, commentDTO)
}

// @Summary GetByTrackID
// @Tags comment
// @Description get comments by track id
// @Accept  json
// @Produce json
// @Param   id   path    string  true  "track id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.TrackDTO
// @Router /tracks/{id}/comments [get]
func (h *CommentHandler) getOnTrack(context *gin.Context) {
	trackID, err := getIdFromPath(context, "track_id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	comments, err := h.s.CommentService.GetCommentsOnTrack(context.Request.Context(), trackID)
	if err != nil {
		errorResponse(context, err)
		return
	}

	commentDTOs := make([]dto.CommentDTO, len(comments))
	for i := range comments {
		commentDTOs[i] = dto.CommentFromDomain(comments[i])
	}

	successResponse(context, commentDTOs)
}

// @Summary GetUsersComments
// @Tags comment
// @Security ApiKeyAuth
// @Description get comments by user id
// @Accept  json
// @Produce json
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.TrackDTO
// @Router /users/me/comments [get]
func (h *CommentHandler) getUsers(context *gin.Context) {
	userID, err := getIdFromRequestContext(context)
	if err != nil {
		errorResponse(context, err)
		return
	}

	comments, err := h.s.CommentService.GetUserComments(context.Request.Context(), userID)
	if err != nil {
		errorResponse(context, err)
		return
	}

	commentDTOs := make([]dto.CommentDTO, len(comments))
	for i := range comments {
		commentDTOs[i] = dto.CommentFromDomain(comments[i])
	}

	successResponse(context, commentDTOs)
}

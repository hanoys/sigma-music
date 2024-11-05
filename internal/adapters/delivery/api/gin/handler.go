package gin

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/hanoys/sigma-music/docs"
	"github.com/hanoys/sigma-music/internal/ports"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router          *gin.Engine
	logger          *zap.Logger
	services        *Services
	albumHandler    *AlbumHandler
	userHandler     *UserHandler
	authHandler     *AuthHandler
	musicianHandler *MusicianHandler
	genreHandler    *GenreHandler
	commentHandler  *CommentHandler
	trackHandler    *TrackHandler
}

func NewHandler(logger *zap.Logger) *Handler {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return &Handler{router: router}
}

func (h *Handler) SetServices(services *Services) {
	h.services = services
}

func (h *Handler) ConfigureHandlers() error {
	if h.services == nil {
		return errors.New("services are not set")
	}

	v1Router := h.router.Group("/api/v1")
	h.authHandler = NewAuthHandler(v1Router, h.logger, h.services)
	h.albumHandler = NewAlbumHandler(v1Router, h.logger, h.services, h.authHandler)
	h.userHandler = NewUserHandler(v1Router, h.logger, h.services)
	h.musicianHandler = NewMusicianHandler(v1Router, h.logger, h.services, h.authHandler)
	h.genreHandler = NewGenreHandler(v1Router, h.logger, h.services, h.authHandler)
	h.commentHandler = NewCommentHandler(v1Router, h.logger, h.services, h.authHandler)
	h.trackHandler = NewTrackHandler(v1Router, h.logger, h.services, h.authHandler)

	return nil
}

func (h *Handler) GetRouter() http.Handler {
	return h.router
}

func getIdFromPath(c *gin.Context, paramName string) (uuid.UUID, error) {
	log.Println("DEBUG: c.Params: ", c.Params)
	idString := c.Param(paramName)
	if idString == "" {
		return uuid.UUID{}, PathIDNotFoundError
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println("PARSE ID ERROR: ", err)
		return uuid.UUID{}, InvalidPathIDError
	}

	return id, nil
}

func getIdFromRequestContext(context *gin.Context) (uuid.UUID, error) {
	id, ok := context.Get("UserID")
	if !ok {
		return uuid.UUID{}, UnauthorizedError
	}

	idParsed, _ := uuid.Parse(id.(string))
	return idParsed, nil
}

func getRoleFromRequestContext(context *gin.Context) (int, error) {
	role, ok := context.Get("UserRole")
	if !ok {
		return 0, UnauthorizedError
	}

	return role.(int), nil
}

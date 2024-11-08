package echo

import (
	"net/http"

	"github.com/hanoys/sigma-music/internal/adapters/delivery/api/dto"
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

type UserHandler struct {
	logger *zap.Logger
	s      *Services
}

func NewUserHandler(e *echo.Echo, logger *zap.Logger, services *Services) *UserHandler {
	userHandler := &UserHandler{
		logger: logger,
		s:      services,
	}

	userGroup := e.Group("/users")
	{
		userGroup.POST("/register", userHandler.register)
		userGroup.GET("/", userHandler.getAll)
	}

	return userHandler
}

func (h *UserHandler) register(context echo.Context) error {
	var registerDTO dto.RegisterUserDTO
	err := context.Bind(&registerDTO)
	if err != nil {
		return context.String(http.StatusBadRequest, "bad request")
	}

	user, err := h.s.UserService.Register(
		context.Request().Context(),
		ports.UserServiceCreateRequest{
			Name:     registerDTO.Name,
			Email:    registerDTO.Email,
			Phone:    registerDTO.Phone,
			Password: registerDTO.Password,
			Country:  registerDTO.Country,
		},
	)
	if err != nil {
		return context.String(http.StatusBadRequest, "can't register user")
	}

	userDTO := dto.UserFromDomain(user)

	return context.JSON(http.StatusOK, userDTO)
}

func (h *UserHandler) getAll(context echo.Context) error {
	users, err := h.s.UserService.GetAll(context.Request().Context())
	if err != nil {
		return context.String(http.StatusInternalServerError, "can't get all users")
	}

	userDTOs := make([]dto.UserDTO, len(users))
	for i := range users {
		userDTOs[i] = dto.UserFromDomain(users[i])
	}

	return context.JSON(http.StatusOK, userDTOs)
}

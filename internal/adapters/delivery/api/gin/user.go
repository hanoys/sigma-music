package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/api/dto"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
)

type UserHandler struct {
	router *gin.RouterGroup
	logger *zap.Logger
	s      *Services
}

func NewUserHandler(router *gin.RouterGroup,
	logger *zap.Logger,
	services *Services,
) *UserHandler {
	userHandler := &UserHandler{
		router: router,
		logger: logger,
		s:      services,
	}

	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", userHandler.register)
		userGroup.GET("/",
			userHandler.getAll)
		userGroup.GET("/:id",
			userHandler.getByID)
	}

	return userHandler
}

// @Summary UserRegister
// @Tags user
// @Description registration
// @Accept  json
// @Produce json
// @Param input body dto.RegisterUserDTO true "user information"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 409 {object} RestErrorConflict
// @Failure 500 {object} RestErrorInternalError
// @Success 201 {string} string "message"
// @Router /users/register [post]
func (h *UserHandler) register(context *gin.Context) {
	var registerDTO dto.RegisterUserDTO
	err := context.ShouldBindJSON(&registerDTO)
	if err != nil {
		errorResponse(context, err)
		return
	}

	user, err := h.s.UserService.Register(
		context.Request.Context(),
		ports.UserServiceCreateRequest{
			Name:     registerDTO.Name,
			Email:    registerDTO.Email,
			Phone:    registerDTO.Phone,
			Password: registerDTO.Password,
			Country:  registerDTO.Country,
		},
	)
	if err != nil {
		errorResponse(context, err)
		return
	}

	userDTO := dto.UserFromDomain(user)
	createdResponse(context, userDTO)
}

// @Summary GetAllUsers
// @Tags user
// @Description get all users
// @Accept  json
// @Produce json
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} []dto.UserDTO
// @Router /users [get]
func (h *UserHandler) getAll(context *gin.Context) {
	users, err := h.s.UserService.GetAll(context.Request.Context())
	if err != nil {
		errorResponse(context, err)
		return
	}

	userDTOs := make([]dto.UserDTO, len(users))
	for i := range users {
		userDTOs[i] = dto.UserFromDomain(users[i])
	}

	successResponse(context, userDTOs)
}

// @Summary GetUserByID
// @Tags user
// @Description get user id
// @Accept  json
// @Produce json
// @Param   id   path    string  true  "user id"
// @Failure 400 {object} RestErrorBadRequest
// @Failure 404 {object} RestErrorNotFound
// @Failure 500 {object} RestErrorInternalError
// @Success 200 {object} dto.UserDTO
// @Router /users/{id} [get]
func (h *UserHandler) getByID(context *gin.Context) {
	id, err := getIdFromPath(context, "id")
	if err != nil {
		errorResponse(context, err)
		return
	}

	user, err := h.s.UserService.GetById(context.Request.Context(), id)
	if err != nil {
		errorResponse(context, err)
		return
	}

	userDTO := dto.UserFromDomain(user)
	successResponse(context, userDTO)
}

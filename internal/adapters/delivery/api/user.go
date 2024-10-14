package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/api/dto"
	"github.com/hanoys/sigma-music/internal/ports"
	"net/http"
)

type UserHandler struct {
	router      *gin.RouterGroup
	userService ports.IUserService
	authService ports.IAuthorizationService
}

func NewUserHandler(router *gin.RouterGroup, service ports.IUserService, authorizationService ports.IAuthorizationService) *UserHandler {
	userHandler := &UserHandler{
		router:      router,
		userService: service,
		authService: authorizationService,
	}

	userGroup := router.Group("/user")
	{
		userGroup.POST("/register", userHandler.register)
		userGroup.POST("/login", userHandler.login)
		userGroup.POST("/logout", userHandler.logout)
	}

	return userHandler
}

func (u *UserHandler) register(c *gin.Context) {
	var registerDTO dto.RegisterUserDTO

	if err := c.ShouldBindJSON(&registerDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad json format"})
		return
	}

	_, err := u.userService.Register(c.Request.Context(), registerDTO.ToServiceRequest())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create user: %v\n", err).Error()})
	}

	c.JSON(http.StatusOK, gin.H{"msg": "user created"})
}

func (u *UserHandler) login(c *gin.Context) {
	var loginDTO dto.LoginUserDTO

	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad json format"})
		return
	}

	tokenPair, err := u.authService.LogIn(c.Request.Context(), loginDTO.ToServiceRequest())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't login user: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LoginUserResponseFromTokenPair(tokenPair))
}

func (u *UserHandler) logout(c *gin.Context) {
	var logoutDTO dto.LogoutUserDTO

	if err := c.ShouldBindJSON(&logoutDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad json format"})
		return
	}

	err := u.authService.LogOut(c.Request.Context(), logoutDTO.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't logout user: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, "user logged out")
}

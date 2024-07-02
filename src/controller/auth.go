package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/infra/entity"
	"github.com/pc-power-api/src/infra/repo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthHandler struct {
	userRepo *repo.UserRepository
}

func NewAuthHandler(e *gin.Engine, jwtMiddleware *jwt.GinJWTMiddleware, userRepo *repo.UserRepository) {
	handler := &AuthHandler{
		userRepo: userRepo,
	}

	group := e.Group("/auth")
	{
		group.POST("/login", jwtMiddleware.LoginHandler)
		group.GET("/refresh_token", jwtMiddleware.RefreshHandler)
		group.POST("/register", handler.register)
	}
}

func (h *AuthHandler) register(c *gin.Context) {
	var credentials *api.RegisterCredentials
	err := c.ShouldBind(&credentials)
	if err != nil {
		c.Error(errors.New(err))
		return
	}

	userUuid := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)

	var newUser = entity.User{
		ID:       userUuid.String(),
		Username: credentials.Username,
		Password: string(hashedPassword),
	}

	aerr := h.userRepo.Create(&newUser)
	if aerr != nil {
		c.Error(aerr)
		return
	}

	c.Status(http.StatusNoContent)
}

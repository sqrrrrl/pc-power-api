package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pc-power-api/src/controller/gateway"
	"github.com/pc-power-api/src/controller/middleware"
	"github.com/pc-power-api/src/infra/repo"
)

type UsersHandler struct {
	userRepo   *repo.UserRepository
	deviceRepo *repo.DeviceRepository
}

func NewUsersHandler(e *gin.Engine, jwtMiddleware *jwt.GinJWTMiddleware, userRepo *repo.UserRepository, deviceRepo *repo.DeviceRepository) {
	handler := &UsersHandler{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
	}

	group := e.Group("/user", jwtMiddleware.MiddlewareFunc())
	{
		group.GET("/gateway", handler.gateway)
	}
}

func (h *UsersHandler) gateway(c *gin.Context) {
	user, err := h.userRepo.GetById(middleware.GetUserIdFromContext(c))
	if err != nil {
		c.Error(err)
		return
	}

	gateway.NewUserClient(c.Writer, c.Request, user)
}

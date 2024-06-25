package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/controller/gateway"
	"github.com/pc-power-api/src/controller/middleware"
	"github.com/pc-power-api/src/exceptions"
	"github.com/pc-power-api/src/infra/repo"
	"net/http"
)

var UserDoesNotOwnDevice = exceptions.NewNoAccess("The user does not own this device")

type DevicesHandler struct {
	deviceGateway *gateway.DeviceGatewayHandler
	deviceRepo    *repo.DeviceRepository
	userRepo      *repo.UserRepository
}

func NewDevicesHandler(e *gin.Engine, jwtMiddleware *jwt.GinJWTMiddleware, deviceRepo *repo.DeviceRepository, userRepo *repo.UserRepository) {
	handler := &DevicesHandler{
		deviceGateway: gateway.NewDeviceGatewayHandler(),
		deviceRepo:    deviceRepo,
		userRepo:      userRepo,
	}

	group := e.Group("/devices")
	{
		group.GET("/gateway", handler.gateway)
		group.POST("/power-switch", jwtMiddleware.MiddlewareFunc(), handler.pressPowerSwitch)
	}
}

func (h *DevicesHandler) gateway(c *gin.Context) {
	var data *api.DeviceIdentify
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(errors.New(err))
		return
	}

	device, aerr := h.deviceRepo.GetByIdAndSecret(data)
	if aerr != nil {
		c.Error(aerr)
		return
	}

	h.deviceGateway.DeviceHandler(c.Writer, c.Request, device.Code)
}

func (h *DevicesHandler) pressPowerSwitch(c *gin.Context) {
	var data *api.UserCommand
	err := c.ShouldBind(&data)
	if err != nil {
		c.Error(errors.New(err))
		return
	}

	jwtUser, _ := c.Get(jwt.IdentityKey)
	user, aerr := h.userRepo.GetById(jwtUser.(*middleware.JwtUser).ID)
	if aerr != nil {
		c.Error(aerr)
		return
	}

	if !user.HasDevice(data.DeviceId) {
		c.Error(errors.New(UserDoesNotOwnDevice))
		return
	}

	aerr = h.deviceGateway.PressPowerSwitch(data.DeviceId, data.Hard)
	if aerr != nil {
		c.Error(aerr)
		return
	}
	c.Status(http.StatusOK)
}

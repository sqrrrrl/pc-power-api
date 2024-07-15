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

const DeviceCodeLength = 6
const DeviceSecretLength = 16
const IdPathParam = "id"

var UserDoesNotOwnDevice = exceptions.NewNoAccess("The user does not own this device")
var DeviceNotConnectedError = exceptions.NewDeviceUnreachable("the device is not online")

type DevicesHandler struct {
	deviceRepo *repo.DeviceRepository
	userRepo   *repo.UserRepository
}

func NewDevicesHandler(e *gin.Engine, jwtMiddleware *jwt.GinJWTMiddleware, deviceRepo *repo.DeviceRepository, userRepo *repo.UserRepository) {
	handler := &DevicesHandler{
		deviceRepo: deviceRepo,
		userRepo:   userRepo,
	}

	group := e.Group("/devices")
	{
		group.GET("/gateway", handler.gateway)
		group.POST("/power-switch", jwtMiddleware.MiddlewareFunc(), handler.pressPowerSwitch)
		group.POST("/reset-switch", jwtMiddleware.MiddlewareFunc(), handler.pressResetSwitch)
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

	gateway.NewDeviceClient(c.Writer, c.Request, device)
}

func (h *DevicesHandler) pressPowerSwitch(c *gin.Context) {
	var data *api.UserCommand
	err := c.ShouldBind(&data)
	if err != nil {
		c.Error(errors.New(err))
		return
	}

	user, aerr := h.userRepo.GetById(middleware.GetUserIdFromContext(c))
	if aerr != nil {
		c.Error(aerr)
		return
	}

	if !user.HasDevice(data.DeviceID) {
		c.Error(errors.New(UserDoesNotOwnDevice))
		return
	}

	if deviceClient, ok := gateway.ConnectedDevices[data.DeviceID]; ok {
		aerr = deviceClient.PressPowerSwitch(data.Hard)
		if aerr != nil {
			c.Error(aerr)
			return
		}
	} else {
		c.Error(errors.New(DeviceNotConnectedError))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DevicesHandler) pressResetSwitch(c *gin.Context) {
	var data *api.UserCommand
	err := c.ShouldBind(&data)
	if err != nil {
		c.Error(errors.New(err))
		return
	}

	user, aerr := h.userRepo.GetById(middleware.GetUserIdFromContext(c))
	if aerr != nil {
		c.Error(aerr)
		return
	}

	if !user.HasDevice(data.DeviceID) {
		c.Error(errors.New(UserDoesNotOwnDevice))
		return
	}

	if deviceClient, ok := gateway.ConnectedDevices[data.DeviceID]; ok {
		aerr = deviceClient.PressResetSwitch()
		if aerr != nil {
			c.Error(aerr)
			return
		}
	} else {
		c.Error(errors.New(DeviceNotConnectedError))
		return
	}
	c.Status(http.StatusNoContent)
}

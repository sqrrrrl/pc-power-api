package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/controller/gateway"
	"github.com/pc-power-api/src/controller/middleware"
	"github.com/pc-power-api/src/exceptions"
	"github.com/pc-power-api/src/infra/entity"
	"github.com/pc-power-api/src/infra/repo"
	"github.com/pc-power-api/src/util"
	"net/http"
)

const DeviceCodeLength = 6
const DeviceSecretLength = 16
const IdPathParam = "id"

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
		group.POST("/reset-switch", jwtMiddleware.MiddlewareFunc(), handler.pressResetSwitch)
		group.POST("/", jwtMiddleware.MiddlewareFunc(), handler.createDevice)
		group.DELETE("/:"+IdPathParam, jwtMiddleware.MiddlewareFunc(), handler.deleteDevice)
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

	user, aerr := h.userRepo.GetById(middleware.GetUserIdFromContext(c))
	if aerr != nil {
		c.Error(aerr)
		return
	}

	if !user.HasDevice(data.DeviceCode) {
		c.Error(errors.New(UserDoesNotOwnDevice))
		return
	}

	aerr = h.deviceGateway.PressPowerSwitch(data.DeviceCode, data.Hard)
	if aerr != nil {
		c.Error(aerr)
		return
	}
	c.Status(http.StatusOK)
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

	if !user.HasDevice(data.DeviceCode) {
		c.Error(errors.New(UserDoesNotOwnDevice))
		return
	}

	aerr = h.deviceGateway.PressResetSwitch(data.DeviceCode)
	if aerr != nil {
		c.Error(aerr)
		return
	}
	c.Status(http.StatusOK)
}

func (h *DevicesHandler) createDevice(c *gin.Context) {
	var deviceInfo *api.DeviceInfo
	err := c.ShouldBind(&deviceInfo)
	if err != nil {
		c.Error(errors.New(err))
		return
	}

	ownerId := middleware.GetUserIdFromContext(c)
	deviceUuid := uuid.New()
	deviceCode := util.GenerateRandomString(DeviceCodeLength)
	deviceSecret := util.GenerateRandomString(DeviceSecretLength)

	var device = entity.Device{
		ID:     deviceUuid.String(),
		Name:   deviceInfo.Name,
		Code:   deviceCode,
		Secret: deviceSecret,
		Status: 0,
		UserID: ownerId,
	}

	aerr := h.deviceRepo.Create(&device)
	if aerr != nil {
		c.Error(aerr)
		return
	}

	c.JSON(200, api.DeviceInfo{
		ID:     device.ID,
		Name:   device.Name,
		Code:   device.Code,
		Secret: device.Secret,
		Status: device.Status,
	})
}

func (h *DevicesHandler) deleteDevice(c *gin.Context) {
	deviceId := c.Param(IdPathParam)
	device, aerr := h.deviceRepo.GetById(deviceId)
	if aerr != nil {
		c.Error(aerr)
		return
	}

	ownerId := middleware.GetUserIdFromContext(c)
	if ownerId != device.UserID {
		c.Error(errors.New(UserDoesNotOwnDevice))
		return
	}

	aerr = h.deviceRepo.Delete(device)
	if aerr != nil {
		c.Error(aerr)
		return
	}

	c.Status(http.StatusOK)
}

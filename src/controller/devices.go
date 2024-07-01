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
		group.GET("/", jwtMiddleware.MiddlewareFunc(), handler.getDevices)
		group.POST("/", jwtMiddleware.MiddlewareFunc(), handler.createDevice)
		group.GET("/:"+IdPathParam, jwtMiddleware.MiddlewareFunc(), handler.getDevice)
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

	gateway.NewDeviceClient(c.Writer, c.Request, device.Code)
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

	if deviceClient, ok := gateway.ConnectedClients[data.DeviceCode]; ok {
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

	if !user.HasDevice(data.DeviceCode) {
		c.Error(errors.New(UserDoesNotOwnDevice))
		return
	}

	if deviceClient, ok := gateway.ConnectedClients[data.DeviceCode]; ok {
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

func (h *DevicesHandler) getDevices(c *gin.Context) {
	userId := middleware.GetUserIdFromContext(c)
	user, aerr := h.userRepo.GetById(userId)
	if aerr != nil {
		c.Error(aerr)
		return
	}

	devicesInfoList := api.DeviceInfoList{
		OnlineDevices:  make([]api.DeviceInfo, 0),
		OfflineDevices: make([]api.DeviceInfo, 0),
	}
	for _, device := range user.Devices {
		if conn, ok := gateway.ConnectedClients[device.Code]; ok {
			devicesInfoList.OnlineDevices = append(devicesInfoList.OnlineDevices, api.DeviceInfo{
				ID:     device.ID,
				Name:   device.Name,
				Code:   device.Code,
				Secret: device.Secret,
				Status: conn.GetStatus(),
				Online: true,
			})
		} else {
			devicesInfoList.OfflineDevices = append(devicesInfoList.OfflineDevices, api.DeviceInfo{
				ID:     device.ID,
				Name:   device.Name,
				Code:   device.Code,
				Secret: device.Secret,
				Status: 0,
				Online: false,
			})
		}
	}

	c.JSON(http.StatusOK, devicesInfoList)
}

func (h *DevicesHandler) getDevice(c *gin.Context) {
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

	status := 0
	online := false
	if conn, ok := gateway.ConnectedClients[device.Code]; ok {
		status = conn.GetStatus()
		online = true
	}
	c.JSON(http.StatusOK, api.DeviceInfo{
		ID:     device.ID,
		Name:   device.Name,
		Code:   device.Code,
		Secret: device.Secret,
		Status: status,
		Online: online,
	})
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
		UserID: ownerId,
	}

	aerr := h.deviceRepo.Create(&device)
	if aerr != nil {
		c.Error(aerr)
		return
	}

	c.JSON(http.StatusOK, api.DeviceInfo{
		ID:     device.ID,
		Name:   device.Name,
		Code:   device.Code,
		Secret: device.Secret,
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

	c.Status(http.StatusNoContent)
}

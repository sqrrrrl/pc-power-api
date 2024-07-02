package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/controller/gateway"
	"github.com/pc-power-api/src/controller/middleware"
	"github.com/pc-power-api/src/infra/entity"
	"github.com/pc-power-api/src/infra/repo"
	"github.com/pc-power-api/src/pubsub"
	"github.com/pc-power-api/src/util"
	"net/http"
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

		deviceGroup := group.Group("/devices")

		deviceGroup.POST("/", handler.createDevice)
		deviceGroup.GET("/", handler.getDevices)
		deviceGroup.GET("/:"+IdPathParam, handler.getDevice)
		deviceGroup.PUT("/:"+IdPathParam, handler.updateDevice)
		deviceGroup.DELETE("/:"+IdPathParam, handler.deleteDevice)
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

func (h *UsersHandler) getDevices(c *gin.Context) {
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
		if conn, ok := gateway.ConnectedDevices[device.ID]; ok {
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

func (h *UsersHandler) getDevice(c *gin.Context) {
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
	if conn, ok := gateway.ConnectedDevices[device.ID]; ok {
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

func (h *UsersHandler) createDevice(c *gin.Context) {
	var deviceInfo *api.DeviceCreateInfo
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
	pubsub.Publish(ownerId, device)

	c.JSON(http.StatusOK, api.DeviceInfo{
		ID:     device.ID,
		Name:   device.Name,
		Code:   device.Code,
		Secret: device.Secret,
	})
}

func (h *UsersHandler) updateDevice(c *gin.Context) {
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

	var deviceInfo *api.DeviceCreateInfo
	err := c.ShouldBind(&deviceInfo)
	if err != nil {
		c.Error(errors.New(err))
		return
	}

	device.Name = deviceInfo.Name
	aerr = h.deviceRepo.Update(device)
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

func (h *UsersHandler) deleteDevice(c *gin.Context) {
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

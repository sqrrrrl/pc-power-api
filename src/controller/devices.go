package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/controller/gateway"
	"github.com/pc-power-api/src/infra/repo"
	"net/http"
)

type DevicesHandler struct {
	deviceGateway *gateway.DeviceGatewayHandler
	deviceRepo    *repo.DeviceRepository
}

func NewDevicesHandler(e *gin.Engine, deviceRepo *repo.DeviceRepository) {
	handler := &DevicesHandler{
		deviceGateway: gateway.NewDeviceGatewayHandler(),
		deviceRepo:    deviceRepo,
	}

	group := e.Group("/devices")
	{
		group.GET("/gateway", handler.gateway)
		group.POST("/power-switch", handler.pressPowerSwitch)
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
	//TODO: check if user has permission to control this device
	var data *api.UserCommand
	err := c.ShouldBind(&data)
	if err != nil {
		c.Error(errors.New(err))
		return
	}
	aerr := h.deviceGateway.PressPowerSwitch(data.DeviceId, data.Hard)
	if aerr != nil {
		c.Error(aerr)
		return
	}
	c.Status(http.StatusOK)
}

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/controller/gateway"
	"net/http"
)

type DevicesHandler struct {
	deviceGateway *gateway.DeviceGatewayHandler
}

func NewDevicesHandler(e *gin.Engine) {
	handler := &DevicesHandler{
		deviceGateway: gateway.NewDeviceGatewayHandler(),
	}

	group := e.Group("/devices")
	{
		group.GET("/gateway", handler.gateway)
		group.POST("/power-switch", handler.pressPowerSwitch)
	}
}

func (h *DevicesHandler) gateway(c *gin.Context) {
	//TODO: get the device id from the authenticated device
	h.deviceGateway.DeviceHandler(c.Writer, c.Request, "1")
}

func (h *DevicesHandler) pressPowerSwitch(c *gin.Context) {
	//TODO: check if user has permission to control this device
	var data *api.UserCommand
	err := c.ShouldBind(&data)
	if err != nil {
		c.Error(errors.New(err))
		return
	}
	err = h.deviceGateway.PressPowerSwitch(data.DeviceId, data.Hard)
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

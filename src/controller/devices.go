package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pc-power-api/src/controller/gateway"
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
	}
}

func (h *DevicesHandler) gateway(c *gin.Context) {
	h.deviceGateway.DeviceHandler(c.Writer, c.Request, "1")
}

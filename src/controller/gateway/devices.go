package gateway

import (
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pc-power-api/src/api/gateway"
	"github.com/pc-power-api/src/exceptions"
	"github.com/pc-power-api/src/util"
	"net/http"
)

var DeviceNotConnectedError = exceptions.NewDeviceUnreachable("the device is not online")
var FailedToCommunicateWithDeviceError = exceptions.NewDeviceUnreachable("the communication with the device failed")

const InvalidMessageTitle = "The message is invalid"
const InvalidMessageDescription = "The message is not valid json or is not following the schema"
const GatewayType = "device"

type DeviceGatewayHandler struct {
	conns map[string]*websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewDeviceGatewayHandler() *DeviceGatewayHandler {
	return &DeviceGatewayHandler{
		conns: make(map[string]*websocket.Conn),
	}
}

func (h *DeviceGatewayHandler) DeviceHandler(w http.ResponseWriter, r *http.Request, deviceId string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.conns[deviceId] = conn
	h.listen(conn)
	delete(h.conns, deviceId)
	conn.Close()
}

func (h *DeviceGatewayHandler) listen(conn *websocket.Conn) {
	for {
		var data gateway.DeviceMessage
		err := conn.ReadJSON(&data)
		if err != nil {
			var closeError *websocket.CloseError
			if errors.As(err, &closeError) {
				break
			}
			h.handleError(errors.New(err), conn)
		}
	}
}

func (h *DeviceGatewayHandler) handleError(err *errors.Error, conn *websocket.Conn) {
	id := uuid.New()
	message := gateway.ErrorMessage{}
	message.SetId(id.String())
	message.SetMessage(err.Error())
	message.SetTitle(InvalidMessageTitle)
	message.SetDescription(InvalidMessageDescription)
	conn.WriteJSON(message)
	util.LogWebsocketError(err, id, conn, GatewayType)
}

func (h *DeviceGatewayHandler) PressPowerSwitch(deviceId string, hardPowerOff bool) *errors.Error {
	if conn, ok := h.conns[deviceId]; ok {
		var op int
		if hardPowerOff {
			op = gateway.HardPowerOffOpcode
		} else {
			op = gateway.PressPowerSwitchOpcode
		}
		message := gateway.CommandMessage{
			Opcode: op,
		}
		err := conn.WriteJSON(message)
		if err != nil {
			return errors.New(FailedToCommunicateWithDeviceError)
		}
	} else {
		return errors.New(DeviceNotConnectedError)
	}
	return nil
}

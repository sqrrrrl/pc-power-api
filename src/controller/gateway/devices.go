package gateway

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pc-power-api/src/api/gateway"
	"github.com/pc-power-api/src/controller/middleware"
	"github.com/pc-power-api/src/exceptions"
	"github.com/pc-power-api/src/util"
	"net"
	"net/http"
	"sync"
	"time"
)

var FailedToCommunicateWithDeviceError = exceptions.NewDeviceUnreachable("the communication with the device failed")

const InvalidMessageTitle = "The message is invalid"
const InvalidMessageDescription = "The message is not valid json or is not following the schema"
const NewSessionOpenedTitle = "Another session has been opened"
const NewSessionOpenedDescription = "Another session has been opened, this one will be closed"
const GatewayType = "device"
const PingPeriod = 2 * time.Minute
const PongWait = PingPeriod + time.Minute

var ConnectedClients = make(map[string]*DeviceClient)

type DeviceClient struct {
	conn       *websocket.Conn
	status     int
	deviceCode string
	writeMu    sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewDeviceClient(w http.ResponseWriter, r *http.Request, deviceCode string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conn.SetReadDeadline(time.Now().Add(PongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })

	client := &DeviceClient{
		conn:       conn,
		status:     0,
		deviceCode: deviceCode,
		writeMu:    sync.Mutex{},
	}
	if ConnectedClients[deviceCode] != nil {
		ConnectedClients[deviceCode].handleError(errors.New(NewSessionOpenedDescription), NewSessionOpenedTitle, NewSessionOpenedDescription)
		ConnectedClients[deviceCode].destroy()
	}
	ConnectedClients[deviceCode] = client

	go client.listen()
	go client.sendPing()
}

func (c *DeviceClient) listen() {
	for c.conn != nil {
		var data gateway.DeviceMessage
		err := c.conn.ReadJSON(&data)
		if err != nil {
			var closeError *websocket.CloseError
			var timeoutError net.Error
			var jsonTypeError *json.UnmarshalTypeError
			var jsonSyntaxError *json.SyntaxError
			if errors.As(err, &closeError) {
				c.destroy()
			} else if errors.As(err, &timeoutError) && timeoutError.Timeout() {
				c.destroy()
			} else if errors.As(err, &jsonTypeError) || errors.As(err, &jsonSyntaxError) {
				c.handleError(errors.New(err), InvalidMessageTitle, InvalidMessageDescription)
			} else {
				c.handleError(errors.New(err))
			}
		} else {
			c.status = data.Status
		}
	}
}

func (c *DeviceClient) sendPing() {
	ticker := time.NewTicker(PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.writeMu.Lock()
			if c.conn == nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.destroy()
			}
			c.writeMu.Unlock()
		}
	}
}

func (c *DeviceClient) GetStatus() int {
	return c.status
}

func (c *DeviceClient) destroy() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	delete(ConnectedClients, c.deviceCode)
}

func (c *DeviceClient) handleError(err *errors.Error, info ...string) {
	id := uuid.New()
	errorTitle := middleware.UnexpectedErrorTitle
	errorDescription := middleware.UnexpectedErrorDescription
	errorMsg := ""
	if len(info) > 0 {
		errorTitle = info[0]
		errorDescription = info[1]
		errorMsg = err.Error()
	}
	message := gateway.ErrorMessage{}
	message.SetId(id.String())
	message.SetMessage(errorMsg)
	message.SetTitle(errorTitle)
	message.SetDescription(errorDescription)
	c.conn.WriteJSON(message)
	util.LogWebsocketError(err, id, c.conn, GatewayType)
}

func (c *DeviceClient) PressPowerSwitch(hardPowerOff bool) *errors.Error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	var op int
	if hardPowerOff {
		op = gateway.HardPowerOffOpcode
	} else {
		op = gateway.PressPowerSwitchOpcode
	}
	message := gateway.CommandMessage{
		Opcode: op,
	}
	err := c.conn.WriteJSON(message)
	if err != nil {
		c.destroy()
		return errors.New(FailedToCommunicateWithDeviceError)
	}
	return nil
}

func (c *DeviceClient) PressResetSwitch() *errors.Error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	message := gateway.CommandMessage{
		Opcode: gateway.PressResetSwitchOpcode,
	}
	err := c.conn.WriteJSON(message)
	if err != nil {
		c.destroy()
		return errors.New(FailedToCommunicateWithDeviceError)
	}
	return nil
}

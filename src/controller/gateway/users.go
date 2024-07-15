package gateway

import (
	"github.com/gorilla/websocket"
	"github.com/pc-power-api/src/infra/entity"
	"github.com/pc-power-api/src/pubsub"
	"net/http"
)

type UserClient struct {
	conn *websocket.Conn
	user *entity.User
}

func NewUserClient(w http.ResponseWriter, r *http.Request, user *entity.User) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &UserClient{
		conn: conn,
		user: user,
	}
	pubsub.Subscribe(client)
	conn.SetCloseHandler(func(code int, text string) error {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, text))
		client.destroy()
		return nil
	})
	go client.listen()
}

// Messages need to be read for the CloseHandler to be called
func (c *UserClient) listen() {
	for c.conn != nil {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			c.destroy()
		}
	}
}

func (c *UserClient) destroy() {
	pubsub.Unsubscribe(c)
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *UserClient) Notify(topic string, data interface{}) {
	if topic == c.user.ID {
		c.user.Devices = append(c.user.Devices, data.(entity.Device))
	} else if c.user.HasDevice(topic) {
		c.conn.WriteJSON(data)
	}
}

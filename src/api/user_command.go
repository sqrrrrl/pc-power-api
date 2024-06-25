package api

type UserCommand struct {
	DeviceCode string `json:"device_id"`
	Hard       bool   `json:"hard"`
}

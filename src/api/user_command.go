package api

type UserCommand struct {
	DeviceId string `json:"device_id"`
	Hard     bool   `json:"hard"`
}

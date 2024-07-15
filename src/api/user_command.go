package api

type UserCommand struct {
	DeviceID string `json:"device_id" binding:"required,uuid"`
	Hard     bool   `json:"hard"`
}

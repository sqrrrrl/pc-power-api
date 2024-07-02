package gateway

type DeviceMessage struct {
	Status int `json:"status" binding:"required,oneof=0 1"`
}

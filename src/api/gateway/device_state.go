package gateway

type DeviceState struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Online bool   `json:"online"`
}

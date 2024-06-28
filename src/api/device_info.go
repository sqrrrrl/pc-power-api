package api

type DeviceInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Secret string `json:"secret"`
	Status int    `json:"status"`
}

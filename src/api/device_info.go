package api

type DeviceInfo struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Secret string `json:"secret"`
	Status int    `json:"status"`
}

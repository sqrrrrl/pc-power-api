package api

type DeviceCreateInfo struct {
	Name string `json:"name" binding:"required,min=1,max=32"`
}

type DeviceInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Secret string `json:"secret"`
	Status int    `json:"status"`
	Online bool   `json:"online"`
}

type DeviceInfoList struct {
	OnlineDevices  []DeviceInfo `json:"online"`
	OfflineDevices []DeviceInfo `json:"offline"`
}

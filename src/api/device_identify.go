package api

type DeviceIdentify struct {
	Code   string `form:"device_id"`
	Secret string `form:"secret"`
}

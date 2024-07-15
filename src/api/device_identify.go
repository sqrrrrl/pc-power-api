package api

type DeviceIdentify struct {
	Code   string `form:"device_id" binding:"required,len=6"`
	Secret string `form:"secret" binding:"required,len=16"`
}

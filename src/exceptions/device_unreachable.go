package exceptions

type DeviceUnreachable struct {
	Message string
}

func NewDeviceUnreachable(message string) *DeviceUnreachable {
	return &DeviceUnreachable{
		Message: message,
	}
}

func (e *DeviceUnreachable) Error() string {
	return e.Message
}

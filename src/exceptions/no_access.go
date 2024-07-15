package exceptions

type NoAccess struct {
	Message string
}

func NewNoAccess(message string) *NoAccess {
	return &NoAccess{
		Message: message,
	}
}

func (e *NoAccess) Error() string {
	return e.Message
}

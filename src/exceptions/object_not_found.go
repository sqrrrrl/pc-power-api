package exceptions

type ObjectNotFound struct {
	Message string
}

func NewObjectNotFound(message string) *ObjectNotFound {
	return &ObjectNotFound{
		Message: message,
	}
}

func (e *ObjectNotFound) Error() string {
	return e.Message
}

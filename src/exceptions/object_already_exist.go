package exceptions

type ObjectAlreadyExist struct {
	Message string
}

func NewObjectAlreadyExist(message string) *ObjectAlreadyExist {
	return &ObjectAlreadyExist{
		Message: message,
	}
}

func (e *ObjectAlreadyExist) Error() string {
	return e.Message
}

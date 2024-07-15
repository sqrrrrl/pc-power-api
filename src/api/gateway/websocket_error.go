package gateway

type ErrorMessage struct {
	Details ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Message     string `json:"message"`
}

func (err *ErrorMessage) SetId(id string) {
	err.Details.Id = id
}

func (err *ErrorMessage) SetTitle(title string) {
	err.Details.Title = title
}

func (err *ErrorMessage) SetDescription(description string) {
	err.Details.Description = description
}

func (err *ErrorMessage) SetMessage(message string) {
	err.Details.Message = message
}

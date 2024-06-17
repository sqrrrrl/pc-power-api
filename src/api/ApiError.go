package api

type ErrorResponse struct {
	Details ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Id          string `json:"id"`
	Status      int    `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Message     string `json:"message"`
	Expected    bool   `json:"expected"`
}

func (err *ErrorResponse) SetId(id string) {
	err.Details.Id = id
}

func (err *ErrorResponse) SetStatus(status int) {
	err.Details.Status = status
}

func (err *ErrorResponse) SetTitle(title string) {
	err.Details.Title = title
}

func (err *ErrorResponse) SetDescription(description string) {
	err.Details.Description = description
}

func (err *ErrorResponse) SetMessage(message string) {
	err.Details.Message = message
}

func (err *ErrorResponse) SetExpected(expected bool) {
	err.Details.Expected = expected
}

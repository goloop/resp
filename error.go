package resp

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code    int    `json:"code"`    // error code
	Message string `json:"message"` // error message
}

// Unpack returns the error code and message.
func (e *ErrorResponse) Unpack() (code int, message string) {
	return e.Code, e.Message
}

// newErrorResponse creates a new errorMessage object with the
// given code and message. If a message is provided, it will be
// used as the error message. Otherwise, the default message
// associated with the given status code will be used.
func newErrorResponse(status int, message ...string) *ErrorResponse {
	msg := statusMessages[status]
	if len(message) > 0 {
		msg = message[0]
	}

	return &ErrorResponse{
		Code:    status,
		Message: msg,
	}
}

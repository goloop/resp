package resp

// errorMessage represents an error response.
type errorMessage struct {
	Code    int    `json:"code"`    // error code
	Message string `json:"message"` // error message
}

// newErrorMessage creates a new errorMessage object with the
// given code and message. If a message is provided, it will be
// used as the error message. Otherwise, the default message
// associated with the given status code will be used.
func newErrorMessage(status int, message ...string) *errorMessage {
	msg := statusMessages[status]
	if len(message) > 0 {
		msg = message[0]
	}

	return &errorMessage{
		Code:    status,
		Message: msg,
	}
}

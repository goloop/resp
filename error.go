package resp

import "errors"

// ErrInvalidJSONPCallback is returned by JSONP when the supplied callback
// name fails validation. It can be matched with errors.Is to distinguish a
// rejected (potentially malicious) callback from an encoding or write error.
var ErrInvalidJSONPCallback = errors.New("invalid JSONP callback name")

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code int `json:"code"` // error code

	// Error is a machine-readable name for what went wrong, such as
	// "rate_limit" or "slug_taken". It is what a client switches on and what
	// a translation key is looked up by; a number cannot do either without a
	// registry that nobody maintains, and "1042" says nothing in a log or a
	// test either.
	//
	// It is omitted when unset, so a response that does not name one is
	// exactly the bytes it was before this field existed. See WithErrorSlug.
	Error string `json:"error,omitempty"`

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

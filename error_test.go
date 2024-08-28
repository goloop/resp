package resp

import "testing"

// TestNewErrorMessage tests the newErrorMessage function.
func TestNewErrorMessage(t *testing.T) {
	tests := []struct {
		status   int
		message  string
		expected *ErrorResponse
	}{
		{
			status:  200,
			message: "OK",
			expected: &ErrorResponse{
				Code:    200,
				Message: "OK",
			},
		},
		{
			status:  400,
			message: "Bad Request",
			expected: &ErrorResponse{
				Code:    400,
				Message: "Bad Request",
			},
		},
		{
			status:  500,
			message: "Internal Server Error",
			expected: &ErrorResponse{
				Code:    500,
				Message: "Internal Server Error",
			},
		},
	}

	for _, test := range tests {
		result := newErrorResponse(test.status, test.message)
		if result.Code != test.expected.Code {
			t.Errorf("newErrorMessage() Code = %d, want %d", result.Code, test.expected.Code)
		}
		if result.Message != test.expected.Message {
			t.Errorf("newErrorMessage() Message = %s, want %s", result.Message, test.expected.Message)
		}
	}
}
func TestErrorResponse_Unpack(t *testing.T) {
	err := &ErrorResponse{
		Code:    200,
		Message: "OK",
	}

	code, message := err.Unpack()

	if code != 200 {
		t.Errorf("Unpack() code = %d, want %d", code, 200)
	}

	if message != "OK" {
		t.Errorf("Unpack() message = %s, want %s", message, "OK")
	}
}

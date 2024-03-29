package resp

import "testing"

// TestNewErrorMessage tests the newErrorMessage function.
func TestNewErrorMessage(t *testing.T) {
	tests := []struct {
		status   int
		message  string
		expected *errorMessage
	}{
		{
			status:  200,
			message: "OK",
			expected: &errorMessage{
				Code:    200,
				Message: "OK",
			},
		},
		{
			status:  400,
			message: "Bad Request",
			expected: &errorMessage{
				Code:    400,
				Message: "Bad Request",
			},
		},
		{
			status:  500,
			message: "Internal Server Error",
			expected: &errorMessage{
				Code:    500,
				Message: "Internal Server Error",
			},
		},
	}

	for _, test := range tests {
		result := newErrorMessage(test.status, test.message)
		if result.Code != test.expected.Code {
			t.Errorf("newErrorMessage() Code = %d, want %d", result.Code, test.expected.Code)
		}
		if result.Message != test.expected.Message {
			t.Errorf("newErrorMessage() Message = %s, want %s", result.Message, test.expected.Message)
		}
	}
}

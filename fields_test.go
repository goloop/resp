package resp

import (
	"reflect"
	"testing"
)

type User struct {
	ID       int
	Email    string
	Password string
	IsActive bool
}

// TestOnlyFields tests the OnlyFields function.
func TestOnlyFields(t *testing.T) {
	user := User{
		ID:       1,
		Email:    "user@example.com",
		Password: "secret",
		IsActive: true,
	}

	expected := R{
		"ID":       1,
		"Email":    "user@example.com",
		"IsActive": true,
	}

	result := OnlyFields(user, "ID", "Email", "IsActive")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OnlyFields() = %v, want %v", result, expected)
	}
}

// TestExcludeFields tests the ExcludeFields function.
func TestExcludeFields(t *testing.T) {
	user := User{
		ID:       1,
		Email:    "user@example.com",
		Password: "secret",
		IsActive: true,
	}

	expected := R{
		"ID":       1,
		"Email":    "user@example.com",
		"IsActive": true,
	}

	result := ExcludeFields(user, "Password")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExcludeFields() = %v, want %v", result, expected)
	}
}

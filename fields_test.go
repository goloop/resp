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

// TestOnlyFieldsSlice tests the OnlyFields function with a slice.
func TestOnlyFieldsSlice(t *testing.T) {
	users := []User{
		{
			ID:       1,
			Email:    "user_a@example.com",
			Password: "secret",
			IsActive: true,
		},
		{
			ID:       2,
			Email:    "user_b@example.com",
			Password: "secret",
			IsActive: true,
		},
	}

	expected := []R{
		{
			"ID":       1,
			"Email":    "user_a@example.com",
			"IsActive": true,
		},
		{
			"ID":       2,
			"Email":    "user_b@example.com",
			"IsActive": true,
		},
	}

	result := OnlyFields(users, "ID", "Email", "IsActive")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OnlyFields() = %v, want %v", result, expected)
	}
}

// TestExcludeFieldsSlice tests the ExcludeFields function with a slice.
func TestExcludeFieldsSlice(t *testing.T) {
	users := []User{
		{
			ID:       1,
			Email:    "user_a@example.com",
			Password: "secret",
			IsActive: true,
		},
		{
			ID:       2,
			Email:    "user_b@example.com",
			Password: "secret",
			IsActive: true,
		},
	}

	expected := []R{
		{
			"ID":       1,
			"Email":    "user_a@example.com",
			"IsActive": true,
		},
		{
			"ID":       2,
			"Email":    "user_b@example.com",
			"IsActive": true,
		},
	}

	result := ExcludeFields(users, "Password")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExcludeFields() = %v, want %v", result, expected)
	}
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

// TestOnlyFieldsPointer tests the OnlyFields function with a pointer to a struct.
func TestOnlyFieldsPointer(t *testing.T) {
	user := &User{
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

// TestOnlyFieldsSlicePointer tests the OnlyFields function with a slice of pointers to structs.
func TestOnlyFieldsSlicePointer(t *testing.T) {
	users := []*User{
		{
			ID:       1,
			Email:    "user_a@example.com",
			Password: "secret",
			IsActive: true,
		},
		{
			ID:       2,
			Email:    "user_b@example.com",
			Password: "secret",
			IsActive: true,
		},
	}

	expected := []R{
		{
			"ID":       1,
			"Email":    "user_a@example.com",
			"IsActive": true,
		},
		{
			"ID":       2,
			"Email":    "user_b@example.com",
			"IsActive": true,
		},
	}

	result := OnlyFields(users, "ID", "Email", "IsActive")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OnlyFields() = %v, want %v", result, expected)
	}
}

// TestOnlyFieldsNonStruct tests the OnlyFields function with non-struct data.
func TestOnlyFieldsNonStruct(t *testing.T) {
	// Test with a string
	inputStr := "not a struct"
	result := OnlyFields(inputStr, "field_a", "field_b")
	if result != inputStr {
		t.Errorf("OnlyFields() with string = %v, want %v", result, inputStr)
	}

	// Test with an int
	inputInt := 42
	result = OnlyFields(inputInt, "field_a", "field_b")
	if result != inputInt {
		t.Errorf("OnlyFields() with int = %v, want %v", result, inputInt)
	}

	// Test with a slice of strings
	inputSlice := []string{"a", "b", "c"}
	result = OnlyFields(inputSlice, "field_a", "field_b")
	if !reflect.DeepEqual(result, inputSlice) {
		t.Errorf("OnlyFields() with slice of strings = %v, want %v", result, inputSlice)
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

// TestExcludeFieldsPointer tests the ExcludeFields function with a pointer to a struct.
func TestExcludeFieldsPointer(t *testing.T) {
	user := &User{
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

// TestExcludeFieldsSlicePointer tests the ExcludeFields function with a slice of pointers to structs.
func TestExcludeFieldsSlicePointer(t *testing.T) {
	users := []*User{
		{
			ID:       1,
			Email:    "user_a@example.com",
			Password: "secret",
			IsActive: true,
		},
		{
			ID:       2,
			Email:    "user_b@example.com",
			Password: "secret",
			IsActive: true,
		},
	}

	expected := []R{
		{
			"ID":       1,
			"Email":    "user_a@example.com",
			"IsActive": true,
		},
		{
			"ID":       2,
			"Email":    "user_b@example.com",
			"IsActive": true,
		},
	}

	result := ExcludeFields(users, "Password")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExcludeFields() = %v, want %v", result, expected)
	}
}

// TestExcludeFieldsNonStruct tests the ExcludeFields function with non-struct data.
func TestExcludeFieldsNonStruct(t *testing.T) {
	// Test with a string
	inputStr := "not a struct"
	result := ExcludeFields(inputStr, "field_a", "field_b")
	if result != inputStr {
		t.Errorf("ExcludeFields() with string = %v, want %v", result, inputStr)
	}

	// Test with an int
	inputInt := 42
	result = ExcludeFields(inputInt, "field_a", "field_b")
	if result != inputInt {
		t.Errorf("ExcludeFields() with int = %v, want %v", result, inputInt)
	}

	// Test with a slice of strings
	inputSlice := []string{"a", "b", "c"}
	result = ExcludeFields(inputSlice, "field_a", "field_b")
	if !reflect.DeepEqual(result, inputSlice) {
		t.Errorf("ExcludeFields() with slice of strings = %v, want %v", result, inputSlice)
	}
}

// TestOnlyFieldsMap tests the OnlyFields function with a map.
func TestOnlyFieldsMap(t *testing.T) {
	inputMap := map[string]any{
		"ID":       1,
		"Email":    "user@example.com",
		"Password": "secret",
		"IsActive": true,
	}

	expected := R{
		"ID":       1,
		"Email":    "user@example.com",
		"IsActive": true,
	}

	result := OnlyFields(inputMap, "ID", "Email", "IsActive")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OnlyFields() = %v, want %v", result, expected)
	}
}

// TestExcludeFieldsMap tests the ExcludeFields function with a map.
func TestExcludeFieldsMap(t *testing.T) {
	inputMap := map[string]any{
		"ID":       1,
		"Email":    "user@example.com",
		"Password": "secret",
		"IsActive": true,
	}

	expected := R{
		"ID":       1,
		"Email":    "user@example.com",
		"IsActive": true,
	}

	result := ExcludeFields(inputMap, "Password")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExcludeFields() = %v, want %v", result, expected)
	}
}

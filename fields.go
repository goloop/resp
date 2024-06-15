package resp

import (
	"reflect"
)

// OnlyFields extracts only the specified fields from the provided
// data and returns them as an `R` map. This function is useful
// for creating JSON responses with a subset of the original data.
//
// Parameters:
//   - data: The input data from which fields will be extracted.
//   - fields: A list of field names to include in the resulting map.
//
// Returns:
//   - An `R` map containing only the specified fields from the input data.
//
// Example Usage:
// The following example demonstrates how to use `OnlyFields` to create
// a JSON response that includes only specific fields from a user object.
//
//	import (
//		"net/http"
//		"github.com/goloop/resp"
//	)
//
//	type User struct {
//		ID       int
//		Email    string
//		Password string
//		IsActive bool
//	}
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//		user := User{
//			ID:       1,
//			Email:    "user@example.com",
//			Password: "secret",
//			IsActive: true,
//		}
//		data := resp.OnlyFields(user, "ID", "Email", "IsActive")
//		if err := resp.JSON(w, data); err != nil {
//			// handle error
//		}
//	}
func OnlyFields(data any, fields ...string) R {
	result := make(R)

	rv := reflect.ValueOf(data)
	rt := rv.Type()

	// for i := 0; i < rv.NumField(); i++ {
	// 	name := rt.Field(i).Name
	// 	if g.In(name, fields...) {
	// 		result[name] = rv.Field(i).Interface()
	// 	}
	// }

	allowed := make(map[string]bool, len(fields))
	for _, field := range fields {
		allowed[field] = true
	}

	for i := 0; i < rv.NumField(); i++ {
		name := rt.Field(i).Name
		if allowed[name] {
			result[name] = rv.Field(i).Interface()
		}
	}

	return result
}

// ExcludeFields removes the specified fields from the provided data
// and returns the remaining fields as an `R` map. This function is
// useful for creating JSON responses without sensitive or unwanted
// fields from the original data.
//
// Parameters:
//   - data: The input data from which fields will be excluded.
//   - fields: A list of field names to exclude from the resulting map.
//
// Returns:
//   - An `R` map containing the fields from the input data except
//     the specified fields.
//
// Example Usage:
// The following example demonstrates how to use `ExcludeFields` to create
// a JSON response that excludes specific fields from a user object.
//
//	import (
//		"net/http"
//		"github.com/goloop/resp"
//	)
//
//	type User struct {
//		ID       int
//		Email    string
//		Password string
//		IsActive bool
//	}
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//		user := User{
//			ID:       1,
//			Email:    "user@example.com",
//			Password: "secret",
//			IsActive: true,
//		}
//		data := resp.ExcludeFields(user, "Password")
//		if err := resp.JSON(w, data); err != nil {
//			// handle error
//		}
//	}
func ExcludeFields(data any, fields ...string) R {
	result := make(R)

	rv := reflect.ValueOf(data)
	rt := rv.Type()

	excluded := make(map[string]bool)
	for _, field := range fields {
		excluded[field] = true
	}

	for i := 0; i < rv.NumField(); i++ {
		name := rt.Field(i).Name
		if !excluded[name] {
			result[name] = rv.Field(i).Interface()
		}
	}

	return result
}

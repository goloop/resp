package resp

import (
	"reflect"
)

// OnlyFields extracts only the specified fields from the provided
// data and returns them as an `R` map. This function is useful
// for creating JSON responses with a subset of the original data.
// The operation can be performed on a single object, a slice of
// objects, an array of objects, or a map. When a slice or an array
// is provided, the function returns a slice of `R` maps. If the data
// is not a struct, slice/array of structs, or map, it returns the
// original data unchanged.
//
// Parameters:
//   - data: The input data from which fields will be extracted.
//   - fields: A list of field names to include in the resulting map.
//
// Returns:
//   - An `R` map containing only the specified fields from the input
//     data, a slice of `R` maps if the input data is a slice or an array,
//     or the original input data unchanged if it is not a struct, a map, or a
//     slice/array of structs.
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
//	func HandlerSingle(w http.ResponseWriter, r *http.Request) {
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
//
//	func HandlerMultiple(w http.ResponseWriter, r *http.Request) {
//		users := []User{
//		    {
//			    ID:       1,
//			    Email:    "user_a@example.com",
//			    Password: "secret",
//			    IsActive: true,
//		    },
//		    {
//			    ID:       2,
//			    Email:    "user_b@example.com",
//			    Password: "secret",
//			    IsActive: true,
//		    },
//		}
//		data := resp.OnlyFields(users, "ID", "Email", "IsActive")
//		if err := resp.JSON(w, data); err != nil {
//			// handle error
//		}
//	}
func OnlyFields(data any, fields ...string) any {
	rv := reflect.ValueOf(data)

	switch rv.Kind() {
	case reflect.Ptr:
		rv = rv.Elem()
		if rv.Kind() == reflect.Struct {
			return onlyFields(rv.Interface(), fields...)
		}
	case reflect.Slice, reflect.Array:
		length := rv.Len()
		if length > 0 {
			elemKind := rv.Index(0).Kind()
			if elemKind == reflect.Ptr {
				elemKind = rv.Index(0).Elem().Kind()
			}
			if elemKind == reflect.Struct {
				result := make([]R, length)
				for i := 0; i < length; i++ {
					elem := rv.Index(i)
					if elem.Kind() == reflect.Ptr {
						elem = elem.Elem()
					}
					result[i] = onlyFields(elem.Interface(), fields...)
				}
				return result
			}
		}
	case reflect.Struct:
		return onlyFields(data, fields...)
	case reflect.Map:
		if rv.Type().Key().Kind() == reflect.String {
			return onlyFieldsMap(data.(map[string]any), fields...)
		}
	}

	return data
}

// ExcludeFields removes the specified fields from the provided data
// and returns the remaining fields as an `R` map. This function is
// useful for creating JSON responses without sensitive or unwanted
// fields from the original data. The operation can be performed on
// a single object, a slice of objects, an array of objects, or a map.
// When a slice or an array is provided, the function returns a slice of
// `R` maps. If the data is not a struct, slice/array of structs, or map,
// it returns the original data unchanged.
//
// Parameters:
//   - data: The input data from which fields will be excluded.
//   - fields: A list of field names to exclude from the resulting map.
//
// Returns:
//   - An `R` map containing the fields from the input data except
//     the specified fields, a slice of `R` maps if the input data
//     is a slice or an array, or the original input data unchanged
//     if it is not a struct, a map, or a slice/array of structs.
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
//	func HandlerSingle(w http.ResponseWriter, r *http.Request) {
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
//
//	func HandlerMultiple(w http.ResponseWriter, r *http.Request) {
//	 	users := []User{
//		    {
//	 		    ID:       1,
//	 		    Email:    "user_a@example.com",
//	 		    Password: "secret",
//	 		    IsActive: true,
//	 	    },
//		    {
//	 		    ID:       2,
//	 		    Email:    "user_b@example.com",
//	 		    Password: "secret",
//	 		    IsActive: true,
//	 	    },
//	 	}
//		data := resp.ExcludeFields(users, "Password")
//		if err := resp.JSON(w, data); err != nil {
//			// handle error
//		}
//	}
func ExcludeFields(data any, fields ...string) any {
	rv := reflect.ValueOf(data)

	switch rv.Kind() {
	case reflect.Ptr:
		rv = rv.Elem()
		if rv.Kind() == reflect.Struct {
			return excludeFields(rv.Interface(), fields...)
		}
	case reflect.Slice, reflect.Array:
		length := rv.Len()
		if length > 0 {
			elemKind := rv.Index(0).Kind()
			if elemKind == reflect.Ptr {
				elemKind = rv.Index(0).Elem().Kind()
			}
			if elemKind == reflect.Struct {
				result := make([]R, length)
				for i := 0; i < length; i++ {
					elem := rv.Index(i)
					if elem.Kind() == reflect.Ptr {
						elem = elem.Elem()
					}
					result[i] = excludeFields(elem.Interface(), fields...)
				}
				return result
			}
		}
	case reflect.Struct:
		return excludeFields(data, fields...)
	case reflect.Map:
		if rv.Type().Key().Kind() == reflect.String {
			return excludeFieldsMap(data.(map[string]any), fields...)
		}
	}

	return data
}

// onlyFields extracts only the specified fields from the provided
// data and returns them as an `R` map.
func onlyFields(data any, fields ...string) R {
	result := make(R)

	rv := reflect.ValueOf(data)
	rt := rv.Type()

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

// onlyFieldsMap extracts only the specified fields from the provided
// map and returns them as an `R` map.
func onlyFieldsMap(data map[string]any, fields ...string) R {
	result := make(R)
	allowed := make(map[string]bool, len(fields))
	for _, field := range fields {
		allowed[field] = true
	}

	for key, value := range data {
		if allowed[key] {
			result[key] = value
		}
	}

	return result
}

// excludeFields removes the specified fields from the provided data
// and returns the remaining fields as an `R` map.
func excludeFields(data any, fields ...string) R {
	result := make(R)

	rv := reflect.ValueOf(data)
	rt := rv.Type()

	excluded := make(map[string]bool, len(fields))
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

// excludeFieldsMap removes the specified fields from the provided
// map and returns the remaining fields as an `R` map.
func excludeFieldsMap(data map[string]any, fields ...string) R {
	result := make(R)
	excluded := make(map[string]bool, len(fields))
	for _, field := range fields {
		excluded[field] = true
	}

	for key, value := range data {
		if !excluded[key] {
			result[key] = value
		}
	}

	return result
}

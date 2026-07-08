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
// is not a struct, slice/array of structs, or map with string keys,
// it returns the original data unchanged.
//
// Parameters:
//   - data: The input data from which fields will be extracted. It
//     can be a struct, a slice/array of structs, or a map with
//     string keys (`map[string]any`).
//   - fields: A list of field names to include in the resulting map.
//
// Returns:
//   - An `R` map containing only the specified fields from the input
//     data, a slice of `R` maps if the input data is a slice or an array,
//     or the original input data unchanged if it is not a struct, a map
//     with string keys, or a slice/array of structs.
//
// Example Usage:
// The following example demonstrates how to use `OnlyFields` to create
// a JSON response that includes only specific fields from a user object.
//
//	import (
//		"net/http"
//		"github.com/goloop/resp/v2"
//	)
//
//	type User struct {
//		ID       int
//		Email    string
//		Password string
//		IsActive bool
//	}
//
//	func UserData(w http.ResponseWriter, r *http.Request) {
//		user := User{
//			ID:       1,
//			Email:    "user_a@example.com",
//			Password: "secret",
//			IsActive: true,
//		}
//		data := resp.OnlyFields(user, "ID", "Email", "IsActive")
//		if err := resp.JSON(w, data); err != nil {
//			// handle error
//		}
//	}
//
//	func UserList(w http.ResponseWriter, r *http.Request) {
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
	return filterFields(data, fields, true)
}

// ExcludeFields removes the specified fields from the provided data
// and returns the remaining fields as an `R` map. This function is
// useful for creating JSON responses without sensitive or unwanted
// fields from the original data. The operation can be performed on
// a single object, a slice of objects, an array of objects, or a map.
// When a slice or an array is provided, the function returns a slice of
// `R` maps. If the data is not a struct, slice/array of structs, or map
// with string keys, it returns the original data unchanged.
//
// Parameters:
//   - data: The input data from which fields will be excluded. It
//     can be a struct, a slice/array of structs, or a map with
//     string keys (`map[string]any`).
//   - fields: A list of field names to exclude from the resulting map.
//
// Returns:
//   - An `R` map containing the fields from the input data except
//     the specified fields, a slice of `R` maps if the input data
//     is a slice or an array, or the original input data unchanged
//     if it is not a struct, a map with string keys, or a slice/array
//     of structs.
//
// Example Usage:
// The following example demonstrates how to use `ExcludeFields` to create
// a JSON response that excludes specific fields from a user object.
//
//	import (
//		"net/http"
//		"github.com/goloop/resp/v2"
//	)
//
//	type User struct {
//		ID       int
//		Email    string
//		Password string
//		IsActive bool
//	}
//
//	func UserData(w http.ResponseWriter, r *http.Request) {
//		user := User{
//			ID:       1,
//			Email:    "user_a@example.com",
//			Password: "secret",
//			IsActive: true,
//		}
//		data := resp.ExcludeFields(user, "Password")
//		if err := resp.JSON(w, data); err != nil {
//			// handle error
//		}
//	}
//
//	func UserList(w http.ResponseWriter, r *http.Request) {
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
	return filterFields(data, fields, false)
}

// filterFields is the shared implementation behind OnlyFields and
// ExcludeFields. The keep flag selects the mode: when true, only the named
// fields are retained; when false, the named fields are dropped. The field
// set is built once and reused across every element of a slice/array, so a
// list of N structs costs O(N) work instead of rebuilding the set N times.
func filterFields(data any, fields []string, keep bool) any {
	set := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		set[field] = struct{}{}
	}

	rv := reflect.ValueOf(data)
	switch rv.Kind() {
	case reflect.Ptr:
		rv = rv.Elem()
		if rv.Kind() == reflect.Struct {
			return pickStruct(rv, set, keep)
		}
	case reflect.Slice, reflect.Array:
		length := rv.Len()
		if length > 0 {
			// Determine the element kind from the slice's element type, not the
			// first value, so a nil pointer at index 0 does not misclassify the
			// slice.
			et := rv.Type().Elem()
			if et.Kind() == reflect.Ptr {
				et = et.Elem()
			}
			if et.Kind() == reflect.Struct {
				result := make([]R, length)
				for i := range length {
					elem := rv.Index(i)
					if elem.Kind() == reflect.Ptr {
						if elem.IsNil() {
							result[i] = nil
							continue
						}
						elem = elem.Elem()
					}
					result[i] = pickStruct(elem, set, keep)
				}
				return result
			}
		}
	case reflect.Struct:
		return pickStruct(rv, set, keep)
	case reflect.Map:
		if rv.Type().Key().Kind() == reflect.String {
			return pickMap(rv, set, keep)
		}
	}

	return data
}

// pickStruct projects the exported fields of a struct value into an R map,
// keeping or dropping the names present in set according to keep. The
// membership test (`_, ok := set[name]; ok == keep`) yields the field for
// the "only" mode when it is in the set and for the "exclude" mode when it
// is not.
func pickStruct(rv reflect.Value, set map[string]struct{}, keep bool) R {
	result := make(R)
	rt := rv.Type()
	for i := range rt.NumField() {
		f := rt.Field(i)
		if !f.IsExported() {
			continue // reflect cannot read unexported fields; skip them
		}
		if _, ok := set[f.Name]; ok == keep {
			result[f.Name] = rv.Field(i).Interface()
		}
	}

	return result
}

// pickMap projects a string-keyed map into an R map, keeping or dropping
// the keys present in set according to keep. It uses reflection so it works
// for any map[string]T, not only map[string]any.
func pickMap(rv reflect.Value, set map[string]struct{}, keep bool) R {
	result := make(R, rv.Len())
	iter := rv.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		if _, ok := set[key]; ok == keep {
			result[key] = iter.Value().Interface()
		}
	}

	return result
}

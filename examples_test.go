package resp_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/goloop/resp/v2"
)

// ExampleJSON sends a JSON object using the R shorthand map.
func ExampleJSON() {
	w := httptest.NewRecorder()

	_ = resp.JSON(w, resp.R{"message": "hello", "ok": true})

	fmt.Print(w.Body.String())
	// Output: {"message":"hello","ok":true}
}

// ExampleError shows the v2 error contract: the first argument is the HTTP
// status (mirrored into the body), and WithErrorCode overrides only the body
// code with a separate business code.
func ExampleError() {
	w := httptest.NewRecorder()

	_ = resp.Error(w, http.StatusNotFound, "User not found", resp.WithErrorCode(1042))

	fmt.Println(w.Code)
	fmt.Print(w.Body.String())
	// Output:
	// 404
	// {"code":1042,"message":"User not found"}
}

// ExampleOnlyFields keeps only the whitelisted fields of a struct — a typical
// way to hide a password before serialization.
func ExampleOnlyFields() {
	type User struct {
		ID       int
		Email    string
		Password string
	}

	u := User{ID: 1, Email: "user@example.com", Password: "secret"}
	data := resp.OnlyFields(u, "ID", "Email")

	out, _ := json.Marshal(data)
	fmt.Println(string(out))
	// Output: {"Email":"user@example.com","ID":1}
}

// ExampleCreated returns 201 Created with a Location header and a JSON body —
// the idiomatic response after creating a REST resource.
func ExampleCreated() {
	w := httptest.NewRecorder()

	_ = resp.Created(w, "/users/42", resp.R{"id": 42})

	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Location"))
	fmt.Print(w.Body.String())
	// Output:
	// 201
	// /users/42
	// {"id":42}
}

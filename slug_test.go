package resp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The slug is what a client switches on and a translation key is looked up by,
// so it has to arrive as its own field rather than be inferred from a number.
func TestWithErrorSlug(t *testing.T) {
	w := httptest.NewRecorder()
	if err := Error(w, http.StatusConflict, "That name is already in use",
		WithErrorSlug("slug_taken")); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want 409", w.Code)
	}

	var body ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not an error response: %v", err)
	}
	if body.Error != "slug_taken" {
		t.Errorf("error = %q, want slug_taken", body.Error)
	}
	if body.Code != http.StatusConflict {
		t.Errorf("code = %d, want the status", body.Code)
	}
	if body.Message != "That name is already in use" {
		t.Errorf("message = %q", body.Message)
	}
}

// The slug and the business code are independent: naming one must not change
// the other.
func TestErrorSlugAndCodeAreIndependent(t *testing.T) {
	w := httptest.NewRecorder()
	if err := Error(w, http.StatusTooManyRequests, "Slow down",
		WithErrorCode(1042), WithErrorSlug("rate_limit")); err != nil {
		t.Fatal(err)
	}

	var body ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code != 1042 || body.Error != "rate_limit" {
		t.Errorf("body = %+v, want both the code and the slug", body)
	}
}

// The compatibility half: a response that names no slug must be exactly the
// bytes it was before the field existed.
func TestErrorWithoutSlugIsUnchanged(t *testing.T) {
	w := httptest.NewRecorder()
	if err := Error(w, http.StatusNotFound, "User not found"); err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(w.Body.String())
	const want = `{"code":404,"message":"User not found"}`
	if got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
	if strings.Contains(got, `"error"`) {
		t.Error("an unset slug still emitted its field")
	}
}

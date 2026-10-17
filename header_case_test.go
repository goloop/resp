package resp

import (
	"net/http/httptest"
	"testing"
)

func TestSingleHeaderCaseInsensitive(t *testing.T) {
	for _, key := range []string{"content-type", "ETAG", "etag", "x-xss-protection"} {
		w := httptest.NewRecorder()
		r := NewResponse(w)
		r.AddHeader(key, "old").AddHeader(key, "new")
		if v := w.Header().Values(key); len(v) != 1 || v[0] != "new" {
			t.Fatalf("%s = %q", key, v)
		}
	}
}

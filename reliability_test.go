package resp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// trackingWriter is an http.ResponseWriter that records whether the status
// line was committed (WriteHeader) and captures the body. It lets a test
// distinguish "the status was sent to the client" from "nothing was written
// yet", which is exactly the guarantee BUG-01 is about.
type trackingWriter struct {
	header      http.Header
	body        strings.Builder
	status      int
	wroteHeader bool
	writeErr    error
}

func newTrackingWriter() *trackingWriter {
	return &trackingWriter{header: http.Header{}}
}

func (w *trackingWriter) Header() http.Header { return w.header }

func (w *trackingWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.status = code
		w.wroteHeader = true
	}
}

func (w *trackingWriter) Write(b []byte) (int, error) {
	if w.writeErr != nil {
		return 0, w.writeErr
	}
	w.wroteHeader = true
	return w.body.WriteString(string(b))
}

// cyclicNode builds a self-referential structure that the JSON encoder
// rejects with a cycle error — a realistic way an encode can fail midway.
type cyclicNode struct {
	Name string      `json:"name"`
	Next *cyclicNode `json:"next"`
}

// -----------------------------------------------------------------------
// BUG-01 — encoding errors must not commit a success status (buffered mode).
// -----------------------------------------------------------------------

func TestJSON_EncodeError_DoesNotCommitStatus(t *testing.T) {
	cases := map[string]any{
		"channel": map[string]any{"bad": make(chan int)},
		"func":    map[string]any{"bad": func() {}},
		"cycle": func() any {
			n := &cyclicNode{Name: "root"}
			n.Next = n
			return n
		}(),
	}

	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			w := newTrackingWriter()
			r := NewResponse(w)

			err := r.JSON(data)
			if err == nil {
				t.Fatalf("JSON(%s) expected an encode error, got nil", name)
			}

			// The whole point of buffered mode: on failure the status line
			// was never sent and no partial body leaked to the client.
			if w.wroteHeader {
				t.Errorf("JSON(%s): status line was committed despite "+
					"encode error (BUG-01 regression)", name)
			}
			if w.body.Len() != 0 {
				t.Errorf("JSON(%s): wrote %d body bytes despite encode "+
					"error, want 0", name, w.body.Len())
			}

			// Because nothing was committed, the handler can still send a
			// proper error response afterwards.
			if err2 := r.Error(http.StatusInternalServerError, "boom"); err2 != nil {
				t.Errorf("follow-up Error() failed: %v", err2)
			}
			if w.status != http.StatusInternalServerError {
				t.Errorf("follow-up status = %d, want %d",
					w.status, http.StatusInternalServerError)
			}
		})
	}
}

func TestJSONP_EncodeError_DoesNotCommitStatus(t *testing.T) {
	w := newTrackingWriter()
	r := NewResponse(w)

	err := r.JSONP(map[string]any{"bad": make(chan int)}, "cb")
	if err == nil {
		t.Fatal("JSONP() expected an encode error, got nil")
	}
	if w.wroteHeader {
		t.Error("JSONP(): status committed despite encode error (BUG-01)")
	}
	if w.body.Len() != 0 {
		t.Errorf("JSONP(): wrote %d body bytes, want 0", w.body.Len())
	}
}

// In streaming mode the caller explicitly accepts that an encode error may
// surface only after the status was committed. We assert the trade-off so
// the documented contract is locked down.
func TestJSON_StreamMode_CommitsStatusThenFails(t *testing.T) {
	w := newTrackingWriter()
	r := NewResponse(w, WithStreamJSON())

	err := r.JSON(map[string]any{"bad": make(chan int)})
	if err == nil {
		t.Fatal("JSON() stream mode: expected encode error, got nil")
	}
	if !w.wroteHeader {
		t.Error("JSON() stream mode: expected status to be committed first")
	}
}

func TestJSON_StreamMode_ValidPayload(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w, WithStreamJSON())

	if err := r.JSON(R{"hello": "world"}); err != nil {
		t.Fatalf("JSON() stream mode returned error: %v", err)
	}
	if got := strings.TrimSpace(w.Body.String()); got != `{"hello":"world"}` {
		t.Errorf("JSON() stream body = %q", got)
	}
	if w.Code != http.StatusOK {
		t.Errorf("JSON() stream status = %d, want 200", w.Code)
	}
}

// -----------------------------------------------------------------------
// BUG-06 — JSONP callback names are validated (reflected-XSS vector).
// -----------------------------------------------------------------------

func TestJSONP_RejectsUnsafeCallbacks(t *testing.T) {
	unsafe := []string{
		"",
		"alert(1);//",
		"<script>alert(1)</script>",
		"a=alert,b",
		"cb;evil()",
		"cb()",
		"cb ",
		" cb",
		"foo.bar()",
		"1cb",   // may not start with a digit
		".cb",   // may not start with a dot
		"cb.",   // may not end with a dot
		"cb..x", // empty path segment
		`cb"`,   // quote
		"cb\n",  // newline
		strings.Repeat("a", maxJSONPCallbackLen+1), // too long
	}

	for _, cb := range unsafe {
		t.Run(strconv.Quote(cb), func(t *testing.T) {
			w := newTrackingWriter()
			r := NewResponse(w)

			err := r.JSONP(R{"a": 1}, cb)
			if !errors.Is(err, ErrInvalidJSONPCallback) {
				t.Fatalf("JSONP(%q) error = %v, want ErrInvalidJSONPCallback",
					cb, err)
			}
			// Nothing must be written for a rejected callback.
			if w.wroteHeader || w.body.Len() != 0 {
				t.Errorf("JSONP(%q) wrote output for a rejected callback", cb)
			}
		})
	}
}

func TestJSONP_AcceptsSafeCallbacks(t *testing.T) {
	safe := []string{
		"cb",
		"callbackFunction",
		"_private",
		"$jsonp",
		"app.handlers.onData",
		"a1.b2.c3",
		"$.cb",
	}

	for _, cb := range safe {
		t.Run(cb, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := NewResponse(w)

			if err := r.JSONP(R{"a": 1}, cb); err != nil {
				t.Fatalf("JSONP(%q) unexpected error: %v", cb, err)
			}
			want := cb + `({"a":1});`
			if got := w.Body.String(); got != want {
				t.Errorf("JSONP(%q) body = %q, want %q", cb, got, want)
			}
		})
	}
}

// -----------------------------------------------------------------------
// BUG-03 — Error: status is the HTTP status, body code mirrors it.
// -----------------------------------------------------------------------

func TestError_StatusSemantics(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		message    string
		opts       []Option
		wantStatus int
		wantBody   string
	}{
		{
			name:       "status mirrored into body",
			status:     http.StatusNotFound,
			message:    "Not Found",
			wantStatus: 404,
			wantBody:   `{"code":404,"message":"Not Found"}`,
		},
		{
			name:       "empty message uses canonical text",
			status:     http.StatusBadRequest,
			message:    "",
			wantStatus: 400,
			wantBody:   `{"code":400,"message":"Bad Request"}`,
		},
		{
			name:       "WithErrorCode overrides only body code",
			status:     http.StatusNotFound,
			message:    "User not found",
			opts:       []Option{WithErrorCode(1042)},
			wantStatus: 404,
			wantBody:   `{"code":1042,"message":"User not found"}`,
		},
		{
			name:       "status arg overrides WithStatus",
			status:     http.StatusNotFound,
			message:    "x",
			opts:       []Option{WithStatus(http.StatusOK)},
			wantStatus: 404,
			wantBody:   `{"code":404,"message":"x"}`,
		},
		{
			name:       "below range falls back to 500",
			status:     7,
			message:    "weird",
			wantStatus: 500,
			wantBody:   `{"code":500,"message":"weird"}`,
		},
		{
			name:       "above range falls back to 500",
			status:     600,
			message:    "weird",
			wantStatus: 500,
			wantBody:   `{"code":500,"message":"weird"}`,
		},
		{
			name:       "negative falls back to 500",
			status:     -1,
			message:    "weird",
			wantStatus: 500,
			wantBody:   `{"code":500,"message":"weird"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if err := Error(w, tt.status, tt.message, tt.opts...); err != nil {
				t.Fatalf("Error() returned error: %v", err)
			}
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if got := strings.TrimSpace(w.Body.String()); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
		})
	}
}

// -----------------------------------------------------------------------
// BUG-04 — Redirect only honours genuine redirect codes.
// -----------------------------------------------------------------------

func TestRedirect_StatusCodes(t *testing.T) {
	tests := []struct {
		set  int
		want int
	}{
		// Genuine redirects are preserved.
		{StatusMultipleChoices, StatusMultipleChoices},     // 300
		{StatusMovedPermanently, StatusMovedPermanently},   // 301
		{StatusFound, StatusFound},                         // 302
		{StatusSeeOther, StatusSeeOther},                   // 303
		{StatusTemporaryRedirect, StatusTemporaryRedirect}, // 307
		{StatusPermanentRedirect, StatusPermanentRedirect}, // 308
		// Non-redirect 3xx codes must fall back to 302.
		{StatusNotModified, StatusFound}, // 304
		{305, StatusFound},               // 305 Use Proxy
		{306, StatusFound},               // 306 (unused)
		// Clearly invalid codes fall back to 302.
		{199, StatusFound},
		{http.StatusOK, StatusFound},
		{http.StatusInternalServerError, StatusFound},
	}

	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.set), func(t *testing.T) {
			w := httptest.NewRecorder()
			r := NewResponse(w, WithStatus(tt.set))
			if err := r.Redirect("https://example.com/x"); err != nil {
				t.Fatalf("Redirect() error: %v", err)
			}
			if w.Code != tt.want {
				t.Errorf("Redirect(set=%d) status = %d, want %d",
					tt.set, w.Code, tt.want)
			}
			if loc := w.Header().Get("Location"); loc != "https://example.com/x" {
				t.Errorf("Redirect() Location = %q", loc)
			}
		})
	}
}

// -----------------------------------------------------------------------
// BUG-02 — ServeFileAsDownload sets length, nosniff and UTF-8 file names.
// -----------------------------------------------------------------------

func TestServeFileAsDownload_HeadersAndUnicode(t *testing.T) {
	tests := []struct {
		name         string
		fileName     string
		data         []byte
		wantContains []string // substrings expected in Content-Disposition
	}{
		{
			name:         "ascii",
			fileName:     "report.pdf",
			data:         []byte("PDF"),
			wantContains: []string{`filename="report.pdf"`, `filename*=UTF-8''report.pdf`},
		},
		{
			name:         "cyrillic",
			fileName:     "звіт.pdf",
			data:         []byte("PDF"),
			wantContains: []string{`filename="____.pdf"`, `filename*=UTF-8''`},
		},
		{
			name:         "quote and backslash sanitized in fallback",
			fileName:     `a"b\c.txt`,
			data:         []byte("x"),
			wantContains: []string{`filename="abc.txt"`, `filename*=UTF-8''`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := NewResponse(w)
			if err := r.ServeFileAsDownload(tt.fileName, tt.data); err != nil {
				t.Fatalf("ServeFileAsDownload() error: %v", err)
			}

			cd := w.Header().Get("Content-Disposition")
			for _, want := range tt.wantContains {
				if !strings.Contains(cd, want) {
					t.Errorf("Content-Disposition %q does not contain %q", cd, want)
				}
			}
			// The quoted ASCII fallback must never contain a bare quote or
			// backslash that would break header parsing.
			quoted := cd[strings.Index(cd, `filename="`)+len(`filename="`):]
			quoted = quoted[:strings.Index(quoted, `"`)]
			if strings.ContainsAny(quoted, "\"\\") {
				t.Errorf("ASCII fallback %q contains unsafe characters", quoted)
			}

			if got, want := w.Header().Get("Content-Length"),
				strconv.Itoa(len(tt.data)); got != want {
				t.Errorf("Content-Length = %q, want %q", got, want)
			}
			if got := w.Header().Get("X-Content-Type-Options"); got != "nosniff" {
				t.Errorf("X-Content-Type-Options = %q, want nosniff", got)
			}
		})
	}
}

// -----------------------------------------------------------------------
// JSONIndent — pretty output, buffered safety, custom encoder ignored.
// -----------------------------------------------------------------------

func TestJSONIndent(t *testing.T) {
	w := httptest.NewRecorder()
	if err := JSONIndent(w, R{"a": 1}, "", "  "); err != nil {
		t.Fatalf("JSONIndent() error: %v", err)
	}
	want := "{\n  \"a\": 1\n}\n"
	if got := w.Body.String(); got != want {
		t.Errorf("JSONIndent() body = %q, want %q", got, want)
	}
	if ct := w.Header().Get("Content-Type"); ct != MIMEApplicationJSONCharsetUTF8 {
		t.Errorf("JSONIndent() Content-Type = %q", ct)
	}
}

func TestJSONIndent_IgnoresCustomEncoder(t *testing.T) {
	// A custom encoder must NOT hijack JSONIndent — indentation is the point.
	enc := func(w io.Writer, v any) error { return errFromEncoder }
	w := httptest.NewRecorder()
	r := NewResponse(w, ApplyJSONEncoder(enc))
	if err := r.JSONIndent(R{"a": 1}, "", "\t"); err != nil {
		t.Fatalf("JSONIndent() unexpectedly used the custom encoder: %v", err)
	}
	if !strings.Contains(w.Body.String(), "\t\"a\": 1") {
		t.Errorf("JSONIndent() body not tab-indented: %q", w.Body.String())
	}
}

func TestJSONIndent_EncodeError_DoesNotCommitStatus(t *testing.T) {
	w := newTrackingWriter()
	r := NewResponse(w)
	if err := r.JSONIndent(map[string]any{"bad": make(chan int)}, "", "  "); err == nil {
		t.Fatal("JSONIndent() expected encode error, got nil")
	}
	if w.wroteHeader || w.body.Len() != 0 {
		t.Error("JSONIndent(): committed output despite encode error")
	}
}

var errFromEncoder = errors.New("custom encoder used")

// -----------------------------------------------------------------------
// RFC 5987 ext-value encoding for Content-Disposition filename*.
// -----------------------------------------------------------------------

func TestRFC5987Encode(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		// attr-chars pass through unencoded.
		{"report.pdf", "report.pdf"},
		{"a-b_c.d~e", "a-b_c.d~e"},
		{"!#$&+-.^_`|~", "!#$&+-.^_`|~"},
		// space and other ASCII are percent-encoded.
		{"my file.txt", "my%20file.txt"},
		{"a(b)c", "a%28b%29c"},
		{`q"x`, "q%22x"},
		{"100%", "100%25"},
		// multi-byte UTF-8 is encoded byte-by-byte, uppercase hex.
		{"é", "%C3%A9"},
		{"звіт", "%D0%B7%D0%B2%D1%96%D1%82"},
	}
	for _, tt := range tests {
		if got := rfc5987Encode(tt.in); got != tt.want {
			t.Errorf("rfc5987Encode(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// -----------------------------------------------------------------------
// New API — XML, Blob, Created, Accepted, getters.
// -----------------------------------------------------------------------

func TestXML(t *testing.T) {
	type note struct {
		XMLName xml.Name `xml:"note"`
		To      string   `xml:"to"`
	}

	w := httptest.NewRecorder()
	if err := XML(w, note{To: "World"}); err != nil {
		t.Fatalf("XML() error: %v", err)
	}
	if ct := w.Header().Get("Content-Type"); ct != MIMEApplicationXMLCharsetUTF8 {
		t.Errorf("XML() Content-Type = %q", ct)
	}
	if got, want := w.Body.String(), `<note><to>World</to></note>`; got != want {
		t.Errorf("XML() body = %q, want %q", got, want)
	}
}

func TestXML_EncodeError_DoesNotCommitStatus(t *testing.T) {
	w := newTrackingWriter()
	r := NewResponse(w)

	// A map has no stable XML representation and fails to marshal.
	err := r.XML(map[string]string{"a": "b"})
	if err == nil {
		t.Fatal("XML() expected encode error for a map, got nil")
	}
	if w.wroteHeader || w.body.Len() != 0 {
		t.Error("XML(): committed output despite encode error")
	}
}

func TestBlob(t *testing.T) {
	t.Run("explicit content type", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := []byte("%PDF-1.7")
		if err := Blob(w, "application/pdf", data); err != nil {
			t.Fatalf("Blob() error: %v", err)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/pdf" {
			t.Errorf("Blob() Content-Type = %q", ct)
		}
		if cl := w.Header().Get("Content-Length"); cl != strconv.Itoa(len(data)) {
			t.Errorf("Blob() Content-Length = %q", cl)
		}
		if w.Body.String() != string(data) {
			t.Errorf("Blob() body = %q", w.Body.String())
		}
	})

	t.Run("empty content type defaults to octet-stream", func(t *testing.T) {
		w := httptest.NewRecorder()
		if err := Blob(w, "", []byte("x")); err != nil {
			t.Fatalf("Blob() error: %v", err)
		}
		if ct := w.Header().Get("Content-Type"); ct != MIMEOctetStream {
			t.Errorf("Blob() Content-Type = %q, want %q", ct, MIMEOctetStream)
		}
	})
}

func TestCreated(t *testing.T) {
	t.Run("location and body", func(t *testing.T) {
		w := httptest.NewRecorder()
		if err := Created(w, "/users/42", R{"id": 42}); err != nil {
			t.Fatalf("Created() error: %v", err)
		}
		if w.Code != StatusCreated {
			t.Errorf("Created() status = %d, want 201", w.Code)
		}
		if loc := w.Header().Get("Location"); loc != "/users/42" {
			t.Errorf("Created() Location = %q", loc)
		}
		if got := strings.TrimSpace(w.Body.String()); got != `{"id":42}` {
			t.Errorf("Created() body = %q", got)
		}
	})

	t.Run("nil body sends only status", func(t *testing.T) {
		w := httptest.NewRecorder()
		if err := Created(w, "", nil); err != nil {
			t.Fatalf("Created() error: %v", err)
		}
		if w.Code != StatusCreated {
			t.Errorf("Created() status = %d, want 201", w.Code)
		}
		if w.Body.Len() != 0 {
			t.Errorf("Created() body = %q, want empty", w.Body.String())
		}
		if loc, ok := w.Header()["Location"]; ok {
			t.Errorf("Created() set Location to %v for empty location", loc)
		}
	})

	t.Run("explicit status is respected", func(t *testing.T) {
		w := httptest.NewRecorder()
		if err := Created(w, "/x", R{"ok": true}, WithStatus(http.StatusOK)); err != nil {
			t.Fatalf("Created() error: %v", err)
		}
		if w.Code != http.StatusOK {
			t.Errorf("Created() status = %d, want 200 (explicit)", w.Code)
		}
	})
}

func TestAccepted(t *testing.T) {
	t.Run("with body", func(t *testing.T) {
		w := httptest.NewRecorder()
		if err := Accepted(w, R{"job_id": "abc"}); err != nil {
			t.Fatalf("Accepted() error: %v", err)
		}
		if w.Code != StatusAccepted {
			t.Errorf("Accepted() status = %d, want 202", w.Code)
		}
		if got := strings.TrimSpace(w.Body.String()); got != `{"job_id":"abc"}` {
			t.Errorf("Accepted() body = %q", got)
		}
	})

	t.Run("nil body", func(t *testing.T) {
		w := httptest.NewRecorder()
		if err := Accepted(w, nil); err != nil {
			t.Fatalf("Accepted() error: %v", err)
		}
		if w.Code != StatusAccepted {
			t.Errorf("Accepted() status = %d, want 202", w.Code)
		}
		if w.Body.Len() != 0 {
			t.Errorf("Accepted() body = %q, want empty", w.Body.String())
		}
	})
}

func TestStatusAndHeaderGetters(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	if got := r.Status(); got != StatusUndefined {
		t.Errorf("Status() initial = %d, want %d", got, StatusUndefined)
	}

	r.SetStatus(http.StatusTeapot)
	if got := r.Status(); got != http.StatusTeapot {
		t.Errorf("Status() = %d, want %d", got, http.StatusTeapot)
	}

	// Header() returns the live header map shared with the writer.
	r.Header().Set("X-Trace", "abc")
	if got := w.Header().Get("X-Trace"); got != "abc" {
		t.Errorf("Header() is not the live writer header: got %q", got)
	}
}

// -----------------------------------------------------------------------
// Security headers — conservative preset + explicit CSP opt-in.
// -----------------------------------------------------------------------

func TestSecurityHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	if err := JSON(w, R{"ok": true}, SecurityHeaders()); err != nil {
		t.Fatalf("JSON() error: %v", err)
	}

	want := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
	}
	for k, v := range want {
		if got := w.Header().Get(k); got != v {
			t.Errorf("SecurityHeaders() %s = %q, want %q", k, got, v)
		}
	}
	// The preset must NOT set a CSP on its own.
	if got := w.Header().Get("Content-Security-Policy"); got != "" {
		t.Errorf("SecurityHeaders() unexpectedly set CSP = %q", got)
	}
}

func TestSecurityHeaders_OverwritesAndIsIdempotent(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)
	// A pre-existing weaker value must be overwritten, and a second apply
	// must not append a duplicate header value.
	r.Header().Set("X-Frame-Options", "SAMEORIGIN")
	SecurityHeaders()(r)
	SecurityHeaders()(r)

	if got := r.Header().Values("X-Frame-Options"); len(got) != 1 || got[0] != "DENY" {
		t.Errorf("X-Frame-Options = %v, want exactly [DENY]", got)
	}
}

func TestWithCSP(t *testing.T) {
	w := httptest.NewRecorder()
	policy := "default-src 'self'; img-src *"
	if err := JSON(w, R{"ok": true}, SecurityHeaders(), WithCSP(policy)); err != nil {
		t.Fatalf("JSON() error: %v", err)
	}
	if got := w.Header().Get("Content-Security-Policy"); got != policy {
		t.Errorf("WithCSP() = %q, want %q", got, policy)
	}
}

// -----------------------------------------------------------------------
// fields — map projection must not depend on the concrete value type and
// must never panic (the legacy code type-asserted map[string]any).
// -----------------------------------------------------------------------

func TestOnlyExcludeFields_NonAnyMap(t *testing.T) {
	// A map[string]int previously panicked on a map[string]any assertion.
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	only := OnlyFields(m, "a", "c")
	got, ok := only.(R)
	if !ok {
		t.Fatalf("OnlyFields returned %T, want R", only)
	}
	if len(got) != 2 || got["a"] != 1 || got["c"] != 3 {
		t.Errorf("OnlyFields(map[string]int) = %v", got)
	}

	excl := ExcludeFields(m, "b").(R)
	if len(excl) != 2 || excl["b"] != nil {
		t.Errorf("ExcludeFields(map[string]int) = %v", excl)
	}
}

// -----------------------------------------------------------------------
// Concurrency — the buffer pool must not leak data between responses.
// -----------------------------------------------------------------------

func TestJSON_ConcurrentBufferPool(t *testing.T) {
	const n = 200
	var wg sync.WaitGroup
	errs := make([]error, n)
	bodies := make([]string, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			w := httptest.NewRecorder()
			// Distinct, differently sized payload per goroutine to maximise
			// the chance of cross-contamination if the pool were misused.
			payload := strings.Repeat("x", i)
			errs[i] = JSON(w, R{"i": i, "p": payload})
			bodies[i] = strings.TrimSpace(w.Body.String())
		}(i)
	}
	wg.Wait()

	for i := 0; i < n; i++ {
		if errs[i] != nil {
			t.Fatalf("JSON() #%d error: %v", i, errs[i])
		}
		want := fmt.Sprintf(`{"i":%d,"p":%q}`, i, strings.Repeat("x", i))
		if bodies[i] != want {
			t.Errorf("JSON() #%d body = %q, want %q", i, bodies[i], want)
		}
	}
}

// A payload larger than maxPooledBuffer must still serialize correctly while
// the oversized buffer is dropped rather than pooled.
func TestJSON_LargePayloadBeyondPoolLimit(t *testing.T) {
	big := strings.Repeat("a", maxPooledBuffer+1024)
	w := httptest.NewRecorder()
	if err := JSON(w, R{"big": big}); err != nil {
		t.Fatalf("JSON() large payload error: %v", err)
	}
	want := fmt.Sprintf(`{"big":%q}`, big)
	if got := strings.TrimSpace(w.Body.String()); got != want {
		t.Errorf("JSON() large payload truncated or corrupted (len %d)", len(got))
	}
}

// -----------------------------------------------------------------------
// Fuzz — adversarial inputs for the security-sensitive surfaces.
// -----------------------------------------------------------------------

// FuzzJSONPCallback asserts the core safety invariant: any callback the
// validator accepts contains only characters that are safe to splice into a
// JavaScript call site — never quotes, brackets, semicolons or whitespace.
func FuzzJSONPCallback(f *testing.F) {
	for _, s := range []string{"cb", "a.b.c", "<script>", "x();", "", "$_"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, cb string) {
		if !validJSONPCallback(cb) {
			return // rejected names need no further guarantees
		}
		if strings.ContainsAny(cb, "()<>\"';, \t\n\r/\\{}[]") {
			t.Fatalf("validJSONPCallback accepted unsafe name %q", cb)
		}
		// An accepted callback must round-trip through JSONP without error.
		w := httptest.NewRecorder()
		if err := JSONP(w, R{"a": 1}, cb); err != nil {
			t.Fatalf("JSONP rejected a validated callback %q: %v", cb, err)
		}
		if want := cb + `({"a":1});`; w.Body.String() != want {
			t.Fatalf("JSONP body = %q, want %q", w.Body.String(), want)
		}
	})
}

// FuzzOnlyFields drives the map projection with arbitrary keys and field
// selections; the only guarantee under test is "never panics and returns a
// subset".
func FuzzOnlyFields(f *testing.F) {
	f.Add("a", "b")
	f.Add("", "x")
	f.Fuzz(func(t *testing.T, k1, field string) {
		m := map[string]any{k1: 1, "fixed": 2}
		out := OnlyFields(m, field)
		r, ok := out.(R)
		if !ok {
			t.Fatalf("OnlyFields returned %T", out)
		}
		for key := range r {
			if key != field {
				t.Fatalf("OnlyFields kept unexpected key %q (want only %q)",
					key, field)
			}
		}
	})
}

// TestServeFile_DetectsContentType guards BUG-03: ServeFile must not preset
// application/octet-stream, so http.ServeFile can detect the real type and a
// browser renders the file instead of downloading it.
func TestServeFile_DetectsContentType(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "page-*.html")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("<html><body>hi</body></html>"); err != nil {
		t.Fatal(err)
	}
	f.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if err := ServeFile(w, req, f.Name()); err != nil {
		t.Fatalf("ServeFile: %v", err)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html...", ct)
	}
}

// TestContentDisposition_NoInjection guards BUG-04: a quote in the filename
// must not break out of the quoted-string to inject extra header parameters.
func TestContentDisposition_NoInjection(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, R{"ok": true},
		AddContentDisposition("attachment", `report.pdf"; dummy="x`))
	got := w.Header().Get("Content-Disposition")
	if strings.Contains(got, `"; dummy=`) || strings.Count(got, `"`) != 2 {
		t.Errorf("Content-Disposition = %q, injection not neutralized", got)
	}
}

// TestError_UnknownStatusHasMessage guards BUG-06: a valid but unnamed status
// still yields a non-empty message.
func TestError_UnknownStatusHasMessage(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, 599, "")
	if !strings.Contains(w.Body.String(), `"message":"Unknown Status"`) {
		t.Errorf("body = %s, want a non-empty message", w.Body.String())
	}
}

// TestHTTPDates_AreGMT guards BUG-05: date headers are IMF-fixdate in GMT.
func TestHTTPDates_AreGMT(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, R{"ok": true},
		AddLastModified(time.Date(2026, 7, 8, 12, 0, 0, 0, time.FixedZone("EEST", 3*3600))))
	if got := w.Header().Get("Last-Modified"); got != "Wed, 08 Jul 2026 09:00:00 GMT" {
		t.Errorf("Last-Modified = %q, want GMT IMF-fixdate", got)
	}
}

// TestContentDisposition_TypeSanitized guards that a disposition type made of
// only invalid characters falls back to "attachment".
func TestContentDisposition_TypeSanitized(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, R{"ok": true}, AddContentDisposition(`";`, "file.pdf"))
	if got := w.Header().Get("Content-Disposition"); !strings.HasPrefix(got, "attachment;") {
		t.Errorf("Content-Disposition = %q, want attachment fallback", got)
	}
}

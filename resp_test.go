package resp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/goloop/g"
)

// TestFuncJSON tests the JSON function.
func TestFuncJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"hello": "world"}

	err := JSON(w, data)

	if err != nil {
		t.Errorf("JSON() returned an error: %v", err)
	}

	// Check that the Content-Type header is set correctly.
	got := w.Header().Get("Content-Type")
	if want := MIMEApplicationJSONCharsetUTF8; got != want {
		t.Errorf("JSON() Content-Type = %v, want %v", got, want)
	}

	// Check the response body.
	expectedBody, _ := json.Marshal(data)
	expected := string(expectedBody)
	res := g.Trim(w.Body.String())
	if res != expected {
		t.Errorf("JSON() body = %v, want %v", res, expected)
	}
}

// TestFuncJSONP tests the JSONP function.
func TestFuncJSONP(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"hello": "world"}
	callback := "callback"

	err := JSONP(w, data, callback)

	if err != nil {
		t.Errorf("JSONP() returned an error: %v", err)
	}

	// Check that the Content-Type header is set correctly.
	got := w.Header().Get("Content-Type")
	if want := MIMEApplicationJavaScriptCharsetUTF8; got != want {
		t.Errorf("JSONP() Content-Type = %v, want %v", got, want)
	}

	// Check the response body.
	expectedBody, _ := json.Marshal(data)
	expected := callback + "(" + string(expectedBody) + ");"
	res := g.Trim(w.Body.String())
	if res != expected {
		t.Errorf("JSONP() body = %v, want %v", res, expected)
	}
}

// TestFuncError tests the Error function.
func TestFuncString(t *testing.T) {
	w := httptest.NewRecorder()
	data := "hello world"

	err := String(w, data)

	if err != nil {
		t.Errorf("String() returned an error: %v", err)
	}

	// Check that the Content-Type header is set correctly.
	got := w.Header().Get("Content-Type")
	if want := MIMETextPlain; got != want {
		t.Errorf("String() Content-Type = %v, want %v", got, want)
	}

	// Check the response body.
	expected := data
	res := g.Trim(w.Body.String())
	if res != expected {
		t.Errorf("String() body = %v, want %v", res, expected)
	}
}

// TestFuncError tests the Error function.
func TestFuncError(t *testing.T) {
	w := httptest.NewRecorder()
	status := 400
	message := "error message"

	err := Error(w, status, message)

	if err != nil {
		t.Errorf("Error() returned an error: %v", err)
	}

	// Check that the Content-Type header is set correctly.
	got := w.Header().Get("Content-Type")
	if want := MIMEApplicationJSONCharsetUTF8; got != want {
		t.Errorf("Error() Content-Type = %v, want %v", got, want)
	}

	// Check the response body.
	expected := `{"code":400,"message":"error message"}`
	res := g.Trim(w.Body.String())
	if res != expected {
		t.Errorf("Error() body = %v, want %v", res, expected)
	}
}

// TestFuncStream tests the Stream function.
func TestFuncStream(t *testing.T) {
	w := httptest.NewRecorder()
	data := "hello world"

	err := Stream(w, strings.NewReader(data))

	if err != nil {
		t.Errorf("Stream() returned an error: %v", err)
	}

	// Check that the Content-Type header is set correctly.
	got := w.Header().Get("Content-Type")
	if want := MIMEOctetStream; got != want {
		t.Errorf("Stream() Content-Type = %v, want %v", got, want)
	}

	// Check the response body.
	expected := data
	res := g.Trim(w.Body.String())
	if res != expected {
		t.Errorf("Stream() body = %v, want %v", res, expected)
	}
}

// TestFuncServeFile tests the ServeFile function.
func TestFuncServeFile(t *testing.T) {
	// Create a temporary file.
	tmpFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some content to the file.
	testContent := "Hello, world!"
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Create a test server.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeFile(w, r, tmpFile.Name()) // use the temp file
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create a request to the test server.
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the content of the file.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if got, want := string(body), testContent; got != want {
		t.Errorf("ServeFile() body = %q, want %q", got, want)
	}
}

// TestFuncServeFileAsDownload tests the ServeFileAsDownload function.
func TestFuncServeFileAsDownload(t *testing.T) {
	// Create a test server.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeFileAsDownload(w, "download.txt", []byte("Hello, download!"))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Request the test server.
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if got, want := string(body), "Hello, download!"; got != want {
		t.Errorf("ServeFileAsDownload() body = %q, want %q", got, want)
	}

	cd := resp.Header.Get("Content-Disposition")
	if got, want := cd, `attachment; filename="download.txt"`; got != want {
		t.Errorf("ServeFileAsDownload() Content-Disposition = %q, want %q",
			got, want)
	}
}

// TestFuncRedirect tests the Redirect function.
func TestFuncRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	url := "http://example.com"

	err := Redirect(w, url)

	if err != nil {
		t.Errorf("Redirect() returned an error: %v", err)
	}

	// Check the response status code.
	got := w.Code
	if want := http.StatusFound; got != want {
		t.Errorf("Redirect() status code = %v, want %v", got, want)
	}

	// Check the Location header.
	gotLocation := w.Header().Get("Location")
	if gotLocation != url {
		t.Errorf("Redirect() Location = %v, want %v", gotLocation, url)
	}
}

// TestFuncNoContent tests the NoContent function.
func TestFuncNoContent(t *testing.T) {
	w := httptest.NewRecorder()

	err := NoContent(w)

	if err != nil {
		t.Errorf("NoContent() returned an error: %v", err)
	}

	// Check the response status code.
	got := w.Code
	if want := http.StatusNoContent; got != want {
		t.Errorf("NoContent() status code = %v, want %v", got, want)
	}
}

// TestFuncHTML tests the HTML function.
func TestFuncHTML(t *testing.T) {
	w := httptest.NewRecorder()
	data := "<html><body>Hello, world!</body></html>"

	err := HTML(w, data)

	if err != nil {
		t.Errorf("HTML() returned an error: %v", err)
	}

	// Check that the Content-Type header is set correctly.
	got := w.Header().Get("Content-Type")
	if want := MIMETextHTMLCharsetUTF8; got != want {
		t.Errorf("HTML() Content-Type = %v, want %v", got, want)
	}

	// Check the response body.
	expected := data
	res := g.Trim(w.Body.String())
	if res != expected {
		t.Errorf("HTML() body = %v, want %v", res, expected)
	}
}

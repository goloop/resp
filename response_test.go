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

// TestNewResponse tests the NewResponse function.
func TestNewResponse(t *testing.T) {
	w := httptest.NewRecorder()
	response := NewResponse(w)

	if response == nil {
		t.Fatal("NewResponse() returned nil")
	}

	if response.httpWriter != w {
		t.Error("NewResponse() did not correctly assign httpWriter")
	}

	if response.statusCode != StatusUndefined {
		t.Errorf("NewResponse() incorrect initial status code, "+
			"got %v, want %v", response.statusCode, StatusUndefined)
	}
}

// TestPrepare tests the prepare method.
func TestPrepare(t *testing.T) {
	w := httptest.NewRecorder()
	response := NewResponse(w)

	// Set the default content type and status code.
	response.prepare(http.StatusOK, MIMEApplicationJSON)

	if w.Header().Get(HeaderContentType) != MIMEApplicationJSON {
		t.Errorf("Prepare() did not set Content-Type header correctly, "+
			"got %v, want %v",
			w.Header().Get(HeaderContentType), MIMEApplicationJSON)
	}

	if response.statusCode != http.StatusOK {
		t.Errorf("Prepare() did not set status code correctly,"+
			"got %v, want %v", response.statusCode, http.StatusOK)
	}
}

// TestSetStatus tests the SetStatus method.
func TestSetStatus(t *testing.T) {
	w := httptest.NewRecorder()
	response := NewResponse(w).SetStatus(http.StatusNotFound)

	if response.statusCode != http.StatusNotFound {
		t.Errorf("SetStatus() did not set the correct status code, "+
			"got %d, want %d", response.statusCode, http.StatusNotFound)
	}
}

// TestSetHeader tests the SetHeader method.
func TestSetHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	r.SetHeader(HeaderContentType, MIMEApplicationJSON)
	if got := w.Header().Get(HeaderContentType); got != MIMEApplicationJSON {
		t.Errorf("SetHeader() = %v, want %v", got, MIMEApplicationJSON)
	}
}

// TestSetHeader_SingleValue tests the SetHeader method for
// a single value header.
func TestSetHeader_SingleValue(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	r.SetHeader(HeaderContentType, MIMEApplicationJSON, MIMETextHTML)
	if got := w.Header().Get(HeaderContentType); got != MIMEApplicationJSON {
		t.Errorf("SetHeader() for single value header = %v, want %v",
			got, MIMEApplicationJSON)
	}
}

// TestSetHeader_MultipleValues tests the SetHeader method for
// a multiple value header.
func TestSetHeader_MultipleValues(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	r.SetHeader(HeaderCacheControl, "no-cache", "no-store")
	if got := w.Header().Get(HeaderCacheControl); got != "no-cache,no-store" {
		t.Errorf("SetHeader() for multiple value header = %v, "+
			"want %v", got, "no-cache,no-store")
	}
}

// TestAddHeader tests the AddHeader method.
func TestAddHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	r.AddHeader("X-Custom-Header", "Value1")
	r.AddHeader("X-Custom-Header", "Value2")

	exp := "Value1,Value2"
	got := w.Header().Values("X-Custom-Header")
	if strings.Join(got, ",") != exp {
		t.Errorf("AddHeader() = %v, want %v", got, exp)
	}
}

// TestAddHeader_SingleValue tests the AddHeader method for
// a single value header.
func TestAddHeader_SingleValue(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	r.AddHeader(HeaderContentLength, "1234", "5678")
	got := w.Header().Values(HeaderContentLength)
	if strings.Join(got, ",") != "1234" {
		t.Errorf("AddHeader() for single value header = %v, "+
			"want %v", got, "1234")
	}
}

// TestAddHeader_MultipleValues tests the AddHeader method for
// a multiple value header.
func TestAddHeader_MultipleValues(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	r.AddHeader("Accept-Encoding", "gzip", "deflate")
	values := w.Header()["Accept-Encoding"]
	if len(values) != 2 || values[0] != "gzip" || values[1] != "deflate" {
		t.Errorf("AddHeader() did not add multiple values correctly, "+
			"got %v", values)
	}
}

// TestDelHeader tests the DelHeader method.
func TestDelHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	r.SetHeader("X-Custom-Header", "Value")
	r.DelHeader("X-Custom-Header")

	if got := w.Header().Get("X-Custom-Header"); got != "" {
		t.Errorf("DelHeader() did not delete the header, got %v, "+
			"want %v", got, "")
	}
}

// TestClearHeaders tests the ClearHeaders method.
func TestClearHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	r.SetHeader("Content-Type", "application/json")
	r.SetHeader("X-Custom-Header", "Value")
	r.ClearHeaders()

	if got, want := len(w.Header()), 0; got != want {
		t.Errorf("ClearHeaders() did not clear the headers, "+
			"got %d headers, want %d", got, want)
	}
}

// TestSetCookie tests the SetCookie method.
func TestSetCookie(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	cookie := &http.Cookie{
		Name:  "test",
		Value: "value",
	}
	r.SetCookie(cookie)

	got := w.Header().Get("Set-Cookie")
	if !strings.Contains(got, "test=value") {
		t.Errorf("SetCookie() did not set the cookie correctly, got %v", got)
	}
}

// TestBindCookie tests the BindCookie method.
func TestBindCookie(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	cookie1 := &http.Cookie{
		Name:  "test",
		Value: "oldValue",
	}
	cookie2 := &http.Cookie{
		Name:  "test",
		Value: "newValue",
	}

	r.SetCookie(cookie1)
	r.BindCookie(cookie2)

	got := w.Header().Get("Set-Cookie")
	if !strings.Contains(got, "newValue") {
		t.Errorf("BindCookie() did not bind the new cookie correctly, "+
			"got %v", got)
	}
}

// TestDelCookie tests the DelCookie method.
func TestDelCookie(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	cookie := &http.Cookie{
		Name:  "test",
		Value: "value",
	}

	r.SetCookie(cookie)
	r.DelCookie("test")

	if got := w.Header().Get("Set-Cookie"); strings.Contains(got, "test=") {
		t.Errorf("DelCookie() did not delete the cookie, got %v", got)
	}
}

// TestDelCookieGlobal tests the DelCookie method with global cookies.
func TestDelCookieGlobal(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	// Додавання кількох cookies до відповіді
	http.SetCookie(w, &http.Cookie{Name: "keep", Value: "1"})
	http.SetCookie(w, &http.Cookie{Name: "delete", Value: "2"})
	http.SetCookie(w, &http.Cookie{Name: "alsoKeep", Value: "3"})

	r.DelCookie("delete")

	cookies := w.Result().Header["Set-Cookie"]
	if len(cookies) != 2 {
		t.Errorf("Expected 2 cookies after deletion, got %d", len(cookies))
	}

	for _, c := range cookies {
		if c == "delete=2" {
			t.Errorf("Cookie 'delete' was not removed")
		}
		if c != "keep=1" && c != "alsoKeep=3" {
			t.Errorf("Unexpected cookie present: %s", c)
		}
	}
}

// TestClearCookies tests the ClearCookies method.
func TestClearCookies(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	cookie1 := &http.Cookie{Name: "test1", Value: "value1"}
	cookie2 := &http.Cookie{Name: "test2", Value: "value2"}

	r.SetCookie(cookie1)
	r.SetCookie(cookie2)

	r.ClearCookies()

	if got := w.Header()["Set-Cookie"]; len(got) != 0 {
		t.Errorf("ClearCookies() did not clear the cookies, got %v", got)
	}
}

// TestExpiredCookie tests the ExpiredCookie method.
func TestExpiredCookie(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	r.ExpiredCookie("test")

	cookie := w.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, "test=deleted") ||
		!strings.Contains(cookie, "Max-Age=0") || // is set to 0 (not -1)
		!strings.Contains(cookie, "Expires=Thu, 01 Jan 1970 00:00:00 GMT") {
		t.Errorf("ExpiredCookie() did not set the cookie as "+
			"expired correctly, got %v", cookie)
	}
}

// TestJSON tests the JSON method.
func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	data := map[string]string{"hello": "world"}
	err := r.JSON(data)
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

// TestJSONP tests the JSONP method.
func TestJSONP(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	data := map[string]string{"hello": "world"}
	callback := "callbackFunction"
	err := r.JSONP(data, callback)
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
	res := w.Body.String()
	if res != expected {
		t.Errorf("JSONP() body = %v, want %v", res, expected)
	}
}

// TestString tests the String method.
func TestString(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	// Send a string response.
	testString := "Hello, world!"
	err := r.String(testString)
	if err != nil {
		t.Errorf("String() returned an error: %v", err)
	}

	// Check that the Content-Type header is set correctly.
	res := w.Body.String()
	if res != testString {
		t.Errorf("String() body = %v, want %v", res, testString)
	}
}

// TestError tests the Error method.
func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w, WithStatus(StatusBadRequest))

	errMessage := "This is an error"
	r.Error(errMessage)

	// Check that the status code is set to StatusInternalServerError
	// and the Content-Type header is set to MIMEApplicationJSONCharsetUTF8
	// and the response body contains the error message.
	if w.Code != StatusBadRequest {
		t.Errorf("Error() status code = %v, want %v",
			w.Code, StatusBadRequest)
	}

	// Check that the Content-Type header is set correctly.
	if !strings.Contains(w.Body.String(), errMessage) {
		t.Errorf("Error() body does not contain error message %v, got %v",
			errMessage, w.Body.String())
	}
}

// TestError_Empty tests the Error method.
func TestError_Empty(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	r.Error()

	// Check that the status code is set to StatusInternalServerError
	// and the Content-Type header is set to MIMEApplicationJSONCharsetUTF8
	// and the response body contains the error message.
	if w.Code != StatusInternalServerError {
		t.Errorf("Error() status code = %v, want %v",
			w.Code, StatusInternalServerError)
	}

	// Check that the Content-Type header is set correctly.
	want := `{"code":500,"message":"Internal Server Error"}`
	got := g.Trim(w.Body.String())
	if want != got {
		t.Errorf("Error() body = %v, want %v", got, want)
	}
}

// TestStream tests the Stream method.
func TestStream(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	// Send a stream response.
	streamData := "Stream data"
	reader := strings.NewReader(streamData)
	err := r.Stream(reader)
	if err != nil {
		t.Errorf("Stream() returned an error: %v", err)
	}

	// Check that the response body contains the stream data.
	respBody := w.Body.String()
	if respBody != streamData {
		t.Errorf("Stream() body = %v, want %v", respBody, streamData)
	}
}

// TestServeFile tests the ServeFile method.
func TestServeFile(t *testing.T) {
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
		resp := NewResponse(w)
		resp.ServeFile(r, tmpFile.Name()) // use the temp file
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

// TestServeFileAsDownload tests the ServeFileAsDownload method.
func TestServeFileAsDownload(t *testing.T) {
	// Create a test server.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := NewResponse(w)
		resp.ServeFileAsDownload("download.txt", []byte("Hello, download!"))
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

// TestRedirect tests the Redirect method.
func TestRedirect(t *testing.T) {
	w := httptest.NewRecorder()

	resp := NewResponse(w)
	err := resp.Redirect("http://example.com/bar")
	if err != nil {
		t.Errorf("Redirect failed: %v", err)
	}

	result := w.Result()
	if result.StatusCode != http.StatusFound {
		t.Errorf("Expected status code %d, got %d",
			http.StatusFound, result.StatusCode)
	}

	location := result.Header.Get("Location")
	if location != "http://example.com/bar" {
		t.Errorf("Expected Location header %s, got %s",
			"http://example.com/bar", location)
	}
}

// TestRedirectWithInvalidStatusCode tests the Redirect method with an
// invalid status code.
func TestRedirectWithInvalidStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	response := NewResponse(w)

	// Set an invalid status code.
	response.SetStatus(199)
	err := response.Redirect("http://example.org")

	if err != nil {
		t.Errorf("Redirect returned error: %v", err)
	}

	// Check that the status code is set to the default value.
	if got := w.Result().StatusCode; got != http.StatusFound {
		t.Errorf("Redirect() status code = %v, want %v",
			got, http.StatusFound)
	}

	// Check that the Location header is set correctly.
	if got := w.Header().Get("Location"); got != "http://example.org" {
		t.Errorf("Redirect() Location header = %v, want %v",
			got, "http://example.org")
	}
}

// TestNoContent tests the NoContent method.
func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()

	resp := NewResponse(w)
	err := resp.NoContent()
	if err != nil {
		t.Errorf("NoContent failed: %v", err)
	}

	result := w.Result()
	if result.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d",
			http.StatusNoContent, result.StatusCode)
	}
}

// TestHTML tests the HTML method.
func TestHTML(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w)

	htmlContent := "<!DOCTYPE html><html><body>Hello, World!</body></html>"
	err := resp.HTML(htmlContent)
	if err != nil {
		t.Errorf("HTML failed: %v", err)
	}

	result := w.Result()
	if result.Header.Get("Content-Type") != MIMETextHTMLCharsetUTF8 {
		t.Errorf("Expected Content-Type %s, got %s",
			MIMETextHTMLCharsetUTF8, result.Header.Get("Content-Type"))
	}

	body, _ := io.ReadAll(result.Body)
	if string(body) != htmlContent {
		t.Errorf("Unexpected response body: %s", body)
	}
}

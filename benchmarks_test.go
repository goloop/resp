package resp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test data structures
type testStruct struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

var (
	// Small payload for testing
	smallData = testStruct{
		Name:    "John Doe",
		Age:     30,
		Email:   "john@example.com",
		Address: "123 Main St",
	}

	// Medium payload for testing
	mediumData = []testStruct{
		{Name: "John Doe", Age: 30, Email: "john@example.com", Address: "123 Main St"},
		{Name: "Jane Doe", Age: 28, Email: "jane@example.com", Address: "456 Oak Ave"},
		{Name: "Bob Smith", Age: 35, Email: "bob@example.com", Address: "789 Pine Rd"},
	}

	// Large payload for testing
	largeData = make([]testStruct, 100)

	// HTML content for testing
	htmlContent = `<!DOCTYPE html>
		<html>
		<head><title>Test Page</title></head>
		<body>
			<h1>Hello, World!</h1>
			<p>This is a test page with some content.</p>
		</body>
		</html>`
)

func init() {
	// Initialize large dataset
	for i := 0; i < 100; i++ {
		largeData[i] = testStruct{
			Name:    "User Name",
			Age:     25 + i,
			Email:   "user@example.com",
			Address: "Street Address",
		}
	}
}

// helperNewRecorder creates a new response recorder for testing
func helperNewRecorder() http.ResponseWriter {
	return httptest.NewRecorder()
}

// BenchmarkJSONSmall benchmarks JSON response with small payload
func BenchmarkJSONSmall(b *testing.B) {
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		JSON(w, smallData)
	}
}

// BenchmarkJSONMedium benchmarks JSON response with medium payload
func BenchmarkJSONMedium(b *testing.B) {
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		JSON(w, mediumData)
	}
}

// BenchmarkJSONLarge benchmarks JSON response with large payload
func BenchmarkJSONLarge(b *testing.B) {
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		JSON(w, largeData)
	}
}

// BenchmarkJSONPSmall benchmarks JSONP response with small payload
func BenchmarkJSONPSmall(b *testing.B) {
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		JSONP(w, smallData, "callback")
	}
}

// BenchmarkHTML benchmarks HTML response
func BenchmarkHTML(b *testing.B) {
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		HTML(w, htmlContent)
	}
}

// BenchmarkString benchmarks String response
func BenchmarkString(b *testing.B) {
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		String(w, "Hello, World! This is a test string response.")
	}
}

// BenchmarkError benchmarks Error response
func BenchmarkError(b *testing.B) {
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Error(w, http.StatusNotFound, "Resource not found")
	}
}

// BenchmarkStream benchmarks Stream response
func BenchmarkStream(b *testing.B) {
	data := bytes.NewReader([]byte("This is test stream data that will be sent to the client"))
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data.Seek(0, io.SeekStart)
		Stream(w, data)
	}
}

// BenchmarkRedirect benchmarks Redirect response
func BenchmarkRedirect(b *testing.B) {
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Redirect(w, "https://example.com")
	}
}

// BenchmarkNoContent benchmarks NoContent response
func BenchmarkNoContent(b *testing.B) {
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		NoContent(w)
	}
}

// BenchmarkJSONWithCustomEncoder benchmarks JSON response with custom encoder
func BenchmarkJSONWithCustomEncoder(b *testing.B) {
	w := helperNewRecorder()
	customEncoder := func(w io.Writer, v interface{}) error {
		return json.NewEncoder(w).Encode(v)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		JSON(w, mediumData, ApplyJSONEncoder(customEncoder))
	}
}

// BenchmarkResponseChaining benchmarks response chaining operations
func BenchmarkResponseChaining(b *testing.B) {
	w := helperNewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		response := NewResponse(w).
			SetStatus(http.StatusOK).
			SetHeader("X-Custom-Header", "value")
		response.JSON(smallData)
	}
}

// BenchmarkServeFileAsDownload benchmarks serving file as download
func BenchmarkServeFileAsDownload(b *testing.B) {
	w := helperNewRecorder()
	data := []byte("This is test file content that will be downloaded")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ServeFileAsDownload(w, "test.txt", data)
	}
}

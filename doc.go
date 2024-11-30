// Package resp provides a powerful and flexible library for handling HTTP
// responses in Go web applications and APIs. It offers a clean, intuitive
// interface while maintaining compatibility with Go's standard net/http
// package.
//
// IMPORTANT: This package, like http.ResponseWriter, is NOT safe for
// concurrent use. ResponseWriter should only be written to from a single
// goroutine, and only until the HTTP handler returns. Writing to
// ResponseWriter after the handler returns or from multiple goroutines
// may result in corrupted responses or other undefined behavior.
//
// For concurrent operations, use channels to coordinate responses:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    resultChan := make(chan Result)
//	    go func() {
//	        // Perform processing...
//	        resultChan <- result
//	    }()
//	    result := <-resultChan
//	    if err := resp.JSON(w, result); err != nil {
//	        log.Printf("Failed to send response: %v", err)
//	    }
//	}
//
// Key Features:
//   - Streamlined response methods (JSON, HTML, String, etc.)
//   - Flexible header and cookie management
//   - Built-in error handling
//   - File streaming and downloads
//   - Support for custom JSON encoders
//   - Chainable response options
//   - Automatic content type handling
//
// Quick Response Methods:
//
// 1. JSON Responses:
//    Send JSON data with automatic content-type headers:
//
//    func HandleJSON(w http.ResponseWriter, r *http.Request) {
//        data := resp.R{"message": "Success", "code": 200}
//        if err := resp.JSON(w, data); err != nil {
//            log.Printf("JSON encoding failed: %v", err)
//        }
//    }
//
// 2. Custom JSON Encoding:
//    Use custom JSON encoders for specific needs:
//
//    import jsoniter "github.com/json-iterator/go"
//
//    func HandleCustomJSON(w http.ResponseWriter, r *http.Request) {
//        encoder := func(w io.Writer, v interface{}) error {
//            return jsoniter.NewEncoder(w).Encode(v)
//        }
//        resp := resp.NewResponse(w, resp.ApplyJSONEncoder(encoder))
//        if err := resp.JSON(data); err != nil {
//            log.Printf("Custom encoding failed: %v", err)
//        }
//    }
//
// 3. Error Handling:
//    Send error responses with custom codes and messages:
//
//    func HandleError(w http.ResponseWriter, r *http.Request) {
//        if err := validateInput(r); err != nil {
//            resp.Error(w, http.StatusBadRequest, "Invalid input")
//            return
//        }
//    }
//
// 4. File Downloads:
//    Serve files with custom headers:
//
//    func HandleDownload(w http.ResponseWriter, r *http.Request) {
//        opts := []resp.Option{
//            resp.AddContentDisposition("attachment", "report.pdf"),
//            resp.WithHeader("X-Download-Type", "report"),
//        }
//        err := resp.ServeFile(w, r, "path/to/file.pdf", opts...)
//        if err != nil {
//            log.Printf("Download failed: %v", err)
//        }
//    }
//
// 5. HTML Responses:
//    Send HTML content with proper content-type:
//
//    func HandleHTML(w http.ResponseWriter, r *http.Request) {
//        template := `
//        <!DOCTYPE html>
//        <html>
//        <head><title>Example</title></head>
//        <body><h1>Hello, World!</h1></body>
//        </html>`
//        if err := resp.HTML(w, template); err != nil {
//            log.Printf("Failed to send HTML: %v", err)
//        }
//    }
//
// Header Management:
//
// 1. Setting Custom Headers:
//
//    response := resp.NewResponse(w).
//        SetHeader("X-Custom-Header", "value").
//        SetHeader("X-Another-Header", "another-value")
//
// 2. Using Predefined Options:
//
//    opts := []resp.Option{
//        resp.WithHeader("Cache-Control", "no-cache"),
//        resp.AddETag("\"123456\""),
//        resp.AsApplicationJSON(),
//    }
//    resp.JSON(w, data, opts...)
//
// Cookie Management:
//
// 1. Setting Cookies:
//
//    cookie := &http.Cookie{
//        Name:     "session",
//        Value:    "abc123",
//        Path:     "/",
//        MaxAge:   3600,
//        HttpOnly: true,
//        Secure:   true,
//    }
//    response := resp.NewResponse(w).SetCookie(cookie)
//
// 2. Deleting Cookies:
//
//    response.ExpiredCookie("session")
//
// Default Behaviors:
//   - Status code defaults to 200 OK if not specified
//   - Content-Type is set automatically based on response type:
//     - JSON: "application/json; charset=utf-8"
//     - HTML: "text/html; charset=utf-8"
//     - Text: "text/plain; charset=utf-8"
//   - Errors are properly wrapped with context
//
// Best Practices:
//   - Always handle errors returned by response methods
//   - Use appropriate content types for responses
//   - Consider using options for fine-grained control
//   - Follow standard HTTP status code conventions
//
// For more examples and detailed documentation, visit:
// https://pkg.go.dev/github.com/goloop/resp

package resp

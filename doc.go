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
//   - Streamlined response methods (JSON, JSONP, XML, HTML, String, Blob)
//   - Safe-by-default encoding: a serialization error never leaves a
//     half-written success response (opt into streaming with WithStreamJSON)
//   - Validated JSONP callbacks (closes the reflected-XSS vector)
//   - Consistent error responses (HTTP status mirrored into the body code)
//   - REST helpers (Created, Accepted, NoContent, Redirect)
//   - Flexible header and cookie management
//   - File streaming and downloads
//   - Support for custom JSON encoders
//   - Chainable response options and state getters (Status, Header)
//   - Automatic content type handling
//   - Zero third-party dependencies
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
//    Plug in any encoder that satisfies JSONEncodeFunc for specific needs:
//
//    func HandleCustomJSON(w http.ResponseWriter, r *http.Request) {
//        encoder := func(w io.Writer, v any) error {
//            enc := json.NewEncoder(w)
//            enc.SetEscapeHTML(false)
//            return enc.Encode(v)
//        }
//        resp := resp.NewResponse(w, resp.ApplyJSONEncoder(encoder))
//        if err := resp.JSON(data); err != nil {
//            log.Printf("Custom encoding failed: %v", err)
//        }
//    }
//
// 3. Error Handling:
//    The first argument to Error is the HTTP status. It is written to the
//    status line and mirrored into the JSON body's "code" field, so the two
//    never disagree. Use WithErrorCode to set a distinct business code:
//
//    func HandleError(w http.ResponseWriter, r *http.Request) {
//        if err := validateInput(r); err != nil {
//            // HTTP 400, body {"code":400,"message":"Invalid input"}.
//        Use WithErrorSlug to add a machine-readable name beside them:
//
//            resp.Error(w, 409, "That name is taken",
//                resp.WithErrorSlug("slug_taken"))
//            // HTTP 409, body {"code":409,"error":"slug_taken","message":...}
//
//        That is the field a client switches on and a translation key is
//        looked up by; a number needs a registry mapping numbers to meanings,
//        and it says nothing in a log line or a failing test. It is omitted
//        when unset, so bodies that do not name one are unchanged.
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
// 6. XML and Binary Responses:
//    Serialize XML like JSON, or send raw bytes with an explicit type:
//
//    func HandleXML(w http.ResponseWriter, r *http.Request) {
//        resp.XML(w, note)
//    }
//
//    func HandlePDF(w http.ResponseWriter, r *http.Request) {
//        // Content-Length is set automatically.
//        resp.Blob(w, "application/pdf", pdfBytes)
//    }
//
// 7. REST Helpers:
//    Express common REST outcomes directly:
//
//    func CreateUser(w http.ResponseWriter, r *http.Request) {
//        user := create(r)
//        resp.Created(w, "/users/"+user.ID, user) // 201 + Location
//    }
//
//    func Enqueue(w http.ResponseWriter, r *http.Request) {
//        resp.Accepted(w, resp.R{"job_id": id}) // 202
//    }
//
// 8. JSONP:
//    The callback name is validated before anything is written; an unsafe
//    name returns ErrInvalidJSONPCallback and produces no output:
//
//    func HandleJSONP(w http.ResponseWriter, r *http.Request) {
//        cb := r.URL.Query().Get("callback")
//        if err := resp.JSONP(w, data, cb); err != nil {
//            resp.Error(w, http.StatusBadRequest, "invalid callback")
//        }
//    }
//
// Encoding safety:
//   - JSON, JSONP and XML buffer the payload and commit the status only
//     after a successful encode, so an encoding error can be reported before
//     the response is written — the handler can still send a proper error.
//   - WithStreamJSON opts into zero-copy streaming when memory matters more
//     than that guarantee.
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
//     - JSON:  "application/json; charset=utf-8"
//     - JSONP: "application/javascript; charset=utf-8"
//     - XML:   "application/xml; charset=utf-8"
//     - HTML:  "text/html; charset=utf-8"
//     - Text:  "text/plain"
//     - Blob:  "application/octet-stream" (unless a type is given)
//   - Errors are properly wrapped with context
//
// Best Practices:
//   - Always handle errors returned by response methods
//   - Use appropriate content types for responses
//   - Consider using options for fine-grained control
//   - Follow standard HTTP status code conventions
//
// For more examples and detailed documentation, visit:
// https://pkg.go.dev/github.com/goloop/resp/v2

package resp

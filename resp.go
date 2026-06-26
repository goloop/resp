package resp

import (
	"io"
	"net/http"
)

// R is a type alias for a map[string]any. It is designed to simplify
// the creation and manipulation of JSON objects in Go, making it easier to
// work with dynamic data structures commonly used in web applications and APIs.
// The use of `any` as the map value allows storing any Go value,
// offering the flexibility needed for JSON serialization.
//
// The alias `R` stands for "Response", emphasizing its utility in preparing
// data for HTTP responses. It streamlines the generation of JSON-encoded
// responses by providing a concise and readable way to construct JSON objects
// without the need for struct definitions or type assertions.
//
// Example Usage:
// The following example demonstrates how to use `R` to send a JSON response
// in an HTTP handler. The `resp.JSON` function leverages the `R` type to
// encode the map into a JSON object and write it to the HTTP response.
// This approach is particularly useful for endpoints that return JSON data,
// as it reduces boilerplate and improves code readability.
//
//	import (
//		"net/http"
//		"time"
//		"github.com/goloop/resp/v2"
//	)
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//		data := resp.R{
//			"Name": "Go Loop",
//			"Created": time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
//		}
//		if err := resp.JSON(w, data); err != nil {
//			// handle error
//		}
//	}
type R map[string]any

// JSON sends a JSON response to the client.
//
// This function wraps the process of setting up a JSON response by
// initializing a new Response object with optional configurations
// and then using it to encode the provided data as JSON.
// It sets the Content-Type header to "application/json" and writes
// the encoded JSON data to the client. This function simplifies the
// process of returning JSON data in your HTTP handlers.
//
// Parameters:
//   - w: The http.ResponseWriter that the JSON response will be written to.
//   - data: The data to be encoded as JSON. This can be any Go data structure,
//     including structs and slices.
//   - opts...: Optional configurations applied to the response. These can be
//     used to set custom headers, status codes, or other response settings.
//
// Returns:
// - An error if encoding the JSON fails. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    data := struct {
//	        Name string `json:"name"`
//	        Age  int    `json:"age"`
//	    }{
//	        Name: "John Doe",
//	        Age:  30,
//	    }
//
//	    // // With custom settings.
//	    // // Send a JSON response with HTTP status 200 OK and custom header.
//	    // if err := resp.JSON(w, data, WithStatus(http.StatusOK),
//	    //              WithHeader("X-Custom-Header", "value")); err != nil {
//	    //     // Handle error...
//	    // }
//
//	    // With default settings.
//	    if err := resp.JSON(w, data); err != nil {
//	        // Handle error...
//	    }
//	}
func JSON(w http.ResponseWriter, data any, opts ...Option) error {
	response := NewResponse(w, opts...)
	return response.JSON(data)
}

// JSONP sends a JSONP (JSON with Padding) response to the client.
//
// This function is particularly useful for serving JSON responses to web pages
// making requests across domain boundaries, where direct JSON responses might
// be blocked due to browser same-origin policies. JSONP wraps the JSON data in
// a callback function specified by the client, enabling cross-domain data
// communication in JavaScript by bypassing same-origin restrictions.
//
// Parameters:
//   - w: The http.ResponseWriter to which the JSONP response is written.
//   - data: The data to be encoded as JSON and wrapped in the
//     callback function. This can be any Go data structure, including
//     structs and slices.
//   - callback: The name of the callback function to wrap the JSON data.
//     This function name is typically provided by the client in the request
//     query string.
//   - opts...: Optional configurations applied to the response. These can
//     be used to set custom headers, status codes, or other response settings.
//
// Returns:
// - An error if encoding the JSON fails. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    data := map[string]string{"hello": "world"}
//
//	    // Extract callback function name from query parameters.
//	    callback := r.URL.Query().Get("callback")
//
//	    // Send a JSONP response.
//	    if err := resp.JSONP(w, data, callback); err != nil {
//	        // Handle error...
//	    }
//	}
func JSONP(
	w http.ResponseWriter,
	data any,
	callback string,
	opts ...Option,
) error {
	response := NewResponse(w, opts...)
	return response.JSONP(data, callback)
}

// String sends a plain text response to the client.
//
// This function simplifies the process of sending text-based responses,
// such as HTML or plain text, by encapsulating the creation of a Response
// object and setting the appropriate Content-Type header. It's particularly
// useful for endpoints that return non-JSON responses, providing a convenient
// way to return text data.
//
// Parameters:
//   - w: The http.ResponseWriter that the text response will be written to.
//   - data:    The string data to be sent as the response body. This can be
//     any text data, including HTML or plain text.
//   - opts...: Optional configurations applied to the response. These can be
//     used to set custom headers, status codes, or other response
//     settings, including changing the Content-Type from its default
//     ("text/plain").
//
// Returns:
// - An error if writing the response fails. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    message := "Hello, World!"
//
//	    // Send a plain text response with HTTP status 200 OK.
//	    if err := resp.String(w, message); err != nil {
//	        // Handle error...
//	    }
//	}
func String(w http.ResponseWriter, data string, opts ...Option) error {
	response := NewResponse(w, opts...)
	return response.String(data)
}

// Error sends a JSON error response.
//
// The status argument is the HTTP status code: it is written to the status
// line and, by default, mirrored into the JSON body's "code" field, so a
// single call keeps the two consistent. The status argument always wins over
// a WithStatus option; an out-of-range value falls back to
// StatusInternalServerError. When message is empty, the canonical text for
// the status code is used.
//
// To return a distinct application/business code in the body while keeping
// the HTTP status, add WithErrorCode.
//
// Parameters:
//   - w: The http.ResponseWriter to which the error response will be written.
//   - status: The HTTP status code (100-599). Also used as the body code
//     unless overridden with WithErrorCode.
//   - message: The error message for the response body. If empty, the
//     canonical message for the status code is used.
//   - opts...: Optional configurations applied to the response.
//
// Returns:
// - An error if writing the response fails. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    // HTTP 404, body {"code":404,"message":"Not Found"}.
//	    if err := resp.Error(w, http.StatusNotFound, ""); err != nil {
//	        // Handle error...
//	    }
//
//	    // HTTP 404, body {"code":1042,"message":"User not found"}.
//	    err := resp.Error(w, http.StatusNotFound, "User not found",
//	        resp.WithErrorCode(1042))
//	    if err != nil {
//	        // Handle error...
//	    }
//	}
func Error(
	w http.ResponseWriter,
	status int,
	message string,
	opts ...Option,
) error {
	response := NewResponse(w, opts...)
	return response.Error(status, message)
}

// Stream sends a stream response to the client.
//
// This function facilitates the sending of streaming data, such as file
// downloads or video streaming, by encapsulating the process of streaming
// from an io.Reader to the http.ResponseWriter. It can be configured with
// various options to set headers or status codes before streaming begins.
// The Content-Type of the response should be set appropriately using options,
// depending on the type of data being streamed.
//
// Parameters:
//   - w:The http.ResponseWriter to which the streaming response will
//     be written.
//   - reader:  An io.Reader from which data will be read and streamed to the
//     response. This could be a file, a buffer, or any other data
//     source implementing io.Reader.
//   - opts...: Optional configurations applied to the response. These can
//     be used to set custom headers, status codes, or other response
//     settings.
//
// Returns:
//   - An error if there's an issue writing the stream to the response.
//     Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    file, err := os.Open("video.mp4")
//	    if err != nil {
//	        resp.Error(w, resp.StatusNotFound)
//	        return
//	    }
//	    defer file.Close()
//
//	    // Stream the file to the response with the correct Content-Type.
//	    err := resp.Stream(w, file, WithHeader("Content-Type", "video/mp4"))
//	    if err != nil {
//	        log.Printf("Failed to stream file: %v", err)
//	    }
//	}
func Stream(
	w http.ResponseWriter,
	reader io.Reader,
	opts ...Option,
) error {
	response := NewResponse(w, opts...)
	return response.Stream(reader)
}

// ServeFile sends a file response to the client.
//
// This function is designed to simplify the process of serving static files
// (e.g., images, documents, media files) to the client over HTTP. It leverages
// http.ServeFile to handle the details of reading and transmitting the file,
// including setting appropriate headers for caching, content type detection,
// and handling range requests for efficient media streaming. The function can
// be configured with various options to set custom headers, status codes, or
// other response settings before the file is served.
//
// Parameters:
//   - w: The http.ResponseWriter to which the file will be written.
//   - r: The *http.Request object that initiated the file request. This is
//     required for handling conditional GET requests and range requests.
//   - filename: The path to the file that will be served. This must be a
//     valid file path accessible by the server.
//   - opts...: Optional configurations applied to the response. These can be
//     used to set custom headers, status codes, or other response settings.
//
// Returns:
//   - An error if there's an issue serving the file. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    // Assuming there's a query parameter 'file' with the filename.
//	    filename := r.URL.Query().Get("file")
//	    if filename == "" {
//	        resp.Error(w, "File not specified", http.StatusBadRequest)
//	        return
//	    }
//
//	    // Serve the file with Content-Disposition header for download.
//	    err := resp.ServeFile(w, r, filepath.Join("static", filename),
//	        resp.AddContentDisposition("attachment", filename))
//	    if err != nil {
//	        log.Printf("Failed to serve file: %v", err)
//	        resp.Error(w, "Failed to serve file", 500)
//	    }
//	}
func ServeFile(
	w http.ResponseWriter,
	r *http.Request,
	filename string,
	opts ...Option,
) error {
	response := NewResponse(w, opts...)
	return response.ServeFile(r, filename)
}

// ServeFileAsDownload sends a file as a download response to the client.
//
// This function is intended for scenarios where you need to serve
// dynamically generated content or files stored in memory as downloads,
// rather than serving files directly from the filesystem. It sets the
// Content-Disposition header to prompt the browser to treat the response
// as a file to be downloaded. The function can be configured with various
// options to set custom headers, status codes, or other response settings
// before the download is initiated.
//
// Parameters:
//   - w: The http.ResponseWriter to which the download response
//     will be written.
//   - filename: The filename to be used in the Content-Disposition header,
//     suggesting the name the browser should use to save the downloaded file.
//   - data: The byte slice containing the file data to be sent as the download.
//   - opts...: Optional configurations applied to the response. These can be
//     used to set custom headers, status codes, or other response settings.
//
// Returns:
//   - An error if there's an issue writing the download response.
//     Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    // Generate or retrieve the file data.
//	    fileData := []byte("Hello, world!")
//	    filename := "greeting.txt"
//
//	    // Send the file data as a download.
//	    err := ServeFileAsDownload(w, filename, fileData)
//	    if err != nil {
//	        log.Printf("Failed to serve file download: %v", err)
//	        resp.Error(w, "Failed to serve file download", 500)
//	    }
//	}
func ServeFileAsDownload(
	w http.ResponseWriter,
	filename string,
	data []byte,
	opts ...Option,
) error {
	response := NewResponse(w, opts...)
	return response.ServeFileAsDownload(filename, data)
}

// Redirect sends a redirect response to the client, instructing the browser
// to navigate to a different URL.
//
// This function simplifies sending HTTP redirects by automatically setting
// the status code to 302 Found, unless otherwise specified in the options.
// It's useful for handling requests where the resource has been moved or
// when the request should be redirected to a different endpoint.
//
// Parameters:
//   - w: The http.ResponseWriter to which the redirect response is written.
//   - url: The URL to which the client will be redirected.
//   - opts...: Optional configurations applied to the response. This can
//     include setting a specific status code using WithStatus if a different
//     type of redirect is required (e.g., 301 Moved Permanently).
//
// Returns:
//   - An error if there's an issue writing the redirect response.
//     Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    // Define the URL to redirect to.
//	    newURL := "https://example.com/new-page"
//
//	    // Redirect the request to the new URL.
//	    if err := Redirect(w, newURL); err != nil {
//	        log.Printf("Failed to redirect: %v", err)
//	    }
//	}
func Redirect(w http.ResponseWriter, url string, opts ...Option) error {
	options := []Option{WithStatus(StatusFound)}
	options = append(options, opts...)
	return NewResponse(w, options...).Redirect(url)
}

// NoContent sends a 204 No Content response to the client.
//
// This function is useful for endpoints that successfully process a request
// but do not need to return any data in the response body, such as a successful
// deletion of a resource or a successful update operation where no confirmation
// is needed. It sets the status code to 204 No Content by default, but this can
// be overridden using options to accommodate any additional requirements.
//
// Parameters:
//   - w: The http.ResponseWriter to which the no content response is written.
//   - opts...: Optional configurations applied to the response. This can be
//     used to set custom headers or other response settings as needed.
//
// Returns:
// - An error if there's an issue setting up the response. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    // Perform the deletion operation...
//
//	    // Respond with no content status to indicate success
//	    if err := resp.NoContent(w); err != nil {
//	        log.Printf("Failed to send no content response: %v", err)
//	    }
//	}
func NoContent(w http.ResponseWriter, opts ...Option) error {
	options := []Option{WithStatus(StatusNoContent)}
	options = append(options, opts...)
	return NewResponse(w, options...).NoContent()
}

// HTML sends an HTML response to the client.
//
// This function simplifies the process of sending HTML content as a response.
// It's particularly useful for serving web pages or HTML fragments. By default,
// it sets the Content-Type header to "text/html". This function can be
// configured with various options to set custom headers, status codes,
// or other response settings, making it versatile for web development needs.
//
// Parameters:
//   - w: The http.ResponseWriter to which the HTML content will be written.
//   - data: The HTML content to be sent as the response body. This should
//     be a valid HTML string.
//   - opts...: Optional configurations applied to the response. These can
//     be used to set custom headers, status codes, or other response settings.
//
// Returns:
//   - An error if there's an issue writing the HTML response. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    template := `
//	        <!DOCTYPE html>
//	        <html>
//	        <head><title>Example Page</title></head>
//	        <body>
//	            <h1>Hello, World!</h1>
//	            <p>This is example page.</p>
//	        </body>
//	        </html>
//	    `
//
//	    // Send the HTML content as a response.
//	    if err := HTML(w, template, WithStatus(http.StatusOK)); err != nil {
//	        // Handle error
//	        log.Printf("Failed to send HTML response: %v", err)
//	        resp.Error(w, "Failed to send HTML response", 500)
//	    }
//	}
func HTML(w http.ResponseWriter, data string, opts ...Option) error {
	return NewResponse(w, opts...).HTML(data)
}

// XML sends an XML response to the client.
//
// It mirrors JSON for XML payloads: the Content-Type defaults to
// "application/xml; charset=utf-8" and the data is encoded with
// encoding/xml. Like JSON, encoding happens into a buffer first, so a
// marshalling error is returned before the status line is committed.
//
// Parameters:
//   - w: The http.ResponseWriter that the XML response will be written to.
//   - data: The data to be encoded as XML (typically a struct with xml tags).
//   - opts...: Optional configurations applied to the response.
//
// Returns:
//   - An error if encoding the XML fails. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    type Note struct {
//	        XMLName xml.Name `xml:"note"`
//	        To      string   `xml:"to"`
//	    }
//	    if err := resp.XML(w, Note{To: "World"}); err != nil {
//	        // Handle error...
//	    }
//	}
func XML(w http.ResponseWriter, data any, opts ...Option) error {
	return NewResponse(w, opts...).XML(data)
}

// Blob sends raw bytes with an explicit content type to the client.
//
// It is the in-memory counterpart to ServeFile for arbitrary binary content
// (PDFs, images, generated documents) that should be displayed rather than
// downloaded. Content-Length is set automatically; if contentType is empty,
// "application/octet-stream" is used.
//
// Parameters:
//   - w: The http.ResponseWriter to which the bytes will be written.
//   - contentType: The MIME type for the data (e.g. resp.MIMEApplicationPDF).
//     If empty, application/octet-stream is used.
//   - data: The raw bytes to send as the response body.
//   - opts...: Optional configurations applied to the response.
//
// Returns:
//   - An error if writing the response fails. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    pdf := renderInvoice()
//	    if err := resp.Blob(w, "application/pdf", pdf); err != nil {
//	        // Handle error...
//	    }
//	}
func Blob(
	w http.ResponseWriter,
	contentType string,
	data []byte,
	opts ...Option,
) error {
	return NewResponse(w, opts...).Blob(contentType, data)
}

// Created sends a 201 Created response with an optional Location header and
// an optional JSON body.
//
// This is the idiomatic helper for REST endpoints that create a resource:
// it sets the status to 201, points Location at the new resource and encodes
// the representation as JSON. A nil body sends only the status and headers.
//
// Parameters:
//   - w: The http.ResponseWriter to which the response will be written.
//   - location: The URL of the newly created resource (Location header).
//     Omitted from the header when empty.
//   - data: The created resource to encode as JSON, or nil for no body.
//   - opts...: Optional configurations applied to the response.
//
// Returns:
//   - An error if encoding or writing the response fails. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    user := createUser(r)
//	    if err := resp.Created(w, "/users/"+user.ID, user); err != nil {
//	        // Handle error...
//	    }
//	}
func Created(
	w http.ResponseWriter,
	location string,
	data any,
	opts ...Option,
) error {
	return NewResponse(w, opts...).Created(location, data)
}

// Accepted sends a 202 Accepted response with an optional JSON body.
//
// It is the idiomatic helper for endpoints that accept a request for
// asynchronous processing: it sets the status to 202 and optionally encodes
// a tracking representation (e.g. a job handle) as JSON. A nil body sends
// only the status.
//
// Parameters:
//   - w: The http.ResponseWriter to which the response will be written.
//   - data: The payload to encode as JSON (e.g. a job descriptor), or nil.
//   - opts...: Optional configurations applied to the response.
//
// Returns:
//   - An error if encoding or writing the response fails. Otherwise, nil.
//
// Example usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    job := enqueue(r)
//	    if err := resp.Accepted(w, resp.R{"job_id": job.ID}); err != nil {
//	        // Handle error...
//	    }
//	}
func Accepted(w http.ResponseWriter, data any, opts ...Option) error {
	return NewResponse(w, opts...).Accepted(data)
}

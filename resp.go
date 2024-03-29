package resp

import (
	"io"
	"net/http"
)

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
//	func MyHandler(w http.ResponseWriter, r *http.Request) {
//	    myData := struct {
//	        Name string `json:"name"`
//	        Age  int    `json:"age"`
//	    }{
//	        Name: "John Doe",
//	        Age:  30,
//	    }
//
//	    // Send a JSON response with HTTP status 200 OK and custom header
//	    if err := resp.JSON(w, myData, WithStatus(http.StatusOK),
//	                 WithHeader("X-Custom-Header", "value")); err != nil {
//	        // Handle error
//	    }
//	}
//
// This example demonstrates how to use the JSON function to send a JSON
// response with a 200 OK status and a custom header. The myData object
// is encoded as JSON and sent to the client.
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
//	func MyJSONPHandler(w http.ResponseWriter, r *http.Request) {
//	    data := map[string]string{"hello": "world"}
//
//	    // Extract callback function name from query parameters
//	    callback := r.URL.Query().Get("callback")
//
//	    // Send a JSONP response
//	    err := resp.JSONP(w, data, callback, WithStatus(http.StatusOK))
//	    if err != nil {
//	        // Handle error
//	        http.Error(w, err.Error(), http.StatusInternalServerError)
//	        return
//	    }
//	}
//
// In this example, the JSONP function is used to send a JSONP response with
// the data `{"hello": "world"}` wrapped in a callback function specified by
// the client.
// The callback name is retrieved from the request's query parameters, allowing
// the client-side JavaScript to process the response data despite cross-domain
// restrictions.
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
//	func MyTextHandler(w http.ResponseWriter, r *http.Request) {
//	    message := "Hello, World!"
//
//	    // Send a plain text response with HTTP status 200 OK
//	    err := resp.String(w, message, WithStatus(http.StatusOK))
//	    if err != nil {
//	        // Handle error
//	        http.Error(w, "Failed to write response",
//	                   http.StatusInternalServerError)
//	    }
//	}
//
// In this example, the String function is used to send a plain text response
// containing "Hello, World!" to the client. The WithStatus option is used to
// explicitly set the HTTP status code to 200 OK. This demonstrates how to use
// the String function to return text responses, making it easy to build
// text-based endpoints in a web application.
func String(w http.ResponseWriter, data string, opts ...Option) error {
	response := NewResponse(w, opts...)
	return response.String(data)
}

// Error sends an error response with a specified HTTP status code and
// error message.
//
// This function is designed to facilitate sending error responses with
// meaningful status codes and messages. It allows for flexible error
// reporting by accepting an optional message parameter, making it suitable
// for endpoints that need to provide more context about an error.
// The Content-Type of the response is set to "text/plain" by default,
// but this can be adjusted using options.
//
// Parameters:
//   - w: The http.ResponseWriter to which the error response will be written.
//   - status: The HTTP status code associated with the error. This should be
//     a valid HTTP status code that indicates an error (4xx or 5xx).
//   - message...: Optional. One or more strings that will be concatenated and
//     sent as the error message in the response body. If no message is
//     provided, a default message based on the status code will be used.
//
// Returns:
// - An error if writing the response fails. Otherwise, nil.
//
// Example usage:
//
//	func MyErrorHandler(w http.ResponseWriter, r *http.Request) {
//	    // Example of handling a not found error
//	    err := resp.Error(w, http.StatusNotFound,
//	                     "The requested resource was not found")
//	    if err != nil {
//	        // Additional error handling logic if needed
//	        log.Printf("Failed to send error response: %v", err)
//	    }
//	}
//
// In this example, the Error function is used to send an error
// response with a 404 Not Found status code and a custom error message.
// This demonstrates how to use the Error function to communicate
// specific error conditions back to the client in a RESTful API
// or web application.
func Error(w http.ResponseWriter, status int, message ...string) error {
	response := NewResponse(w, WithStatus(status))
	return response.Error(message...)
}

// Stream sends a stream response to the client.
//
// This function facilitates the sending of streaming data, such as file downloads or
// video streaming, by encapsulating the process of streaming from an io.Reader to the
// http.ResponseWriter. It can be configured with various options to set headers or
// status codes before streaming begins. The Content-Type of the response should be set
// appropriately using options, depending on the type of data being streamed.
//
// Parameters:
//   - w: The http.ResponseWriter to which the streaming response will be written.
//   - reader: An io.Reader from which data will be read and streamed to the response. This
//     could be a file, a buffer, or any other data source implementating io.Reader.
//   - opts...: Optional configurations applied to the response. These can be used to set
//     custom headers, status codes, or other response settings.
//
// Returns:
// - An error if there's an issue writing the stream to the response. Otherwise, nil.
//
// Example usage:
//
//	func FileStreamHandler(w http.ResponseWriter, r *http.Request) {
//	    file, err := os.Open("video.mp4")
//	    if err != nil {
//	        // Handle error (file not found, etc.)
//	        http.Error(w, "File not found", http.StatusNotFound)
//	        return
//	    }
//	    defer file.Close()
//
//	    // Stream the file to the response with the correct Content-Type
//	    if err := resp.Stream(w, file, WithHeader("Content-Type", "video/mp4")); err != nil {
//	        // Handle error streaming file
//	        log.Printf("Failed to stream file: %v", err)
//	    }
//	}
//
// This example demonstrates how to use the Stream function to serve a video file directly
// to the client. It sets the Content-Type header to "video/mp4" to ensure the client handles
// the response correctly. The function reads data from the file and streams it to the response,
// enabling efficient serving of large files or real-time data.
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
// This function is designed to simplify the process of serving static files (e.g., images,
// documents, media files) to the client over HTTP. It leverages http.ServeFile to handle
// the details of reading and transmitting the file, including setting appropriate headers
// for caching, content type detection, and handling range requests for efficient media
// streaming. The function can be configured with various options to set custom headers,
// status codes, or other response settings before the file is served.
//
// Parameters:
//   - w: The http.ResponseWriter to which the file will be written.
//   - r: The *http.Request object that initiated the file request. This is required for
//     handling conditional GET requests and range requests.
//   - filename: The path to the file that will be served. This must be a valid file path
//     accessible by the server.
//   - opts...: Optional configurations applied to the response. These can be used to set
//     custom headers, status codes, or other response settings.
//
// Returns:
// - An error if there's an issue serving the file. Otherwise, nil.
//
// Example usage:
//
//	func MyFileHandler(w http.ResponseWriter, r *http.Request) {
//	    // Assuming there's a query parameter 'file' with the filename
//	    filename := r.URL.Query().Get("file")
//	    if filename == "" {
//	        http.Error(w, "File not specified", http.StatusBadRequest)
//	        return
//	    }
//
//	    // Serve the requested file with Content-Disposition header for download
//	    if err := resp.ServeFile(w, r, filepath.Join("static", filename), WithHeader("Content-Disposition", "attachment; filename=\""+filename+"\"")); err != nil {
//	        // Handle error (file not found, permissions issue, etc.)
//	        log.Printf("Failed to serve file: %v", err)
//	        http.Error(w, "Failed to serve file", http.StatusInternalServerError)
//	    }
//	}
//
// In this example, the ServeFile function is used to send a requested file to the client,
// with an additional header to suggest that the browser treat the response as a file to
// download. The file is expected to be located within a "static" directory, and its name
// is retrieved from a query parameter. This setup is useful for dynamically serving files
// based on client requests, including handling downloads.
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
// This function is intended for scenarios where you need to serve dynamically generated
// content or files stored in memory as downloads, rather than serving files directly from
// the filesystem. It sets the Content-Disposition header to prompt the browser to treat
// the response as a file to be downloaded. The function can be configured with various
// options to set custom headers, status codes, or other response settings before the
// download is initiated.
//
// Parameters:
//   - w: The http.ResponseWriter to which the download response will be written.
//   - filename: The filename to be used in the Content-Disposition header, suggesting the
//     name the browser should use to save the downloaded file.
//   - data: The byte slice containing the file data to be sent as the download.
//   - opts...: Optional configurations applied to the response. These can be used to set
//     custom headers, status codes, or other response settings.
//
// Returns:
// - An error if there's an issue writing the download response. Otherwise, nil.
//
// Example usage:
//
//	func MyDownloadHandler(w http.ResponseWriter, r *http.Request) {
//	    // Generate or retrieve the file data
//	    fileData := []byte("Hello, world!")
//	    filename := "greeting.txt"
//
//	    // Send the file data as a download
//	    if err := ServeFileAsDownload(w, filename, fileData, WithStatus(http.StatusOK)); err != nil {
//	        // Handle error
//	        log.Printf("Failed to serve file download: %v", err)
//	        http.Error(w, "Failed to serve file download", http.StatusInternalServerError)
//	    }
//	}
//
// In this example, the ServeFileAsDownload function is used to send dynamically generated
// file data as a download to the client. The filename "greeting.txt" is suggested to the
// client for saving the file. Additional response configurations, such as setting the
// HTTP status code, can be applied through the options parameter.
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
//   - opts...: Optional configurations applied to the response. This can include
//     setting a specific status code using WithStatus if a different
//     type of redirect is required (e.g., 301 Moved Permanently).
//
// Returns:
// - An error if there's an issue writing the redirect response. Otherwise, nil.
//
// Example usage:
//
//	func MyRedirectHandler(w http.ResponseWriter, r *http.Request) {
//	    // Define the URL to redirect to
//	    newURL := "https://example.com/new-page"
//
//	    // Redirect the request to the new URL
//	    if err := Redirect(w, newURL); err != nil {
//	        // Handle potential error
//	        log.Printf("Failed to redirect: %v", err)
//	    }
//	}
//
// In this example, the Redirect function is used within an HTTP handler to redirect
// the client to a new URL. By default, it uses the 302 Found status code for the
// redirect, but this can be changed by passing an additional option, like
// WithStatus(http.StatusMovedPermanently), to specify a different redirect status code.
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
//   - opts...: Optional configurations applied to the response. This can be used to
//     set custom headers or other response settings as needed.
//
// Returns:
// - An error if there's an issue setting up the response. Otherwise, nil.
//
// Example usage:
//
//	func MyDeleteHandler(w http.ResponseWriter, r *http.Request) {
//	    // Perform the deletion operation...
//
//	    // Respond with no content status to indicate success
//	    if err := resp.NoContent(w); err != nil {
//	        // Handle potential error
//	        log.Printf("Failed to send no content response: %v", err)
//	    }
//	}
//
// This example demonstrates how to use the NoContent function within an HTTP handler
// that processes deletion requests. After successfully deleting a resource, it sends
// a 204 No Content response to the client to indicate that the request was successfully
// processed, but there is no additional content to return in the response body.
func NoContent(w http.ResponseWriter, opts ...Option) error {
	options := []Option{WithStatus(StatusNoContent)}
	options = append(options, opts...)
	return NewResponse(w, options...).NoContent()
}

// HTML sends an HTML response to the client.
//
// This function simplifies the process of sending HTML content as a response. It's
// particularly useful for serving web pages or HTML fragments. By default, it sets
// the Content-Type header to "text/html". This function can be configured with
// various options to set custom headers, status codes, or other response settings,
// making it versatile for web development needs.
//
// Parameters:
//   - w: The http.ResponseWriter to which the HTML content will be written.
//   - data: The HTML content to be sent as the response body. This should be a valid
//     HTML string.
//   - opts...: Optional configurations applied to the response. These can be used to
//     set custom headers, status codes, or other response settings.
//
// Returns:
// - An error if there's an issue writing the HTML response. Otherwise, nil.
//
// Example usage:
//
//	func MyPageHandler(w http.ResponseWriter, r *http.Request) {
//	    // Define the HTML content
//	    htmlContent := `
//	        <!DOCTYPE html>
//	        <html>
//	        <head><title>My Page</title></head>
//	        <body>
//	            <h1>Hello, World!</h1>
//	            <p>This is my page.</p>
//	        </body>
//	        </html>
//	    `
//
//	    // Send the HTML content as a response
//	    if err := HTML(w, htmlContent, WithStatus(http.StatusOK)); err != nil {
//	        // Handle error
//	        log.Printf("Failed to send HTML response: %v", err)
//	        http.Error(w, "Failed to send HTML response", http.StatusInternalServerError)
//	    }
//	}
//
// In this example, the HTML function is used to send an HTML document as a response to
// the client. The document includes a simple message wrapped in basic HTML structure.
// The WithStatus option is used to explicitly set the HTTP status code to 200 OK.
func HTML(w http.ResponseWriter, data string, opts ...Option) error {
	return NewResponse(w, opts...).HTML(data)
}

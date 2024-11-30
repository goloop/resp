package resp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/goloop/g"
)

// JSONEncodeFunc represents a function that encodes the provided data
// into JSON and writes it to the provided io.Writer.
// This allows for custom JSON encoding strategies.
//
// Example Usage:
//
//	// Using a custom JSON encoder with jsoniter
//	import jsoniter "github.com/json-iterator/go"
//
//	customEncoder := func(w io.Writer, v interface{}) error {
//	    return jsoniter.NewEncoder(w).Encode(v)
//	}
//
//	resp.JSON(w, data, resp.ApplyJSONEncoder(customEncoder))
type JSONEncodeFunc func(w io.Writer, v interface{}) error

// Response represents an HTTP response.
// It provides methods for setting headers, cookies, and writing data
// to the response body. It can be customized using various options.
//
// Example Usage:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	    response := resp.NewResponse(w, resp.WithStatus(http.StatusOK))
//	    response.SetHeader("X-Custom-Header", "value")
//	    response.JSON(resp.R{"message": "Hello, World!"})
//	}
type Response struct {
	httpWriter     http.ResponseWriter
	statusCode     int
	jsonEncodeFunc JSONEncodeFunc
}

// NewResponse creates a new instance of Response with the provided
// http.ResponseWriter and options. It applies the provided options
// to the response and returns the pointer to the created response.
//
// Example Usage:
//
//	response := resp.NewResponse(w, resp.WithStatus(http.StatusOK),
//	    resp.AsApplicationJSON(),
//	    resp.ApplyJSONEncoder(customEncoder))
func NewResponse(w http.ResponseWriter, opts ...Option) *Response {
	// Create a new response with the provided http.ResponseWriter.
	response := &Response{
		httpWriter:     w,
		statusCode:     StatusUndefined,
		jsonEncodeFunc: nil,
	}

	// Apply the provided options to the response.
	for _, opt := range opts {
		response = opt(response)
	}

	return response
}

// prepare prepares the response by setting the default status
// code and content type.
//
// If the content type is not already set and a default content type
// is provided, it sets the default content type for the response.
// If the status code is not already set, it sets the default status
// code for the response.
func (r *Response) prepare(defStatus int, defContentType ...string) {
	// Set the default content type if it is not already set.
	_, ok := r.httpWriter.Header()[HeaderContentType]
	if !ok && len(defContentType) > 0 {
		r.httpWriter.Header().Set(HeaderContentType, defContentType[0])
	}

	// Set the default status code if it is not already set.
	if r.statusCode == StatusUndefined {
		r.statusCode = defStatus
	}
}

// SetJSONEncoder sets the custom JSON encoder function for the response
// and returns pointer to the modified response object.
func (r *Response) SetJSONEncoder(f JSONEncodeFunc) *Response {
	r.jsonEncodeFunc = f
	return r
}

// GetJSONEncoder returns the JSON encoder function of the response.
// If the JSON encoder function is not set, it returns nil.
func (r *Response) GetJSONEncoder() JSONEncodeFunc {
	return r.jsonEncodeFunc
}

// SetStatus sets the status code of the response and returns
// the modified response.
func (r *Response) SetStatus(code int) *Response {
	r.statusCode = code
	return r
}

// SetHeader sets the header with the provided key and value(s) and
// returns the modified response.
func (r *Response) SetHeader(key string, value ...string) *Response {
	// If the header can contain only one value, use first value only.
	if g.In(key, singleHeaders...) && len(value) > 0 {
		r.httpWriter.Header().Set(key, value[0])
		return r
	}

	// Set the header with the provided key and value(s).
	r.httpWriter.Header().Set(key, strings.Join(value, ","))
	return r
}

// AddHeader adds into header with the provided key and value(s) and
// returns the modified response.
func (r *Response) AddHeader(key string, value ...string) *Response {
	// If the header can contain only one value, use first value only.
	if g.In(key, singleHeaders...) && len(value) > 0 {
		r.SetHeader(key, value[0])
		return r
	}

	// Add the header with the provided key and value(s).
	for _, v := range value {
		r.httpWriter.Header().Add(key, v)
	}

	return r
}

// DelHeader deletes the header with the provided key from the response
// and returns the modified response.
func (r *Response) DelHeader(key string) *Response {
	r.httpWriter.Header().Del(key)
	return r
}

// ClearHeaders deletes all headers from the response and returns the
// modified response.
func (r *Response) ClearHeaders() *Response {
	for k := range r.httpWriter.Header() {
		r.httpWriter.Header().Del(k)
	}
	return r
}

// SetCookie sets a cookie in the response and returns the modified response.
func (r *Response) SetCookie(cookie *http.Cookie) *Response {
	http.SetCookie(r.httpWriter, cookie)
	return r
}

// BindCookie binds a cookie to the response. If a cookie with the same
// name already exists, it replaces the existing cookie with the new one.
//
// If a cookie already exists, it will be deleted and a new one will be re-set.
// If there are multiple cookies with the same name, they will all be deleted.
func (r *Response) BindCookie(cookie *http.Cookie) *Response {
	// Add the new one.
	r.DelCookie(cookie.Name)
	http.SetCookie(r.httpWriter, cookie)
	return r
}

// DelCookie deletes a cookie with the specified name from the response.
// It removes the cookie from the response's header and returns the
// modified response.
//
// Pay attention. All cookies with this name will be deleted.
func (r *Response) DelCookie(name string) *Response {
	// Get all existing cookies.
	// Filter cookies by name, without the one we want to add.
	filteredCookies := []string{}
	for _, c := range r.httpWriter.Header().Values(HeaderSetCookie) {
		if !strings.HasPrefix(c, name+"=") {
			filteredCookies = append(filteredCookies, c)
		}
	}

	// Remove all existing cookies, and add the filtered ones back.
	r.httpWriter.Header().Del(HeaderSetCookie)
	for _, c := range filteredCookies {
		r.httpWriter.Header().Add(HeaderSetCookie, c)
	}

	return r
}

// ClearCookies deletes all cookies from the response and returns the
// modified response.
func (r *Response) ClearCookies() *Response {
	r.httpWriter.Header().Del(HeaderSetCookie)
	return r
}

// ExpiredCookie expires a cookie with the specified name from the response.
func (r *Response) ExpiredCookie(name string) *Response {
	expiredCookie := &http.Cookie{
		Name:    name,
		Value:   "deleted",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}

	http.SetCookie(r.httpWriter, expiredCookie)
	return r
}

// JSON sends a JSON response.
// If the status code is not set - StatusOK will be set.
// If ContentType isn't defined - MIMEApplicationJSON will be used by default.
func (r *Response) JSON(data any) error {
	r.prepare(StatusOK, MIMEApplicationJSONCharsetUTF8)
	r.httpWriter.WriteHeader(r.statusCode)

	if r.jsonEncodeFunc != nil {
		if err := r.jsonEncodeFunc(r.httpWriter, data); err != nil {
			return fmt.Errorf("custom JSON encoder failed: %w", err)
		}
		return nil
	}

	if err := json.NewEncoder(r.httpWriter).Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON response: %w", err)
	}
	return nil
}

// JSONP sends a JSONP response.
// If the status code is not set - StatusOK will be set.
// If ContentType isn't defined - MIMEApplicationJavaScript will
// be used by default.
func (r *Response) JSONP(data any, callback string) error {
	r.prepare(StatusOK, MIMEApplicationJavaScriptCharsetUTF8)
	r.httpWriter.WriteHeader(r.statusCode)

	var buf bytes.Buffer

	var err error
	if r.jsonEncodeFunc != nil {
		err = r.jsonEncodeFunc(&buf, data)
		if err != nil {
			return fmt.Errorf("custom JSON encoder failed in JSONP: %w", err)
		}
	} else {
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			return fmt.Errorf("failed to encode JSONP data: %w", err)
		}
	}

	// Remove the trailing newline character if present.
	jsonData := buf.Bytes()
	if len(jsonData) > 0 && jsonData[len(jsonData)-1] == '\n' {
		jsonData = jsonData[:len(jsonData)-1]
	}

	// Write the JSONP response.
	_, err = fmt.Fprintf(r.httpWriter, "%s(%s);", callback, jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSONP response: %w", err)
	}

	return nil
}

// String sends a string response.
// If the status code is not set - StatusOK will be set.
// If ContentType isn't defined - MIMETextPlain will be used by default.
func (r *Response) String(data string) error {
	r.prepare(StatusOK, MIMETextPlain)
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := r.httpWriter.Write([]byte(data))
	return err
}

// Error sends an error response.
// If no error description is passed, it will be generated from the
// status code from the response. If more than one message is sent,
// only the first one will be used.
//
// If the status code isn't set - StatusInternalServerError will be set.
func (r *Response) Error(code int, message string) error {
	if r.statusCode == StatusUndefined {
		r.statusCode = StatusInternalServerError
	}

	return r.JSON(newErrorResponse(code, message))
}

// Stream sends a stream response.
func (r *Response) Stream(data io.Reader) error {
	r.prepare(StatusOK, MIMEOctetStream)
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := io.Copy(r.httpWriter, data)
	return err
}

// File sends a file response.
func (r *Response) ServeFile(req *http.Request, file string) error {
	r.prepare(StatusOK, MIMEOctetStream)

	// The http.ServeFile function from the net/http package independently
	// sets the response headers and status code before starting the file
	// transfer, no need: r.httpWriter.WriteHeader(r.statusCode)
	http.ServeFile(r.httpWriter, req, file)
	return nil
}

// ServeFileAsDownload sends a file as download response.
func (r *Response) ServeFileAsDownload(fileName string, data []byte) error {
	r.httpWriter.Header().Set(
		HeaderContentDisposition,
		"attachment; filename=\""+fileName+"\"",
	)

	r.prepare(StatusOK, MIMEOctetStream)
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := r.httpWriter.Write(data)
	return err
}

// Redirect sends an HTTP redirect to the specified URL.
func (r *Response) Redirect(url string) error {
	r.prepare(StatusFound)
	s := r.statusCode

	if s < StatusMultipleChoices || s > StatusPermanentRedirect {
		s = StatusFound // default to 302 if an invalid status code is provided
	}

	r.httpWriter.Header().Set("Location", url) // redirect to the specified URL
	r.httpWriter.WriteHeader(s)
	return nil
}

// NoContent sends a 204 No Content response.
func (r *Response) NoContent() error {
	r.SetStatus(StatusNoContent)
	r.prepare(StatusNoContent)
	r.httpWriter.WriteHeader(http.StatusNoContent)
	return nil
}

// HTML sends an HTML response.
func (r *Response) HTML(html string) error {
	r.prepare(http.StatusOK, MIMETextHTMLCharsetUTF8)
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := r.httpWriter.Write([]byte(html))
	return err
}

package resp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// JSONEncodeFunc represents a function that encodes the provided data
// into JSON and writes it to the provided io.Writer.
// This allows for custom JSON encoding strategies.
//
// Example Usage:
//
//	// Using a custom JSON encoder.
//	customEncoder := func(w io.Writer, v any) error {
//	    return json.NewEncoder(w).Encode(v)
//	}
//
//	resp.JSON(w, data, resp.ApplyJSONEncoder(customEncoder))
type JSONEncodeFunc func(w io.Writer, v any) error

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
	streamJSON     bool
	errorCode      *int
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

// encodeJSON serializes data into w using the custom encoder when one is
// configured, otherwise the standard library encoder. The returned error
// is already wrapped with context and supports errors.Is/As.
func (r *Response) encodeJSON(w io.Writer, data any) error {
	if r.jsonEncodeFunc != nil {
		if err := r.jsonEncodeFunc(w, data); err != nil {
			return fmt.Errorf("custom JSON encoder failed: %w", err)
		}
		return nil
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON response: %w", err)
	}
	return nil
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

// Status returns the status code currently configured for the response.
// Before any status is set (via SetStatus, WithStatus or a terminal write
// method) it returns StatusUndefined (0). This is primarily useful for
// middleware and tests that need to inspect the intended status.
func (r *Response) Status() int {
	return r.statusCode
}

// Header returns the underlying response header map, allowing callers to
// read or modify headers directly. It is the same map returned by the
// wrapped http.ResponseWriter.Header.
func (r *Response) Header() http.Header {
	return r.httpWriter.Header()
}

// SetHeader sets the header with the provided key and value(s) and
// returns the modified response.
func (r *Response) SetHeader(key string, value ...string) *Response {
	// If the header can contain only one value, use first value only.
	if _, single := singleHeaders[key]; single && len(value) > 0 {
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
	if _, single := singleHeaders[key]; single && len(value) > 0 {
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
	clear(r.httpWriter.Header())
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
	prefix := name + "="
	existing := r.httpWriter.Header().Values(HeaderSetCookie)
	filtered := make([]string, 0, len(existing))
	for _, c := range existing {
		if !strings.HasPrefix(c, prefix) {
			filtered = append(filtered, c)
		}
	}

	// Remove all existing cookies, and add the filtered ones back.
	r.httpWriter.Header().Del(HeaderSetCookie)
	for _, c := range filtered {
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
//
// By default the payload is encoded into an in-memory buffer first, so that
// a serialization failure is reported BEFORE the status line is written and
// the calling handler can still respond with a proper error. Opt into
// zero-copy streaming with WithStreamJSON when memory matters more than this
// guarantee (see WithStreamJSON for the trade-offs).
func (r *Response) JSON(data any) error {
	r.prepare(StatusOK, MIMEApplicationJSONCharsetUTF8)

	// Streaming mode (opt-in): encode straight into the ResponseWriter.
	// An encoding error may surface only after the status line and part of
	// the body were already flushed — the caller accepts that explicitly.
	if r.streamJSON {
		r.httpWriter.WriteHeader(r.statusCode)
		return r.encodeJSON(r.httpWriter, data)
	}

	// Buffered mode (default): commit the status only on a successful
	// encode, so a failure leaves the response writer untouched.
	buf := getBuf()
	defer putBuf(buf)
	if err := r.encodeJSON(buf, data); err != nil {
		return err
	}

	r.httpWriter.WriteHeader(r.statusCode)
	_, err := r.httpWriter.Write(buf.Bytes())
	return err
}

// JSONIndent sends an indented (pretty-printed) JSON response — handy for
// human-facing or debug endpoints. Unlike JSON it always uses the standard
// library encoder with the given prefix and indent, ignoring any custom
// encoder set via ApplyJSONEncoder (indentation is the whole point). Encoding
// is buffered, so an error is reported before the status line is committed.
func (r *Response) JSONIndent(data any, prefix, indent string) error {
	r.prepare(StatusOK, MIMEApplicationJSONCharsetUTF8)

	buf := getBuf()
	defer putBuf(buf)
	enc := json.NewEncoder(buf)
	enc.SetIndent(prefix, indent)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON response: %w", err)
	}

	r.httpWriter.WriteHeader(r.statusCode)
	_, err := r.httpWriter.Write(buf.Bytes())
	return err
}

// JSONP sends a JSONP response.
// If the status code is not set - StatusOK will be set.
// If ContentType isn't defined - MIMEApplicationJavaScript will
// be used by default.
//
// The callback name is validated against a strict JavaScript-identifier
// whitelist before anything is written. An unsafe name (containing quotes,
// brackets, markup, etc.) yields ErrInvalidJSONPCallback and no output,
// which prevents the callback padding from becoming a reflected-XSS vector.
func (r *Response) JSONP(data any, callback string) error {
	if !validJSONPCallback(callback) {
		return fmt.Errorf("%w: %q", ErrInvalidJSONPCallback, callback)
	}

	r.prepare(StatusOK, MIMEApplicationJavaScriptCharsetUTF8)

	// Encode into a buffer first: this both lets us trim the trailing
	// newline and keeps the status uncommitted until encoding succeeds.
	buf := getBuf()
	defer putBuf(buf)
	if err := r.encodeJSON(buf, data); err != nil {
		return err
	}
	jsonData := bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})

	// Write the JSONP response.
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := fmt.Fprintf(r.httpWriter, "%s(%s);", callback, jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSONP response: %w", err)
	}

	return nil
}

// XML sends an XML response.
// If the status code is not set - StatusOK will be set.
// If ContentType isn't defined - MIMEApplicationXMLCharsetUTF8 is used.
//
// Like JSON, the payload is encoded into a buffer first so a marshalling
// error is reported before the status line is committed.
func (r *Response) XML(data any) error {
	r.prepare(StatusOK, MIMEApplicationXMLCharsetUTF8)

	buf := getBuf()
	defer putBuf(buf)
	if err := xml.NewEncoder(buf).Encode(data); err != nil {
		return fmt.Errorf("failed to encode XML response: %w", err)
	}

	r.httpWriter.WriteHeader(r.statusCode)
	_, err := r.httpWriter.Write(buf.Bytes())
	return err
}

// Blob sends raw bytes with an explicit content type. If contentType is
// empty, MIMEOctetStream is used. Content-Length is set automatically.
// This is the building block behind in-memory file/binary responses
// (PDF, images, pre-rendered documents) without a Content-Disposition.
func (r *Response) Blob(contentType string, data []byte) error {
	if contentType != "" {
		r.httpWriter.Header().Set(HeaderContentType, contentType)
	}
	r.prepare(StatusOK, MIMEOctetStream)
	r.httpWriter.Header().Set(HeaderContentLength, strconv.Itoa(len(data)))

	r.httpWriter.WriteHeader(r.statusCode)
	_, err := r.httpWriter.Write(data)
	return err
}

// String sends a string response.
// If the status code is not set - StatusOK will be set.
// If ContentType isn't defined - MIMETextPlain will be used by default.
func (r *Response) String(data string) error {
	r.prepare(StatusOK, MIMETextPlain)
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := io.WriteString(r.httpWriter, data)
	return err
}

// Error sends an error response as JSON.
//
// The status argument is always the HTTP status code: it is written to the
// response status line AND, by default, mirrored into the JSON body's
// "code" field. Use WithErrorCode to override only the body code with a
// distinct application/business code while keeping the HTTP status intact.
//
// The status argument takes precedence over any WithStatus option. If it is
// outside the valid HTTP range (100-599) it falls back to
// StatusInternalServerError. When message is empty, the canonical text for
// the status code is used.
func (r *Response) Error(status int, message string) error {
	if status < 100 || status > 599 {
		status = StatusInternalServerError
	}
	r.statusCode = status

	code := status
	if r.errorCode != nil {
		code = *r.errorCode
	}

	if message == "" {
		message = statusMessages[status]
	}

	return r.JSON(&ErrorResponse{Code: code, Message: message})
}

// Stream sends a stream response.
func (r *Response) Stream(data io.Reader) error {
	r.prepare(StatusOK, MIMEOctetStream)
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := io.Copy(r.httpWriter, data)
	return err
}

// ServeFile sends a file response.
//
// It delegates to http.ServeFile, which writes the status code, headers and
// body itself (including 404 handling, range requests and content sniffing).
// Because http.ServeFile manages the whole exchange, this method always
// returns nil; inspect the written response (e.g. via a recorder) rather
// than the return value to detect failures.
func (r *Response) ServeFile(req *http.Request, file string) error {
	r.prepare(StatusOK, MIMEOctetStream)

	// The http.ServeFile function from the net/http package independently
	// sets the response headers and status code before starting the file
	// transfer, no need: r.httpWriter.WriteHeader(r.statusCode)
	http.ServeFile(r.httpWriter, req, file)
	return nil
}

// ServeFileAsDownload sends in-memory data as a downloadable file.
//
// It sets Content-Disposition (with an RFC 5987 UTF-8 form so non-ASCII
// file names survive), Content-Length from len(data), and
// X-Content-Type-Options: nosniff to stop browsers from re-sniffing the
// type. The content type defaults to MIMEOctetStream but is preserved if a
// caller set one explicitly via options.
func (r *Response) ServeFileAsDownload(fileName string, data []byte) error {
	r.httpWriter.Header().Set(
		HeaderContentDisposition,
		attachmentDisposition(fileName),
	)
	r.httpWriter.Header().Set(HeaderContentLength, strconv.Itoa(len(data)))
	r.httpWriter.Header().Set(HeaderXContentTypeOptions, "nosniff")

	r.prepare(StatusOK, MIMEOctetStream)
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := r.httpWriter.Write(data)
	return err
}

// Redirect sends an HTTP redirect to the specified URL.
//
// Only genuine redirect status codes are honoured: 300, 301, 302, 303, 307
// and 308. Any other configured status (including the non-redirect 304, 305
// and 306) falls back to StatusFound (302). This method always returns nil.
func (r *Response) Redirect(url string) error {
	r.prepare(StatusFound)

	switch r.statusCode {
	case StatusMultipleChoices, StatusMovedPermanently, StatusFound,
		StatusSeeOther, StatusTemporaryRedirect, StatusPermanentRedirect:
		// A real redirect status — keep it.
	default:
		r.statusCode = StatusFound // default to 302 for non-redirect codes
	}

	r.httpWriter.Header().Set(HeaderLocation, url)
	r.httpWriter.WriteHeader(r.statusCode)
	return nil
}

// NoContent sends a 204 No Content response. It always returns nil.
func (r *Response) NoContent() error {
	r.SetStatus(StatusNoContent)
	r.prepare(StatusNoContent)
	r.httpWriter.WriteHeader(http.StatusNoContent)
	return nil
}

// Created sends a 201 Created response with an optional Location header and
// an optional JSON body. If location is non-empty it is written to the
// Location header. If data is nil, only the status (and headers) are sent;
// otherwise data is encoded as JSON. An explicitly configured status code
// (via WithStatus/SetStatus) is respected, otherwise 201 is used.
func (r *Response) Created(location string, data any) error {
	if location != "" {
		r.httpWriter.Header().Set(HeaderLocation, location)
	}
	if r.statusCode == StatusUndefined {
		r.statusCode = StatusCreated
	}

	if data == nil {
		r.prepare(StatusCreated)
		r.httpWriter.WriteHeader(r.statusCode)
		return nil
	}

	return r.JSON(data)
}

// Accepted sends a 202 Accepted response with an optional JSON body. If
// data is nil, only the status is sent; otherwise data is encoded as JSON.
// An explicitly configured status code is respected, otherwise 202 is used.
func (r *Response) Accepted(data any) error {
	if r.statusCode == StatusUndefined {
		r.statusCode = StatusAccepted
	}

	if data == nil {
		r.prepare(StatusAccepted)
		r.httpWriter.WriteHeader(r.statusCode)
		return nil
	}

	return r.JSON(data)
}

// HTML sends an HTML response.
// If the status code is not set - StatusOK will be set.
// If ContentType isn't defined - MIMETextHTMLCharsetUTF8 is used.
func (r *Response) HTML(html string) error {
	r.prepare(http.StatusOK, MIMETextHTMLCharsetUTF8)
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := io.WriteString(r.httpWriter, html)
	return err
}

// attachmentDisposition builds a Content-Disposition value for an
// attachment. It emits a sanitized ASCII fallback in the quoted filename
// (dropping control characters, quotes and backslashes that would break the
// header) and always appends an RFC 5987 filename* form so UTF-8 names are
// delivered intact to compliant clients.
func attachmentDisposition(fileName string) string {
	var ascii strings.Builder
	ascii.Grow(len(fileName))
	for _, r := range fileName {
		switch {
		case r < 0x20 || r == 0x7f || r == '"' || r == '\\':
			// Skip characters unsafe for a quoted-string.
		case r > 0x7f:
			ascii.WriteByte('_') // non-ASCII placeholder in the fallback
		default:
			ascii.WriteRune(r)
		}
	}

	return `attachment; filename="` + ascii.String() +
		`"; filename*=UTF-8''` + rfc5987Encode(fileName)
}

package resp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/goloop/g"
)

// Response represents an HTTP response.
// It provides methods for setting headers, cookies, and writing data
// to the response body.
type Response struct {
	httpWriter http.ResponseWriter
	statusCode int
}

// NewResponse creates a new instance of Response with the provided
// http.ResponseWriter and options. It applies the provided options
// to the response and returns the pointer of created response.
func NewResponse(w http.ResponseWriter, opts ...Option) *Response {
	// Create a new response with the provided http.ResponseWriter.
	response := &Response{
		httpWriter: w,
		statusCode: StatusUndefined,
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
	// We use strings.Builder instead of concatenation.
	var b strings.Builder
	b.Grow(len(cookie.Name) + len(cookie.Value) + 50)

	b.WriteString(cookie.String())
	r.httpWriter.Header().Add("Set-Cookie", b.String())
	return r
	// http.SetCookie(r.httpWriter, cookie)
	// return r
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
	buf := jsonBufPool.Get().(*bytes.Buffer)
	buf.Reset()

	defer jsonBufPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return err
	}

	r.prepare(StatusOK, MIMEApplicationJSONCharsetUTF8)
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := io.Copy(r.httpWriter, buf)

	return err
}

// StreamJSON sends a JSON stream response.
func (r *Response) StreamJSON(data any) error {
	r.prepare(StatusOK, MIMEApplicationJSONCharsetUTF8)
	r.httpWriter.WriteHeader(r.statusCode)

	enc := json.NewEncoder(r.httpWriter)
	enc.SetEscapeHTML(false) // reduce the number of allocations

	return enc.Encode(data)
}

// JSONP sends a JSONP response.
// If the status code is not set - StatusOK will be set.
// If ContentType isn't defined - MIMEApplicationJavaScript will
// be used by default.
func (r *Response) JSONP(data any, callback string) error {
	// resp, err := json.Marshal(data)
	// if err != nil {
	// 	return err
	// }

	// r.prepare(StatusOK, MIMEApplicationJavaScriptCharsetUTF8)
	// r.httpWriter.WriteHeader(r.statusCode)
	// _, err = r.httpWriter.Write([]byte(callback + "(" + string(resp) + ");"))
	// return err

	buf := jsonpBufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer jsonpBufPool.Put(buf)

	buf.WriteString(callback)
	buf.WriteByte('(')

	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return err
	}

	// Remove trailing newline from encoder.
	buf.Truncate(buf.Len() - 1)
	buf.WriteString(");")

	r.prepare(StatusOK, MIMEApplicationJavaScriptCharsetUTF8)
	r.httpWriter.WriteHeader(r.statusCode)
	_, err := io.Copy(r.httpWriter, buf)
	return err
}

// String sends a string response.
// If the status code is not set - StatusOK will be set.
// If ContentType isn't defined - MIMETextPlain will be used by default.
func (r *Response) String(data string) error {
	if len(data) > 32*1024 {
		buf := largeBufPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer largeBufPool.Put(buf)

		buf.WriteString(data)
		r.prepare(StatusOK, MIMETextPlain)
		r.httpWriter.WriteHeader(r.statusCode)
		_, err := io.Copy(r.httpWriter, buf)
		return err
	}

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
	// For large files, use buffered copy.
	buf := make([]byte, 32*1024)

	r.prepare(StatusOK, MIMEOctetStream)
	r.httpWriter.WriteHeader(r.statusCode)

	_, err := io.CopyBuffer(r.httpWriter, data, buf)
	return err

	// r.prepare(StatusOK, MIMEOctetStream)
	// r.httpWriter.WriteHeader(r.statusCode)
	// _, err := io.Copy(r.httpWriter, data)
	// return err
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

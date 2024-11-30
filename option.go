package resp

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Option represents a response option.
type Option func(*Response) *Response

// WarningHeader represents a Warning header.
type WarningHeader struct {
	Code  int
	Agent string
	Text  string
	Date  time.Time
}

// LinkHeader represents a Link header.
type LinkHeader struct {
	URI   string
	Rel   string
	Type  string
	Title string
}

// WithHeader adds the provided header key-value pair to the response.
func WithHeader(key string, values ...string) Option {
	return func(r *Response) *Response {
		return r.AddHeader(key, values...)
	}
}

// WithStatus sets the status code of the response.
func WithStatus(code int) Option {
	return func(r *Response) *Response {
		return r.SetStatus(code)
	}
}

// WithCookie sets the provided cookie in the response.
func WithCookie(cookie *http.Cookie) Option {
	return func(r *Response) *Response {
		return r.SetCookie(cookie)
	}
}

// WithStatusContinue sets the status code to 100.
func WithStatusContinue() Option {
	return WithStatus(StatusContinue)
}

// WithStatusSwitchingProtocols sets the status code to 101.
func WithStatusSwitchingProtocols() Option {
	return WithStatus(StatusSwitchingProtocols)
}

// WithStatusProcessing sets the status code to 102.
func WithStatusProcessing() Option {
	return WithStatus(StatusProcessing)
}

// WithStatusEarlyHints sets the status code to 103.
func WithStatusEarlyHints() Option {
	return WithStatus(StatusEarlyHints)
}

// WithStatusOK sets the status code to 200.
func WithStatusOK() Option {
	return WithStatus(StatusOK)
}

// WithStatusCreated sets the status code to 201.
func WithStatusCreated() Option {
	return WithStatus(StatusCreated)
}

// WithStatusAccepted sets the status code to 202.
func WithStatusAccepted() Option {
	return WithStatus(StatusAccepted)
}

// WithStatusNonAuthoritativeInfo sets the status code to 203.
func WithStatusNonAuthoritativeInfo() Option {
	return WithStatus(StatusNonAuthoritativeInfo)
}

// WithStatusNoContent sets the status code to 204.
func WithStatusNoContent() Option {
	return WithStatus(StatusNoContent)
}

// WithStatusResetContent sets the status code to 205.
func WithStatusResetContent() Option {
	return WithStatus(StatusResetContent)
}

// WithStatusPartialContent sets the status code to 206.
func WithStatusPartialContent() Option {
	return WithStatus(StatusPartialContent)
}

// WithStatusMultiStatus sets the status code to 207.
func WithStatusMultiStatus() Option {
	return WithStatus(StatusMultiStatus)
}

// WithStatusAlreadyReported sets the status code to 208.
func WithStatusAlreadyReported() Option {
	return WithStatus(StatusAlreadyReported)
}

// WithStatusIMUsed sets the status code to 226.
func WithStatusIMUsed() Option {
	return WithStatus(StatusIMUsed)
}

// WithStatusMultipleChoices sets the status code to 300.
func WithStatusMultipleChoices() Option {
	return WithStatus(StatusMultipleChoices)
}

// WithStatusMovedPermanently sets the status code to 301.
func WithStatusMovedPermanently() Option {
	return WithStatus(StatusMovedPermanently)
}

// WithStatusFound sets the status code to 302.
func WithStatusFound() Option {
	return WithStatus(StatusFound)
}

// WithStatusSeeOther sets the status code to 303.
func WithStatusSeeOther() Option {
	return WithStatus(StatusSeeOther)
}

// WithStatusNotModified sets the status code to 304.
func WithStatusNotModified() Option {
	return WithStatus(StatusNotModified)
}

// WithStatusUseProxy sets the status code to 305.
func WithStatusUseProxy() Option {
	return WithStatus(StatusUseProxy)
}

// WithStatusTemporaryRedirect sets the status code to 307.
func WithStatusTemporaryRedirect() Option {
	return WithStatus(StatusTemporaryRedirect)
}

// WithStatusPermanentRedirect sets the status code to 308.
func WithStatusPermanentRedirect() Option {
	return WithStatus(StatusPermanentRedirect)
}

// WithStatusBadRequest sets the status code to 400.
func WithStatusBadRequest() Option {
	return WithStatus(StatusBadRequest)
}

// WithStatusUnauthorized sets the status code to 401.
func WithStatusUnauthorized() Option {
	return WithStatus(StatusUnauthorized)
}

// WithStatusPaymentRequired sets the status code to 402.
func WithStatusPaymentRequired() Option {
	return WithStatus(StatusPaymentRequired)
}

// WithStatusForbidden sets the status code to 403.
func WithStatusForbidden() Option {
	return WithStatus(StatusForbidden)
}

// WithStatusNotFound sets the status code to 404.
func WithStatusNotFound() Option {
	return WithStatus(StatusNotFound)
}

// WithStatusMethodNotAllowed sets the status code to 405.
func WithStatusMethodNotAllowed() Option {
	return WithStatus(StatusMethodNotAllowed)
}

// WithStatusNotAcceptable sets the status code to 406.
func WithStatusNotAcceptable() Option {
	return WithStatus(StatusNotAcceptable)
}

// WithStatusProxyAuthRequired sets the status code to 407.
func WithStatusProxyAuthRequired() Option {
	return WithStatus(StatusProxyAuthRequired)
}

// WithStatusRequestTimeout sets the status code to 408.
func WithStatusRequestTimeout() Option {
	return WithStatus(StatusRequestTimeout)
}

// WithStatusConflict sets the status code to 409.
func WithStatusConflict() Option {
	return WithStatus(StatusConflict)
}

// WithStatusGone sets the status code to 410.
func WithStatusGone() Option {
	return WithStatus(StatusGone)
}

// WithStatusLengthRequired sets the status code to 411.
func WithStatusLengthRequired() Option {
	return WithStatus(StatusLengthRequired)
}

// WithStatusPreconditionFailed sets the status code to 412.
func WithStatusPreconditionFailed() Option {
	return WithStatus(StatusPreconditionFailed)
}

// WithStatusRequestEntityTooLarge sets the status code to 413.
func WithStatusRequestEntityTooLarge() Option {
	return WithStatus(StatusRequestEntityTooLarge)
}

// WithStatusRequestURITooLong sets the status code to 414.
func WithStatusRequestURITooLong() Option {
	return WithStatus(StatusRequestURITooLong)
}

// WithStatusUnsupportedMediaType sets the status code to 415.
func WithStatusUnsupportedMediaType() Option {
	return WithStatus(StatusUnsupportedMediaType)
}

// WithStatusRequestedRangeNotSatisfiable sets the status code to 416.
func WithStatusRequestedRangeNotSatisfiable() Option {
	return WithStatus(StatusRequestedRangeNotSatisfiable)
}

// WithStatusExpectationFailed sets the status code to 417.
func WithStatusExpectationFailed() Option {
	return WithStatus(StatusExpectationFailed)
}

// WithStatusTeapot sets the status code to 418.
func WithStatusTeapot() Option {
	return WithStatus(StatusTeapot)
}

// WithStatusMisdirectedRequest sets the status code to 421.
func WithStatusMisdirectedRequest() Option {
	return WithStatus(StatusMisdirectedRequest)
}

// WithStatusUnprocessableEntity sets the status code to 422.
func WithStatusUnprocessableEntity() Option {
	return WithStatus(StatusUnprocessableEntity)
}

// WithStatusLocked sets the status code to 423.
func WithStatusLocked() Option {
	return WithStatus(StatusLocked)
}

// WithStatusFailedDependency sets the status code to 424.
func WithStatusFailedDependency() Option {
	return WithStatus(StatusFailedDependency)
}

// WithStatusTooEarly sets the status code to 425.
func WithStatusTooEarly() Option {
	return WithStatus(StatusTooEarly)
}

// WithStatusUpgradeRequired sets the status code to 426.
func WithStatusUpgradeRequired() Option {
	return WithStatus(StatusUpgradeRequired)
}

// WithStatusPreconditionRequired sets the status code to 428.
func WithStatusPreconditionRequired() Option {
	return WithStatus(StatusPreconditionRequired)
}

// WithStatusTooManyRequests sets the status code to 429.
func WithStatusTooManyRequests() Option {
	return WithStatus(StatusTooManyRequests)
}

// WithStatusRequestHeaderFieldsTooLarge sets the status code to 431.
func WithStatusRequestHeaderFieldsTooLarge() Option {
	return WithStatus(StatusRequestHeaderFieldsTooLarge)
}

// WithStatusUnavailableForLegalReasons sets the status code to 451.
func WithStatusUnavailableForLegalReasons() Option {
	return WithStatus(StatusUnavailableForLegalReasons)
}

// WithStatusInternalServerError sets the status code to 500.
func WithStatusInternalServerError() Option {
	return WithStatus(StatusInternalServerError)
}

// WithStatusNotImplemented sets the status code to 501.
func WithStatusNotImplemented() Option {
	return WithStatus(StatusNotImplemented)
}

// WithStatusBadGateway sets the status code to 502.
func WithStatusBadGateway() Option {
	return WithStatus(StatusBadGateway)
}

// WithStatusServiceUnavailable sets the status code to 503.
func WithStatusServiceUnavailable() Option {
	return WithStatus(StatusServiceUnavailable)
}

// WithStatusGatewayTimeout sets the status code to 504.
func WithStatusGatewayTimeout() Option {
	return WithStatus(StatusGatewayTimeout)
}

// WithStatusHTTPVersionNotSupported sets the status code to 505.
func WithStatusHTTPVersionNotSupported() Option {
	return WithStatus(StatusHTTPVersionNotSupported)
}

// WithStatusVariantAlsoNegotiates sets the status code to 506.
func WithStatusVariantAlsoNegotiates() Option {
	return WithStatus(StatusVariantAlsoNegotiates)
}

// WithStatusInsufficientStorage sets the status code to 507.
func WithStatusInsufficientStorage() Option {
	return WithStatus(StatusInsufficientStorage)
}

// WithStatusLoopDetected sets the status code to 508.
func WithStatusLoopDetected() Option {
	return WithStatus(StatusLoopDetected)
}

// WithStatusNotExtended sets the status code to 510.
func WithStatusNotExtended() Option {
	return WithStatus(StatusNotExtended)
}

// WithStatusNetworkAuthenticationRequired sets the status code to 511.
func WithStatusNetworkAuthenticationRequired() Option {
	return WithStatus(StatusNetworkAuthenticationRequired)
}

// AddContentType sets the Content-Type header.
func AddContentType(value string) Option {
	return WithHeader(HeaderContentType, value)
}

// AddETag sets the ETag header.
func AddETag(value string) Option {
	return WithHeader(HeaderETag, value)
}

// AddLastModified sets the Last-Modified header.
func AddLastModified(t time.Time) Option {
	return WithHeader(HeaderLastModified, t.Format(time.RFC1123))
}

// AddContentLength sets the Content-Length header.
func AddContentLength(length int64) Option {
	return WithHeader(HeaderContentLength, strconv.FormatInt(length, 10))
}

// AddUserAgent sets the User-Agent header.
func AddUserAgent(value string) Option {
	return WithHeader(HeaderUserAgent, value)
}

// AddHost sets the Host header.
func AddHost(value string) Option {
	return WithHeader(HeaderHost, value)
}

// AddReferer sets the Referer header.
func AddReferer(value string) Option {
	return WithHeader(HeaderReferer, value)
}

// AddServer sets the Server header.
func AddServer(value string) Option {
	return WithHeader(HeaderServer, value)
}

// AddDate sets the Date header.
func AddDate(date time.Time) Option {
	return WithHeader(HeaderDate, date.Format(time.RFC1123))
}

// AddLocation sets the Location header.
func AddLocation(value string) Option {
	return WithHeader(HeaderLocation, value)
}

// // AddRetryAfter sets the Retry-After header.
// func AddRetryAfter(seconds int) Option {
// 	return WithHeader(HeaderRetryAfter, strconv.Itoa(seconds))
// }

// AddRetryAfter sets the Retry-After header.
func AddRetryAfter[T int | time.Time | time.Duration](value T) Option {
	return func(r *Response) *Response {
		var stringValue string
		switch v := any(value).(type) {
		case int:
			stringValue = strconv.Itoa(v)
		case time.Time:
			stringValue = v.Format(time.RFC1123)
		case time.Duration:
			stringValue = strconv.Itoa(int(v.Seconds()))
		}

		return WithHeader(HeaderRetryAfter, stringValue)(r)
	}
}

// AddContentDisposition sets the Content-Disposition header.
func AddContentDisposition(
	dispositionType,
	filename string,
	useUTF8Encoding ...bool,
) Option {
	return func(r *Response) *Response {
		// Check if UTF-8 encoding is needed for the filename.
		if len(useUTF8Encoding) > 0 && useUTF8Encoding[0] {
			// Encode the filename using URL encoding.
			encodedFilename := url.PathEscape(filename)
			value := fmt.Sprintf(
				`%s; filename*=UTF-8''%s`,
				dispositionType,
				encodedFilename,
			)
			return WithHeader(HeaderContentDisposition, value)(r)
		} else {
			// Standard encoding.
			value := fmt.Sprintf(`%s; filename="%s"`, dispositionType, filename)
			return WithHeader(HeaderContentDisposition, value)(r)
		}
	}
}

// AddContentEncoding sets the Content-Encoding header.
func AddContentEncoding(value string) Option {
	return WithHeader(HeaderContentEncoding, value)
}

// AddContentLanguage sets the Content-Language header.
func AddContentLanguage(value string) Option {
	return WithHeader(HeaderContentLanguage, value)
}

// AddContentLocation sets the Content-Location header.
func AddContentLocation(value string) Option {
	return WithHeader(HeaderContentLocation, value)
}

// AddWWWAuthenticate sets the WWW-Authenticate header.
func AddWWWAuthenticate(value ...string) Option {
	return WithHeader(HeaderWWWAuthenticate, value...)
}

// AddAuthorization sets the Authorization header.
func AddAuthorization(value ...string) Option {
	return WithHeader(HeaderAuthorization, value...)
}

// AddProxyAuthenticate sets the Proxy-Authenticate header.
func AddProxyAuthenticate(value ...string) Option {
	return WithHeader(HeaderProxyAuthenticate, value...)
}

// AddProxyAuthorization sets the Proxy-Authorization header.
func AddProxyAuthorization(value ...string) Option {
	return WithHeader(HeaderProxyAuthorization, value...)
}

// AddIfMatch sets the If-Match header.
func AddIfMatch(value ...string) Option {
	return WithHeader(HeaderIfMatch, value...)
}

// AddIfNoneMatch sets the If-None-Match header.
func AddIfNoneMatch(value ...string) Option {
	return WithHeader(HeaderIfNoneMatch, value...)
}

// AddIfModifiedSince sets the If-Modified-Since header.
func AddIfModifiedSince(t time.Time) Option {
	return WithHeader(HeaderIfModifiedSince, t.Format(time.RFC1123))
}

// AddIfUnmodifiedSince sets the If-Unmodified-Since header.
func AddIfUnmodifiedSince(t time.Time) Option {
	return WithHeader(HeaderIfUnmodifiedSince, t.Format(time.RFC1123))
}

// AddIfRange sets the If-Range header.
func AddIfRange(value string) Option {
	return WithHeader(HeaderIfRange, value)
}

// AddContentSecurityPolicy sets the Content-Security-Policy header.
func AddContentSecurityPolicy(value ...string) Option {
	return WithHeader(HeaderContentSecurityPolicy, value...)
}

// AddContentSecurityPolicyReportOnly sets the
// Content-Security-Policy-Report-Only header.
func AddContentSecurityPolicyReportOnly(value ...string) Option {
	return WithHeader(HeaderContentSecurityPolicyReportOnly, value...)
}

// AddStrictTransportSecurity sets the Strict-Transport-Security header.
// The maxAgeSeconds parameter is the number of seconds that the browser
// should remember that this site is only to be accessed using HTTPS.
//
// If includeSubDomains is true, this rule applies to all of the site's
// subdomains.
// If preload is true, the site is included in the HSTS preload list.
//
// The val parameter is optional and can be used to set the includeSubDomains
// and preload values.
//
// For example:
//
//	AddStrictTransportSecurity(31536000)
//	AddStrictTransportSecurity(31536000, true, true)
//	AddStrictTransportSecurity(31536000, false, true)
//	AddStrictTransportSecurity(31536000, true)
func AddStrictTransportSecurity(maxAgeSeconds int, val ...bool) Option {
	includeSubDomains := false
	preload := false

	if len(val) > 0 {
		includeSubDomains = val[0]
	}

	if len(val) > 1 {
		preload = val[1]
	}

	value := fmt.Sprintf("max-age=%d", maxAgeSeconds)
	if includeSubDomains {
		value += "; includeSubDomains"
	}

	if preload {
		value += "; preload"
	}

	return WithHeader(HeaderStrictTransportSecurity, value)
}

// AddReferrerPolicy sets the Referrer-Policy header.
func AddReferrerPolicy(value string) Option {
	return WithHeader(HeaderReferrerPolicy, value)
}

// AddUpgradeInsecureRequests sets the Upgrade-Insecure-Requests header.
func AddUpgradeInsecureRequests[T int | bool | string](enable T) Option {
	value := "0"
	switch v := any(enable).(type) {
	case int:
		if v > 0 {
			value = "1"
		}
	case bool:
		if v {
			value = "1"
		}
	case string:
		value = v
	}

	return WithHeader(HeaderUpgradeInsecureRequests, value)
}

// AddXContentTypeOptions sets the X-Content-Type-Options header.
func AddXContentTypeOptions(value string) Option {
	return WithHeader(HeaderXContentTypeOptions, value)
}

// AddXFrameOptions sets the X-Frame-Options header.
func AddXFrameOptions(value string) Option {
	return WithHeader(HeaderXFrameOptions, value)
}

// AddXXSSProtection sets the X-XSS-Protection header.
func AddXXSSProtection(value string) Option {
	return WithHeader(HeaderXXSSProtection, value)
}

// AddContentDPR sets the Content-DPR header.
func AddContentDPR(value float64) Option {
	return WithHeader(HeaderContentDPR, strconv.FormatFloat(value, 'f', -1, 64))
}

// AddDPR sets the DPR header.
func AddDPR(value float64) Option {
	return WithHeader(HeaderDPR, strconv.FormatFloat(value, 'f', -1, 64))
}

// AddViewportWidth sets the Viewport-Width header.
func AddViewportWidth(value int) Option {
	return WithHeader(HeaderViewportWidth, strconv.Itoa(value))
}

// AddWidth sets the Width header.
func AddWidth(value int) Option {
	return WithHeader(HeaderWidth, strconv.Itoa(value))
}

// AddContentRange sets the Content-Range header.
func AddContentRange(start, end, total int) Option {
	return func(r *Response) *Response {
		value := fmt.Sprintf("bytes %d-%d/%d", start, end, total)
		r.httpWriter.Header().Set(HeaderContentRange, value)
		return r
	}
}

// AddAcceptRanges sets the Accept-Ranges header.
func AddAccept(value ...string) Option {
	return WithHeader(HeaderAccept, value...)
}

// AddAcceptCharset sets the Accept-Charset header.
func AddAcceptCharset(value ...string) Option {
	return WithHeader(HeaderAcceptCharset, value...)
}

// AddAcceptEncoding sets the Accept-Encoding header.
func AddAcceptEncoding(value ...string) Option {
	return WithHeader(HeaderAcceptEncoding, value...)
}

// AddAcceptLanguage sets the Accept-Language header.
func AddAcceptLanguage(value ...string) Option {
	return WithHeader(HeaderAcceptLanguage, value...)
}

// AddCacheControl sets the Cache-Control header.
func AddCacheControl(value ...string) Option {
	return WithHeader(HeaderCacheControl, value...)
}

// AddPragma sets the Pragma header.
func AddPragma(value ...string) Option {
	return WithHeader(HeaderPragma, value...)
}

// AddWarning sets the Warning header.
func AddWarning(warnings ...WarningHeader) Option {
	return func(r *Response) *Response {
		for _, warning := range warnings {
			dateStr := ""
			if !warning.Date.IsZero() {
				dateStr = warning.Date.Format(time.RFC1123)
			}

			value := fmt.Sprintf("%d", warning.Code)

			if warning.Agent != "" {
				value += " " + warning.Agent
			}

			if warning.Text != "" {
				value += " \"" + warning.Text + "\""
			}

			if dateStr != "" {
				value += " \"" + dateStr + "\""
			}

			r.httpWriter.Header().Add(HeaderWarning, value)
		}
		return r
	}
}

// AddVary sets the Vary header.
func AddVary(value ...string) Option {
	return WithHeader(HeaderVary, value...)
}

// AddConnection sets the Connection header.
func AddConnection(value ...string) Option {
	return WithHeader(HeaderConnection, value...)
}

// AddTransferEncoding sets the Transfer-Encoding header.
func AddTransferEncoding(value ...string) Option {
	return WithHeader(HeaderTransferEncoding, value...)
}

// AddAccessControlAllowHeaders sets the Access-Control-Allow-Headers header.
func AddAccessControlAllowHeaders(value ...string) Option {
	return WithHeader(HeaderAccessControlAllowHeaders, value...)
}

// AddAccessControlAllowMethods sets the Access-Control-Allow-Methods header.
func AddAccessControlAllowMethods(value ...string) Option {
	return WithHeader(HeaderAccessControlAllowMethods, value...)
}

// AddAccessControlExposeHeaders sets the Access-Control-Expose-Headers header.
func AddAccessControlExposeHeaders(value ...string) Option {
	return WithHeader(HeaderAccessControlExposeHeaders, value...)
}

// AddLink sets the Link header.
func AddLink(links ...LinkHeader) Option {
	return func(r *Response) *Response {
		for _, link := range links {
			linkValue := fmt.Sprintf("<%s>; rel=\"%s\"", link.URI, link.Rel)
			if link.Type != "" {
				linkValue += fmt.Sprintf("; type=\"%s\"", link.Type)
			}
			if link.Title != "" {
				linkValue += fmt.Sprintf("; title=\"%s\"", link.Title)
			}

			r.httpWriter.Header().Add(HeaderLink, linkValue)
		}
		return r
	}
}

// AddAccessControlAllowCredentials sets the
// Access-Control-Allow-Credentials header.
func AddAccessControlAllowCredentials(enable bool) Option {
	return WithHeader(HeaderAccessControlAllowCredentials, strconv.FormatBool(enable))
}

// AddAccessControlAllowOrigin sets the Access-Control-Allow-Origin header.
func AddAccessControlAllowOrigin(value ...string) Option {
	return WithHeader(HeaderAccessControlAllowOrigin, value...)
}

// AddAccessControlRequestHeaders sets the
// Access-Control-Request-Headers header.
func AddAccessControlRequestHeaders(value ...string) Option {
	return WithHeader(HeaderAccessControlRequestHeaders, value...)
}

// AddAccessControlRequestMethod sets the Access-Control-Request-Method header.
func AddAccessControlRequestMethod(value ...string) Option {
	return WithHeader(HeaderAccessControlRequestMethod, value...)
}

// AddOrigin sets the Origin header.
func AddOrigin(value ...string) Option {
	return WithHeader(HeaderOrigin, value...)
}

// AsTextXML sets the Content-Type header to text/xml.
func AsTextXML() Option {
	return AddContentType(MIMETextXML)
}

// AsTextHTML sets the Content-Type header to text/html.
func AsTextHTML() Option {
	return AddContentType(MIMETextHTML)
}

// AsTextPlain sets the Content-Type header to text/plain.
func AsTextPlain() Option {
	return AddContentType(MIMETextPlain)
}

// AsTextJavaScript sets the Content-Type header to text/javascript.
func AsTextJavaScript() Option {
	return AddContentType(MIMETextJavaScript)
}

// AsApplicationXML sets the Content-Type header to application/xml.
func AsApplicationXML() Option {
	return AddContentType(MIMEApplicationXML)
}

// AsApplicationJSON sets the Content-Type header to application/json.
func AsApplicationJSON() Option {
	return AddContentType(MIMEApplicationJSON)
}

// AsApplicationJavaScript sets the Content-Type header
// to application/javascript.
func AsApplicationJavaScript() Option {
	return AddContentType(MIMEApplicationJavaScript)
}

// AsApplicationForm sets the Content-Type header
// to application/x-www-form-urlencoded.
func AsApplicationForm() Option {
	return AddContentType(MIMEApplicationForm)
}

// AsOctetStream sets the Content-Type header to application/octet-stream.
func AsOctetStream() Option {
	return AddContentType(MIMEOctetStream)
}

// AsMultipartForm sets the Content-Type header to multipart/form-data.
func AsMultipartForm() Option {
	return AddContentType(MIMEMultipartForm)
}

// AsTextXMLCharsetUTF8 sets the Content-Type header
// to text/xml; charset=utf-8.
func AsTextXMLCharsetUTF8() Option {
	return AddContentType(MIMETextXMLCharsetUTF8)
}

// AsTextHTMLCharsetUTF8 sets the Content-Type header
// to text/html; charset=utf-8.
func AsTextHTMLCharsetUTF8() Option {
	return AddContentType(MIMETextHTMLCharsetUTF8)
}

// AsTextPlainCharsetUTF8 sets the Content-Type header
// to text/plain; charset=utf-8.
func AsTextPlainCharsetUTF8() Option {
	return AddContentType(MIMETextPlainCharsetUTF8)
}

// AsTextJavaScriptCharsetUTF8 sets the Content-Type header
// to text/javascript; charset=utf-8.
func AsTextJavaScriptCharsetUTF8() Option {
	return AddContentType(MIMETextJavaScriptCharsetUTF8)
}

// AsApplicationXMLCharsetUTF8 sets the Content-Type header
// to application/xml; charset=utf-8.
func AsApplicationXMLCharsetUTF8() Option {
	return AddContentType(MIMEApplicationXMLCharsetUTF8)
}

// AsApplicationJSONCharsetUTF8 sets the Content-Type header
// to application/json; charset=utf-8.
func AsApplicationJSONCharsetUTF8() Option {
	return AddContentType(MIMEApplicationJSONCharsetUTF8)
}

// AsApplicationJavaScriptCharsetUTF8 sets the Content-Type header
// to application/javascript; charset=utf-8.
func AsApplicationJavaScriptCharsetUTF8() Option {
	return AddContentType(MIMEApplicationJavaScriptCharsetUTF8)
}

// ApplyJSONEncoder sets the custom JSON encoder function.
// This allows us to use a different JSON encoding library
// or customize encoding.
//
// Example Usage:
//
//	import jsoniter "github.com/json-iterator/go"
//
//	customEncoder := func(w io.Writer, v interface{}) error {
//	    return jsoniter.NewEncoder(w).Encode(v)
//	}
//
//	response := resp.NewResponse(w, resp.ApplyJSONEncoder(customEncoder))
func ApplyJSONEncoder(encodeFunc JSONEncodeFunc) Option {
	return func(r *Response) *Response {
		r.jsonEncodeFunc = encodeFunc
		return r
	}
}

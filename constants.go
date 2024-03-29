package resp

import "net/http"

// MIME types that are commonly used
const (
	MIMETextXML                          = "text/xml"
	MIMETextHTML                         = "text/html"
	MIMETextPlain                        = "text/plain"
	MIMETextJavaScript                   = "text/javascript"
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEOctetStream                      = "application/octet-stream"
	MIMEMultipartForm                    = "multipart/form-data"
	MIMETextXMLCharsetUTF8               = "text/xml; charset=utf-8"
	MIMETextHTMLCharsetUTF8              = "text/html; charset=utf-8"
	MIMETextPlainCharsetUTF8             = "text/plain; charset=utf-8"
	MIMETextJavaScriptCharsetUTF8        = "text/javascript; charset=utf-8"
	MIMEApplicationXMLCharsetUTF8        = "application/xml; charset=utf-8"
	MIMEApplicationJSONCharsetUTF8       = "application/json; charset=utf-8"
	MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript; charset=utf-8"
)

// HTTP Headers were copied from net/http.
const (
	HeaderAuthorization                      = "Authorization"
	HeaderProxyAuthenticate                  = "Proxy-Authenticate"
	HeaderProxyAuthorization                 = "Proxy-Authorization"
	HeaderWWWAuthenticate                    = "WWW-Authenticate"
	HeaderAge                                = "Age"
	HeaderCacheControl                       = "Cache-Control"
	HeaderClearSiteData                      = "Clear-Site-Data"
	HeaderExpires                            = "Expires"
	HeaderPragma                             = "Pragma"
	HeaderWarning                            = "Warning"
	HeaderAcceptCH                           = "Accept-CH"
	HeaderAcceptCHLifetime                   = "Accept-CH-Lifetime"
	HeaderContentDPR                         = "Content-DPR"
	HeaderDPR                                = "DPR"
	HeaderEarlyData                          = "Early-Data"
	HeaderSaveData                           = "Save-Data"
	HeaderViewportWidth                      = "Viewport-Width"
	HeaderWidth                              = "Width"
	HeaderETag                               = "ETag"
	HeaderIfMatch                            = "If-Match"
	HeaderIfModifiedSince                    = "If-Modified-Since"
	HeaderIfNoneMatch                        = "If-None-Match"
	HeaderIfUnmodifiedSince                  = "If-Unmodified-Since"
	HeaderLastModified                       = "Last-Modified"
	HeaderVary                               = "Vary"
	HeaderConnection                         = "Connection"
	HeaderKeepAlive                          = "Keep-Alive"
	HeaderAccept                             = "Accept"
	HeaderAcceptCharset                      = "Accept-Charset"
	HeaderAcceptEncoding                     = "Accept-Encoding"
	HeaderAcceptLanguage                     = "Accept-Language"
	HeaderCookie                             = "Cookie"
	HeaderExpect                             = "Expect"
	HeaderMaxForwards                        = "Max-Forwards"
	HeaderSetCookie                          = "Set-Cookie"
	HeaderAccessControlAllowCredentials      = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowHeaders          = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowMethods          = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowOrigin           = "Access-Control-Allow-Origin"
	HeaderAccessControlExposeHeaders         = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge                = "Access-Control-Max-Age"
	HeaderAccessControlRequestHeaders        = "Access-Control-Request-Headers"
	HeaderAccessControlRequestMethod         = "Access-Control-Request-Method"
	HeaderOrigin                             = "Origin"
	HeaderTimingAllowOrigin                  = "Timing-Allow-Origin"
	HeaderXPermittedCrossDomainPolicies      = "X-Permitted-Cross-Domain-Policies"
	HeaderDNT                                = "DNT"
	HeaderTk                                 = "Tk"
	HeaderContentDisposition                 = "Content-Disposition"
	HeaderContentEncoding                    = "Content-Encoding"
	HeaderContentLanguage                    = "Content-Language"
	HeaderContentLength                      = "Content-Length"
	HeaderContentLocation                    = "Content-Location"
	HeaderContentType                        = "Content-Type"
	HeaderForwarded                          = "Forwarded"
	HeaderVia                                = "Via"
	HeaderXForwardedFor                      = "X-Forwarded-For"
	HeaderXForwardedHost                     = "X-Forwarded-Host"
	HeaderXForwardedProto                    = "X-Forwarded-Proto"
	HeaderXForwardedProtocol                 = "X-Forwarded-Protocol"
	HeaderXForwardedSsl                      = "X-Forwarded-Ssl"
	HeaderXUrlScheme                         = "X-Url-Scheme"
	HeaderLocation                           = "Location"
	HeaderFrom                               = "From"
	HeaderHost                               = "Host"
	HeaderReferer                            = "Referer"
	HeaderReferrerPolicy                     = "Referrer-Policy"
	HeaderUserAgent                          = "User-Agent"
	HeaderAllow                              = "Allow"
	HeaderServer                             = "Server"
	HeaderAcceptRanges                       = "Accept-Ranges"
	HeaderContentRange                       = "Content-Range"
	HeaderIfRange                            = "If-Range"
	HeaderRange                              = "Range"
	HeaderContentSecurityPolicy              = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly    = "Content-Security-Policy-Report-Only"
	HeaderCrossOriginResourcePolicy          = "Cross-Origin-Resource-Policy"
	HeaderExpectCT                           = "Expect-CT"
	HeaderPermissionsPolicy                  = "Permissions-Policy"
	HeaderPublicKeyPins                      = "Public-Key-Pins"
	HeaderPublicKeyPinsReportOnly            = "Public-Key-Pins-Report-Only"
	HeaderStrictTransportSecurity            = "Strict-Transport-Security"
	HeaderUpgradeInsecureRequests            = "Upgrade-Insecure-Requests"
	HeaderXContentTypeOptions                = "X-Content-Type-Options"
	HeaderXDownloadOptions                   = "X-Download-Options"
	HeaderXFrameOptions                      = "X-Frame-Options"
	HeaderXPoweredBy                         = "X-Powered-By"
	HeaderXXSSProtection                     = "X-XSS-Protection"
	HeaderLastEventID                        = "Last-Event-ID"
	HeaderNEL                                = "NEL"
	HeaderPingFrom                           = "Ping-From"
	HeaderPingTo                             = "Ping-To"
	HeaderReportTo                           = "Report-To"
	HeaderTE                                 = "TE"
	HeaderTrailer                            = "Trailer"
	HeaderTransferEncoding                   = "Transfer-Encoding"
	HeaderSecWebSocketAccept                 = "Sec-WebSocket-Accept"
	HeaderSecWebSocketExtensions             = "Sec-WebSocket-Extensions"
	HeaderSecWebSocketKey                    = "Sec-WebSocket-Key"
	HeaderSecWebSocketProtocol               = "Sec-WebSocket-Protocol"
	HeaderSecWebSocketVersion                = "Sec-WebSocket-Version"
	HeaderAcceptPatch                        = "Accept-Patch"
	HeaderAcceptPushPolicy                   = "Accept-Push-Policy"
	HeaderAcceptSignature                    = "Accept-Signature"
	HeaderAltSvc                             = "Alt-Svc"
	HeaderDate                               = "Date"
	HeaderIndex                              = "Index"
	HeaderLargeAllocation                    = "Large-Allocation"
	HeaderLink                               = "Link"
	HeaderPushPolicy                         = "Push-Policy"
	HeaderRetryAfter                         = "Retry-After"
	HeaderServerTiming                       = "Server-Timing"
	HeaderSignature                          = "Signature"
	HeaderSignedHeaders                      = "Signed-Headers"
	HeaderSourceMap                          = "SourceMap"
	HeaderUpgrade                            = "Upgrade"
	HeaderXDNSPrefetchControl                = "X-DNS-Prefetch-Control"
	HeaderXPingback                          = "X-Pingback"
	HeaderXRequestID                         = "X-Request-ID"
	HeaderXRequestedWith                     = "X-Requested-With"
	HeaderXRobotsTag                         = "X-Robots-Tag"
	HeaderXUACompatible                      = "X-UA-Compatible"
	HeaderAccessControlAllowPrivateNetwork   = "Access-Control-Allow-Private-Network"
	HeaderAccessControlRequestPrivateNetwork = "Access-Control-Request-Private-Network"
)

// HTTP status codes.
const (
	// StatusUndefined is a special status code that indicates that the status
	// code has not been set yet.
	StatusUndefined = 0

	// StatusContinue indicates that the initial part of a request has been
	// received and has not yet been rejected by the server.
	// The client SHOULD continue sending the remainder of the request or,
	// if the request has already been completed, ignore this response.
	// RFC 9110, Section 15.2.1
	StatusContinue = http.StatusContinue // 100

	// StatusSwitchingProtocols indicates that the server understands
	// and is willing to comply with the client's request, via the
	// Upgrade header field, for a change in the application protocol
	// being used on this connection.
	// RFC 9110, Section 15.2.2
	StatusSwitchingProtocols = http.StatusSwitchingProtocols // 101

	// StatusProcessing is used to inform the client that the server has
	// accepted the complete request but has not yet completed it.
	// This status code SHOULD only be sent when the server has a reasonable
	// expectation that the request will take significant time to complete.
	// RFC 2518, Section 10.1 (WebDAV)
	StatusProcessing = http.StatusProcessing // 102

	// StatusEarlyHints is primarily intended to be used with the Link header
	// to allow the user agent to start preloading resources while the server
	// prepares a response.
	// RFC 8297
	StatusEarlyHints = http.StatusEarlyHints // 103

	// StatusOK indicates that the request has succeeded.
	// RFC 9110, Section 15.3.1
	StatusOK = http.StatusOK // 200

	// StatusCreated indicates that the request has been fulfilled and has
	// resulted in one or more new resources being created.
	// RFC 9110, Section 15.3.2
	StatusCreated = http.StatusCreated // 201

	// StatusAccepted indicates that the request has been accepted for
	// processing, but the processing has not been completed.
	// RFC 9110, Section 15.3.3
	StatusAccepted = http.StatusAccepted // 202

	// StatusNonAuthoritativeInfo indicates that the request was
	// successful but the enclosed payload has been modified from that of
	// the origin server's 200 (OK) response by a transforming proxy.
	// RFC 9110, Section 15.3.4
	StatusNonAuthoritativeInfo = http.StatusNonAuthoritativeInfo // 203

	// StatusNoContent indicates that the server has successfully fulfilled
	// the request and that there is no additional content to send in the
	// response payload body.
	// RFC 9110, Section 15.3.5
	StatusNoContent = http.StatusNoContent // 204

	// StatusResetContent indicates that the server has fulfilled the request
	// and desires that the user agent reset the "document view", which caused
	// the request to be sent, to its original state as received from the
	// origin server.
	// RFC 9110, Section 15.3.6
	StatusResetContent = http.StatusResetContent // 205

	// StatusPartialContent indicates that the server is successfully
	// fulfilling a range request for the target resource by transferring
	// one or more parts of the selected representation that correspond to
	// the satisfiable ranges found in the request's Range header field.
	// RFC 9110, Section 15.3.7
	StatusPartialContent = http.StatusPartialContent // 206

	// StatusMultiStatus provides status for multiple independent
	// operations (see WebDAV).
	// RFC 4918, Section 11.1
	StatusMultiStatus = http.StatusMultiStatus // 207

	// StatusAlreadyReported is used inside a DAV: propstat response element
	// to avoid repeatedly enumerating the internal members of multiple
	// bindings to the same collection.
	// RFC 5842, Section 7.1
	StatusAlreadyReported = http.StatusAlreadyReported // 208

	// StatusIMUsed indicates that the server has fulfilled a GET request
	// for the resource, and the response is a representation of the result
	// of one or more instance-manipulations applied to the current instance.
	// RFC 3229, Section 10.4.1
	StatusIMUsed = http.StatusIMUsed // 226

	// StatusMultipleChoices indicates that the target resource has more than
	// one representation, each with its own more specific identifier, and
	// information about the alternatives is being provided so that the user
	// (or user agent) can select a preferred representation by redirecting
	// its request to one or more of those identifiers.
	// RFC 9110, Section 15.4.1
	StatusMultipleChoices = http.StatusMultipleChoices // 300

	// StatusMovedPermanently indicates that the target resource has been
	// assigned a new permanent URI and any future references to this resource
	// ought to use one of the enclosed URIs.
	// RFC 9110, Section 15.4.2
	StatusMovedPermanently = http.StatusMovedPermanently // 301

	// StatusFound indicates that the target resource resides temporarily
	// under a different URI.
	// RFC 9110, Section 15.4.3
	StatusFound = http.StatusFound // 302

	// StatusSeeOther indicates that the server is redirecting the user agent
	// to a different resource, as indicated by a URI in the Location header
	// field, which is intended to provide an indirect response to the original
	// request.
	// RFC 9110, Section 15.4.4
	StatusSeeOther = http.StatusSeeOther // 303

	// StatusNotModified indicates that a conditional GET request has been
	// received and would have resulted in a 200 OK response if it were not
	// for the fact that the condition evaluated to false.
	// RFC 9110, Section 15.4.5
	StatusNotModified = http.StatusNotModified // 304

	// StatusUseProxy is deprecated due to security concerns regarding in-band
	// configuration of a proxy.
	// RFC 9110, Section 15.4.6
	StatusUseProxy = http.StatusUseProxy // 305

	// StatusSwitchProxy was defined in a previous version of the HTTP
	// specification to indicate that a proxy should switch protocols.
	// RFC 9110 does not define StatusSwitchProxy and this code
	// is no longer in use.
	StatusSwitchProxy = 306 // no longer used

	// StatusTemporaryRedirect indicates that the target resource resides
	// temporarily under a different URI and the user agent MUST NOT change
	// the request method if it performs an automatic redirection to that URI.
	// RFC 9110, Section 15.4.8
	StatusTemporaryRedirect = http.StatusTemporaryRedirect // 307

	// StatusPermanentRedirect indicates that the target resource has been
	// assigned a new permanent URI and any future references to this resource
	// ought to use one of the enclosed URIs.
	// RFC 9110, Section 15.4.9
	StatusPermanentRedirect = http.StatusPermanentRedirect // 308

	// StatusBadRequest indicates that the server cannot or will not process
	// the request due to something that is perceived to be a client error
	// (e.g., malformed request syntax, invalid request message framing, or
	// deceptive request routing).
	// RFC 9110, Section 15.5.1
	StatusBadRequest = http.StatusBadRequest // 400

	// StatusUnauthorized indicates that the request has not been applied
	// because it lacks valid authentication credentials for the target
	// resource.
	// RFC 9110, Section 15.5.2
	StatusUnauthorized = http.StatusUnauthorized // 401

	// StatusPaymentRequired is reserved for future use.
	// RFC 9110, Section 15.5.3
	StatusPaymentRequired = http.StatusPaymentRequired // 402

	// StatusForbidden indicates that the server understood the request but
	// refuses to authorize it.
	// RFC 9110, Section 15.5.4
	StatusForbidden = http.StatusForbidden // 403

	// StatusNotFound indicates that the origin server did not find a current
	// representation for the target resource or is not willing to disclose
	// that one exists.
	// RFC 9110, Section 15.5.5
	StatusNotFound = http.StatusNotFound // 404

	// StatusMethodNotAllowed indicates that the method specified in the
	// request-line is known by the origin server but not supported by
	// the target resource.
	// RFC 9110, Section 15.5.6
	StatusMethodNotAllowed = http.StatusMethodNotAllowed // 405

	// StatusNotAcceptable indicates that the target resource does not have
	// a current representation that would be acceptable to the user agent,
	// according to the proactive negotiation header fields received in the
	// request, and the server is unwilling to supply a default representation.
	// RFC 9110, Section 15.5.7
	StatusNotAcceptable = http.StatusNotAcceptable // 406

	// StatusProxyAuthRequired is similar to 401 Unauthorized, but it indicates
	// that the client needs to authenticate itself in order to use a proxy.
	// RFC 9110, Section 15.5.8
	StatusProxyAuthRequired = http.StatusProxyAuthRequired // 407

	// StatusRequestTimeout indicates that the server did not receive
	// a complete request message within the time that it was prepared to wait.
	// RFC 9110, Section 15.5.9
	StatusRequestTimeout = http.StatusRequestTimeout // 408

	// StatusConflict indicates that the request could not be completed due
	// to a conflict with the current state of the target resource.
	// RFC 9110, Section 15.5.10
	StatusConflict = http.StatusConflict // 409

	// StatusGone indicates that access to the target resource is no longer
	// available at the origin server and that this condition is likely
	// to be permanent.
	// RFC 9110, Section 15.5.11
	StatusGone = http.StatusGone // 410

	// StatusLengthRequired indicates that the server refuses to accept
	// the request without a defined Content-Length.
	// RFC 9110, Section 15.5.12
	StatusLengthRequired = http.StatusLengthRequired // 411

	// StatusPreconditionFailed indicates that one or more conditions
	// given in the request header fields evaluated to false when tested
	// on the server.
	// RFC 9110, Section 15.5.13
	StatusPreconditionFailed = http.StatusPreconditionFailed // 412

	// StatusRequestEntityTooLarge is deprecated in favor of
	// StatusPayloadTooLarge.
	// The server is refusing to process a request because the request
	// payload is larger than the server is willing or able to process.
	// RFC 9110, Section 15.5.14
	StatusRequestEntityTooLarge = http.StatusRequestEntityTooLarge // 413

	// StatusRequestURITooLong indicates that the server is refusing
	// to service the request because the request-target is longer
	// than the server is willing to interpret.
	// RFC 9110, Section 15.5.15
	StatusRequestURITooLong = http.StatusRequestURITooLong // 414

	// StatusUnsupportedMediaType indicates that the origin server is
	// refusing to service the request because the payload is in a format
	// not supported by this method on the target resource.
	// RFC 9110, Section 15.5.16
	StatusUnsupportedMediaType = http.StatusUnsupportedMediaType // 415

	// StatusRequestedRangeNotSatisfiable indicates that none of the
	// ranges in the request's Range header field overlap the current
	// extent of the selected resource or that the set of ranges requested
	// has been rejected due to invalid ranges or an excessive request of
	// small or overlapping ranges.
	// RFC 9110, Section 15.5.17
	StatusRequestedRangeNotSatisfiable = http.StatusRequestedRangeNotSatisfiable // 416

	// StatusExpectationFailed indicates that the expectation given in the
	// request's Expect header field could not be met by at least one of
	// the inbound servers.
	// RFC 9110, Section 15.5.18
	StatusExpectationFailed = http.StatusExpectationFailed // 417

	// StatusTeapot is an April Fools' joke from RFC 2324, Hyper Text Coffee
	// Pot Control Protocol. It's not expected to be implemented by actual
	// HTTP servers.
	// RFC 2324, Section 2.3.2
	StatusTeapot = http.StatusTeapot // 418

	// StatusMisdirectedRequest indicates that the request was directed at
	// a server that is not able to produce a response. This can be sent by
	// a server that is not configured to produce responses for the combination
	// of scheme and authority that are included in the request URI.
	// RFC 9110, Section 15.5.20
	StatusMisdirectedRequest = http.StatusMisdirectedRequest // 421

	// StatusUnprocessableEntity indicates that the server understands the
	// content type of the request entity, and the syntax of the request entity
	// is correct, but it was unable to process the contained instructions.
	// RFC 4918, Section 11.2 (WebDAV)
	StatusUnprocessableEntity = http.StatusUnprocessableEntity // 422

	// StatusLocked indicates that the source or destination resource of a
	// method is locked.
	// RFC 4918, Section 11.3 (WebDAV)
	StatusLocked = http.StatusLocked // 423

	// StatusFailedDependency indicates that the method could not be performed
	// on the resource because the requested action depended on another action
	// and that action failed.
	// RFC 4918, Section 11.4 (WebDAV)
	StatusFailedDependency = http.StatusFailedDependency // 424

	// StatusTooEarly indicates that the server is unwilling to risk processing
	// a request that might be replayed.
	// RFC 8470, Section 5.2
	StatusTooEarly = http.StatusTooEarly // 425

	// StatusUpgradeRequired indicates that the server refuses to perform the
	// request using the current protocol but might be willing to do so after
	// the client upgrades to a different protocol.
	// RFC 9110, Section 15.5.22
	StatusUpgradeRequired = http.StatusUpgradeRequired // 426

	// StatusPreconditionRequired indicates that the origin server requires
	// the request to be conditional.
	// RFC 6585, Section 3
	StatusPreconditionRequired = http.StatusPreconditionRequired // 428

	// StatusTooManyRequests indicates the user has sent too many requests
	// in a given amount of time ("rate limiting").
	// RFC 6585, Section 4
	StatusTooManyRequests = http.StatusTooManyRequests // 429

	// StatusRequestHeaderFieldsTooLarge indicates that the server is unwilling
	// to process the request because its header fields are too large.
	// The request MAY be resubmitted after reducing the size of the request
	// header fields.
	// RFC 6585, Section 5
	StatusRequestHeaderFieldsTooLarge = http.StatusRequestHeaderFieldsTooLarge // 431

	// StatusUnavailableForLegalReasons indicates that the server is denying
	// access to the resource as a consequence of a legal demand.
	// RFC 7725, Section 3
	StatusUnavailableForLegalReasons = http.StatusUnavailableForLegalReasons // 451

	// StatusInternalServerError indicates that the server encountered an
	// unexpected condition that prevented it from fulfilling the request.
	// RFC 9110, Section 15.6.1
	StatusInternalServerError = http.StatusInternalServerError // 500

	// StatusNotImplemented indicates that the server does not support the
	// functionality required to fulfill the request.
	// RFC 9110, Section 15.6.2
	StatusNotImplemented = http.StatusNotImplemented // 501

	// StatusBadGateway indicates that the server, while acting as a gateway
	// or proxy, received an invalid response from an inbound server it
	// accessed while attempting to fulfill the request.
	// RFC 9110, Section 15.6.3
	StatusBadGateway = http.StatusBadGateway // 502

	// StatusServiceUnavailable indicates that the server is currently
	// unable to handle the request due to a temporary overload or scheduled
	// maintenance, which will likely be alleviated after some delay.
	// RFC 9110, Section 15.6.4
	StatusServiceUnavailable = http.StatusServiceUnavailable // 503

	// StatusGatewayTimeout indicates that the server, while acting as
	// a gateway or proxy, did not receive a timely response from an
	// upstream server it needed to access in order to complete the request.
	// RFC 9110, Section 15.6.5
	StatusGatewayTimeout = http.StatusGatewayTimeout // 504

	// StatusHTTPVersionNotSupported indicates that the server does not
	// support, or refuses to support, the major version of HTTP
	// that was used in the request message.
	// RFC 9110, Section 15.6.6
	StatusHTTPVersionNotSupported = http.StatusHTTPVersionNotSupported // 505

	// StatusVariantAlsoNegotiates indicates that the server has an internal
	// configuration error: the chosen variant resource is configured to
	// engage in transparent content negotiation itself, and is therefore
	// not a proper end point in the negotiation process.
	// RFC 2295, Section 8.1
	StatusVariantAlsoNegotiates = http.StatusVariantAlsoNegotiates // 506

	// StatusInsufficientStorage indicates that the server is unable to store
	// the representation needed to complete the request.
	// RFC 4918, Section 11.5 (WebDAV)
	StatusInsufficientStorage = http.StatusInsufficientStorage // 507

	// StatusLoopDetected indicates that the server detected an infinite loop
	// while processing the request (sent instead of 208 Already Reported).
	// RFC 5842, Section 7.2 (WebDAV Binding Extensions)
	StatusLoopDetected = http.StatusLoopDetected // 508

	// StatusNotExtended further extensions to the request are required for
	// the server to fulfill it.
	// RFC 2774, Section 7
	StatusNotExtended = http.StatusNotExtended // 510

	// StatusNetworkAuthenticationRequired indicates that the client needs
	// to authenticate to gain network access.
	// RFC 6585, Section 6
	StatusNetworkAuthenticationRequired = http.StatusNetworkAuthenticationRequired // 511
)

// statusMessages is a map that maps HTTP status codes to their corresponding
// status messages. The keys of the map are integers representing the status
// codes, and the values are strings representing the status messages.
var statusMessages = map[int]string{
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "Switching Protocols",
	StatusProcessing:         "Processing",
	StatusEarlyHints:         "Early Hints",

	StatusOK:                   "OK",
	StatusCreated:              "Created",
	StatusAccepted:             "Accepted",
	StatusNonAuthoritativeInfo: "Non-Authoritative Information",
	StatusNoContent:            "No Content",
	StatusResetContent:         "Reset Content",
	StatusPartialContent:       "Partial Content",
	StatusMultiStatus:          "Multi-Status",
	StatusAlreadyReported:      "Already Reported",
	StatusIMUsed:               "IM Used",

	StatusMultipleChoices:   "Multiple Choices",
	StatusMovedPermanently:  "Moved Permanently",
	StatusFound:             "Found",
	StatusSeeOther:          "See Other",
	StatusNotModified:       "Not Modified",
	StatusUseProxy:          "Use Proxy",
	StatusSwitchProxy:       "Switch Proxy",
	StatusTemporaryRedirect: "Temporary Redirect",
	StatusPermanentRedirect: "Permanent Redirect",

	StatusBadRequest:                   "Bad Request",
	StatusUnauthorized:                 "Unauthorized",
	StatusPaymentRequired:              "Payment Required",
	StatusForbidden:                    "Forbidden",
	StatusNotFound:                     "Not Found",
	StatusMethodNotAllowed:             "Method Not Allowed",
	StatusNotAcceptable:                "Not Acceptable",
	StatusProxyAuthRequired:            "Proxy Authentication Required",
	StatusRequestTimeout:               "Request Timeout",
	StatusConflict:                     "Conflict",
	StatusGone:                         "Gone",
	StatusLengthRequired:               "Length Required",
	StatusPreconditionFailed:           "Precondition Failed",
	StatusRequestEntityTooLarge:        "Request Entity Too Large",
	StatusRequestURITooLong:            "Request URI Too Long",
	StatusUnsupportedMediaType:         "Unsupported Media Type",
	StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	StatusExpectationFailed:            "Expectation Failed",
	StatusTeapot:                       "I'm a teapot",
	StatusMisdirectedRequest:           "Misdirected Request",
	StatusUnprocessableEntity:          "Unprocessable Entity",
	StatusLocked:                       "Locked",
	StatusFailedDependency:             "Failed Dependency",
	StatusTooEarly:                     "Too Early",
	StatusUpgradeRequired:              "Upgrade Required",
	StatusPreconditionRequired:         "Precondition Required",
	StatusTooManyRequests:              "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	StatusInternalServerError:           "Internal Server Error",
	StatusNotImplemented:                "Not Implemented",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "Service Unavailable",
	StatusGatewayTimeout:                "Gateway Timeout",
	StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	StatusVariantAlsoNegotiates:         "Variant Also Negotiates",
	StatusInsufficientStorage:           "Insufficient Storage",
	StatusLoopDetected:                  "Loop Detected",
	StatusNotExtended:                   "Not Extended",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

// singleHeaders is a slice of strings that contains all the HTTP headers that
// are not lists of values.
var singleHeaders = []string{
	HeaderContentType,
	HeaderETag,
	HeaderLastModified,
	HeaderContentLength,
	HeaderUserAgent,
	HeaderHost,
	HeaderReferer,
	HeaderServer,
	HeaderDate,
	HeaderLocation,
	HeaderRetryAfter,
	HeaderContentDisposition,
	HeaderContentEncoding,
	HeaderContentLanguage,
	HeaderContentLocation,
	// HeaderWWWAuthenticate,
	// HeaderAuthorization,
	// HeaderProxyAuthenticate,
	// HeaderProxyAuthorization,
	// HeaderIfMatch,
	// HeaderIfNoneMatch,
	HeaderIfModifiedSince,
	HeaderIfUnmodifiedSince,
	HeaderIfRange,
	// HeaderContentSecurityPolicy,
	HeaderStrictTransportSecurity,
	HeaderUpgradeInsecureRequests,
	HeaderXContentTypeOptions,
	HeaderXFrameOptions,
	HeaderXXSSProtection,
	HeaderContentDPR,
	HeaderDPR,
	HeaderViewportWidth,
	HeaderWidth,
	HeaderContentRange,
}

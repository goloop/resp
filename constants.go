package resp

import "net/http"

// MIME types that are commonly used
const (
	// MIMETextXML is the MIME type for XML documents.
	MIMETextXML = "text/xml"

	// MIMETextHTML is the MIME type for HTML documents.
	MIMETextHTML = "text/html"

	// MIMETextPlain is the MIME type for plain text.
	MIMETextPlain = "text/plain"

	// MIMETextJavaScript is the MIME type for JavaScript code.
	MIMETextJavaScript = "text/javascript"

	// MIMEApplicationXML is the MIME type for XML documents,
	// typically used for APIs or web services.
	MIMEApplicationXML = "application/xml"

	// MIMEApplicationJSON is the MIME type for JSON formatted data.
	MIMEApplicationJSON = "application/json"

	// MIMEApplicationJavaScript is the MIME type for JavaScript code,
	// often used for APIs serving JavaScript.
	MIMEApplicationJavaScript = "application/javascript"

	// MIMEApplicationForm is the MIME type for URL-encoded form data.
	MIMEApplicationForm = "application/x-www-form-urlencoded"

	// MIMEOctetStream is the MIME type for arbitrary binary data.
	MIMEOctetStream = "application/octet-stream"

	// MIMEMultipartForm is the MIME type for multipart form data,
	// used for form submissions that include file uploads.
	MIMEMultipartForm = "multipart/form-data"

	// MIMETextXMLCharsetUTF8 is the MIME type for XML documents
	// using UTF-8 character encoding.
	MIMETextXMLCharsetUTF8 = "text/xml; charset=utf-8"

	// MIMETextHTMLCharsetUTF8 is the MIME type for HTML documents
	// using UTF-8 character encoding.
	MIMETextHTMLCharsetUTF8 = "text/html; charset=utf-8"

	// MIMETextPlainCharsetUTF8 is the MIME type for plain text
	// using UTF-8 character encoding.
	MIMETextPlainCharsetUTF8 = "text/plain; charset=utf-8"

	// MIMETextJavaScriptCharsetUTF8 is the MIME type for JavaScript
	// code using UTF-8 character encoding.
	MIMETextJavaScriptCharsetUTF8 = "text/javascript; charset=utf-8"

	// MIMEApplicationXMLCharsetUTF8 is the MIME type for XML documents
	// using UTF-8 character encoding.
	MIMEApplicationXMLCharsetUTF8 = "application/xml; charset=utf-8"

	// MIMEApplicationJSONCharsetUTF8 is the MIME type for JSON formatted
	// data using UTF-8 character encoding.
	MIMEApplicationJSONCharsetUTF8 = "application/json; charset=utf-8"

	// MIMEApplicationJavaScriptCharsetUTF8 is the MIME type for JavaScript
	// code using UTF-8 character encoding.
	MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript; charset=utf-8"
)

// HTTP Headers were copied from net/http.
const (
	// HeaderAuthorization is the HTTP header that represents the credentials
	// that the client uses to authenticate itself when making a request
	// to the server.
	HeaderAuthorization = "Authorization"

	// HeaderProxyAuthenticate is the HTTP header that represents the
	// credentials that the proxy can accept for authenticating the client.
	HeaderProxyAuthenticate = "Proxy-Authenticate"

	// HeaderProxyAuthorization is the HTTP header that represents the
	// credentials that the client uses to authenticate itself to a proxy.
	HeaderProxyAuthorization = "Proxy-Authorization"

	// HeaderWWWAuthenticate is the HTTP header that represents the method
	// that should be used to authenticate the client to the server.
	HeaderWWWAuthenticate = "WWW-Authenticate"

	// HeaderAge is the HTTP header that represents the time in seconds
	// the object has been in a proxy cache.
	HeaderAge = "Age"

	// HeaderCacheControl is the HTTP header that represents directives
	// for caching mechanisms in both requests and responses.
	HeaderCacheControl = "Cache-Control"

	// HeaderClearSiteData is the HTTP header that tells the browser to
	// clear various types of cached data (cookies, storage, etc.).
	HeaderClearSiteData = "Clear-Site-Data"

	// HeaderExpires is the HTTP header that represents the date/time
	// after which the response is considered stale.
	HeaderExpires = "Expires"

	// HeaderPragma is the HTTP header that represents implementation-specific
	// directives that might apply to any agent along the request/response
	// chain.
	HeaderPragma = "Pragma"

	// HeaderWarning is the HTTP header that represents general warning
	// information about possible problems.
	HeaderWarning = "Warning"

	// HeaderAcceptCH is the HTTP header that represents the client hints
	// that the server is willing to accept.
	HeaderAcceptCH = "Accept-CH"

	// HeaderAcceptCHLifetime is the HTTP header that represents the lifetime
	// of the client hints that the server supports.
	HeaderAcceptCHLifetime = "Accept-CH-Lifetime"

	// HeaderContentDPR is the HTTP header that represents the device pixel
	// ratio (DPR) of the client's device.
	HeaderContentDPR = "Content-DPR"

	// HeaderDPR is the HTTP header that represents the client's device pixel
	// ratio for content negotiation.
	HeaderDPR = "DPR"

	// HeaderEarlyData is the HTTP header that represents the indication that
	// the request has been made with early data.
	HeaderEarlyData = "Early-Data"

	// HeaderSaveData is the HTTP header that represents the client's
	// preference for data saving, useful for clients in data-saving mode.
	HeaderSaveData = "Save-Data"

	// HeaderViewportWidth is the HTTP header that represents the width of
	// the viewport in pixels, a hint to help the server select optimized
	// content for the device.
	HeaderViewportWidth = "Viewport-Width"

	// HeaderWidth is the HTTP header that represents the width of the
	// requested resource, a hint that allows servers to serve different
	// versions based on the content width.
	HeaderWidth = "Width"

	// HeaderETag is the HTTP header that represents the entity tag of the
	// resource, a mechanism for cache validation and conditional requests.
	HeaderETag = "ETag"

	// HeaderIfMatch is the HTTP header that makes the request method
	// conditional on the resource's ETag matching one of the listed ETags.
	HeaderIfMatch = "If-Match"

	// HeaderIfModifiedSince is the HTTP header that makes the request
	// method conditional on the resource being modified after the
	// specified date.
	HeaderIfModifiedSince = "If-Modified-Since"

	// HeaderIfNoneMatch is the HTTP header that makes the request method
	// conditional on the resource's ETag not matching any of the listed ETags.
	HeaderIfNoneMatch = "If-None-Match"

	// HeaderIfUnmodifiedSince is the HTTP header that makes the request
	// method conditional on the resource not being modified after the
	// specified date.
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"

	// HeaderLastModified is the HTTP header that represents the date and
	// time at which the server believes the resource was last modified.
	HeaderLastModified = "Last-Modified"

	// HeaderVary is the HTTP header that indicates the set of request-header
	// fields that fully determines the response sent by the server.
	HeaderVary = "Vary"

	// HeaderConnection is the HTTP header that controls whether the network
	// connection stays open after the current transaction finishes.
	HeaderConnection = "Connection"

	// HeaderKeepAlive is the HTTP header that specifies directives for
	// managing the persistence of the connection.
	HeaderKeepAlive = "Keep-Alive"

	// HeaderAccept is the HTTP header that specifies the media types which
	// are acceptable for the response.
	HeaderAccept = "Accept"

	// HeaderAcceptCharset is the HTTP header that specifies the character
	// sets that are acceptable for the response.
	HeaderAcceptCharset = "Accept-Charset"

	// HeaderAcceptEncoding is the HTTP header that specifies the content
	// encodings that are acceptable in the response.
	HeaderAcceptEncoding = "Accept-Encoding"

	// HeaderAcceptLanguage is the HTTP header that specifies the natural
	// languages that are preferred as a response to the request.
	HeaderAcceptLanguage = "Accept-Language"

	// HeaderCookie is the HTTP header that contains stored HTTP cookies
	// previously sent by the server with the Set-Cookie header.
	HeaderCookie = "Cookie"

	// HeaderExpect is the HTTP header that indicates that particular
	// server behaviors are required by the client.
	HeaderExpect = "Expect"

	// HeaderMaxForwards is the HTTP header that represents the maximum
	// number of times that the request can be forwarded through proxies
	// or gateways.
	HeaderMaxForwards = "Max-Forwards"

	// HeaderSetCookie is the HTTP header that represents the instruction
	// sent by the server to the client to store a cookie.
	HeaderSetCookie = "Set-Cookie"

	// HeaderAccessControlAllowCredentials is the HTTP header that indicates
	// whether the response to the request can be exposed when the credentials
	// flag is true.
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"

	// HeaderAccessControlAllowHeaders is the HTTP header that specifies
	// the headers that can be used during the actual request.
	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"

	// HeaderAccessControlAllowMethods is the HTTP header that specifies
	// the methods allowed when accessing the resource in response to
	// a preflight request.
	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"

	// HeaderAccessControlAllowOrigin is the HTTP header that specifies
	// the origins that are allowed to access the resource.
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"

	// HeaderAccessControlExposeHeaders is the HTTP header that lets
	// a server whitelist headers that browsers are allowed to access.
	HeaderAccessControlExposeHeaders = "Access-Control-Expose-Headers"

	// HeaderAccessControlMaxAge is the HTTP header that indicates how
	// long the results of a preflight request can be cached.
	HeaderAccessControlMaxAge = "Access-Control-Max-Age"

	// HeaderAccessControlRequestHeaders is the HTTP header used in preflight
	// requests to indicate which HTTP headers can be used during the actual
	// request.
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"

	// HeaderAccessControlRequestMethod is the HTTP header used in preflight
	// requests to indicate which HTTP method can be used during the actual
	// request.
	HeaderAccessControlRequestMethod = "Access-Control-Request-Method"

	// HeaderOrigin is the HTTP header that indicates the origin of the
	// request or the page making the request.
	HeaderOrigin = "Origin"

	// HeaderTimingAllowOrigin is the HTTP header that specifies which
	// origins are allowed to include timing resources in performance data.
	HeaderTimingAllowOrigin = "Timing-Allow-Origin"

	// HeaderXPermittedCrossDomainPolicies is the HTTP header that specifies
	// the policy that on-site or cross-site policy files can be used to grant
	// cross-domain loading of content.
	HeaderXPermittedCrossDomainPolicies = "X-Permitted-Cross-Domain-Policies"

	// HeaderDNT is the HTTP header that represents the "Do Not Track"
	// request header, indicating the user's tracking preference.
	HeaderDNT = "DNT"

	// HeaderTk is the HTTP header that represents the tracking status sent
	// in response to a DNT request.
	HeaderTk = "Tk"

	// HeaderContentDisposition is the HTTP header that represents directives
	// for content disposition,
	// such as inline or attachment.
	HeaderContentDisposition = "Content-Disposition"

	// HeaderContentEncoding is the HTTP header that represents the encoding
	// transformations that have been applied to the content.
	HeaderContentEncoding = "Content-Encoding"

	// HeaderContentLanguage is the HTTP header that represents the natural
	// language(s) of the intended audience for the content.
	HeaderContentLanguage = "Content-Language"

	// HeaderContentLength is the HTTP header that represents the size of the
	// entity-body in bytes.
	HeaderContentLength = "Content-Length"

	// HeaderContentLocation is the HTTP header that represents an alternate
	// location for the returned content.
	HeaderContentLocation = "Content-Location"

	// HeaderContentType is the HTTP header that represents the media type
	// of the content.
	HeaderContentType = "Content-Type"

	// HeaderForwarded is the HTTP header that represents the details of
	// the forwarding mechanism, including the originating IP.
	HeaderForwarded = "Forwarded"

	// HeaderVia is the HTTP header that represents intermediate protocols
	// and recipients between the user agent and the server on requests,
	// and between the origin server and the client on responses.
	HeaderVia = "Via"

	// HeaderXForwardedFor is the HTTP header that represents the originating
	// IP addresses of a client connecting to a web server through an HTTP
	// proxy or a load balancer.
	HeaderXForwardedFor = "X-Forwarded-For"

	// HeaderXForwardedHost is the HTTP header that represents the original
	// host requested by the client in the Host HTTP request header.
	HeaderXForwardedHost = "X-Forwarded-Host"

	// HeaderXForwardedProto is the HTTP header that represents the protocol
	// (HTTP or HTTPS) that a client used to connect to your proxy or load
	// balancer.
	HeaderXForwardedProto = "X-Forwarded-Proto"

	// HeaderXForwardedProtocol is similar to X-Forwarded-Proto, and it
	// represents the protocol used by the client.
	HeaderXForwardedProtocol = "X-Forwarded-Protocol"

	// HeaderXForwardedSSL is the HTTP header that represents the usage of
	// SSL (HTTPS) by the client.
	HeaderXForwardedSSL = "X-Forwarded-Ssl"

	// HeaderXUrlScheme is the HTTP header that represents the scheme
	// (HTTP/HTTPS) part of the URL requested by the client.
	HeaderXUrlScheme = "X-Url-Scheme"

	// HeaderLocation is the HTTP header that represents the URL to
	// redirect a page to.
	HeaderLocation = "Location"

	// HeaderFrom is the HTTP header that represents the email address
	// of the user making the request.
	HeaderFrom = "From"

	// HeaderHost is the HTTP header that represents the domain name of
	// the server (for virtual hosting), and optionally the TCP port number
	// on which the server is listening.
	HeaderHost = "Host"

	// HeaderReferer is the HTTP header that represents the address of
	// the previous web page from which a link to the currently requested
	// page was followed.
	HeaderReferer = "Referer"

	// HeaderReferrerPolicy is the HTTP header that represents the policy
	// governing which referrer information sent in the Referer header should
	// be included with requests made.
	HeaderReferrerPolicy = "Referrer-Policy"

	// HeaderUserAgent is the HTTP header that represents the string
	// identifying the user agent originating the request.
	HeaderUserAgent = "User-Agent"

	// HeaderAllow is the HTTP header that represents the list of HTTP
	// request methods supported by a resource.
	HeaderAllow = "Allow"

	// HeaderServer is the HTTP header that represents the web server
	// software being used on the responding server.
	HeaderServer = "Server"

	// HeaderAcceptRanges is the HTTP header that represents the range
	// of bytes that the server accepts in the Range header of a request.
	HeaderAcceptRanges = "Accept-Ranges"

	// HeaderContentRange is the HTTP header that represents the part of
	// a document that the server is returning in a partial request response.
	HeaderContentRange = "Content-Range"

	// HeaderIfRange is the HTTP header that represents a condition for
	// sending a partial copy of a resource only if a previous copy has
	// not been modified.
	HeaderIfRange = "If-Range"

	// HeaderRange is the HTTP header that represents the specific part
	// of a document requested by a client.
	HeaderRange = "Range"

	// HeaderContentSecurityPolicy is the HTTP header that represents the
	// policy that specifies valid sources for content on a webpage.
	HeaderContentSecurityPolicy = "Content-Security-Policy"

	// HeaderContentSecurityPolicyReportOnly is the HTTP header that represents
	// the policy that specifies valid sources for content on a webpage but is
	// reported only.
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"

	// HeaderCrossOriginResourcePolicy is the HTTP header that represents the
	// policy that allows a server to specify the origins that are allowed to
	// load its resources.
	HeaderCrossOriginResourcePolicy = "Cross-Origin-Resource-Policy"

	// HeaderExpectCT is the HTTP header that represents the policy that allows
	// sites to opt in to reporting and/or enforcement of Certificate
	// Transparency requirements.
	HeaderExpectCT = "Expect-CT"

	// HeaderPermissionsPolicy is the HTTP header that represents the
	// permissions policy for a web page, specifying which features are
	// allowed to be used by the page.
	HeaderPermissionsPolicy = "Permissions-Policy"

	// HeaderPublicKeyPins is the HTTP header that represents the Public
	// Key Pinning Extension for HTTP, associating a specific cryptographic
	// public key with a web server to decrease the risk of MITM attacks with
	// forged certificates.
	HeaderPublicKeyPins = "Public-Key-Pins"

	// HeaderPublicKeyPinsReportOnly is the HTTP header that represents the
	// Public Key Pinning Extension for HTTP for report-only mode, which
	// reports pin validation failures without enforcing them.
	HeaderPublicKeyPinsReportOnly = "Public-Key-Pins-Report-Only"

	// HeaderStrictTransportSecurity is the HTTP header that represents the
	// policy that enforces secure (HTTPS) connections to the server.
	HeaderStrictTransportSecurity = "Strict-Transport-Security"

	// HeaderUpgradeInsecureRequests is the HTTP header that represents the
	// request that the server upgrades the connection to HTTPS.
	HeaderUpgradeInsecureRequests = "Upgrade-Insecure-Requests"

	// HeaderXContentTypeOptions is the HTTP header that represents the
	// directive for the browser to disable MIME type sniffing and strictly
	// follow the declared content type.
	HeaderXContentTypeOptions = "X-Content-Type-Options"

	// HeaderXDownloadOptions is the HTTP header that represents the directive
	// for Internet Explorer 8+ to not open downloads directly in the browser.
	HeaderXDownloadOptions = "X-Download-Options"

	// HeaderXFrameOptions is the HTTP header that represents the directive
	// that specifies whether a browser should be allowed to render a page
	// in a <frame>, <iframe>, <embed>, or <object>.
	HeaderXFrameOptions = "X-Frame-Options"

	// HeaderXPoweredBy is the HTTP header that represents information
	// about the technology supporting the web server (often removed for
	// security reasons).
	HeaderXPoweredBy = "X-Powered-By"

	// HeaderXXSSProtection is the HTTP header that represents the directive
	// for cross-site scripting filter built into most modern web browsers.
	HeaderXXSSProtection = "X-XSS-Protection"

	// HeaderLastEventID is the HTTP header that represents the ID of the
	// last event in an EventSource / Server-Sent Events connection, allowing
	// resuming of event streams after a disconnection.
	HeaderLastEventID = "Last-Event-ID"

	// HeaderNEL is the HTTP header that represents the Network Error Logging
	// policy, which allows a server to collect network error reports from
	// user agents.
	HeaderNEL = "NEL"

	// HeaderPingFrom is the HTTP header that represents the URL of the Web
	// page that initiated a ping request during hyperlink auditing.
	HeaderPingFrom = "Ping-From"

	// HeaderPingTo is the HTTP header that represents the URL to which a ping
	// request during hyperlink auditing is sent.
	HeaderPingTo = "Ping-To"

	// HeaderReportTo is the HTTP header that represents the endpoint
	// groups for the Reporting API, allowing the browser to report errors
	// to the server.
	HeaderReportTo = "Report-To"

	// HeaderTE is the HTTP header that represents the transfer encodings
	// the user agent is willing to accept: chunked, compress, deflate, gzip,
	// identity.
	HeaderTE = "TE"

	// HeaderTrailer is the HTTP header that represents the header fields
	// present in the trailer of a message encoded with chunked transfer coding.
	HeaderTrailer = "Trailer"

	// HeaderTransferEncoding is the HTTP header that represents the form of
	// encoding used to safely transfer the payload body to the user.
	HeaderTransferEncoding = "Transfer-Encoding"

	// HeaderSecWebSocketAccept is the HTTP header that represents the
	// server's acceptance of a WebSocket handshake request to establish
	// a WebSocket connection.
	HeaderSecWebSocketAccept = "Sec-WebSocket-Accept"

	// HeaderSecWebSocketExtensions is the HTTP header that represents the
	// accepted extensions for a WebSocket connection.
	HeaderSecWebSocketExtensions = "Sec-WebSocket-Extensions"

	// HeaderSecWebSocketKey is the HTTP header that represents the encoded
	// key for the server to accept the WebSocket connection.
	HeaderSecWebSocketKey = "Sec-WebSocket-Key"

	// HeaderSecWebSocketProtocol is the HTTP header that represents the
	// agreed-upon protocol during a WebSocket connection handshake.
	HeaderSecWebSocketProtocol = "Sec-WebSocket-Protocol"

	// HeaderSecWebSocketVersion is the HTTP header that represents
	// the WebSocket protocol version being used by the client.
	HeaderSecWebSocketVersion = "Sec-WebSocket-Version"

	// HeaderAcceptPatch is the HTTP header that represents the patch
	// document formats accepted by the server.
	HeaderAcceptPatch = "Accept-Patch"

	// HeaderAcceptPushPolicy is the HTTP header that represents the
	// server's preferences for HTTP/2 server push.
	HeaderAcceptPushPolicy = "Accept-Push-Policy"

	// HeaderAcceptSignature is the HTTP header that represents the
	// client's support for the HTTP Signatures.
	HeaderAcceptSignature = "Accept-Signature"

	// HeaderAltSvc is the HTTP header that represents an alternative
	// service for accessing a resource.
	HeaderAltSvc = "Alt-Svc"

	// HeaderDate is the HTTP header that represents the date and time
	// at which the message was originated.
	HeaderDate = "Date"

	// HeaderIndex is the HTTP header that represents index resources
	// for a collection of resources.
	HeaderIndex = "Index"

	// HeaderLargeAllocation is the HTTP header that hints to the browser
	// that a large allocation will be made.
	HeaderLargeAllocation = "Large-Allocation"

	// HeaderLink is the HTTP header that represents relationships between
	// the current document and an external resource.
	HeaderLink = "Link"

	// HeaderPushPolicy is the HTTP header that represents the server's
	// policy for HTTP/2 server push.
	HeaderPushPolicy = "Push-Policy"

	// HeaderRetryAfter is the HTTP header that represents the amount of
	// time the client should wait before making a follow-up request.
	HeaderRetryAfter = "Retry-After"

	// HeaderServerTiming is the HTTP header that represents the server
	// timing for performance tracking.
	HeaderServerTiming = "Server-Timing"

	// HeaderSignature is the HTTP header that represents the digital
	// signature for the message content for verification.
	HeaderSignature = "Signature"

	// HeaderSignedHeaders is the HTTP header that represents the list
	// of headers that are included in the digital signature.
	HeaderSignedHeaders = "Signed-Headers"

	// HeaderSourceMap is the HTTP header that provides a link to a source
	// map for debugging purposes.
	HeaderSourceMap = "SourceMap"

	// HeaderUpgrade is the HTTP header that requests the client to switch
	// to a different protocol.
	HeaderUpgrade = "Upgrade"

	// HeaderXDNSPrefetchControl is the HTTP header that controls DNS
	// prefetching, allowing browsers to resolve domain names before
	// resources are requested.
	HeaderXDNSPrefetchControl = "X-DNS-Prefetch-Control"

	// HeaderXPingback is the HTTP header that specifies the pingback
	// URL for the resource.
	HeaderXPingback = "X-Pingback"

	// HeaderXRequestID is the HTTP header that provides a unique identifier
	// for the request, facilitating tracing and debugging.
	HeaderXRequestID = "X-Request-ID"

	// HeaderXRequestedWith is the HTTP header that identifies the request
	// as being made with a particular technology, often used to identify
	// Ajax requests.
	HeaderXRequestedWith = "X-Requested-With"

	// HeaderXRobotsTag is the HTTP header that provides directives to search
	// engines for indexing and serving the content.
	HeaderXRobotsTag = "X-Robots-Tag"

	// HeaderXUACompatible is the HTTP header that advises the web browser
	// to use a specific rendering engine (e.g., IE=edge).
	HeaderXUACompatible = "X-UA-Compatible"

	// HeaderAccessControlAllowPrivateNetwork is the HTTP header that grants
	// web applications on the public internet access to resources on the
	// user’s private network.
	HeaderAccessControlAllowPrivateNetwork = "Access-Control-Allow-Private-Network"

	// HeaderAccessControlRequestPrivateNetwork is the HTTP header used in
	// preflight requests to indicate access to the user’s private network
	// is requested by the web application.
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

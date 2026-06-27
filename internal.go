package resp

import (
	"bytes"
	"regexp"
	"strings"
	"sync"
)

// maxPooledBuffer is the upper bound (in bytes) on the capacity of a
// *bytes.Buffer that may be returned to bufPool. Buffers that grew past
// this size while encoding a large response are dropped instead of being
// pooled, so the pool does not pin oversized allocations in memory.
const maxPooledBuffer = 1 << 20 // 1 MiB

// maxJSONPCallbackLen is the maximum accepted length of a JSONP callback
// name. It is a defensive bound: legitimate callbacks are short, while an
// attacker-supplied value could be arbitrarily long.
const maxJSONPCallbackLen = 128

// bufPool recycles *bytes.Buffer values used to render JSON/JSONP/XML
// payloads before they are flushed to the http.ResponseWriter. Buffering
// lets us detect an encoding error *before* the status line is committed
// (see Response.JSON), while the pool keeps the per-call allocation count
// on par with direct streaming.
var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

// getBuf returns a reset buffer ready for use.
func getBuf() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

// putBuf returns buf to the pool unless it grew beyond maxPooledBuffer,
// in which case it is left for the garbage collector.
func putBuf(buf *bytes.Buffer) {
	if buf.Cap() > maxPooledBuffer {
		return
	}
	buf.Reset()
	bufPool.Put(buf)
}

// jsonpCallbackRe compiles, lazily and exactly once, the whitelist that a
// JSONP callback name must match. The grammar mirrors a JavaScript member
// expression: an identifier, optionally followed by dotted identifiers
// (e.g. "cb", "app.handlers.onData"). Anything outside this set — quotes,
// parentheses, angle brackets, semicolons, whitespace — is rejected, which
// closes the reflected-XSS / script-injection vector around JSONP padding.
var jsonpCallbackRe = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*(\.[A-Za-z_$][A-Za-z0-9_$]*)*$`)
})

// validJSONPCallback reports whether callback is a safe JSONP callback name.
func validJSONPCallback(callback string) bool {
	if callback == "" || len(callback) > maxJSONPCallbackLen {
		return false
	}
	return jsonpCallbackRe().MatchString(callback)
}

// rfc5987Encode percent-encodes s into an RFC 5987 ext-value, the form
// required by the Content-Disposition filename* parameter. Only RFC 5987
// "attr-char" bytes are emitted literally; every other byte of the UTF-8
// representation is written as %XX, so arbitrary Unicode file names round
// trip to compliant clients.
func rfc5987Encode(s string) string {
	const upperhex = "0123456789ABCDEF"
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if c := s[i]; isAttrChar(c) {
			b.WriteByte(c)
		} else {
			b.WriteByte('%')
			b.WriteByte(upperhex[c>>4])
			b.WriteByte(upperhex[c&0x0f])
		}
	}
	return b.String()
}

// isAttrChar reports whether c is an RFC 5987 "attr-char" that may appear
// literally in an ext-value.
func isAttrChar(c byte) bool {
	switch {
	case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9':
		return true
	}
	switch c {
	case '!', '#', '$', '&', '+', '-', '.', '^', '_', '`', '|', '~':
		return true
	}
	return false
}

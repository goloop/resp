# resp — reference

The full reference for the `resp` package: the mental model, the response
helpers, the `Response` builder, options, the safe-by-default encoding
guarantee, security helpers, errors and practical recipes.

Ukrainian version: **[DOC.UK.md](DOC.UK.md)**.

## Contents

- [Mental model](#mental-model)
- [Response helpers](#response-helpers)
- [The Response builder](#the-response-builder)
- [Options](#options)
- [Error responses](#error-responses)
- [Safe-by-default encoding](#safe-by-default-encoding)
- [JSONP](#jsonp)
- [Custom JSON encoder](#custom-json-encoder)
- [Field projection](#field-projection)
- [Files and downloads](#files-and-downloads)
- [Security headers](#security-headers)
- [Concurrency](#concurrency)
- [Errors](#errors)
- [Recipes and tips](#recipes-and-tips)

## Mental model

`resp` crafts HTTP responses on top of the standard `net/http`. Instead of
setting a content type, writing a status and encoding a body by hand, you call
one helper that does all three consistently — and reports an encoding error
before anything is committed.

Two shapes cover every use:

1. **Package helpers** — `resp.JSON(w, data, opts...)` and friends: one call, a
   `http.ResponseWriter`, optional functional options.
2. **A `Response` builder** — `resp.NewResponse(w).SetStatus(...).SetHeader(...)`:
   chainable configuration, then a terminal method (`JSON`, `HTML`, …).

Both accept the same `Option` values, so anything you can do inline you can also
do on a builder.

```go
import "github.com/goloop/resp/v2"
```

The convenience type `resp.R` is `map[string]any`, handy for ad-hoc JSON:

```go
resp.JSON(w, resp.R{"message": "Hello", "status": "success"})
```

## Response helpers

Each helper takes `w`, its payload, and optional `...Option`:

| Helper | Sends |
|--------|-------|
| `JSON(w, data)`            | `application/json` body |
| `JSONIndent(w, data, prefix, indent)` | pretty-printed JSON |
| `JSONP(w, data, callback)` | validated JSONP |
| `XML(w, data)`             | `application/xml` body |
| `HTML(w, html)`            | `text/html` body |
| `String(w, text)`          | `text/plain` body |
| `Blob(w, contentType, bytes)` | raw bytes, `Content-Length` set |
| `Stream(w, reader)`        | streamed body from an `io.Reader` |
| `Error(w, status, message)` | JSON error body (see below) |
| `Created(w, location, data)` | `201` + `Location` + JSON body |
| `Accepted(w, data)`        | `202` + JSON body |
| `NoContent(w)`             | `204` |
| `Redirect(w, url)`         | redirect (valid 3xx only) |
| `ServeFile(w, r, filename)`| serve a file from disk |
| `ServeFileAsDownload(w, filename, data)` | in-memory download |

Default content types are applied automatically (`application/json;
charset=utf-8`, etc.), and the default status is `200 OK`.

## The Response builder

```go
func NewResponse(w http.ResponseWriter, opts ...Option) *Response
```

The builder mirrors every helper as a method and adds chainable configuration:

```go
resp.NewResponse(w).
    SetStatus(http.StatusOK).
    SetHeader("X-Custom-Header", "value").
    JSON(resp.R{"status": "success"})
```

Configuration methods (`SetStatus`, `SetHeader`/`AddHeader`/`DelHeader`,
`SetCookie`/`DelCookie`/`ExpiredCookie`, `SetJSONEncoder`, …) return `*Response`
for chaining; the terminal methods (`JSON`, `HTML`, `XML`, `String`, `Blob`,
`Stream`, `Error`, `Created`, `Accepted`, `NoContent`, `Redirect`, `ServeFile`)
write the response and return `error`. Getters `Status()` and `Header()` read
back the current state.

## Options

Options are `func(*Response) *Response` values accepted by every helper and by
`NewResponse`. They fall into a few families:

- **Status** — `WithStatus(code)` plus a named constant for every status
  (`WithStatusCreated()`, `WithStatusNotFound()`, …).
- **Headers** — `WithHeader(key, values...)`, plus typed `Add…` helpers for the
  full set of standard headers (`AddETag`, `AddLastModified`, `AddCacheControl`,
  `AddContentDisposition(type, filename)`, CORS `AddAccessControl…`, `AddLink`,
  `AddWarning`, and many more).
- **Content type** — `As…` shortcuts (`AsApplicationJSON()`, `AsTextHTML()`,
  `AsOctetStream()`, …).
- **Cookies** — `WithCookie(cookie)`.
- **JSON behaviour** — `WithStreamJSON()`, `ApplyJSONEncoder(fn)`.
- **Errors** — `WithErrorCode(n)`.
- **Security** — `SecurityHeaders()`, `WithCSP(policy)`.

```go
resp.JSON(w, data,
    resp.WithStatus(http.StatusOK),
    resp.WithHeader("X-Custom-Header", "value"),
    resp.AddETag(`"123456"`),
)
```

## Error responses

```go
func Error(w http.ResponseWriter, status int, message string, opts ...Option) error
```

The **first argument is always the HTTP status**. It is written to the status
line and mirrored into the JSON body's `code` field, so the two never disagree:

```go
// HTTP 404, body {"code":404,"message":"Not Found"}.
resp.Error(w, http.StatusNotFound, "")

// HTTP 404, body {"code":1042,"message":"User not found"}.
resp.Error(w, http.StatusNotFound, "User not found", resp.WithErrorCode(1042))
```

`WithErrorCode(n)` sets a distinct business code in the body while keeping the
HTTP status. The status argument takes precedence over a `WithStatus` option,
and an out-of-range value falls back to `500 Internal Server Error`.

## Safe-by-default encoding

By default `JSON`/`JSONP`/`XML` encode into a pooled buffer first, so a
serialization error is returned **before** the status line is committed — the
handler can still send a proper error response:

```go
if err := resp.JSON(w, data); err != nil {
    // Nothing has been written yet, so a fallback still works.
    resp.Error(w, http.StatusInternalServerError, "Internal server error")
    return
}
```

For very large payloads where memory matters more than that guarantee, opt into
zero-copy streaming:

```go
resp.JSON(w, largeData, resp.WithStreamJSON())
```

A pooled buffer keeps the default (safe) path at the same allocation count as
direct streaming.

## JSONP

The callback name is validated against a strict JavaScript-identifier whitelist
**before** anything is written; an unsafe name returns
[`ErrInvalidJSONPCallback`](#errors) and writes nothing, closing the
reflected-XSS vector:

```go
callback := r.URL.Query().Get("callback")
if err := resp.JSONP(w, data, callback); err != nil {
    // errors.Is(err, resp.ErrInvalidJSONPCallback) for a rejected name.
    resp.Error(w, http.StatusBadRequest, "invalid callback")
}
```

## Custom JSON encoder

Plug in any encoder satisfying `JSONEncodeFunc` (`func(w io.Writer, v any) error`):

```go
enc := func(w io.Writer, v any) error {
    e := json.NewEncoder(w)
    e.SetEscapeHTML(false)
    return e.Encode(v)
}

resp.NewResponse(w, resp.ApplyJSONEncoder(enc)).JSON(resp.R{"message": "ok"})
```

`JSONIndent(w, data, prefix, indent)` produces pretty-printed JSON for
human-facing or debug endpoints.

## Field projection

Shape the payload before serialization:

```go
data := resp.OnlyFields(user, "id", "email")     // keep only these
data = resp.ExcludeFields(user, "password_hash") // or drop these
resp.JSON(w, data)
```

Both follow `json` tags: a tagged field is projected under its `json` name, a
field may be selected by either its `json` name or its Go field name, and a
`json:"-"` field is never emitted (so it cannot leak through the redaction).
Structs without json tags are keyed by their Go field names.

## Files and downloads

```go
// Serve a static file, forcing a download name.
resp.ServeFile(w, r, "path/to/file.pdf",
    resp.AddContentDisposition("attachment", "document.pdf"))

// Serve in-memory content as a download.
resp.ServeFileAsDownload(w, "file.txt", []byte("Dynamic content"))
```

`ServeFileAsDownload` sets `Content-Length`, `X-Content-Type-Options: nosniff`,
and an RFC 5987 `filename*` form for non-ASCII names.

## Security headers

`SecurityHeaders()` applies a conservative, broadly-safe preset (`nosniff`,
`X-Frame-Options: DENY`, `Referrer-Policy: strict-origin-when-cross-origin`).
CSP is an explicit opt-in, because a useful policy is application-specific:

```go
resp.JSON(w, data,
    resp.SecurityHeaders(),
    resp.WithCSP("default-src 'self'"),
)
```

## Concurrency

Like `http.ResponseWriter`, a `resp` response is **not** safe for concurrent
use. Write from a single goroutine, and only until the handler returns. To
combine work from several goroutines, collect results over channels and write in
the handler's goroutine:

```go
select {
case err := <-errChan:
    resp.Error(w, http.StatusInternalServerError, err.Error())
case result := <-resultChan:
    resp.JSON(w, result)
case <-time.After(5 * time.Second):
    resp.Error(w, http.StatusGatewayTimeout, "Processing timeout")
}
```

## Errors

```go
var ErrInvalidJSONPCallback = errors.New("invalid JSONP callback name")
```

Returned by `JSONP` for a callback name that fails the identifier whitelist.
Match it with `errors.Is`. All helpers return the encode/write error (or `nil`)
so you can decide how to react.

## Recipes and tips

**Always check the error — it's safe to recover.** In the default buffered mode,
a failed `JSON` has written nothing, so a fallback `Error` response still works.

**Keep status and body in sync.** Use `Error` (status mirrored into `code`)
instead of hand-writing a status plus a separate body code; add `WithErrorCode`
only when the business code differs from the HTTP status.

**Reach for the builder when configuring a lot.** Chain `SetStatus`/`SetHeader`/
`SetCookie` on `NewResponse` rather than passing a long `...Option` list.

**Stream the big stuff.** Use `Stream` for large files and `WithStreamJSON()`
for very large JSON, trading the buffered safety net for lower memory.

**Hide sensitive fields.** Run `ExcludeFields`/`OnlyFields` before `JSON` so a
password or token never reaches the encoder.

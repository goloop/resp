[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/resp/v2)](https://goreportcard.com/report/github.com/goloop/resp/v2) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/resp/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://pkg.go.dev/github.com/goloop/resp/v2) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# resp

A lightweight Go package for crafting HTTP responses with ease. Built on top of Go's standard `net/http` package, `resp` simplifies common response patterns while maintaining flexibility and performance. It has **no third-party dependencies**.

## What's new in v2

Version 2 is a focused correctness and modernization release. The module path is now `github.com/goloop/resp/v2`.

- **Safe-by-default JSON.** `JSON`/`JSONP`/`XML` encode into a pooled buffer first, so a serialization error is returned **before** the status line is committed — the handler can still respond with a proper error. Opt into zero-copy streaming with `WithStreamJSON()` when memory matters more than this guarantee.
- **Hardened JSONP.** The callback name is validated against a strict JavaScript-identifier whitelist; an unsafe name returns `ErrInvalidJSONPCallback` and writes nothing, closing the reflected-XSS vector.
- **Clear `Error` semantics.** The first argument is the HTTP status, mirrored into the body `code` by default. Use `WithErrorCode(n)` to set a distinct business code in the body.
- **Strict redirects.** `Redirect` honours only genuine redirect codes (300, 301, 302, 303, 307, 308); anything else falls back to 302.
- **Safer downloads.** `ServeFileAsDownload` sets `Content-Length`, `X-Content-Type-Options: nosniff`, and an RFC 5987 `filename*` form for non-ASCII names.
- **New helpers.** `XML`, `Blob`, `Created`, `Accepted`, plus `(*Response).Status()` and `(*Response).Header()` getters.
- **Modern baseline.** `go 1.24`, `any` instead of `interface{}`, zero dependencies.

## Important Note

This package, like `http.ResponseWriter`, is NOT safe for concurrent use. The `ResponseWriter` should only be written to from a single goroutine, and only until the HTTP handler returns. Writing to `ResponseWriter` after the handler returns or from multiple goroutines may result in corrupted responses or other undefined behavior.

For safe concurrent operations, use channels to collect results and write to the ResponseWriter in the main goroutine:

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    resultChan := make(chan MyData)
    errChan := make(chan error)

    go func() {
        data, err := processData()
        if err != nil {
            errChan <- err
            return
        }
        resultChan <- data
    }()

    select {
    case err := <-errChan:
        resp.Error(w, http.StatusInternalServerError, err.Error())
    case result := <-resultChan:
        resp.JSON(w, result)
    case <-time.After(5 * time.Second):
        resp.Error(w, http.StatusGatewayTimeout, "Processing timeout")
    }
}
```

## Features

- 🚀 Quick and intuitive response methods (JSON, JSONP, XML, HTML, String, Blob, etc.)
- 🛡️ Safe-by-default encoding: errors never leave a half-written success response
- 🔒 Validated JSONP callbacks
- 🛠️ Flexible header and status code management
- 🍪 Simple cookie handling
- 📁 File streaming and downloads
- ⚙️ Custom JSON encoder support
- 🔧 Chainable response options
- 📦 Zero third-party dependencies

## Installation

```bash
go get -u github.com/goloop/resp/v2
```

## Basic Usage

```go
import "github.com/goloop/resp/v2"

// Simple JSON response.
func HandleJSON(w http.ResponseWriter, r *http.Request) {
    data := resp.R{
        "message": "Hello, World!",
        "status":  "success",
    }

    if err := resp.JSON(w, data); err != nil {
        log.Printf("Failed to send response: %v", err)
    }
}

// HTML response with status code.
func HandleHTML(w http.ResponseWriter, r *http.Request) {
    html := `<h1>Hello World</h1>`
    if err := resp.HTML(w, html, resp.WithStatus(http.StatusOK)); err != nil {
        log.Printf("Failed to send response: %v", err)
    }
}

// Error response: HTTP 404 with body {"code":404,"message":"Resource not found"}.
func HandleError(w http.ResponseWriter, r *http.Request) {
    resp.Error(w, http.StatusNotFound, "Resource not found")
}
```

## Error Responses

In v2 the first argument to `Error` is always the HTTP status code. It is written to the status line and mirrored into the JSON body's `code` field, so the two never disagree:

```go
// HTTP 404, body {"code":404,"message":"Not Found"}.
resp.Error(w, http.StatusNotFound, "")

// HTTP 404, body {"code":1042,"message":"User not found"}.
resp.Error(w, http.StatusNotFound, "User not found", resp.WithErrorCode(1042))
```

The status argument takes precedence over a `WithStatus` option, and an out-of-range value falls back to `500 Internal Server Error`.

## REST Helpers

```go
// 201 Created with a Location header and a JSON body.
resp.Created(w, "/users/42", user)

// 202 Accepted for asynchronous processing.
resp.Accepted(w, resp.R{"job_id": jobID})

// 204 No Content.
resp.NoContent(w)
```

## Binary and XML Responses

```go
// Raw bytes with an explicit content type (Content-Length set automatically).
resp.Blob(w, "application/pdf", pdfBytes)

// XML serialization, mirroring JSON.
resp.XML(w, note)
```

## Security Headers

`SecurityHeaders()` applies a conservative, broadly-safe preset; CSP stays an explicit opt-in because a useful policy is application-specific:

```go
resp.JSON(w, data,
    resp.SecurityHeaders(),               // nosniff, X-Frame-Options: DENY,
                                          // Referrer-Policy: strict-origin-when-cross-origin
    resp.WithCSP("default-src 'self'"))   // optional, explicit CSP
```

## JSONP

The callback name is validated before anything is written:

```go
callback := r.URL.Query().Get("callback")
if err := resp.JSONP(w, data, callback); err != nil {
    // errors.Is(err, resp.ErrInvalidJSONPCallback) for a rejected name.
    resp.Error(w, http.StatusBadRequest, "invalid callback")
}
```

## Advanced Features

### Custom JSON Encoder

You can plug in any encoder that satisfies `JSONEncodeFunc` (`func(w io.Writer, v any) error`):

```go
func HandleCustomJSON(w http.ResponseWriter, r *http.Request) {
    customEncoder := func(w io.Writer, v any) error {
        enc := json.NewEncoder(w)
        enc.SetEscapeHTML(false)
        return enc.Encode(v)
    }

    response := resp.NewResponse(w, resp.ApplyJSONEncoder(customEncoder))
    data := resp.R{"message": "Using custom encoder"}

    if err := response.JSON(data); err != nil {
        log.Printf("Encoding failed: %v", err)
    }
}
```

### Streaming Large JSON

By default `JSON` buffers the payload so an encode error can be reported before the status is committed. For very large responses where memory matters more than that guarantee, opt into direct streaming:

```go
resp.JSON(w, largeData, resp.WithStreamJSON())
```

### Pretty-Printed JSON

For human-facing or debug endpoints, `JSONIndent` formats the output:

```go
resp.JSONIndent(w, data, "", "  ")
```

### Field Projection

```go
// Keep only selected fields (e.g. hide a password before serialization).
data := resp.OnlyFields(user, "ID", "Email")

// Or drop selected fields.
data := resp.ExcludeFields(user, "Password")

resp.JSON(w, data)
```

### File Downloads

```go
func HandleFileDownload(w http.ResponseWriter, r *http.Request) {
    // Serve a static file.
    if err := resp.ServeFile(w, r, "path/to/file.pdf",
        resp.AddContentDisposition("attachment", "document.pdf")); err != nil {
        log.Printf("File serving failed: %v", err)
    }

    // Serve in-memory content as a download.
    data := []byte("Dynamic content")
    if err := resp.ServeFileAsDownload(w, "file.txt", data); err != nil {
        log.Printf("Download failed: %v", err)
    }
}
```

### Response Chaining

```go
func HandleChainedResponse(w http.ResponseWriter, r *http.Request) {
    response := resp.NewResponse(w).
        SetStatus(http.StatusOK).
        SetHeader("X-Custom-Header", "value")

    if err := response.JSON(resp.R{"status": "success"}); err != nil {
        log.Printf("Response failed: %v", err)
    }
}
```

### Header Management

```go
func HandleHeaders(w http.ResponseWriter, r *http.Request, data any) {
    opts := []resp.Option{
        resp.WithHeader("X-Custom-Header", "value"),
        resp.AddETag("\"123456\""),
        resp.AsApplicationJSON(),
    }

    if err := resp.JSON(w, data, opts...); err != nil {
        log.Printf("Response failed: %v", err)
    }
}
```

## Default Behaviors

- If no status code is set, `200 OK` is used.
- Content-Type headers are automatically set based on the response type:
  - JSON: `application/json; charset=utf-8`
  - JSONP: `application/javascript; charset=utf-8`
  - XML: `application/xml; charset=utf-8`
  - HTML: `text/html; charset=utf-8`
  - String: `text/plain`

## Error Handling

Always check for errors when sending responses:

```go
func Handler(w http.ResponseWriter, r *http.Request, data any) {
    if err := resp.JSON(w, data); err != nil {
        log.Printf("Failed to send response: %v", err)
        // Safe in buffered (default) mode: nothing has been written yet,
        // so a fallback error response still works.
        resp.Error(w, http.StatusInternalServerError, "Internal server error")
        return
    }
}
```

## Benchmarks

Benchmark results are machine-dependent and provided for illustration only:

```txt
goos: linux
goarch: arm64
pkg: github.com/goloop/resp/v2
BenchmarkJSONSmall-8               	 2383800	       531.6 ns/op	     280 B/op	       2 allocs/op
BenchmarkJSONMedium-8              	 1000000	      1152 ns/op	     608 B/op	       2 allocs/op
BenchmarkJSONLarge-8               	   41572	     68665 ns/op	   32359 B/op	       2 allocs/op
BenchmarkJSONPSmall-8              	 1827565	       680.9 ns/op	     372 B/op	       4 allocs/op
BenchmarkHTML-8                    	 5446101	      1248 ns/op	     639 B/op	       1 allocs/op
BenchmarkString-8                  	 7752152	       130.1 ns/op	     186 B/op	       1 allocs/op
BenchmarkError-8                   	 4737940	       354.2 ns/op	     185 B/op	       2 allocs/op
BenchmarkStream-8                  	12930627	       239.6 ns/op	     214 B/op	       1 allocs/op
BenchmarkRedirect-8                	11419272	       104.2 ns/op	      80 B/op	       3 allocs/op
BenchmarkNoContent-8               	22821877	        48.30 ns/op	      64 B/op	       2 allocs/op
BenchmarkJSONWithCustomEncoder-8   	  671696	      1781 ns/op	     871 B/op	       2 allocs/op
BenchmarkResponseChaining-8        	 3198278	       453.1 ns/op	     379 B/op	       3 allocs/op
BenchmarkServeFileAsDownload-8     	 3580976	       469.9 ns/op	     317 B/op	       6 allocs/op
PASS
ok  	github.com/goloop/resp/v2	29.021s
```

A pooled buffer keeps the default (safe) JSON path at the same 2 allocations per operation as direct streaming.

### Best Practices for Performance

1. Use the appropriate response type for your data:
   - `String()` for simple text responses.
   - `Stream()` for large files or real-time data.
   - `JSON()` with `WithStreamJSON()` or a custom encoder for very large payloads.

2. Consider response size:
   - Enable compression for large responses.
   - Use pagination for large datasets.
   - Consider streaming for very large data.

3. Memory efficiency:
   - Use `Stream()` for large files instead of loading them into memory.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

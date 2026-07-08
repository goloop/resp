![Go Report Card](https://img.shields.io/badge/go%20report-retired-lightgrey) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/resp/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://pkg.go.dev/github.com/goloop/resp/v2) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)

# resp

`resp` crafts HTTP responses on top of Go's standard `net/http`. Instead of
setting a content type, writing a status and encoding a body by hand, you call
one helper that does all three consistently — and reports an encoding error
before anything is committed. It has **no third-party dependencies**.

Use it two ways: package helpers (`resp.JSON(w, data, opts...)`) or a chainable
`Response` builder (`resp.NewResponse(w).SetStatus(...).SetHeader(...)`). Both
accept the same functional options.

## Features

- Quick response methods — JSON, JSONP, XML, HTML, String, Blob, Stream,
  Created/Accepted/NoContent, Redirect, file downloads.
- **Safe-by-default encoding** — a JSON/XML error never leaves a half-written
  success response (opt into streaming with `WithStreamJSON()`).
- **Validated JSONP** — callback names checked against a strict whitelist.
- Flexible headers and status via composable options; simple cookie handling.
- Field projection (`OnlyFields`/`ExcludeFields`), custom JSON encoders, and a
  conservative `SecurityHeaders()` preset.

## Installation

```bash
go get -u github.com/goloop/resp/v2
```

```go
import "github.com/goloop/resp/v2"
```

Requires Go 1.24 or newer.

## Quick start

```go
// Simple JSON.
func HandleJSON(w http.ResponseWriter, r *http.Request) {
    if err := resp.JSON(w, resp.R{"message": "Hello, World!", "status": "success"}); err != nil {
        log.Printf("response failed: %v", err)
    }
}

// Error: HTTP 404, body {"code":404,"message":"Resource not found"}.
func HandleError(w http.ResponseWriter, r *http.Request) {
    resp.Error(w, http.StatusNotFound, "Resource not found")
}

// Chainable builder with headers and status.
func HandleChained(w http.ResponseWriter, r *http.Request) {
    resp.NewResponse(w).
        SetStatus(http.StatusOK).
        SetHeader("X-Custom-Header", "value").
        JSON(resp.R{"status": "success"})
}
```

Because JSON is buffered by default, a failed encode has written nothing — a
fallback error response still works:

```go
if err := resp.JSON(w, data); err != nil {
    resp.Error(w, http.StatusInternalServerError, "Internal server error")
    return
}
```

> Like `http.ResponseWriter`, a `resp` response is **not** safe for concurrent
> use — write from a single goroutine until the handler returns.

## Documentation

- Full reference and recipes: [DOC.md](DOC.md) · [DOC.UK.md](DOC.UK.md)
- Package API: [pkg.go.dev/github.com/goloop/resp/v2](https://pkg.go.dev/github.com/goloop/resp/v2)
- Changes between versions: [CHANGELOG.md](CHANGELOG.md)

## Contributing

Contributions are welcome. Please run `go test ./...`, `go vet ./...` and
`gofmt -l .` before submitting a pull request.

## License

`resp` is released under the MIT License. See [LICENSE](LICENSE).

[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/resp)](https://goreportcard.com/report/github.com/goloop/resp) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/resp/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/resp) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# resp

A lightweight Go package for crafting HTTP responses with ease. Built on top of Go's standard `net/http` package, `resp` simplifies common response patterns while maintaining flexibility and performance.

## ⚠️ Important Note

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

- 🚀 Quick and intuitive response methods (JSON, HTML, String, etc.)
- 🛠️ Flexible header and status code management
- 🍪 Simple cookie handling
- 📁 File streaming and downloads
- ⚙️ Custom JSON encoder support
- 🔧 Chainable response options

## Installation

```bash
go get -u github.com/goloop/resp
```

## Basic Usage

```go
import "github.com/goloop/resp"

// Simple JSON response
func HandleJSON(w http.ResponseWriter, r *http.Request) {
    data := resp.R{
        "message": "Hello, World!",
        "status": "success",
    }

    if err := resp.JSON(w, data); err != nil {
        log.Printf("Failed to send response: %v", err)
    }
}

// HTML response with status code
func HandleHTML(w http.ResponseWriter, r *http.Request) {
    html := `<h1>Hello World</h1>`
    if err := resp.HTML(w, html, resp.WithStatus(http.StatusOK)); err != nil {
        log.Printf("Failed to send response: %v", err)
    }
}

// Error response
func HandleError(w http.ResponseWriter, r *http.Request) {
    resp.Error(w, http.StatusNotFound, "Resource not found")
}
```

## Advanced Features

### Custom JSON Encoder

You can use any JSON encoder that satisfies the `JSONEncodeFunc` interface:

```go
import jsoniter "github.com/json-iterator/go"

func HandleCustomJSON(w http.ResponseWriter, r *http.Request) {
    customEncoder := func(w io.Writer, v interface{}) error {
        return jsoniter.NewEncoder(w).Encode(v)
    }

    response := resp.NewResponse(w, resp.ApplyJSONEncoder(customEncoder))
    data := resp.R{"message": "Using custom encoder"}

    if err := response.JSON(data); err != nil {
        log.Printf("Encoding failed: %v", err)
    }
}
```

### File Downloads

```go
func HandleFileDownload(w http.ResponseWriter, r *http.Request) {
    // Serve static file
    if err := resp.ServeFile(w, r, "path/to/file.pdf",
        resp.AddContentDisposition("attachment", "document.pdf")); err != nil {
        log.Printf("File serving failed: %v", err)
    }

    // Serve dynamic content as file
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
func HandleHeaders(w http.ResponseWriter, r *http.Request) {
    opts := []resp.Option{
        resp.WithHeader("X-Custom-Header", "value"),
        resp.AddContentType("application/json"),
        resp.AddETag("\"123456\""),
        resp.AsApplicationJSON(),
    }

    if err := resp.JSON(w, data, opts...); err != nil {
        log.Printf("Response failed: %v", err)
    }
}
```

## Default Behaviors

- If no status code is set, `200 OK` is used
- Content-Type headers are automatically set based on the response type:
  - JSON: `application/json; charset=utf-8`
  - HTML: `text/html; charset=utf-8`
  - String: `text/plain; charset=utf-8`

## Error Handling

Always check for errors when sending responses:

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    if err := resp.JSON(w, data); err != nil {
        log.Printf("Failed to send response: %v", err)
        // Optionally send a fallback error response
        resp.Error(w, http.StatusInternalServerError, "Internal server error")
        return
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

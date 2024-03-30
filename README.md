[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/resp)](https://goreportcard.com/report/github.com/goloop/resp) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/resp/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/resp) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# resp

The resp package offers a streamlined interface for crafting a variety of HTTP responses in your Go web applications and APIs. Leveraging the power of Go's standard net/http library, resp adds a layer of convenience that simplifies common tasks such as sending JSON and HTML responses, managing cookies, handling file downloads, and streaming content.

## Features

 - Simplified Response Handling: Easily send JSON, HTML, and plain text responses with minimal code.
 - Flexible Header and Status Code Management: Set custom HTTP headers and status codes effortlessly.
 - Cookie Management: Intuitive methods for setting, deleting, and expiring cookies.
 - File Downloads and Streaming: Tools for sending files as downloads or streaming content directly to the client.
 - Seamless Integration: Works out of the box with Go's HTTP server and middleware ecosystem, requiring no additional setup.

## Getting Started

To begin using resp, install the package using `go get`:

```
$ go get github.com/goloop/resp
```

Then, import it into your project:

```
import (
	"github.com/goloop/resp"
)
```

### Example Usage

Below is a quick example to get you started:

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goloop/resp"
)

// Item represents a single some item.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// items is a slice of Item.
// We use it as a test database.
var items []Item

// GetItems sends all items as a JSON response.
func GetItems(w http.ResponseWriter, r *http.Request) {
	// Send the items as a `application/json; charset=utf-8` response
	// with a 200 status code.
	resp.JSON(w, items)
}

// GetItem finds and sends a single item by ID as a JSON response.
func GetItem(w http.ResponseWriter, r *http.Request) {
	// Try to get the id from the URL query.
	id := r.PathValue("id")

	// Find the item by ID.
	for _, item := range items {
		if item.ID == id {
			// Extended response.
			// Send the item as a `application/json` response
			// with a 200 status code.
			resp.JSON(w, item, resp.AsApplicationJSON(), resp.WithStatusOK())
			return
		}
	}

	// The error message is taken from the status code.
	resp.Error(w, http.StatusNotFound)
}

// CreateItem adds a new item from the request body to the items slice.
func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		//  Sets the custom error message.
		resp.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	items = append(items, item)
	resp.JSON(w, items)
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /items/{id}", GetItem)
	router.HandleFunc("GET /items", GetItems)
	router.HandleFunc("POST /item", CreateItem)

	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
```

This code snippet demonstrates basic CRUD operations for an in-memory list of items, showcasing how the resp package can be used to simplify response handling in your Go web applications.

Test usage:

```
$ curl -X POST -H "Content-Type: application/json" -d "{\"id\": \"10\", \"name\": \"john\"}" http://localhost:8000/item
[{"id":"10","name":"john"}]

$ curl -X POST -H "Content-Type: application/json" -d "{\"id\": \"11\", \"name\": \"bob\"}" http://localhost:8000/item
[{"id":"10","name":"john"},{"id":"11","name":"bob"}]

$ curl localhost:8000/items
[{"id":"10","name":"john"},{"id":"11","name":"bob"}]

$ curl localhost:8000/items/11
{"id":"11","name":"bob"}
```

## Hot features

Quick call functions.

 - **JSON** - function initializes a new HTTP response with optional settings, sets the content type to `application/json; charset=utf-8`, and encodes the given data into JSON format for sending to the client.
 - **JSONP** - function encodes the provided data as JSON, wraps it in a client-specified callback for cross-domain requests, and sends this JSONP response to the client.
 - **String** - function sends a text response, such as plain text, to the client with optional configurations for headers and status codes, simplifying the delivery of text-based content.
 - **Error** - function sends an HTTP error response with a customizable status code and message, offering flexibility in error reporting by allowing either a status code or a message as the primary argument, with optional detailed status codes and messages for enhanced context.
 - **Stream** - function enables the delivery of streaming content, like file downloads or live video, directly to the client by reading from an io.Reader and writing to the http.ResponseWriter, with support for custom headers and status codes to tailor the response as needed.
 - **ServeFile** - function streamlines the delivery of static files to the client, utilizing http.ServeFile for automatic content type detection, range request handling, and caching, while offering the flexibility to set custom response headers and status codes.
 - **ServeFileAsDownload** - function facilitates serving in-memory data or dynamically generated content as a file download by setting the Content-Disposition header, with the capability to customize the response through various options for headers, status codes, and other settings.
 - **Redirect** - function orchestrates HTTP redirects, setting the status to 302 Found by default, and allows for flexible redirection to a specified URL with options for customizing the response, such as changing the redirect status code.
 - **NoContent** - function conveys a 204 No Content response to clients, ideal for operations that require no return data, with the flexibility to modify the response via optional configurations for headers or other settings.
 - **HTML** - function facilitates the delivery of HTML content to the client, defaulting to a "text/html" Content-Type, and offers customization through optional configurations for headers and status codes, ideal for serving web pages or HTML snippets.

## Types

The provided types allow us to manage responses more conveniently.

 - **R** is a type alias for `map[string]interface{}` in Go, designed to facilitate the construction and manipulation of JSON objects for HTTP responses. It offers a flexible and concise way to handle dynamic data structures, reducing boilerplate and enhancing code readability in web applications and APIs.
 - **Response** - type encapsulates an HTTP response, offering methods to manipulate headers, cookies, and the response body for streamlined server-side communication. Utilizing `NewResponse` for its creation ensures that any provided options, such as custom headers or status codes, are applied correctly from the start, establishing a robust foundation for handling HTTP responses with enhanced flexibility and control.
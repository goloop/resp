// Package resp is a Go library designed to simplify the creation and management
// of HTTP responses within web applications and APIs.
//
// The resp package provides a set of utilities and helper functions that
// streamline the process of sending various types of responses, including
// JSON, HTML, String, and more. It is built with ease of use and flexibility
// in mind, allowing developers to efficiently implement common response
// patterns while maintaining readability and reducing boilerplate code.
//
// Features:
//
//   - Easy-to-use functions for sending JSON, JSONP, HTML, and text responses.
//   - Support for setting custom HTTP headers and status codes.
//   - Convenient methods for managing cookies, including setting, deleting,
//     and marking cookies as expired.
//   - Utilities for handling file downloads and streaming responses.
//   - Integration with standard net/http interfaces, making it compatible
//     with Go's HTTP server and middleware ecosystems.
//
// Example Usage:
//
// To send a JSON response with the resp package, you can use the JSON
// function like so:
//
//	import (
//		"net/http"
//		"github.com/goloop/resp"
//	)
//
//	func JSONHandler(w http.ResponseWriter, r *http.Request) {
//		data := resp.R{
//			"hello": "world",
//		}
//		if err := resp.JSON(w, data); err != nil {
//			// handle error
//		}
//	}
//
// This is a simple example showing how to send an HTML response:
//
//	func HTMLHandler(w http.ResponseWriter, r *http.Request) {
//		content := "<!DOCTYPE html><html><body>Hello, World!</body></html>"
//		if err := resp.HTML(w, content); err != nil {
//			// handle error
//		}
//	}
package resp

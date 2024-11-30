package resp

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// Test data structures
type testUser struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
	IsActive bool      `json:"is_active"`
}

var (
	testUserData = testUser{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Created:  time.Now(),
		IsActive: true,
	}

	testLargeJSON = generateLargeJSON()
)

// Helper functions
func generateLargeJSON() []testUser {
	users := make([]testUser, 1000)
	for i := range users {
		users[i] = testUser{
			ID:       i + 1,
			Name:     "User " + string(rune(i)),
			Email:    "user" + string(rune(i)) + "@example.com",
			Created:  time.Now(),
			IsActive: i%2 == 0,
		}
	}
	return users
}

func newTestResponseWriter() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// // Benchmarks for JSON responses
// func BenchmarkJSON(b *testing.B) {
// 	w := newTestResponseWriter()

// 	b.Run("SmallJSON", func(b *testing.B) {
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			JSON(w, testUserData)
// 		}
// 	})

// 	b.Run("LargeJSON", func(b *testing.B) {
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			StreamJSON(w, testLargeJSON)
// 		}
// 	})

// 	b.Run("JSONWithHeaders", func(b *testing.B) {
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			JSON(w, testUserData,
// 				WithHeader("X-Custom", "value"),
// 				WithHeader("Cache-Control", "no-cache"),
// 			)
// 		}
// 	})
// }

// Benchmarks for String responses
func BenchmarkString(b *testing.B) {
	w := newTestResponseWriter()
	smallStr := "Hello, World!"
	largeStr := strings.Repeat("Lorem ipsum dolor sit amet ", 1000)

	b.Run("SmallString", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			String(w, smallStr)
		}
	})

	b.Run("LargeString", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			String(w, largeStr)
		}
	})
}

// Benchmarks for Error responses
func BenchmarkError(b *testing.B) {
	w := newTestResponseWriter()

	b.Run("SimpleError", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Error(w, 400, "Bad Request")
		}
	})

	b.Run("ErrorWithCustomStatus", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Error(w, 400, "Bad Request", WithStatus(http.StatusBadRequest))
		}
	})
}

// Benchmarks for Stream responses
func BenchmarkStream(b *testing.B) {
	w := newTestResponseWriter()
	smallData := bytes.NewReader([]byte("Hello, World!"))
	largeData := bytes.NewReader([]byte(strings.Repeat("Lorem ipsum dolor sit amet ", 1000)))

	b.Run("SmallStream", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Stream(w, smallData)
			smallData.Seek(0, io.SeekStart)
		}
	})

	b.Run("LargeStream", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Stream(w, largeData)
			largeData.Seek(0, io.SeekStart)
		}
	})
}

// Benchmarks for HTML responses
func BenchmarkHTML(b *testing.B) {
	w := newTestResponseWriter()
	smallHTML := "<h1>Hello, World!</h1>"
	largeHTML := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
		</head>
		<body>
			<div>` + strings.Repeat("<p>Lorem ipsum dolor sit amet</p>", 100) + `</div>
		</body>
		</html>
	`

	b.Run("SmallHTML", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			HTML(w, smallHTML)
		}
	})

	b.Run("LargeHTML", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			HTML(w, largeHTML)
		}
	})
}

// Benchmarks for Cookie operations
func BenchmarkCookie(b *testing.B) {
	w := newTestResponseWriter()
	cookie := &http.Cookie{
		Name:     "test",
		Value:    "value",
		Path:     "/",
		Domain:   "example.com",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}

	b.Run("SetCookie", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			NewResponse(w).SetCookie(cookie)
		}
	})

	b.Run("DeleteCookie", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			NewResponse(w).DelCookie("test")
		}
	})
}

// Benchmarks for Header operations
func BenchmarkHeaders(b *testing.B) {
	w := newTestResponseWriter()

	b.Run("SetSingleHeader", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			NewResponse(w).SetHeader("X-Custom", "value")
		}
	})

	b.Run("SetMultipleHeaders", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			NewResponse(w,
				WithHeader("X-Custom1", "value1"),
				WithHeader("X-Custom2", "value2"),
				WithHeader("X-Custom3", "value3"),
			)
		}
	})
}

// Benchmarks for Response creation
func BenchmarkResponseCreation(b *testing.B) {
	w := newTestResponseWriter()

	b.Run("SimpleResponse", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			NewResponse(w)
		}
	})

	b.Run("ResponseWithOptions", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			NewResponse(w,
				WithStatus(http.StatusOK),
				WithHeader("Content-Type", "application/json"),
				WithHeader("X-Custom", "value"),
			)
		}
	})
}

// // Benchmarks for JSONP responses
// func BenchmarkJSONP(b *testing.B) {
// 	w := newTestResponseWriter()

// 	b.Run("SmallJSONP", func(b *testing.B) {
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			JSONP(w, testUserData, "callback")
// 		}
// 	})

// 	b.Run("LargeJSONP", func(b *testing.B) {
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			JSONP(w, testLargeJSON, "callback")
// 		}
// 	})
// }

func BenchmarkJSONEncoders(b *testing.B) {
	w := newTestResponseWriter()
	jsoniterAPI := jsoniter.ConfigCompatibleWithStandardLibrary
	jsoniterFast := jsoniter.ConfigFastest

	// Тести для малого об'єкта
	b.Run("SmallJSON/stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			JSON(w, testUserData)
		}
	})

	b.Run("SmallJSON/jsoniter-compat", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
			jsoniterAPI.NewEncoder(w).Encode(testUserData)
		}
	})

	b.Run("SmallJSON/jsoniter-fast", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
			jsoniterFast.NewEncoder(w).Encode(testUserData)
		}
	})

	// Тести для великого об'єкта
	b.Run("LargeJSON/stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			JSON(w, testLargeJSON)
		}
	})

	b.Run("LargeJSON/jsoniter-compat", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
			jsoniterAPI.NewEncoder(w).Encode(testLargeJSON)
		}
	})

	b.Run("LargeJSON/jsoniter-fast", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
			jsoniterFast.NewEncoder(w).Encode(testLargeJSON)
		}
	})

	// Додамо тест з дуже великим об'єктом
	veryLargeJSON := make([]testUser, 10000)
	copy(veryLargeJSON, testLargeJSON)

	b.Run("VeryLargeJSON/stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			JSON(w, veryLargeJSON)
		}
	})

	b.Run("VeryLargeJSON/jsoniter-compat", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
			jsoniterAPI.NewEncoder(w).Encode(veryLargeJSON)
		}
	})

	b.Run("VeryLargeJSON/jsoniter-fast", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
			jsoniterFast.NewEncoder(w).Encode(veryLargeJSON)
		}
	})
}

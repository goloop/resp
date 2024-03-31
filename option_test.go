package resp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// TestWithHeader tests the WithHeader function.
func TestWithHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	// Call the WithHeader function.
	option := WithHeader("Content-Type", "application/json")
	option(r)

	// Check that the header is set correctly.
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("WithHeader() did not set the header correctly, got %v", got)
	}
}

// TestWithHeader_MultipleValues tests the WithHeader
// function with multiple values.
func TestWithHeader_MultipleValues(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	// Call the WithHeader function with multiple values.
	option := WithHeader("Accept-Encoding", "gzip", "deflate")
	option(r)

	// Check that the header is set correctly.
	values := w.Header().Values(HeaderAcceptEncoding)
	if len(values) != 2 || values[0] != "gzip" || values[1] != "deflate" {
		t.Errorf("WithHeader() did not set the multiple values "+
			"correctly, got %v", values)
	}
}

// TestWithHeader_EmptyValues tests the WithHeader function with empty values.
func TestWithHeader_EmptyValues(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	// Call the WithHeader function with empty values.
	option := WithHeader("X-Custom-Header")
	option(r)

	// Check that the header is set correctly.
	if got := w.Header().Get("X-Custom-Header"); got != "" {
		t.Errorf("WithHeader() did not set the header correctly, got %v", got)
	}
}

// TestWithStatus tests the WithStatus function.
func TestWithStatus(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	// Call the WithStatus function.
	option := WithStatus(200)
	option(r)

	// Check that the status code is set correctly.
	if got := w.Code; got != 200 {
		t.Errorf("WithStatus() did not set the status code correctly, got %v", got)
	}
}

// TestWithCookie tests the WithCookie function.
func TestWithCookie(t *testing.T) {
	// Create a test cookie.
	testCookie := &http.Cookie{
		Name:  "test",
		Value: "cookie_value",
	}

	// Create a new recorder.
	w := httptest.NewRecorder()
	_ = NewResponse(w, WithCookie(testCookie))

	result := w.Result()
	defer result.Body.Close()

	cookies := result.Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == testCookie.Name && cookie.Value == testCookie.Value {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("WithCookie() did not set the expected cookie: %v", testCookie)
	}
}

// TestWithStatusContinue tests the WithStatusContinue function.
func TestWithStatusContinue(t *testing.T) {
	// Create a new HTTP response recorder.
	// Apply WithStatusContinue option.
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusContinue())

	// Dummy call to trigger response writing.
	resp.httpWriter.WriteHeader(resp.statusCode)

	// Check if status code is set to 100.
	if w.Code != StatusContinue {
		t.Errorf("WithStatusContinue() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusContinue)
	}
}

// TestWithStatusSwitchingProtocols tests the
// WithStatusSwitchingProtocols function.
func TestWithStatusSwitchingProtocols(t *testing.T) {
	// Create a new HTTP response recorder.
	// Apply WithStatusSwitchingProtocols option.
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusSwitchingProtocols())

	// Dummy call to trigger response writing.
	resp.httpWriter.WriteHeader(resp.statusCode)

	// Check if status code is set to 101.
	if w.Code != StatusSwitchingProtocols {
		t.Errorf("WithStatusSwitchingProtocols() did not set "+
			"the correct status code: got %v, want %v",
			w.Code, StatusSwitchingProtocols)
	}
}

// TestWithStatusProcessing tests the WithStatusProcessing function.
func TestWithStatusProcessing(t *testing.T) {
	// Create a new HTTP response recorder.
	w := httptest.NewRecorder()
	// Apply WithStatusProcessing option.
	resp := NewResponse(w, WithStatusProcessing())

	// Dummy call to trigger response writing.
	resp.httpWriter.WriteHeader(resp.statusCode)

	// Check if status code is set to 102.
	if w.Code != StatusProcessing {
		t.Errorf("WithStatusProcessing() did not set "+
			"the correct status code: got %v, want %v",
			w.Code, StatusProcessing)
	}
}

// TestWithStatusEarlyHints tests the WithStatusEarlyHints function.
func TestWithStatusEarlyHints(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusEarlyHints())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusEarlyHints {
		t.Errorf("WithStatusEarlyHints() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusEarlyHints)
	}
}

// TestWithStatusOK tests the WithStatusOK function.
func TestWithStatusOK(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusOK {
		t.Errorf("WithStatusOK() did not set the correct status code: "+
			"got %v, want %v", w.Code, StatusOK)
	}
}

// TestWithStatusCreated tests the WithStatusCreated function.
func TestWithStatusCreated(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusCreated())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusCreated {
		t.Errorf("WithStatusCreated() did not set the correct status code: "+
			"got %v, want %v", w.Code, StatusCreated)
	}
}

// TestWithStatusAccepted tests the WithStatusAccepted function.
func TestWithStatusAccepted(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusAccepted())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusAccepted {
		t.Errorf("WithStatusAccepted() did not set the correct status code: "+
			"got %v, want %v", w.Code, StatusAccepted)
	}
}

// TestWithStatusNonAuthoritativeInfo tests the
// WithStatusNonAuthoritativeInfo function.
func TestWithStatusNonAuthoritativeInfo(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusNonAuthoritativeInfo())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusNonAuthoritativeInfo {
		t.Errorf("WithStatusNonAuthoritativeInfo() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusNonAuthoritativeInfo)
	}
}

// TestWithStatusNoContent tests the WithStatusNoContent function.
func TestWithStatusNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusNoContent())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusNoContent {
		t.Errorf("WithStatusNoContent() did not set the correct status code: "+
			"got %v, want %v", w.Code, StatusNoContent)
	}
}

// TestWithStatusResetContent tests the WithStatusResetContent function.
func TestWithStatusResetContent(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusResetContent())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusResetContent {
		t.Errorf("WithStatusResetContent() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusResetContent)
	}
}

// TestWithStatusPartialContent tests the WithStatusPartialContent function.
func TestWithStatusPartialContent(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusPartialContent())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusPartialContent {
		t.Errorf("WithStatusPartialContent() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusPartialContent)
	}
}

// TestWithStatusMultiStatus tests the WithStatusMultiStatus function.
func TestWithStatusMultiStatus(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusMultiStatus())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusMultiStatus {
		t.Errorf("WithStatusMultiStatus() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusMultiStatus)
	}
}

// TestWithStatusAlreadyReported tests the WithStatusAlreadyReported function.
func TestWithStatusAlreadyReported(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusAlreadyReported())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusAlreadyReported {
		t.Errorf("WithStatusAlreadyReported() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusAlreadyReported)
	}
}

// TestWithStatusIMUsed tests the WithStatusIMUsed function.
func TestWithStatusIMUsed(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusIMUsed())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusIMUsed {
		t.Errorf("WithStatusIMUsed() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusIMUsed)
	}
}

// TestWithStatusMultipleChoices tests the WithStatusMultipleChoices function.
func TestWithStatusMultipleChoices(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusMultipleChoices())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusMultipleChoices {
		t.Errorf("WithStatusMultipleChoices() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusMultipleChoices)
	}
}

// TestWithStatusMovedPermanently tests the WithStatusMovedPermanently function.
func TestWithStatusMovedPermanently(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusMovedPermanently())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusMovedPermanently {
		t.Errorf("WithStatusMovedPermanently() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusMovedPermanently)
	}
}

// TestWithStatusFound tests the WithStatusFound function.
func TestWithStatusFound(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusFound())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusFound {
		t.Errorf("WithStatusFound() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusFound)
	}
}

// TestWithStatusSeeOther tests the WithStatusSeeOther function.
func TestWithStatusSeeOther(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusSeeOther())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusSeeOther {
		t.Errorf("WithStatusSeeOther() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusSeeOther)
	}
}

// TestWithStatusNotModified tests the WithStatusNotModified function.
func TestWithStatusNotModified(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusNotModified())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusNotModified {
		t.Errorf("WithStatusNotModified() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusNotModified)
	}
}

// TestWithStatusUseProxy tests the WithStatusUseProxy function.
func TestWithStatusUseProxy(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusUseProxy())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusUseProxy {
		t.Errorf("WithStatusUseProxy() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusUseProxy)
	}
}

// TestWithStatusTemporaryRedirect tests the
// WithStatusTemporaryRedirect function.
func TestWithStatusTemporaryRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusTemporaryRedirect())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusTemporaryRedirect {
		t.Errorf("WithStatusTemporaryRedirect() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusTemporaryRedirect)
	}
}

// TestWithStatusPermanentRedirect tests the
// WithStatusPermanentRedirect function.
func TestWithStatusPermanentRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusPermanentRedirect())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusPermanentRedirect {
		t.Errorf("WithStatusPermanentRedirect() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusPermanentRedirect)
	}
}

// TestWithStatusBadRequest tests the WithStatusBadRequest function.
func TestWithStatusBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusBadRequest())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusBadRequest {
		t.Errorf("WithStatusBadRequest() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusBadRequest)
	}
}

// TestWithStatusUnauthorized tests the WithStatusUnauthorized function.
func TestWithStatusUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusUnauthorized())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusUnauthorized {
		t.Errorf("WithStatusUnauthorized() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusUnauthorized)
	}
}

// TestWithStatusPaymentRequired tests the WithStatusPaymentRequired function.
func TestWithStatusPaymentRequired(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusPaymentRequired())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusPaymentRequired {
		t.Errorf("WithStatusPaymentRequired() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusPaymentRequired)
	}
}

// TestWithStatusForbidden tests the WithStatusForbidden function.
func TestWithStatusForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusForbidden())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusForbidden {
		t.Errorf("WithStatusForbidden() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusForbidden)
	}
}

// TestWithStatusNotFound tests the WithStatusNotFound function.
func TestWithStatusNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusNotFound())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusNotFound {
		t.Errorf("WithStatusNotFound() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusNotFound)
	}
}

// TestWithStatusMethodNotAllowed tests the WithStatusMethodNotAllowed function.
func TestWithStatusMethodNotAllowed(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusMethodNotAllowed())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusMethodNotAllowed {
		t.Errorf("WithStatusMethodNotAllowed() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusMethodNotAllowed)
	}
}

// TestWithStatusNotAcceptable tests the WithStatusNotAcceptable function.
func TestWithStatusNotAcceptable(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusNotAcceptable())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusNotAcceptable {
		t.Errorf("WithStatusNotAcceptable() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusNotAcceptable)
	}
}

// TestWithStatusProxyAuthRequired tests the WithStatusProxyAuthRequired function.
func TestWithStatusProxyAuthRequired(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusProxyAuthRequired())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusProxyAuthRequired {
		t.Errorf("WithStatusProxyAuthRequired() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusProxyAuthRequired)
	}
}

// TestWithStatusRequestTimeout tests the WithStatusRequestTimeout function.
func TestWithStatusRequestTimeout(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusRequestTimeout())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusRequestTimeout {
		t.Errorf("WithStatusRequestTimeout() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusRequestTimeout)
	}
}

// TestWithStatusConflict tests the WithStatusConflict function.
func TestWithStatusConflict(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusConflict())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusConflict {
		t.Errorf("WithStatusConflict() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusConflict)
	}
}

// TestWithStatusGone tests the WithStatusGone function.
func TestWithStatusGone(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusGone())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusGone {
		t.Errorf("WithStatusGone() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusGone)
	}
}

// TestWithStatusLengthRequired tests the WithStatusLengthRequired function.
func TestWithStatusLengthRequired(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusLengthRequired())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusLengthRequired {
		t.Errorf("WithStatusLengthRequired() did not set the correct "+
			"status code: got %v, want %v", w.Code, StatusLengthRequired)
	}
}

// TestWithStatusPreconditionFailed tests the
// WithStatusPreconditionFailed function.
func TestWithStatusPreconditionFailed(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusPreconditionFailed())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusPreconditionFailed {
		t.Errorf("WithStatusPreconditionFailed() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusPreconditionFailed)
	}
}

// TestWithStatusRequestEntityTooLarge tests the
// WithStatusRequestEntityTooLarge function.
func TestWithStatusRequestEntityTooLarge(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusRequestEntityTooLarge())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusRequestEntityTooLarge {
		t.Errorf("WithStatusRequestEntityTooLarge() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusRequestEntityTooLarge)
	}
}

// TestWithStatusRequestURITooLong tests the
// WithStatusRequestURITooLong function.
func TestWithStatusRequestURITooLong(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusRequestURITooLong())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusRequestURITooLong {
		t.Errorf("WithStatusRequestURITooLong() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusRequestURITooLong)
	}
}

// TestWithStatusUnsupportedMediaType tests the
// WithStatusUnsupportedMediaType function.
func TestWithStatusUnsupportedMediaType(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusUnsupportedMediaType())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusUnsupportedMediaType {
		t.Errorf("WithStatusUnsupportedMediaType() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusUnsupportedMediaType)
	}
}

// TestWithStatusRequestedRangeNotSatisfiable tests the
// WithStatusRequestedRangeNotSatisfiable function.
func TestWithStatusRequestedRangeNotSatisfiable(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusRequestedRangeNotSatisfiable())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusRequestedRangeNotSatisfiable {
		t.Errorf("WithStatusRequestedRangeNotSatisfiable() did "+
			"not set the correct status code: got %v, want %v",
			w.Code, StatusRequestedRangeNotSatisfiable)
	}
}

// TestWithStatusExpectationFailed tests the
// WithStatusExpectationFailed function.
func TestWithStatusExpectationFailed(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusExpectationFailed())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusExpectationFailed {
		t.Errorf("WithStatusExpectationFailed() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusExpectationFailed)
	}
}

// TestWithStatusTeapot tests the WithStatusTeapot function.
func TestWithStatusTeapot(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusTeapot())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusTeapot {
		t.Errorf("WithStatusTeapot() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusTeapot)
	}
}

// TestWithStatusMisdirectedRequest tests the
// WithStatusMisdirectedRequest function.
func TestWithStatusMisdirectedRequest(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusMisdirectedRequest())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusMisdirectedRequest {
		t.Errorf("WithStatusMisdirectedRequest() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusMisdirectedRequest)
	}
}

// TestWithStatusUnprocessableEntity tests the
// WithStatusUnprocessableEntity function.
func TestWithStatusUnprocessableEntity(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusUnprocessableEntity())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusUnprocessableEntity {
		t.Errorf("WithStatusUnprocessableEntity() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusUnprocessableEntity)
	}
}

// TestWithStatusLocked tests the WithStatusLocked function.
func TestWithStatusLocked(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusLocked())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusLocked {
		t.Errorf("WithStatusLocked() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusLocked)
	}
}

// TestWithStatusFailedDependency tests the
// WithStatusFailedDependency function.
func TestWithStatusFailedDependency(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusFailedDependency())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusFailedDependency {
		t.Errorf("WithStatusFailedDependency() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusFailedDependency)
	}
}

// TestWithStatusTooEarly tests the WithStatusTooEarly function.
func TestWithStatusTooEarly(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusTooEarly())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusTooEarly {
		t.Errorf("WithStatusTooEarly() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusTooEarly)
	}
}

// TestWithStatusUpgradeRequired tests the WithStatusUpgradeRequired function.
func TestWithStatusUpgradeRequired(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusUpgradeRequired())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusUpgradeRequired {
		t.Errorf("WithStatusUpgradeRequired() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusUpgradeRequired)
	}
}

// TestWithStatusPreconditionRequired tests the
// WithStatusPreconditionRequired function.
func TestWithStatusPreconditionRequired(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusPreconditionRequired())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusPreconditionRequired {
		t.Errorf("WithStatusPreconditionRequired() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusPreconditionRequired)
	}
}

// TestWithStatusTooManyRequests tests the WithStatusTooManyRequests function.
func TestWithStatusTooManyRequests(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusTooManyRequests())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusTooManyRequests {
		t.Errorf("WithStatusTooManyRequests() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusTooManyRequests)
	}
}

// TestWithStatusRequestHeaderFieldsTooLarge tests the
// WithStatusRequestHeaderFieldsTooLarge function.
func TestWithStatusRequestHeaderFieldsTooLarge(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusRequestHeaderFieldsTooLarge())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusRequestHeaderFieldsTooLarge {
		t.Errorf("WithStatusRequestHeaderFieldsTooLarge() "+
			"did not set the correct status code: got %v, want %v",
			w.Code, StatusRequestHeaderFieldsTooLarge)
	}
}

// TestWithStatusUnavailableForLegalReasons tests the
// WithStatusUnavailableForLegalReasons function.
func TestWithStatusUnavailableForLegalReasons(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusUnavailableForLegalReasons())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusUnavailableForLegalReasons {
		t.Errorf("WithStatusUnavailableForLegalReasons() "+
			"did not set the correct status code: got %v, want %v",
			w.Code, StatusUnavailableForLegalReasons)
	}
}

// TestWithStatusInternalServerError tests the
// WithStatusInternalServerError function.
func TestWithStatusInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusInternalServerError())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusInternalServerError {
		t.Errorf("WithStatusInternalServerError() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusInternalServerError)
	}
}

// TestWithStatusNotImplemented tests the WithStatusNotImplemented function.
func TestWithStatusNotImplemented(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusNotImplemented())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusNotImplemented {
		t.Errorf("WithStatusNotImplemented() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusNotImplemented)
	}
}

// TestWithStatusBadGateway tests the WithStatusBadGateway function.
func TestWithStatusBadGateway(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusBadGateway())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusBadGateway {
		t.Errorf("WithStatusBadGateway() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusBadGateway)
	}
}

// TestWithStatusServiceUnavailable tests the
// WithStatusServiceUnavailable function.
func TestWithStatusServiceUnavailable(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusServiceUnavailable())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusServiceUnavailable {
		t.Errorf("WithStatusServiceUnavailable() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusServiceUnavailable)
	}
}

// TestWithStatusGatewayTimeout tests the WithStatusGatewayTimeout function.
func TestWithStatusGatewayTimeout(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusGatewayTimeout())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusGatewayTimeout {
		t.Errorf("WithStatusGatewayTimeout() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusGatewayTimeout)
	}
}

// TestWithStatusHTTPVersionNotSupported tests the
// WithStatusHTTPVersionNotSupported function.
func TestWithStatusHTTPVersionNotSupported(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusHTTPVersionNotSupported())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusHTTPVersionNotSupported {
		t.Errorf("WithStatusHTTPVersionNotSupported() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusHTTPVersionNotSupported)
	}
}

// TestWithStatusVariantAlsoNegotiates tests the
// WithStatusVariantAlsoNegotiates function.
func TestWithStatusVariantAlsoNegotiates(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusVariantAlsoNegotiates())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusVariantAlsoNegotiates {
		t.Errorf("WithStatusVariantAlsoNegotiates() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusVariantAlsoNegotiates)
	}
}

// TestWithStatusInsufficientStorage tests the
// WithStatusInsufficientStorage function.
func TestWithStatusInsufficientStorage(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusInsufficientStorage())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusInsufficientStorage {
		t.Errorf("WithStatusInsufficientStorage() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusInsufficientStorage)
	}
}

// TestWithStatusLoopDetected tests the WithStatusLoopDetected function.
func TestWithStatusLoopDetected(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusLoopDetected())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusLoopDetected {
		t.Errorf("WithStatusLoopDetected() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusLoopDetected)
	}
}

// TestWithStatusNotExtended tests the WithStatusNotExtended function.
func TestWithStatusNotExtended(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusNotExtended())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusNotExtended {
		t.Errorf("WithStatusNotExtended() did not set the correct "+
			"status code: got %v, want %v",
			w.Code, StatusNotExtended)
	}
}

// TestWithStatusNetworkAuthenticationRequired tests the
// WithStatusNetworkAuthenticationRequired function.
func TestWithStatusNetworkAuthenticationRequired(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, WithStatusNetworkAuthenticationRequired())

	resp.httpWriter.WriteHeader(resp.statusCode)

	if w.Code != StatusNetworkAuthenticationRequired {
		t.Errorf("WithStatusNetworkAuthenticationRequired() did "+
			"not set the correct status code: got %v, want %v",
			w.Code, StatusNetworkAuthenticationRequired)
	}
}

// TestAddContentType tests the AddContentType function.
func TestAddContentType(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddContentType("application/json"), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("AddContentType() did not set the correct "+
			"Content-Type header: got %v, want %v",
			contentType, "application/json")
	}
}

// TestAddContentType_ExistingHeader tests the AddContentType function
// with an existing Content-Type header.
func TestAddContentType_ExistingHeader(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		WithHeader("Content-Type", "text/html"),
		AddContentType("application/json"), // update the Content-Type header
		WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("AddContentType() did not set the correct "+
			"Content-Type header: got %v, want %v",
			contentType, "application/json")
	}
}

// TestAddContentType_EmptyValue tests the AddContentType function
// with an empty value.
func TestAddContentType_EmptyValue(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddContentType(""), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	contentType := w.Header().Get("Content-Type")
	if contentType != "" {
		t.Errorf("AddContentType() did not set the correct "+
			"Content-Type header: got %v, want %v",
			contentType, "")
	}
}

// TestAddETag tests the AddETag function.
func TestAddETag(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddETag("123456"), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	eTag := w.Header().Get("ETag")
	if eTag != "123456" {
		t.Errorf("AddETag() did not set the correct ETag header: "+
			"got %v, want %v", eTag, "123456")
	}
}

// TestAddLastModified tests the AddLastModified function.
func TestAddLastModified(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddLastModified(time.Now()), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	lastModified := w.Header().Get("Last-Modified")
	if lastModified == "" {
		t.Errorf("AddLastModified() did not set the Last-Modified header")
	}
}

// TestAddContentLength tests the AddContentLength function.
func TestAddContentLength(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddContentLength(1024), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	contentLength := w.Header().Get("Content-Length")
	if contentLength != "1024" {
		t.Errorf("AddContentLength() did not set the correct "+
			"Content-Length header: got %v, want %v", contentLength, "1024")
	}
}

// TestAddUserAgent tests the AddUserAgent function.
func TestAddUserAgent(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddUserAgent("test-agent"), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	userAgent := w.Header().Get("User-Agent")
	if userAgent != "test-agent" {
		t.Errorf("AddUserAgent() did not set the correct User-Agent header: "+
			"got %v, want %v", userAgent, "test-agent")
	}
}

// TestAddHost tests the AddHost function.
func TestAddHost(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddHost("example.com"), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	host := w.Header().Get("Host")
	if host != "example.com" {
		t.Errorf("AddHost() did not set the correct Host header: "+
			"got %v, want %v", host, "example.com")
	}
}

// TestAddReferer tests the AddReferer function.
func TestAddReferer(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddReferer("http://example.com"), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	referer := w.Header().Get("Referer")
	if referer != "http://example.com" {
		t.Errorf("AddReferer() did not set the correct Referer header: "+
			"got %v, want %v", referer, "http://example.com")
	}
}

// TestAddServer tests the AddServer function.
func TestAddServer(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddServer("test-server"), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	server := w.Header().Get("Server")
	if server != "test-server" {
		t.Errorf("AddServer() did not set the correct Server header: "+
			"got %v, want %v", server, "test-server")
	}
}

// TestAddDate tests the AddDate function.
func TestAddDate(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddDate(time.Now()), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	date := w.Header().Get("Date")
	if date == "" {
		t.Errorf("AddDate() did not set the Date header")
	}
}

// TestAddLocation tests the AddLocation function.
func TestAddLocation(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddLocation("/path/to/resource"), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	location := w.Header().Get("Location")
	if location != "/path/to/resource" {
		t.Errorf("AddLocation() did not set the correct Location header: "+
			"got %v, want %v", location, "/path/to/resource")
	}
}

// TestAddLocation_EmptyValue tests the AddLocation function
// with an empty value.
func TestAddLocation_EmptyValue(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddLocation(""), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	location := w.Header().Get("Location")
	if location != "" {
		t.Errorf("AddLocation() did not set the correct Location header: "+
			"got %v, want %v", location, "")
	}
}

// TestAddRetryAfter tests the AddRetryAfter function.
func TestAddRetryAfter(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddRetryAfter(time.Now().Add(5*time.Minute)),
		WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	retryAfter := w.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Errorf("AddRetryAfter() did not set the Retry-After header")
	}
}

// TestAddRetryAfter_Seconds tests the AddRetryAfter function
// with a time.Duration value.
func TestAddRetryAfter_Seconds(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddRetryAfter(5*time.Minute),
		WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	retryAfter := w.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Errorf("AddRetryAfter() did not set the Retry-After header")
	}
}

// TestAddRetryAfter_Int tests the AddRetryAfter function
// with a time.Duration value.
func TestAddRetryAfter_Int(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddRetryAfter(5), WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	retryAfter := w.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Errorf("AddRetryAfter() did not set the Retry-After header")
	}
}

// TestAddContentDisposition tests the AddContentDisposition function.
func TestAddContentDisposition(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddContentDisposition("attachment", "example.txt"),
		WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	want := `attachment; filename="example.txt"`
	contentDisposition := w.Header().Get("Content-Disposition")
	if contentDisposition != want {
		t.Errorf("AddContentDisposition() did not set the correct "+
			"Content-Disposition header: got %v, want %v",
			contentDisposition, want)
	}
}

// TestAddContentDisposition_UTF8 tests the AddContentDisposition function
// with a UTF-8 filename.
func TestAddContentDisposition_UTF8(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddContentDisposition("attachment", "ロシア人はテロリストだ.txt", true),
		WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	want := `attachment; filename*=UTF-8''%E3%83%AD%E3%82%B7%E3%82%A2%E4%BA` +
		`%BA%E3%81%AF%E3%83%86%E3%83%AD%E3%83%AA%E3%82%B9%E3%83%88%E3%81%A0.txt`
	contentDisposition := w.Header().Get("Content-Disposition")
	if contentDisposition != want {
		t.Errorf("AddContentDisposition() did not set the correct "+
			"Content-Disposition header: got %v, want %v",
			contentDisposition, want)
	}
}

// TestAddContentEncoding tests the AddContentEncoding function.
func TestAddContentEncoding(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddContentEncoding("gzip"),
		WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	contentEncoding := w.Header().Get("Content-Encoding")
	if contentEncoding != "gzip" {
		t.Errorf("AddContentEncoding() did not set the correct "+
			"Content-Encoding header: got %v, want %v",
			contentEncoding, "gzip")
	}
}

// TestAddContentLanguage tests the AddContentLanguage function.
func TestAddContentLanguage(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddContentLanguage("en"),
		WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	contentLanguage := w.Header().Get("Content-Language")
	if contentLanguage != "en" {
		t.Errorf("AddContentLanguage() did not set the correct "+
			"Content-Language header: got %v, want %v",
			contentLanguage, "en")
	}
}

// TestAddContentLocation tests the AddContentLocation function.
func TestAddContentLocation(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddContentLocation("/path/to/resource"),
		WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	contentLocation := w.Header().Get("Content-Location")
	if contentLocation != "/path/to/resource" {
		t.Errorf("AddContentLocation() did not set the correct "+
			"Content-Location header: got %v, want %v",
			contentLocation, "/path/to/resource")
	}
}

// TestAddWWWAuthenticate tests the AddWWWAuthenticate function.
func TestAddWWWAuthenticate(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddWWWAuthenticate("challenge0", "challenge1"),
		WithStatusUnauthorized())

	resp.httpWriter.WriteHeader(resp.statusCode)

	got := w.Header().Values(HeaderWWWAuthenticate)
	want := []string{"challenge0", "challenge1"}

	// Перевіряємо, чи збігаються зрізи
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddWWWAuthenticate() did not set the correct "+
			"WWW-Authenticate headers: got %v, want %v", got, want)
	}
}

// TestAddAuthorization tests the AddAuthorization function.
func TestAddAuthorization(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddAuthorization("Bearer eyJhbGciOiJIUzI.eyJzdWIiOiIx5IyfQ.SflKxwRJSM"),
		WithStatusOK())

	resp.httpWriter.WriteHeader(resp.statusCode)

	Authorization := w.Header().Get("Authorization")
	if Authorization == "" {
		t.Errorf("AddAuthorization() did not set the Authorization header")
	}
}

// TestAddProxyAuthenticate tests the AddProxyAuthenticate function.
func TestAddProxyAuthenticate(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddProxyAuthenticate("challenge0", "challenge1"),
		WithStatusUnauthorized())

	resp.httpWriter.WriteHeader(resp.statusCode)

	got := w.Header().Values(HeaderProxyAuthenticate)
	want := []string{"challenge0", "challenge1"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddProxyAuthenticate() did not set the correct "+
			"Proxy-Authenticate headers: got %v, want %v", got, want)
	}
}

// TestAddProxyAuthorizationSingleValue tests the AddProxyAuthorization
// function with a single value.
func TestAddProxyAuthorizationSingleValue(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddProxyAuthorization("Basic QWxhZGRjpPcGVuYW1l"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := []string{"Basic QWxhZGRjpPcGVuYW1l"}
	got := w.Header().Values(HeaderProxyAuthorization)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddProxyAuthorization() single value test failed: "+
			"got %v, want %v", got, want)
	}
}

// TestAddProxyAuthorizationMultipleValues tests the AddProxyAuthorization
// function with multiple values.
func TestAddProxyAuthorizationMultipleValues(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddProxyAuthorization("Basic QWxhZGRjpPcGVuYW1l", "Bearer sometoken"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := []string{"Basic QWxhZGRjpPcGVuYW1l", "Bearer sometoken"}
	got := w.Header().Values(HeaderProxyAuthorization)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddProxyAuthorization() multiple values test failed: "+
			"got %v, want %v", got, want)
	}
}

// TestAddIfMatchSingleValue tests the AddIfMatch function
// with a single value.
func TestAddIfMatchSingleValue(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddIfMatch(`"123456"`))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := []string{`"123456"`}
	got := w.Header().Values("If-Match")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddIfMatch() single value test failed: "+
			"got %v, want %v", got, want)
	}
}

// TestAddIfMatchMultipleValues tests the AddIfMatch function
// with multiple values.
func TestAddIfMatchMultipleValues(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddIfMatch(`"123456"`, `"abcdef"`))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := []string{`"123456"`, `"abcdef"`}
	got := w.Header().Values("If-Match")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddIfMatch() multiple values test failed: "+
			"got %v, want %v", got, want)
	}
}

// TestAddIfNoneMatchSingleValue tests the AddIfNoneMatch function
// with a single value.
func TestAddIfNoneMatchSingleValue(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddIfNoneMatch(`"123456"`))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := []string{`"123456"`}
	got := w.Header().Values("If-None-Match")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddIfNoneMatch() single value test failed: "+
			"got %v, want %v", got, want)
	}
}

// TestAddIfNoneMatchMultipleValues tests the AddIfNoneMatch function
// with multiple values.
func TestAddIfNoneMatchMultipleValues(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddIfNoneMatch(`"123456"`, `"abcdef"`))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := []string{`"123456"`, `"abcdef"`}
	got := w.Header().Values("If-None-Match")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddIfNoneMatch() multiple values test failed: "+
			"got %v, want %v", got, want)
	}
}

// TestAddIfModifiedSince tests the AddIfModifiedSince function.
func TestAddIfModifiedSince(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddIfModifiedSince(time.Now()))

	resp.httpWriter.WriteHeader(http.StatusOK)

	ifModifiedSince := w.Header().Get("If-Modified-Since")
	if ifModifiedSince == "" {
		t.Errorf("AddIfModifiedSince() did not set the " +
			"If-Modified-Since header")
	}
}

// TestAddIfUnmodifiedSince tests the AddIfUnmodifiedSince function.
func TestAddIfUnmodifiedSince(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddIfUnmodifiedSince(time.Now()))

	resp.httpWriter.WriteHeader(http.StatusOK)

	ifUnmodifiedSince := w.Header().Get("If-Unmodified-Since")
	if ifUnmodifiedSince == "" {
		t.Errorf("AddIfUnmodifiedSince() did not set the " +
			"If-Unmodified-Since header")
	}
}

// TestAddIfRange tests the AddIfRange function.
func TestAddIfRange(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddIfRange("123456"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	ifRange := w.Header().Get("If-Range")
	if ifRange != "123456" {
		t.Errorf("AddIfRange() did not set the correct If-Range header: "+
			"got %v, want %v", ifRange, "123456")
	}
}

// TestAddContentSecurityPolicy tests the AddContentSecurityPolicy function.
func TestAddContentSecurityPolicy(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddContentSecurityPolicy("default-src 'self'"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	contentSecurityPolicy := w.Header().Get("Content-Security-Policy")
	if contentSecurityPolicy != "default-src 'self'" {
		t.Errorf("AddContentSecurityPolicy() did not set the correct "+
			"Content-Security-Policy header: got %v, want %v",
			contentSecurityPolicy, "default-src 'self'")
	}
}

// TestAddContentSecurityPolicyReportOnly tests the
// AddContentSecurityPolicyReportOnly function.
func TestAddContentSecurityPolicyReportOnly(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddContentSecurityPolicyReportOnly("default-src 'self'"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	got := w.Header().Get("Content-Security-Policy-Report-Only")
	if got != "default-src 'self'" {
		t.Errorf("AddContentSecurityPolicyReportOnly() did not "+
			"set the correct Content-Security-Policy-Report-Only header: "+
			"got %v, want %v", got, "default-src 'self'")
	}
}

// TestAddStrictTransportSecurity tests the AddStrictTransportSecurity function.
func TestAddStrictTransportSecurity(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddStrictTransportSecurity(31536000, true, true))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "max-age=31536000; includeSubDomains; preload"
	got := w.Header().Get("Strict-Transport-Security")
	if got != want {
		t.Errorf("AddStrictTransportSecurity() did not set the correct "+
			"Strict-Transport-Security header: got %v, want %v", got, want)
	}
}

// TestAddReferrerPolicy tests the AddReferrerPolicy function.
func TestAddReferrerPolicy(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddReferrerPolicy("no-referrer"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "no-referrer"
	got := w.Header().Get("Referrer-Policy")
	if got != want {
		t.Errorf("AddReferrerPolicy() did not set the correct "+
			"Referrer-Policy header: got %v, want %v", got, want)
	}
}

// TestAddUpgradeInsecureRequests_AsString tests the
// AddUpgradeInsecureRequests function.
func TestAddUpgradeInsecureRequests_AsString(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddUpgradeInsecureRequests("1"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "1"
	got := w.Header().Get("Upgrade-Insecure-Requests")
	if got != want {
		t.Errorf("AddUpgradeInsecureRequests() did not set the correct "+
			"Upgrade-Insecure-Requests header: got %v, want %v", got, want)
	}
}

// TestAddUpgradeInsecureRequests_AsInt tests the
// AddUpgradeInsecureRequests function.
func TestAddUpgradeInsecureRequests_AsInt(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddUpgradeInsecureRequests(1))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "1"
	got := w.Header().Get("Upgrade-Insecure-Requests")
	if got != want {
		t.Errorf("AddUpgradeInsecureRequests() did not set the correct "+
			"Upgrade-Insecure-Requests header: got %v, want %v", got, want)
	}
}

// TestAddUpgradeInsecureRequests_AsBool tests the
// AddUpgradeInsecureRequests function.
func TestAddUpgradeInsecureRequests_AsBool(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddUpgradeInsecureRequests(true))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "1"
	got := w.Header().Get("Upgrade-Insecure-Requests")
	if got != want {
		t.Errorf("AddUpgradeInsecureRequests() did not set the correct "+
			"Upgrade-Insecure-Requests header: got %v, want %v", got, want)
	}
}

// TestAddXContentTypeOptions tests the AddXContentTypeOptions function.
func TestAddXContentTypeOptions(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddXContentTypeOptions("nosniff"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "nosniff"
	got := w.Header().Get(HeaderXContentTypeOptions)
	if got != want {
		t.Errorf("AddXContentTypeOptions() did not set the correct "+
			"X-Content-Type-Options header: got %v, want %v", got, want)
	}
}

// TestAddXFrameOptions tests the AddXFrameOptions function.
func TestAddXFrameOptions(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddXFrameOptions("DENY"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "DENY"
	got := w.Header().Get(HeaderXFrameOptions)
	if got != want {
		t.Errorf("AddXFrameOptions() did not set the correct "+
			"X-Frame-Options header: got %v, want %v", got, want)
	}
}

// TestAddXXSSProtection tests the AddXXSSProtection function.
func TestAddXXSSProtection(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w,
		AddXXSSProtection("1; mode=block"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "1; mode=block"
	got := w.Header().Get(HeaderXXSSProtection)
	if got != want {
		t.Errorf("AddXXSSProtection() did not set the correct "+
			"X-XSS-Protection header: got %v, want %v", got, want)
	}
}

// TestAddContentDPR tests the AddContentDPR function.
func TestAddContentDPR(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddContentDPR(2.0))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "2"
	got := w.Header().Get(HeaderContentDPR)
	if got != want {
		t.Errorf("AddContentDPR() did not set the correct "+
			"Content-DPR header: got %v, want %v", got, want)
	}
}

// TestAddDPR tests the AddDPR function.
func TestAddDPR(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddDPR(2.0))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "2"
	got := w.Header().Get(HeaderDPR)
	if got != want {
		t.Errorf("AddDPR() did not set the correct DPR header: "+
			"got %v, want %v", got, want)
	}
}

// TestAddViewportWidth tests the AddViewportWidth function.
func TestAddViewportWidth(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddViewportWidth(1024))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "1024"
	got := w.Header().Get(HeaderViewportWidth)
	if got != want {
		t.Errorf("AddViewportWidth() did not set the correct "+
			"Viewport-Width header: got %v, want %v", got, want)
	}
}

// TestAddWidth tests the AddWidth function.
func TestAddWidth(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddWidth(1024))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "1024"
	got := w.Header().Get(HeaderWidth)
	if got != want {
		t.Errorf("AddWidth() did not set the correct Width header: "+
			"got %v, want %v", got, want)
	}
}

// TestAddContentRange tests the AddContentRange function.
func TestAddContentRange(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddContentRange(0, 100, 200))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "bytes 0-100/200"
	got := w.Header().Get(HeaderContentRange)
	if got != want {
		t.Errorf("AddContentRange() did not set the correct "+
			"Content-Range header: got %v, want %v", got, want)
	}
}

// TestAddAccept tests the AddAccept function.
func TestAddAccept(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAccept("application/json"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "application/json"
	got := w.Header().Get(HeaderAccept)
	if got != want {
		t.Errorf("AddAccept() did not set the correct Accept header: "+
			"got %v, want %v", got, want)
	}
}

// TestAddAcceptCharset tests the AddAcceptCharset function.
func TestAddAcceptCharset(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAcceptCharset("utf-8"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "utf-8"
	got := w.Header().Get(HeaderAcceptCharset)
	if got != want {
		t.Errorf("AddAcceptCharset() did not set the correct "+
			"Accept-Charset header: got %v, want %v", got, want)
	}
}

// TestAddAcceptEncoding tests the AddAcceptEncoding function.
func TestAddAcceptEncoding(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAcceptEncoding("gzip"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "gzip"
	got := w.Header().Get(HeaderAcceptEncoding)
	if got != want {
		t.Errorf("AddAcceptEncoding() did not set the correct "+
			"Accept-Encoding header: got %v, want %v", got, want)
	}
}

// TestAddAcceptLanguage tests the AddAcceptLanguage function.
func TestAddAcceptLanguage(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAcceptLanguage("en"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "en"
	got := w.Header().Get(HeaderAcceptLanguage)
	if got != want {
		t.Errorf("AddAcceptLanguage() did not set the correct "+
			"Accept-Language header: got %v, want %v", got, want)
	}
}

// TestAddCacheControl tests the AddCacheControl function.
func TestAddCacheControl(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddCacheControl("no-cache"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "no-cache"
	got := w.Header().Get(HeaderCacheControl)
	if got != want {
		t.Errorf("AddCacheControl() did not set the correct "+
			"Cache-Control header: got %v, want %v", got, want)
	}
}

// TestAddPragma tests the AddPragma function.
func TestAddPragma(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddPragma("no-cache"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "no-cache"
	got := w.Header().Get(HeaderPragma)
	if got != want {
		t.Errorf("AddPragma() did not set the correct Pragma header: "+
			"got %v, want %v", got, want)
	}
}

// TestAddWarning tests the AddWarning function.
func TestAddWarning(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddWarning(WarningHeader{
		Code: 110,
		Text: "Response is stale",
		// Agent: "Server",
		// Date:  time.Now(),
	}))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := `110 "Response is stale"`
	got := w.Header().Get(HeaderWarning)
	if got != want {
		t.Errorf("AddWarning() did not set the correct Warning header: "+
			"got %v, want %v", got, want)
	}
}

// TestAddWarningWithDateAndAgent tests the AddWarning function
// with a date and agent.
func TestAddWarningWithDateAndAgent(t *testing.T) {
	w := httptest.NewRecorder()
	warningDate := time.Date(2022, time.March, 25, 0, 0, 0, 0, time.UTC)
	warning := WarningHeader{
		Code:  299,
		Agent: "TestAgent",
		Text:  "Deprecated Feature",
		Date:  warningDate,
	}

	resp := NewResponse(w, AddWarning(warning))
	resp.httpWriter.WriteHeader(http.StatusOK)

	expectedDateStr := warningDate.Format(time.RFC1123)
	expectedValue := fmt.Sprintf("%d %s \"%s\" \"%s\"",
		warning.Code, warning.Agent, warning.Text, expectedDateStr)

	got := w.Header().Get("Warning")
	if got != expectedValue {
		t.Errorf("AddWarning() = %v, want %v", got, expectedValue)
	}
}

// TestAddVary tests the AddVary function.
func TestAddVary(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddVary("Accept-Encoding"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "Accept-Encoding"
	got := w.Header().Get(HeaderVary)
	if got != want {
		t.Errorf("AddVary() did not set the correct Vary header: "+
			"got %v, want %v", got, want)
	}
}

// TestAddConnection tests the AddConnection function.
func TestAddConnection(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddConnection("close"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "close"
	got := w.Header().Get(HeaderConnection)
	if got != want {
		t.Errorf("AddConnection() did not set the correct Connection header: "+
			"got %v, want %v", got, want)
	}
}

// TestAddTransferEncoding tests the AddTransferEncoding function.
func TestAddTransferEncoding(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddTransferEncoding("chunked"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "chunked"
	got := w.Header().Get(HeaderTransferEncoding)
	if got != want {
		t.Errorf("AddTransferEncoding() did not set the correct "+
			"Transfer-Encoding header: got %v, want %v", got, want)
	}
}

// TestAddAccessControlAllowHeaders tests the AddAccessControlAllowHeaders
// function.
func TestAddAccessControlAllowHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAccessControlAllowHeaders("Content-Type"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "Content-Type"
	got := w.Header().Get(HeaderAccessControlAllowHeaders)
	if got != want {
		t.Errorf("AddAccessControlAllowHeaders() did not set the correct "+
			"Access-Control-Allow-Headers header: got %v, want %v", got, want)
	}
}

// TestAddAccessControlAllowMethods tests the AddAccessControlAllowMethods
// function.
func TestAddAccessControlAllowMethods(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAccessControlAllowMethods("GET"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "GET"
	got := w.Header().Get(HeaderAccessControlAllowMethods)
	if got != want {
		t.Errorf("AddAccessControlAllowMethods() did not set the correct "+
			"Access-Control-Allow-Methods header: got %v, want %v", got, want)
	}
}

// TestAddAccessControlExposeHeaders tests the AddAccessControlExposeHeaders
// function.
func TestAddAccessControlExposeHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAccessControlExposeHeaders("Content-Type"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "Content-Type"
	got := w.Header().Get(HeaderAccessControlExposeHeaders)
	if got != want {
		t.Errorf("AddAccessControlExposeHeaders() did not set the correct "+
			"Access-Control-Expose-Headers header: got %v, want %v", got, want)
	}
}

// TestAddLink tests the AddLink function.
func TestAddLink(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddLink(LinkHeader{
		URI: "https://example.com/page",
		Rel: "canonical",
	}))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := `<https://example.com/page>; rel="canonical"`
	got := w.Header().Get(HeaderLink)
	if got != want {
		t.Errorf("AddLink() did not set the correct Link header: "+
			"got %v, want %v", got, want)
	}
}

// TestAddLinkWithType tests the AddLink function with a type.
func TestAddLinkWithType(t *testing.T) {
	w := httptest.NewRecorder()
	link := LinkHeader{
		URI:  "https://example.com",
		Rel:  "stylesheet",
		Type: "text/css",
	}

	resp := NewResponse(w, AddLink(link))
	resp.httpWriter.WriteHeader(http.StatusOK)

	expectedValue := `<https://example.com>; rel="stylesheet"; type="text/css"`
	got := w.Header().Get("Link")
	if got != expectedValue {
		t.Errorf("AddLink() with type = %v, want %v", got, expectedValue)
	}
}

// TestAddLinkWithTitle tests the AddLink function with a title.
func TestAddLinkWithTitle(t *testing.T) {
	w := httptest.NewRecorder()
	link := LinkHeader{
		URI:   "https://example.com/help",
		Rel:   "help",
		Title: "Help Page",
	}

	resp := NewResponse(w, AddLink(link))
	resp.httpWriter.WriteHeader(http.StatusOK)

	expectedValue := `<https://example.com/help>; rel="help"; title="Help Page"`
	got := w.Header().Get("Link")
	if got != expectedValue {
		t.Errorf("AddLink() with title = %v, want %v", got, expectedValue)
	}
}

// TestAddAccessControlAllowCredentials tests the
// AddAccessControlAllowCredentials function.
func TestAddAccessControlAllowCredentials(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAccessControlAllowCredentials(true))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "true"
	got := w.Header().Get(HeaderAccessControlAllowCredentials)
	if got != want {
		t.Errorf("AddAccessControlAllowCredentials() did not set the correct "+
			"Access-Control-Allow-Credentials header:"+
			" got %v, want %v", got, want)
	}
}

// TestAddAccessControlAllowOrigin tests the AddAccessControlAllowOrigin
// function.
func TestAddAccessControlAllowOrigin(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAccessControlAllowOrigin("https://example.com"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "https://example.com"
	got := w.Header().Get(HeaderAccessControlAllowOrigin)
	if got != want {
		t.Errorf("AddAccessControlAllowOrigin() did not set the correct "+
			"Access-Control-Allow-Origin header: got %v, want %v", got, want)
	}
}

// TestAddAccessControlRequestHeaders tests the AddAccessControlRequestHeaders
// function.
func TestAddAccessControlRequestHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAccessControlRequestHeaders("Content-Type"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "Content-Type"
	got := w.Header().Get(HeaderAccessControlRequestHeaders)
	if got != want {
		t.Errorf("AddAccessControlRequestHeaders() did not set the correct "+
			"Access-Control-Request-Headers header: got %v, want %v", got, want)
	}
}

// TestAddAccessControlRequestMethod tests the AddAccessControlRequestMethod
// function.
func TestAddAccessControlRequestMethod(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddAccessControlRequestMethod("GET"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "GET"
	got := w.Header().Get(HeaderAccessControlRequestMethod)
	if got != want {
		t.Errorf("AddAccessControlRequestMethod() did not set the correct "+
			"Access-Control-Request-Method header: got %v, want %v", got, want)
	}
}

// TestAddOrigin tests the AddOrigin function.
func TestAddOrigin(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AddOrigin("https://example.com"))

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := "https://example.com"
	got := w.Header().Get(HeaderOrigin)
	if got != want {
		t.Errorf("AddOrigin() did not set the correct Origin header: "+
			"got %v, want %v", got, want)
	}
}

// TestAsTextXML tests the AsTextXML function.
func TestAsTextXML(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsTextXML())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMETextXML
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsTextXML() did not set the correct Content-Type header: "+
			"got %v, want %v", contentType, want)
	}
}

// TestAsTextHTML tests the AsTextHTML function.
func TestAsTextHTML(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsTextHTML())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMETextHTML
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsTextHTML() did not set the correct Content-Type header: "+
			"got %v, want %v", contentType, want)
	}
}

// TestAsTextPlain tests the AsTextPlain function.
func TestAsTextPlain(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsTextPlain())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMETextPlain
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsTextPlain() did not set the correct Content-Type header: "+
			"got %v, want %v", contentType, want)
	}
}

// TestAsTextJavaScript tests the AsTextJavaScript function.
func TestAsTextJavaScript(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsTextJavaScript())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMETextJavaScript
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsTextJavaScript() did not set the "+
			"correct Content-Type header: "+
			"got %v, want %v", contentType, want)
	}
}

// TestAsApplicationXML tests the AsApplicationXML function.
func TestAsApplicationXML(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsApplicationXML())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMEApplicationXML
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsApplicationXML() did not set the "+
			"correct Content-Type header: "+
			"got %v, want %v", contentType, want)
	}
}

// TestAsApplicationJSON tests the AsApplicationJSON function.
func TestAsApplicationJSON(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsApplicationJSON())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMEApplicationJSON
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsApplicationJSON() did not set the "+
			"correct Content-Type header: "+
			"got %v, want %v", contentType, want)
	}
}

// TestAsApplicationJavaScript tests the AsApplicationJavaScript function.
func TestAsApplicationJavaScript(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsApplicationJavaScript())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMEApplicationJavaScript
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsApplicationJavaScript() did not set the "+
			"correct Content-Type header: "+
			"got %v, want %v", contentType, want)
	}
}

// TestAsApplicationForm tests the AsApplicationForm function.
func TestAsApplicationForm(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsApplicationForm())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMEApplicationForm
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsApplicationForm() did not set the "+
			"correct Content-Type header: "+
			"got %v, want %v", contentType, want)
	}
}

// TestAsOctetStream tests the AsOctetStream function.
func TestAsOctetStream(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsOctetStream())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMEOctetStream
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsOctetStream() did not set the "+
			"correct Content-Type header: "+
			"got %v, want %v", contentType, want)
	}
}

// TestAsMultipartForm tests the AsMultipartForm function.
func TestAsMultipartForm(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsMultipartForm())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMEMultipartForm
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsMultipartForm() did not set the "+
			"correct Content-Type header: "+
			"got %v, want %v", contentType, want)
	}
}

// TestAsTextXMLCharsetUTF8 tests the AsTextXMLCharsetUTF8 function.
func TestAsTextXMLCharsetUTF8(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsTextXMLCharsetUTF8())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMETextXMLCharsetUTF8
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsTextXMLCharsetUTF8() did not set the correct "+
			"Content-Type header: got %v, want %v", contentType, want)
	}
}

// TestAsTextHTMLCharsetUTF8 tests the AsTextHTMLCharsetUTF8 function.
func TestAsTextHTMLCharsetUTF8(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsTextHTMLCharsetUTF8())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMETextHTMLCharsetUTF8
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsTextHTMLCharsetUTF8() did not set the correct "+
			"Content-Type header: got %v, want %v", contentType, want)
	}
}

// TestAsTextPlainCharsetUTF8 tests the AsTextPlainCharsetUTF8 function.
func TestAsTextPlainCharsetUTF8(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsTextPlainCharsetUTF8())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMETextPlainCharsetUTF8
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsTextPlainCharsetUTF8() did not set the correct "+
			"Content-Type header: got %v, want %v", contentType, want)
	}
}

// TestAsTextJavaScriptCharsetUTF8 tests the
// AsTextJavaScriptCharsetUTF8 function.
func TestAsTextJavaScriptCharsetUTF8(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsTextJavaScriptCharsetUTF8())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMETextJavaScriptCharsetUTF8
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsTextJavaScriptCharsetUTF8() did not set the correct "+
			"Content-Type header: got %v, want %v", contentType, want)
	}
}

// TestAsApplicationXMLCharsetUTF8 tests the
// AsApplicationXMLCharsetUTF8 function.
func TestAsApplicationXMLCharsetUTF8(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsApplicationXMLCharsetUTF8())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMEApplicationXMLCharsetUTF8
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsApplicationXMLCharsetUTF8() did not set the correct "+
			"Content-Type header: got %v, want %v", contentType, want)
	}
}

// TestAsApplicationJSONCharsetUTF8 tests the
// AsApplicationJSONCharsetUTF8 function.
func TestAsApplicationJSONCharsetUTF8(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsApplicationJSONCharsetUTF8())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMEApplicationJSONCharsetUTF8
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsApplicationJSONCharsetUTF8() did not set the correct "+
			"Content-Type header: got %v, want %v", contentType, want)
	}
}

// TestAsApplicationJavaScriptCharsetUTF8 tests the
// AsApplicationJavaScriptCharsetUTF8 function.
func TestAsApplicationJavaScriptCharsetUTF8(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w, AsApplicationJavaScriptCharsetUTF8())

	resp.httpWriter.WriteHeader(http.StatusOK)

	want := MIMEApplicationJavaScriptCharsetUTF8
	contentType := w.Header().Get(HeaderContentType)
	if contentType != want {
		t.Errorf("AsApplicationJavaScriptCharsetUTF8() did not "+
			"set the correct Content-Type header: got %v, want %v",
			contentType, want)
	}
}

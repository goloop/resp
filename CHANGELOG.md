# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0]

Module path changed to `github.com/goloop/resp/v2`.

### Added
- `WithStreamJSON()` option to opt into zero-copy JSON streaming.
- `WithErrorCode(int)` option to set a distinct business code in error bodies.
- `XML`, `Blob`, `Created`, `Accepted`, `JSONIndent` response helpers.
- `(*Response).Status()` and `(*Response).Header()` getters.
- `ErrInvalidJSONPCallback` sentinel error.

### Changed
- **Breaking:** `JSON`, `JSONP` and `XML` now buffer the payload and commit the
  status only after a successful encode, so an encoding error no longer leaves a
  half-written success response. Direct streaming is available via
  `WithStreamJSON()`.
- **Breaking:** `Error(w, status, message, ...)` — the first argument is the HTTP
  status, written to the status line and mirrored into the body `code`. It takes
  precedence over `WithStatus`; out-of-range values fall back to `500`.
- **Breaking:** `JSONP` validates the callback name and returns
  `ErrInvalidJSONPCallback` for unsafe input instead of reflecting it verbatim.
- **Breaking:** `Redirect` honours only genuine redirect codes (300, 301, 302,
  303, 307, 308); other codes fall back to `302`.
- `ServeFileAsDownload` now sets `Content-Length`, `X-Content-Type-Options:
  nosniff`, and a strict RFC 5987 `filename*` form for non-ASCII file names.
- Raised the language baseline to `go 1.24` and removed all third-party
  dependencies (the package is now zero-dependency).
- `JSONEncodeFunc` and the public API use `any` instead of `interface{}`.

### Fixed
- Field projection (`OnlyFields`/`ExcludeFields`) no longer panics on
  string-keyed maps whose value type is not `any` (e.g. `map[string]int`).

# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.3.2] - 2026-10-18

Recognize single-value headers regardless of letter case.

## [2.3.1]

Patch release.

### Documentation
- The reference covers `WithErrorSlug` and the `error` body field, which until
  now lived only in the godoc and the changelog.

## [2.3.0]

Minor release: a machine-readable name for an error.

### Added
- `WithErrorSlug` adds an `error` field to the JSON error body, beside the
  numeric code and the human message: `{"code":409,"error":"slug_taken",
  "message":"..."}`. That is the field a client switches on and a translation
  key is looked up by. A numeric code cannot serve either purpose without a
  registry mapping numbers to meanings - which is the thing that stops being
  maintained - and it says nothing in a log line or a failing test either.
  Without it, applications wrap `Error` in their own helper and the package's
  own error path goes unused.
- The field is omitted when unset, so every response that does not name a slug
  is byte for byte what it was before. `WithErrorCode` and the mirroring of the
  HTTP status into `code` are unchanged.

## [2.2.0]

Minor release: `OnlyFields`/`ExcludeFields` now respect `json` tags.

### Changed
- `OnlyFields` and `ExcludeFields` project a struct the way `encoding/json`
  would: a tagged field is emitted under its `json` name (not the Go field
  name), and a field can be selected by either its `json` name or its Go field
  name. This fixes redaction on tagged structs: `ExcludeFields(u,
  "password_hash")` now drops a field declared as `PasswordHash string`
  with tag `json:"password_hash"` instead of leaking it under a `PasswordHash`
  key. Structs without json tags are unaffected (still keyed by Go field name).

### Fixed
- A field tagged `json:"-"` is now never emitted by `OnlyFields`/`ExcludeFields`,
  so a field excluded from JSON serialization can no longer leak through them.

## [2.1.0]

Minor release: correctness and safety fixes. Fully backward compatible.

### Fixed
- `OnlyFields`/`ExcludeFields` no longer panic on a struct with unexported
  fields (they are skipped) or on a nil pointer element within a slice.
- `ServeFile` no longer forces `application/octet-stream`, so `http.ServeFile`
  detects the real content type and files render instead of downloading. A
  caller-set Content-Type is still respected.
- `AddContentDisposition` sanitizes the filename and disposition type, so a
  quote or control character can no longer inject extra header parameters; the
  UTF-8 form now also emits a sanitized ASCII fallback.
- Date headers (`AddLastModified`, `AddDate`, `AddRetryAfter`,
  `AddIfModifiedSince`, `AddIfUnmodifiedSince`, `AddWarning`) are formatted as
  IMF-fixdate in GMT, as required for HTTP dates.
- `Error` with a valid but unnamed status code now returns a non-empty message,
  falling back to `http.StatusText` and then "Unknown Status".

### Documentation
- `WithStreamJSON` no longer claims to affect JSONP (JSONP is always buffered).
- Fixed outdated godoc examples that used the old `Error` argument order.

## [2.0.0]

Module path changed to `github.com/goloop/resp/v2`.

### Added
- `WithStreamJSON()` option to opt into zero-copy JSON streaming.
- `WithErrorCode(int)` option to set a distinct business code in error bodies.
- `XML`, `Blob`, `Created`, `Accepted`, `JSONIndent` response helpers.
- `(*Response).Status()` and `(*Response).Header()` getters.
- `SecurityHeaders()` conservative security-header preset and `WithCSP()`
  option for an explicit Content-Security-Policy.
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

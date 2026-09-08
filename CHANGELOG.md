# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- **Breaking:** `ErrCodeInsufficientPrivileges` changed from `-32401` to `-32504`
  to match the official Payme Merchant API specification. Invalid or missing
  HTTP Basic Auth credentials now return the correct JSON-RPC error code.
- Non-POST HTTP methods now return `-32300` (`ErrCodeInvalidHTTPMethod`) instead
  of `-32600` (`ErrCodeInvalidRequest`), per the Payme Merchant API specification.
- HTTP Basic Auth credential comparison now uses `crypto/subtle.ConstantTimeCompare`
  to prevent timing attacks. All auth failure variants (missing, malformed, wrong
  login, wrong password) produce indistinguishable responses.
- Basic Auth scheme prefix matching is now case-insensitive per RFC 7617.
- `buildResponse` no longer leaks internal error details (SQL, stack traces) in
  JSON-RPC error messages for non-`*rpc.Error` handler errors. Generic handler
  errors now return a generic "Internal error" message with code `-32603`.

### Added

- `rpc.ErrCodeInvalidHTTPMethod` (`-32300`) and `rpc.ErrInvalidHTTPMethod()`
  constructor for non-POST HTTP method errors.
- Contract tests (`contract_test.go`) verifying Merchant API protocol compliance:
  HTTP 200 for all responses, correct JSON-RPC error codes for auth/transport
  failures, null id for pre-dispatch errors, no internal detail leakage, and
  case-insensitive Basic Auth prefix matching.

## [1.0.0] - 2025-01-01

### Added

- Payme Business SDK for Go, initial release
- Merchant API: JSON-RPC 2.0 HTTP handler with `MerchantHandler` interface
  - `CheckPerformTransaction`, `CreateTransaction`, `PerformTransaction`, `CancelTransaction`, `CheckTransaction`, `GetStatement`
  - HTTP Basic Auth validation
  - Error constructors with localized messages (ru, uz, en)
- Subscribe API: HTTP client for Payme Business endpoints
  - Card methods: `CardsCreate`, `CardsGetVerifyCode`, `CardsVerify`, `CardsCheck`, `CardsRemove`
  - Receipt methods: `ReceiptsCreate`, `ReceiptsPay`, `ReceiptsSend`, `ReceiptsCancel`, `ReceiptsCheck`, `ReceiptsGet`, `ReceiptsGetAll`
  - Partial auth (merchantID only) for client-side card methods
  - Full auth (merchantID:key) for server-side methods
- Payment initialization: URL and HTML form generation for Payme checkout
  - `GeneratePaymentURL` with GET method (base64 format)
  - `GeneratePaymentHTMLForm` with POST method
  - Options: callback URL, language, description, detail, test mode
- Shared types: `Account`, `Transaction`, `Receiver`, `TransactionState`, `CancelReason`
- Error handling: `rpc.Error` type with JSON-RPC 2.0 error codes
- Test data constants: test cards (Uzcard, Humo), SMS code, test merchant ID
- Unit tests: 88%+ coverage across all packages
- CI/CD: GitHub Actions for testing and linting, GoReleaser for releases
- Documentation: GoDoc comments, README, examples

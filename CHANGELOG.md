# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-01-01

### Added

- Initial release of the Payme Business SDK for Go.
- **Merchant API**: JSON-RPC 2.0 HTTP handler with `MerchantHandler` interface
  - `CheckPerformTransaction`, `CreateTransaction`, `PerformTransaction`, `CancelTransaction`, `CheckTransaction`, `GetStatement`
  - HTTP Basic Auth validation
  - Pre-defined error constructors with localized messages (ru, uz, en)
- **Subscribe API**: HTTP client for Payme Business endpoints
  - Card methods: `CardsCreate`, `CardsGetVerifyCode`, `CardsVerify`, `CardsCheck`, `CardsRemove`
  - Receipt methods: `ReceiptsCreate`, `ReceiptsPay`, `ReceiptsSend`, `ReceiptsCancel`, `ReceiptsCheck`, `ReceiptsGet`, `ReceiptsGetAll`
  - Partial auth (merchantID only) for client-side card methods
  - Full auth (merchantID:key) for server-side methods
- **Payment Initialization**: URL and HTML form generation for Payme checkout
  - `GeneratePaymentURL` with GET method
  - `GeneratePaymentHTMLForm` with POST method
  - Options: callback URL, language, description, detail, test mode
- **Shared types**: `Account`, `Transaction`, `Receiver`, `TransactionState`, `CancelReason`
- **Error handling**: `rpc.Error` type with JSON-RPC 2.0 error codes
- **Unit tests**: 88%+ coverage across all packages
- **CI/CD**: GitHub Actions for testing and linting, GoReleaser for releases
- **Documentation**: Full GoDoc comments, README, examples

package rpc

import "fmt"

// Error represents a JSON-RPC 2.0 error object returned by the Payme API.
//
// The Message field contains localized error messages keyed by language code
// ("ru", "uz", "en"). The Data field provides additional context about the error,
// typically the name of the field that caused the error.
type Error struct {
	Code    int               `json:"code"`
	Message map[string]string `json:"message"`
	Data    interface{}       `json:"data,omitempty"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	if msg, ok := e.Message["en"]; ok {
		return fmt.Sprintf("payme error [%d]: %s", e.Code, msg)
	}
	if msg, ok := e.Message["ru"]; ok {
		return fmt.Sprintf("payme error [%d]: %s", e.Code, msg)
	}
	return fmt.Sprintf("payme error [%d]", e.Code)
}

// NewError creates a new RPC Error with the given code, messages, and optional data.
func NewError(code int, message map[string]string, data interface{}) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewErrorSimple creates a new RPC Error with a single English message and optional data.
func NewErrorSimple(code int, enMessage string, data interface{}) *Error {
	return &Error{
		Code: code,
		Message: map[string]string{
			"en": enMessage,
		},
		Data: data,
	}
}

// Common JSON-RPC 2.0 error codes.
const (
	// ErrCodeParseError indicates invalid JSON was received.
	ErrCodeParseError = -32700

	// ErrCodeInvalidRequest indicates the JSON sent is not a valid Request object.
	ErrCodeInvalidRequest = -32600

	// ErrCodeMethodNotFound indicates the method does not exist or is not available.
	ErrCodeMethodNotFound = -32601

	// ErrCodeInvalidParams indicates invalid method parameters.
	ErrCodeInvalidParams = -32602

	// ErrCodeInternal indicates an internal JSON-RPC error.
	ErrCodeInternal = -32603

	// ErrCodeHTTPError indicates a non-200 HTTP status was returned.
	ErrCodeHTTPError = -32400

	// ErrCodeInsufficientPrivileges indicates insufficient privileges to perform the operation.
	ErrCodeInsufficientPrivileges = -32401

	// ErrCodeAccountNotFound indicates the account was not found (range -31050 to -31099).
	ErrCodeAccountNotFound = -31050
)

// Pre-defined error constructors for common JSON-RPC errors.

// ErrParseError returns a parse error (-32700).
func ErrParseError() *Error {
	return NewErrorSimple(ErrCodeParseError, "Parse error", nil)
}

// ErrInvalidRequest returns an invalid request error (-32600).
func ErrInvalidRequest() *Error {
	return NewErrorSimple(ErrCodeInvalidRequest, "Invalid Request", nil)
}

// ErrMethodNotFound returns a method not found error (-32601).
func ErrMethodNotFound(method string) *Error {
	return NewErrorSimple(ErrCodeMethodNotFound, "Method not found", method)
}

// ErrInvalidParams returns an invalid params error (-32602).
func ErrInvalidParams(data interface{}) *Error {
	return NewErrorSimple(ErrCodeInvalidParams, "Invalid params", data)
}

// ErrInternal returns an internal error (-32603).
func ErrInternal() *Error {
	return NewErrorSimple(ErrCodeInternal, "Internal error", nil)
}

// ErrHTTPError returns an HTTP error (-32400).
func ErrHTTPError() *Error {
	return NewErrorSimple(ErrCodeHTTPError, "HTTP error", nil)
}

// ErrInsufficientPrivileges returns an insufficient privileges error (-32401).
func ErrInsufficientPrivileges() *Error {
	return NewErrorSimple(ErrCodeInsufficientPrivileges, "Insufficient privileges", nil)
}

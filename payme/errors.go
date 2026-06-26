package payme

import "github.com/xusk947/payme-sdk-go/payme/rpc"

// RPCError is an alias for rpc.Error, representing a JSON-RPC 2.0 error object.
//
// This alias is provided for convenience so users can reference payme.RPCError
// without importing the rpc package directly.
type RPCError = rpc.Error

// NewRPCError creates a new RPCError with the given code, messages, and optional data.
func NewRPCError(code int, message map[string]string, data interface{}) *rpc.Error {
	return rpc.NewError(code, message, data)
}

// NewRPCErrorSimple creates a new RPCError with a single English message and optional data.
func NewRPCErrorSimple(code int, enMessage string, data interface{}) *rpc.Error {
	return rpc.NewErrorSimple(code, enMessage, data)
}

// Re-export common error codes from the rpc package.
const (
	ErrCodeParseError             = rpc.ErrCodeParseError
	ErrCodeInvalidRequest         = rpc.ErrCodeInvalidRequest
	ErrCodeMethodNotFound         = rpc.ErrCodeMethodNotFound
	ErrCodeInvalidParams          = rpc.ErrCodeInvalidParams
	ErrCodeInternal               = rpc.ErrCodeInternal
	ErrCodeHTTPError              = rpc.ErrCodeHTTPError
	ErrCodeInsufficientPrivileges = rpc.ErrCodeInsufficientPrivileges
	ErrCodeAccountNotFound        = rpc.ErrCodeAccountNotFound
)

// Re-export error constructors from the rpc package.
var (
	ErrParseError             = rpc.ErrParseError
	ErrInvalidRequest         = rpc.ErrInvalidRequest
	ErrMethodNotFound         = rpc.ErrMethodNotFound
	ErrInvalidParams          = rpc.ErrInvalidParams
	ErrInternal               = rpc.ErrInternal
	ErrHTTPError              = rpc.ErrHTTPError
	ErrInsufficientPrivileges = rpc.ErrInsufficientPrivileges
)
